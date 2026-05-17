# Mentat

Mentat is a simple tool for managing task via markdown files.

It lets you create daily task notes and give you list of the pending task.

![Mentat Recording](./mentat.gif)

## Setup

Create your config:

```sh
mentat config init
```

Mentat asks for your notes vault path and preferred editor. The config lives at:

```text
~/.config/mentat/config.toml
```

## Usage

```sh
mentat dl                 # open today's daily note
mentat dl -d 2026-05-16   # open a specific daily note
mentat dl -m              # open this month's note
mentat dl -y              # open this year's note
mentat status             # view pending tasks
mentat clean              # rename completed daily notes with -X
mentat config show        # show config
mentat config open        # edit config
mentat -h                 # find help
```

Aliases:

```sh
mentat dl    # daily-note
mentat st    # status
mentat cln   # clean
mentat cfg   # config
```

## Task Format

```md
- [ ] Write README #docs
  Add install and usage notes.
  Link related note [[Mentat]]

- [x] Finished task
```

Mentat treats `- [ ]` as pending and `- [x]` / `- [X]` as done. It also reads
tags like `#docs` and linked notes like `[[Mentat]]`.

