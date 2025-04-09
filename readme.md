# DemoProject

This is a minimal terminal-based email application built using [Bubble Tea](https://github.com/charmbracelet/bubbletea) in Go. I configured it to run on a pseudo-backend which is a SQLite database (`emails.db`). Were this to be converted to a real email client, the touch-points for the database (contained within db) would be altered instead to interface with an IMAP server, and carry out the specific logic related to that.

## Features

This is about as simple as it gets for an email client, and would not feel out of place on an ARPANET mainframe
in a dusty basement somewhere in Berkeley.

- Folders: Inbox, Sent, Archive, Deleted, Spam
- Context-specific actions (Delete, Archive, Mark as Spam, etc.)
- Compose emails (saved to the `Sent` folder)
- Basic key bindings for navigation

## Code Structure
1. The code is organized into:
   - `cmd/demoproject/` – the main entry point.
   - `internal/db/` – logic for loading, inserting, updating, deleting emails in SQLite.
   - `internal/model/` – core application data structures and state (folders, current mode, etc.).
   - `internal/ui/` – Bubble Tea update and view logic.

2. When the application starts, it reads `emails.db`, loads the emails into memory, and displays a terminal interface that you navigate with keyboard commands.  

3. You can compose emails (saved to the `Sent` folder), move and delete emails, or refresh the entire list from the database.

## Built With
- **Go** (1.20+)
- **Bubble Tea** (github.com/charmbracelet/bubbletea)
- **SQLite** (github.com/mattn/go-sqlite3)
- CGO (for building go-sqlite3)

## Dockerfile Purpose
The included Dockerfile allows you to build and run the entire project in a container. It:
1. Uses the official `golang:1.20` image.
2. Installs necessary build tools for CGO and SQLite.
3. Copies the project source code into the container.
4. Builds the `demoproject` binary.
5. Defines a default `CMD` to run the compiled binary.

## Local Usage
  
Ensure you have Docker installed and run:
```bash
just run
```
That will spin up a docker container with everything installed, and bind your interactive
terminal to the container. I picked this approach to make it clean, and easy to deploy and interact
with this, without having to fiddle around with versions, and installing anything other than Just.

If you don't want to use just, that's cool too, you can instead copy the commands in Just since they're
POSIX compliant already.

