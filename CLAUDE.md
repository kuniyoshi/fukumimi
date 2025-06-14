# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

fukumimi is a CLI tool for retrieving and tracking radio show episodes from a login-protected fan club website. It manages local read/unread status for each broadcast episode using Markdown format for Git-friendly tracking.

## Development Status

This is a newly initialized project. The implementation has not begun yet, but the design specifies:
- Single binary Go CLI application
- Cookie-based session management
- Local state tracking in Markdown
- Commands: `login`, `fetch`, `merge`

## Architecture Notes

The intended architecture (from README.md):
- Written in Go for portability
- Minimal external dependencies
- Cookie storage for session persistence
- Plain text (Markdown) output for version control compatibility
- No GUI required - pure CLI interaction

## Development Guidelines

Since this is a greenfield Go project:
1. Create `go.mod` when initializing the Go module
2. Follow standard Go project layout conventions
3. Implement the three core commands as described in the README
4. Use standard Go testing practices once code is written