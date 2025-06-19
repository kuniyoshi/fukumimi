# fukumimi

[![Go](https://github.com/kuniyoshi/fukumimi/actions/workflows/go.yml/badge.svg)](https://github.com/kuniyoshi/fukumimi/actions/workflows/go.yml)

`fukumimi` is a CLI tool for retrieving and tracking radio show episodes from a fan club website.
It helps you manage local read/unread status for each broadcast.

## ✨ Features

- Fetch radio show episodes from public website
- Track listened/unlistened status
- Output in Markdown format
- Disk-based caching for improved performance
- Concurrent URL processing for faster fetching

## 🚀 Usage

### Fetch radio show episodes

Retrieve all radio show episodes from the fan club website.

```bash
% fukumimi fetch
- [ ] [06/11](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55) (#38)
- [ ] [05/28](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54) (#37)
- [ ] [05/09](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53) (#36)
...
```

### Merge listened/unlistened

Merge new radio show episodes to local state manage.

```bash
% fukumimi fetch
- [ ] [06/11](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55) (#38)
- [ ] [05/28](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54) (#37)
- [ ] [05/09](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53) (#36)
...
% cat local
- [x] [05/28](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54) (#37)
- [x] [05/09](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53) (#36)
% fukumimi fetch | fukumimi merge local
- [ ] [06/11](https://kitoakari-fc.com/special_contents/?contents_id=1&id=55) (#38)
- [x] [05/28](https://kitoakari-fc.com/special_contents/?contents_id=1&id=54) (#37)
- [x] [05/09](https://kitoakari-fc.com/special_contents/?contents_id=1&id=53) (#36)
```

## Design

- Built as a single binary CLI tool for portability and simplicity
- Designed to minimize dependencies and require no GUI interaction
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

### `fukumimi fetch`
Fetch all radio show episodes from the fan club website. Outputs episode list in markdown format to stdout.

Options:
- `-i, --ignore-cache`: Ignore cache and fetch fresh data (still updates cache)

### `fukumimi merge <filename>`
Merge fetched episodes with a local file containing listened status. Preserves `[x]` marks for previously listened episodes.

Options:
- `-r, --replace`: Update the local file in-place instead of outputting to stdout

## Implementation Details

### Episode Format
Episodes are output in the following format:
```
- [ ] [MM/DD](URL) (#NNN)
- [x] [MM/DD](URL) (#NNN)
```
- `[ ]` indicates an unlistened episode
- `[x]` indicates a listened episode
- Episodes are sorted by episode number (newest first)

### Merge Behavior
The merge command reads new episodes from stdin and preserves the listened status from the local file based on episode numbers.

### Caching
fukumimi uses disk-based caching to improve performance and reduce load on the fan club website:
- Cache files are stored in `.fukumimi-cache/` directory
- Cache persists indefinitely until manually cleared
- Use `--ignore-cache` flag to bypass cache reading while still updating it
- Cache is automatically created and managed

## Development

```bash
make build      # Build the binary
make test       # Run tests
make fmt        # Format code
make vet        # Run go vet
make check      # Run all quality checks
```

