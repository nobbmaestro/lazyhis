# User Config

Default path for the global configuration file:

- `~/.config/lazyhis/lazyhis.yml`

## Default

```yml
# Config relating to the database
db:
  # List of regex for excluding commands
  # See https://github.com/nobbmaestro/lazyhis/blob/master/docs/config.md#exclude-commands-from-database
  excludeCommands: []

# Config relating to the GUI
gui:
  # List of GUI Columns
  # See https://github.com/nobbmaestro/lazyhis/blob/master/docs/config.md#gui-columns
  columnLayout:
    - EXIT_CODE
    - EXECUTED_AT
    - COMMAND

  # If true, GUI will display only unique commands
  showUniqueCommands: true

# Config relating to things outside of LazyHis like how sessions are obtain etc
os:
  # Command for retrieving current session
  # See https://github.com/nobbmaestro/lazyhis/blob/master/docs/config.md#custom-session-providers
  fetchCurrentSessionCmd: "tmux display-message -p '#S'"
```

## Exclude Commands from Database

In case you want to exclude commands from being registered into the database, you can use
`excludeCommands` for listing custom regex expressions.

Example:

```yml
db:
  excludeCommands:
    - ^clear
    - ^nvim
    - ^ls\s*$
```

## Custom Column Layout

`LazyHis` stores various metadata associated to the history record. In order to preserve screen real
estate, this metadata is hidden from GUI by default. The visibility and order of the columns can be
easily customized by modifying the `columnLayout` configuration, allowing you to choose which
columns to display and in what order they should appear in the GUI.

Available columns:

| Column      | Description                        |
| ----------- | ---------------------------------- |
| COMMAND     | The command                        |
| EXIT_CODE   | Exit code of the command           |
| EXECUTED_AT | Execution timestamp of the command |
| EXECUTED_IN | Execution duration of the command  |
| PATH        | Path context of the command        |
| SESSION     | Session context of the command     |

## Custom Session Providers

By default, 'LazyHis' assumes sessions to be provided by [tmux](https://github.com/tmux/tmux).
However, this can simply configured by overwriting the `fetchCurrentSessionCmd`.

Example:

```yml
os:
  fetchCurrentSessionCmd: echo 'Hello World'
```
