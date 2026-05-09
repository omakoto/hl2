[![Build Status](https://travis-ci.org/omakoto/hl2.svg?branch=master)](https://travis-ci.org/omakoto/hl2)
# hl
Highlighter: versatile coloring filter

## Installation

```bash
sudo apt-get install libpcre3-dev
go get -u github.com/omakoto/hl2/src/cmd/hl
```

## Usage

```bash
hl [ -r RULE_TOML ] [OPTIONS] [ FILTER-SPEC... ] <FILE
```

### Examples

- `hl "^Error" @red /var/log/syslog`: Highlights lines starting with "Error" in red.
- `hl "WARN" @yellow "ERROR" @red`: Multiple highlight patterns.
- `hl -n "start" , "end"`: Only show lines between "start" and "end" (inclusive).
- `hl "ERROR" @/red@red`: Highlights "ERROR" with red background, and also colors the whole line red.
- `hl -c "tail -f /var/log/syslog" "^.*ERROR.*$" @red`: Run a command and highlight its output.

### Command Line Arguments

`FILTER-SPEC` is a list of patterns and optional color specifications:

- `PATTERN [ COLOR-SPEC ]`
- `PATTERN [ COLOR-SPEC ] ',' PATTERN [ COLOR-SPEC ]` (Range match)

#### Range Matching
When a range (using `,`) is specified, `hl` will show all lines from the first pattern match to the second pattern match. Range patterns automatically imply `-n` (hide by default).

#### Color Specification (CLI)
In the CLI, colors are specified using the `@` prefix:
`@ [MATCH-COLOR] [ @ [LINE-COLOR] ]`

Example: `@red@b/white` colors the matched text red and the whole line bold with a white background.
See [Color Specification](TOML_SYNTAX.md#color-specification) in `TOML_SYNTAX.md` for details on the color format itself.

### Advanced Regex
Patterns support prefixes in curly braces:
- `{!}`: Negate match.
- `{#}`: Ignore unescaped spaces in the pattern.

For more complex rules and stateful highlighting, use a TOML configuration file with the `-r` option. See [TOML Rule Syntax](TOML_SYNTAX.md) for more information.

## Options

- `-r, --rule FILE`: Specify TOML rule file.
- `-n, --hide`: Hide all lines by default (useful for filtering).
- `-A, --after N`: Show N lines after a match.
- `-B, --before N`: Show N lines before a match.
- `-C, --context N`: Show N lines before and after a match.
- `-i, --ignore-case`: Perform case-insensitive match.
- `-c, --command`: Execute a command and highlight its output.
- `-2, --stderr`: Process stderr as well when using `-c`.
- `-w, --width`: Set terminal width for pre/post lines.
- `-S, --no-skip-marker`: Suppress the `---` skip markers between hidden lines.
- `-h, --help`: Show help.
