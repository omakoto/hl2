# hl TOML Rule File Syntax

`hl` can load coloring rules from a TOML file via the `-r` flag:

```sh
hl -r rules.toml [OPTIONS] [PATTERN...]
```

## File Structure

A rule file contains an array of rule tables. Each `[[rule]]` block defines one rule, applied in order:

```toml
[[rule]]
pattern = 'ERROR'
color = 'bred'

[[rule]]
pattern = 'WARN'
line_color = '550'
```

## Self-Executable Rule Files

A TOML rule file can double as a shell script using this header trick:

```sh
#!/bin/sh
IGNORE=''''
exec hl -r "$0" "${@}"
'''
```

- The shell executes `hl -r "$0"` (passing the file itself as the rule file) then exits.
- The TOML parser sees `#!/bin/sh` as a comment and `IGNORE=''''...'''` as a multi-line literal string, so the shell header is silently ignored.
- Everything after the closing `'''` is parsed as normal TOML `[[rule]]` blocks.

## Fields Reference

All fields are optional except `pattern`.

| Field | Type | Description |
|---|---|---|
| `pattern` | string | **Required.** PCRE regex to match against each input line. See [Pattern Syntax](#pattern-syntax). |
| `when` | string | Pre-condition pattern. The rule is skipped unless this pattern also matches the line (checked before `pattern`). |
| `color` | string | Color for matched text. If the pattern has no capture groups, colors the entire match; otherwise colors only the captured portions. See [Color Format](#color-format). |
| `line_color` | string | Color applied to the entire line when the pattern matches. |
| `pre_line` | string | A string (typically a single character) repeated to fill the terminal width and printed as a decorative line *before* the matching line. |
| `pre_line_color` | string | Color for `pre_line`. |
| `post_line` | string | Same as `pre_line`, but printed *after* the matching line. |
| `post_line_color` | string | Color for `post_line`. |
| `show` | bool | Force this line to be shown (useful with `-n` / `hide = true` default). |
| `hide` | bool | Suppress this line from output. Cannot be combined with `before` or `after`. |
| `stop` | bool | Stop evaluating further rules for this line once this rule matches. |
| `before` | int | Number of context lines to show before a matching line (overrides the global `-B` value). |
| `after` | int | Number of context lines to show after a matching line (overrides the global `-A` value). |
| `states` | array of strings | States in which this rule is active. Omit (or leave empty) to apply in all states. See [State Machine](#state-machine). |
| `next_state` | string | Transition to this state when this rule matches. |

## Pattern Syntax

Patterns are PCRE regular expressions by default (use `--no-pcre` / `-N` to fall back to Go's regexp engine).

### Pattern Prefix Flags

A pattern may start with a `{...}` prefix containing one or more flag characters:

| Flag | Effect |
|---|---|
| `!` | **Negate**: match lines that do *not* match the rest of the pattern. |
| `#` | **Strip spaces**: unescaped spaces in the pattern are removed before compiling, allowing readable verbose-style regex. Use `\ ` (backslash-space) to include a literal space. |

Flags can be combined in any order:

```toml
pattern = '{!}ERROR'       # lines NOT containing ERROR
pattern = '{#}a  b  c'    # equivalent to pattern 'abc' (spaces stripped)
pattern = '{!#}a  b  c'   # lines not matching 'abc'
pattern = '{#}foo\ bar'   # matches "foo bar" (escaped space preserved)
```

### Capture Groups

- **No capture groups**: `color` is applied to the entire match.
- **One or more capture groups**: `color` is applied only to the captured substrings.

```toml
[[rule]]
pattern = 'level=(ERROR|WARN)'   # captures only the level word
color = 'bred'                   # only "ERROR" or "WARN" is colored red
```

## Color Format

Color strings follow this format (all parts optional, case-insensitive):

```
[ATTRS] [FG-COLOR] [/ BG-COLOR]
```

### Attributes

Any combination of these letters placed before the color:

| Letter | Effect |
|---|---|
| `b` | Bold / intense |
| `i` | Italic |
| `f` | Faint |
| `u` | Underline |
| `s` | Strikethrough |

### Color Values

Three formats are supported for both foreground and background:

**Named colors** (8 standard terminal colors):

```
black  red  green  yellow  blue  magenta  cyan  white
```

**Xterm 216-color** — three digits, each in the range `0`–`5`, representing R, G, B levels:

```
500   # bright red
050   # bright green
005   # bright blue
550   # yellow
555   # bright white
000   # black
```

Each digit maps linearly: `0`→0, `1`→51, `2`→102, `3`→153, `4`→204, `5`→255.

**24-bit true color** — six hex digits `RRGGBB`:

```
ff0000   # red
00ff00   # green
ffffff   # white
444444   # dark gray
```

### Color Examples

```toml
color = 'red'           # named red foreground
color = 'bred'          # bold red
color = 'bured'         # bold + underline red
color = '550'           # xterm yellow foreground
color = '550/200'       # yellow foreground, dark-red background
color = 'b555/500'      # bold bright-white on bright-red
color = '/111'          # no foreground, dark-gray background
color = 'ff0000'        # 24-bit red foreground
color = 'b00ff00/000000' # bold bright-green on black
```

## State Machine

Rules can be conditioned on a named state, and can trigger state transitions. This allows multi-line or context-aware highlighting.

- Processing begins in the **initial state**, which is the empty string `""`.
- A rule with no `states` field (or `states = []`) applies in **all** states.
- A rule with `states = ['']` applies only in the initial empty-string state.
- When a matching rule has `next_state = 'FOO'`, the current state becomes `"FOO"`.
- From that point on, only rules that include `"FOO"` in their `states` list (or rules with no `states` restriction) are evaluated.

### Example: Highlighting fatal log blocks

```toml
# Entering the fatal state: match a Fatal line while in the initial state.
[[rule]]
pattern = '''(?:\d F |\bF[\/\(])'''
states = ['']
pre_line = '#'
pre_line_color = 'bred'
next_state = 'in_fatal'

# While in fatal state: keep coloring fatal lines.
[[rule]]
pattern = '''(?:\d F |\bF[\/\(])'''
states = ['in_fatal']
line_color = 'bred/550'
stop = true

# While in fatal state: exit when a non-fatal line is seen.
[[rule]]
pattern = '''{!}(?:\d F |\bF[\/\(])'''
states = ['in_fatal']
pre_line = '#'
pre_line_color = 'bred'
next_state = 'back_to_normal'
```

## Rule Evaluation

For each input line, rules are evaluated from top to bottom:

1. If the rule's `states` list does not include the current state, the rule is skipped.
2. If the rule has a `when` field, the line must match it for the rule to proceed.
3. If the rule's `pattern` does not match the line, the rule is skipped.
4. If the rule matches:
   - Colors are recorded.
   - `show`/`hide` may override the line's visibility.
   - `next_state` updates the current state.
   - If `stop = true`, no further rules are checked for this line.

Multiple rules can match the same line (when `stop` is absent). For overlapping regions, earlier rules' `color` takes visual priority; `line_color` from earlier rules takes priority over later ones; and any `color` (match color) overrides `line_color` within matched regions.

## Complete Example

```toml
# Color ERROR lines red, bold.
[[rule]]
pattern = 'ERROR'
line_color = 'bred'
stop = true

# Color WARN lines yellow.
[[rule]]
pattern = 'WARN'
line_color = '550'
stop = true

# Highlight only the SQL keyword in SQLite log lines.
[[rule]]
pattern = '''(?i)(SELECT|INSERT|UPDATE|DELETE)'''
when = 'SQLiteConnection: execute'
color = 'b055'
line_color = '/022'

# Use verbose regex (spaces stripped) to match a ms timing value.
[[rule]]
pattern = '''{#} ( \d{1,3} \. \d+ \s* ms\b )'''
color = '050'

# Hide noisy debug lines.
[[rule]]
pattern = 'verbose_tag'
hide = true
```
