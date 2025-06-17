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

## Specification

*Login URL*
https://kitoakari-fc.com/slogin.php

*Radio show episodes*

https://kitoakari-fc.com/special_contents/?category_id=4&page=1

The URL has pagenation.

