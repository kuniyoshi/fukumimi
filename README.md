# fukumimi

[![Go](https://github.com/kuniyoshi/fukumimi/actions/workflows/go.yml/badge.svg)](https://github.com/kuniyoshi/fukumimi/actions/workflows/go.yml)

`fukumimi` is a CLI tool for retrieving and tracking radio show episodes from a login-protected fan club website.
It helps you manage local read/unread status for each broadcast.

## ✨ Features

- Login form automation (cookie-based session management)
- Fetch and store radio show episodes lists
- Track listened/unlistened status
- Output in Markdown format

## 🚀 Usage

### Initial login

You will be prompted for login credentials on first run.
Once authenticated, session cookies are saved and reused for future commands.

```bash
% fukumimi login
```

### Fetch radio show episodes

Retrieve the radio show episodes from the fan club site and store them locally.

```bash
% fukumimi fetch
[ ] 06/11 [#38](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55)
[ ] 05/28 [#37](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54)
[ ] 05/09 [#36](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53)
...
```

### Merge listened/unlistened

Merge new radio show episodes to local state manage.

```bash
% fukumimi fetch
[ ] 06/11 [#38](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55)
[ ] 05/28 [#37](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54)
[ ] 05/09 [#36](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53)
...
% cat local
[x] 05/28 [#37](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54)
[x] 05/09 [#36](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53)
% fukumimi fetch | fukumimi merge local
[ ] 06/11 [#38](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55)
[x] 05/28 [#37](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54)
[x] 05/09 [#36](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53)
```

## Design

- Built as a single binary CLI tool for portability and simplicity
- Designed to minimize dependencies and require no GUI interaction
- Stores session cookies locally to avoid repeated logins
- Uses browser-like form submission to support login with auto-filled credentials
- Tracks state using plain text (Markdown) for readability and Git-friendly diffs
- Written in Go for fast execution and easy distribution

## Installation

Install using Go:

```bash
go install github.com/kuniyoshi/fukumimi@latest
```

Or clone and build locally:

```bash
git clone https://github.com/kuniyoshi/fukumimi.git
cd fukumimi
make build
```

## Commands

### `fukumimi login`
Authenticate with the fan club website. Credentials are requested interactively and session cookies are stored in `~/.fukumimi_cookies`.

### `fukumimi fetch`
Fetch all radio show episodes from the fan club website. Outputs episode list in markdown format to stdout.

### `fukumimi merge <filename>`
Merge fetched episodes with a local file containing listened status. Preserves `[x]` marks for previously listened episodes.

Options:
- `-r, --replace`: Update the local file in-place instead of outputting to stdout

## Implementation Details

### Cookie Storage
Session cookies are stored in `~/.fukumimi_cookies` as JSON. The login session is reused automatically for subsequent commands.

### Episode Format
Episodes are output in the following format:
```
[ ] MM/DD [#NNN](URL)
[x] MM/DD [#NNN](URL)
```
- `[ ]` indicates an unlistened episode
- `[x]` indicates a listened episode
- Episodes are sorted by episode number (newest first)

### Merge Behavior
The merge command reads new episodes from stdin and preserves the listened status from the local file based on episode numbers.

## Development

```bash
make build      # Build the binary
make test       # Run tests
make fmt        # Format code
make vet        # Run go vet
make check      # Run all quality checks
```

