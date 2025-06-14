# fukumimi

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
[ ] 2025-06-11 #38
[ ] 2025-05-28 #37
[ ] 2025-05-09 #36
...
```

### Merge listened/unlistened

Merge new radio show episodes to local state manage.

```bash
% fukumimi fetch
[ ] 2025-06-11 #38
[ ] 2025-05-28 #37
[ ] 2025-05-09 #36
...
% cat local
[x] 2025-05-28 #37
[x] 2025-05-09 #36
% fukumimi fetch | fukumimi merge local
[ ] 2025-06-11 #38
[x] 2025-05-28 #37
[x] 2025-05-09 #36
```

## Design

- Built as a single binary CLI tool for portability and simplicity
- Designed to minimize dependencies and require no GUI interaction
- Stores session cookies locally to avoid repeated logins
- Uses browser-like form submission to support login with auto-filled credentials
- Tracks state using plain text (Markdown) for readability and Git-friendly diffs
- Written in Go for fast execution and easy distribution
