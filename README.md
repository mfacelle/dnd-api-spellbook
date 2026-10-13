# dnd-api-spellbook

A simple Go command-line application for looking up D&D 2024 spells through the [D&D 5e API](https://www.dnd5eapi.co/).

Stores information about spells when they are retrieved from the API and "writes" them into the spellbook, so future searches don't require an additional API call.


## Usage

### Build
Build the application with Go, from the root of the projecet directory

```sh
go build .
```

### Run
Run the application and enter a spell name when prompted:

```sh
go run .
```

### Commands when Running
This application is an interactive command line interface and displays basic commands.  Current options are:
- `<spell name>` enter a spell name to search for
- `mine` display spells currently written in this spellbook
- `all` displays all possible spells
- `exit` close the spellbook

## Project structure

- `main.go` main application, handles basic input/output
- `dndapi/` handles communication with the REST API and data types retrieved from it
- `spellbook/` handles managing spells and formatting of output strings

The HTTP client times out requests after 10 seconds and caches successful, public GET responses in memory for up to one hour. Responses marked `no-store`, `no-cache`, or `private` are not cached.
