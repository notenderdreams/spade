"""Exercise installer downloads, platform selection, and failure isolation."""

import hashlib
import io
import os
from pathlib import Path
import subprocess
import sys
import tarfile
import tempfile
import unittest

INSTALLER = Path(__file__).resolve().parents[1] / "installer.sh"


class InstallerTest(unittest.TestCase):
    def setUp(self):
        self.directory = tempfile.TemporaryDirectory(prefix="spade-installer-")
        self.addCleanup(self.directory.cleanup)
        self.root = Path(self.directory.name)
        self.fixtures = self.root / "fixtures"
        self.fixtures.mkdir()
        self.bin = self.root / "bin"
        self.bin.mkdir()
        self.destination = self.root / "install with spaces"
        self.env = dict(os.environ, PATH=str(self.bin) + os.pathsep + os.environ["PATH"],
                        HOME=str(self.root), FIXTURES=str(self.fixtures),
                        TEST_OS="Linux", TEST_ARCH="x86_64")
        self.wrapper("uname", '''import os, sys
print(os.environ["TEST_OS" if sys.argv[1] == "-s" else "TEST_ARCH"])
''')
        self.wrapper("curl", '''import os, shutil, sys
from pathlib import Path
args = sys.argv[1:]
if "-w" in args:
    print("https://github.com/notenderdreams/spade/releases/tag/v0.1.1", end="")
else:
    filename = args[-1].rsplit("/", 1)[-1]
    source = Path(os.environ["FIXTURES"]) / filename
    if not source.exists():
        sys.exit(22)
    shutil.copyfile(source, args[args.index("-o") + 1])
''')
        checksums = []
        for operating_system in ("linux", "darwin"):
            for arch in ("amd64", "arm64"):
                name = f"spade_0.1.1_{operating_system}_{arch}.tar.gz"
                payload = b"#!/bin/sh\nprintf 'spade version 0.1.1\\n'\n"
                with tarfile.open(self.fixtures / name, "w:gz") as archive:
                    member = tarfile.TarInfo("spd")
                    member.size = len(payload)
                    member.mode = 0o755
                    archive.addfile(member, io.BytesIO(payload))
                checksum = hashlib.sha256((self.fixtures / name).read_bytes()).hexdigest()
                checksums.append(f"{checksum}  {name}\n")
        (self.fixtures / "spade_0.1.1_checksums.txt").write_text("".join(checksums))

    def wrapper(self, name, source):
        path = self.bin / name
        path.write_text(f"#!{sys.executable}\n" + source)
        path.chmod(0o755)

    def run_installer(self, *args):
        return subprocess.run(["sh", str(INSTALLER), *args], env=self.env,
                              capture_output=True, text=True)

    def test_latest_platforms_and_default_directory(self):
        for operating_system, arch in (("Linux", "x86_64"), ("Linux", "aarch64"),
                                        ("Darwin", "x86_64"), ("Darwin", "arm64")):
            with self.subTest(os=operating_system, arch=arch):
                self.env.update(TEST_OS=operating_system, TEST_ARCH=arch)
                result = self.run_installer()
                self.assertEqual(result.returncode, 0, result.stderr)
                binary = self.root / ".local/bin/spd"
                self.assertTrue(os.access(binary, os.X_OK))
                self.assertEqual(subprocess.check_output([str(binary)], text=True),
                                 "spade version 0.1.1\n")

    def test_pinned_version_and_custom_directory(self):
        result = self.run_installer("--version", "v0.1.1", "--install-dir", str(self.destination))
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertTrue((self.destination / "spd").is_file())

    def test_checksum_failure_preserves_existing_install(self):
        self.destination.mkdir()
        binary = self.destination / "spd"
        binary.write_text("existing executable")
        (self.fixtures / "spade_0.1.1_linux_amd64.tar.gz").write_bytes(b"corrupted archive")
        result = self.run_installer("--install-dir", str(self.destination))
        self.assertNotEqual(result.returncode, 0)
        self.assertIn("Checksum mismatch", result.stderr)
        self.assertEqual(binary.read_text(), "existing executable")
        self.assertEqual(list(self.destination.iterdir()), [binary])

    def test_missing_download_does_not_install(self):
        (self.fixtures / "spade_0.1.1_checksums.txt").unlink()
        result = self.run_installer("--install-dir", str(self.destination))
        self.assertNotEqual(result.returncode, 0)
        self.assertFalse(self.destination.exists())

    def test_unsupported_platform_and_invalid_options(self):
        self.env["TEST_OS"] = "Windows_NT"
        self.assertNotEqual(self.run_installer().returncode, 0)
        self.env["TEST_OS"] = "Linux"
        for args in (("--version", "../invalid"), ("--version",), ("--unknown",)):
            with self.subTest(args=args):
                self.assertNotEqual(self.run_installer(*args).returncode, 0)
        self.assertFalse((self.root / ".local/bin/spd").exists())


if __name__ == "__main__":
    unittest.main()
