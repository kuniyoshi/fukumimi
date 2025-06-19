# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

fukumimi is a CLI tool for retrieving and tracking radio show episodes from a login-protected fan club website. It manages local read/unread status for each broadcast episode using Markdown format for Git-friendly tracking.

## Development Status

The project now has basic CLI structure with:
- Cobra-based command framework
- Login command implemented with cookie persistence
- HTTP client with session management

Remaining commands to implement: `fetch`, `merge`

## Architecture Notes

The intended architecture (from README.md):
- Written in Go for portability
- Minimal external dependencies
- Cookie storage for session persistence
- Plain text (Markdown) output for version control compatibility
- No GUI required - pure CLI interaction

## Development Guidelines

1. Follow standard Go project layout conventions
2. Commands are in `cmd/` directory
3. Internal packages in `internal/` directory
4. Use `make build` to build the binary
5. Use `make run-login` to test the login command
6. Cookie files are stored in user's home directory as `.fukumimi_cookies`
7. After finished an work, run go fmt to tidy up code base

## Build and Test Commands

```bash
make build      # Build the binary
make run        # Run the application
make test       # Run tests (when available)
make fmt        # Format code
```