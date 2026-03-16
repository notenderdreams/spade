# Shell Completion

Spade supports shell completion for commands, flags and their aliases.

## Zsh

```bash
echo 'source <(spd completion zsh)' >> ~/.zshrc
source ~/.zshrc
```

## Bash

```bash
echo 'source <(spd completion bash)' >> ~/.bashrc
source ~/.bashrc
```

## Fish
```bash
spd completion fish > ~/.config/fish/completions/spd.fish
```
