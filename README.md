Spade is a lightweight CLI to save, manage, and run reusable command shortcuts.

```bash
./spd sit --version
spade version 0.1.0

./spd sit -- --version
sit 0.1.0
```

## Template 
```bash
spd add greet "echo Hello {name} you are {age} years old"

spd greet Batman 21                    # positional
spd greet name=Batman age=21           # named
spd greet Batman age=21                # mixed
spd greet name=Batman 21               # mixed
```