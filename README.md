[![Build Status](https://travis-ci.org/omakoto/hl2.svg?branch=master)](https://travis-ci.org/omakoto/hl2)
# hl
Highlighter: versatile coloring filter

## Installation

```bash
go get -u github.com/omakoto/hl2/src/cmd/hl
```

(`sudo apt-get install libpcre3-dev` is no longer needed)

## Usage

```
hl [-r RULE_TOML] [OPTIONS] [FILTER-SPEC...] < FILE
hl -f [-r RULE_TOML] [OPTIONS] FILE... [, FILTER-SPEC...]
hl -c [-2] [-r RULE_TOML] [OPTIONS] COMMAND [ARG...] [, FILTER-SPEC...]
```

- Default mode reads from stdin.
- `-f` reads one or more files.
- `-c` executes a command and colors its stdout. Add `-2` to also process stderr.

### Filter specs

A `FILTER-SPEC` is one or more of:

```
PATTERN [COLOR-SPEC]
PATTERN [COLOR-SPEC] , PATTERN [COLOR-SPEC]
```

The first form colors lines matching `PATTERN`. The second form (with `,`) defines a **range**: only lines between the two patterns are shown (implies `-n`).

If no `COLOR-SPEC` is given for a pattern, a color is chosen automatically from a built-in palette.

`PATTERN` is a PCRE regular expression. The `{!}` and `{#}` prefix flags also apply — see [Pattern Syntax](TOML_SYNTAX.md#pattern-syntax).

### Color spec (command line)

On the command line, a color spec starts with `@` and optionally includes a second `@` for the line color:

```
@[ATTRS][FG-COLOR][/BG-COLOR][@[ATTRS][LINE-FG-COLOR][/LINE-BG-COLOR]]
```

For the color and attribute format, see [Color Format](TOML_SYNTAX.md#color-format).

Examples:

```sh
# Color "ERROR" in bold red, with a dark-red background on the whole line.
hl 'ERROR' @bred@/200

# Show only lines between "BEGIN" and "END", coloring them cyan.
hl 'BEGIN' @bcyan , 'END' @bcyan < file.log
```

### Options

| Flag | Description |
|---|---|
| `-r FILE` | Load coloring rules from a TOML file. See [TOML Rule Files](TOML_SYNTAX.md). |
| `-n` | Hide all lines by default (show only matching lines). |
| `-i` | Case-insensitive matching. |
| `-A N` | Show N lines of context after each match. |
| `-B N` | Show N lines of context before each match. |
| `-C N` | Shorthand for `-A N -B N`. |
| `-S` | Suppress the `---` skip marker printed between hidden sections. |
| `-w N` | Set terminal width (used for `pre_line`/`post_line` decorations). |
| `-s SEP` | Change the range separator (default: `,`). |
| `-N` | Disable PCRE; use Go's regexp engine instead. |
| `-c` | Treat remaining arguments as a command to execute. |
| `-2` | With `-c`: also process the command's stderr. |
| `-f` | Treat arguments before `,` as input files. |
| `-q` | Suppress the "waiting for stdin" warning. |

## TOML Rule Files

For complex or reusable coloring rules, write a TOML rule file and load it with `-r`:

```sh
hl -r my-rules.toml < input.log
```

See **[TOML_SYNTAX.md](TOML_SYNTAX.md)** for the full syntax reference, including color formats, pattern flags, decorative lines, and the state machine.
