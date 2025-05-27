# envui

A terminal user interface for viewing and copying environment variables.

>[!NOTE]
> This is not the tool you want. Instead use `env | fzf | pbc` (read as "fuzzy find in env and copy match to clipboard").

## Motivation

I wanted to build a minimal TUI and used this as a learning exercise. Thanks to [tview](https://github.com/rivo/tview) this didn't even take 200 lines of code.

But it's not as practical as using `env | fzf | pbc`.

## Features

- View all environment variables in an interactive list
- Search and filter environment variables in real-time
- Copy environment variables to clipboard
- Parse and display variables from `.env` files
- Vim-style keyboard navigation (j/k/↑/↓)
- Clean, responsive terminal interface

## Installation

```bash
go install github.com/rtzll/envui
```

## Usage

### View system environment variables

```bash
envui
```

### View variables from a .env file

```bash
envui /path/to/.env
```

## Keyboard Shortcuts

### Navigation
- `j` / `↓`: Move down
- `k` / `↑`: Move up

### Actions
- `y`: Copy selected variable to clipboard
- `/`: Enter search mode
- `q` / `Ctrl+C`: Quit

### Search Mode
- `Enter`: Accept search and exit search mode
- `Esc`: Clear search and exit search mode
- Type to filter variables in real-time

## Dependencies

- [tview](https://github.com/rivo/tview) - Terminal UI library
- [tcell](https://github.com/gdamore/tcell) - Terminal handling
- [clipboard](https://github.com/atotto/clipboard) - Clipboard operations
- [godotenv](https://github.com/joho/godotenv) - .env file parsing
