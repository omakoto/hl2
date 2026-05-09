# HL TOML Rule Syntax

The `hl` command reads highlighting rules from a TOML file. A rule file consists of one or more `[[rule]]` sections.

## Rule Structure

Each rule is defined in a `[[rule]]` block.

```toml
[[rule]]
pattern = 'regex'
when = 'regex'
color = 'color_spec'
line_color = 'color_spec'
pre_line = 'marker'
pre_line_color = 'color_spec'
post_line = 'marker'
post_line_color = 'color_spec'
show = true
hide = true
stop = true
next_state = 'STATE_NAME'
states = ['STATE1', 'STATE2']
after = 5
before = 2
```

### Fields

| Field | Type | Description |
| :--- | :--- | :--- |
| `pattern` | String | Regular expression to match. Supports prefixes like `{!}` for negation and `{#}` to ignore unescaped spaces. |
| `when` | String | (Optional) Secondary regex that must also match for the rule to apply. |
| `color` | String | Color specification for the matched text. If the regex has capturing groups, only the captured parts are colored. |
| `line_color` | String | Color specification for the entire line if the pattern matches. |
| `pre_line` | String | A marker string printed before the line. Repeated to fill the terminal width. |
| `pre_line_color` | String | Color for the `pre_line` marker. |
| `post_line` | String | A marker string printed after the line. Repeated to fill the terminal width. |
| `post_line_color` | String | Color for the `post_line` marker. |
| `show` | Boolean | Forces the line to be displayed. Overrides `hide`. |
| `hide` | Boolean | Hides the line from output. |
| `stop` | Boolean | Stops processing subsequent rules for the current line if this rule matches. |
| `next_state` | String | Transitions the highlighter to this state for subsequent lines. |
| `states` | List of Strings | Limits the rule to specific states. Default state is `INIT`. |
| `after` | Integer | Shows N additional lines following a match (useful with `show`). |
| `before` | Integer | Shows N preceding lines before a match (useful with `show`). |

## Regular Expression Prefixes

Patterns can be prefixed with a special syntax in curly braces:

- `{!}`: **Negate**. The rule matches if the pattern does *not* match.
- `{#}`: **Ignore Spaces**. Removes all unescaped spaces from the pattern (useful for multi-line or formatted regex).
- `{!#}`: Both negation and space-ignoring.

Example: `pattern = '{#} ^ \d+ \s+ ERROR'` is equivalent to `^ \d+\s+ERROR`.

## Color Specification

Colors are specified using: `[attributes] [foreground] [/ [background]]`

### Attributes
- `b`: Intense (Bold)
- `i`: Italic
- `f`: Faint
- `u`: Underline
- `s`: Strike

### Color Formats
1.  **Names**: `black`, `red`, `green`, `yellow`, `blue`, `magenta`, `cyan`, `white`.
2.  **216-color (RGB666)**: Three digits (0-5), e.g., `511` for bright red.
3.  **RGB888 (Hex)**: Six hex digits, e.g., `ff0000`.

Multiple color rules can apply to the same text; colors are layered.

## State Management

`hl` is a stateful highlighter. It starts in `INIT`. 

- A rule with `states = ["MYSTATE"]` only triggers when the current state is `MYSTATE`.
- `next_state = "NEWSTATE"` changes the state for the *next* line.

This allows for complex multi-line highlighting (e.g., highlighting everything between a "START" and "END" tag).

## Processing Logic

For each line:
1.  All rules are checked in the order they appear in the TOML file.
2.  Rules filtered by `states` or `when` are skipped.
3.  If `stop = true` is encountered, further rules are ignored for this line.
4.  If any matching rule has `hide = true`, the line is marked hidden.
5.  If any matching rule has `show = true`, the line is marked shown (overrides `hide`).
6.  If `-n` is used, the default is hidden; otherwise, the default is shown.
7.  `before`/`after` context is applied to shown lines.
8.  A skip marker (`---`) is printed when one or more lines are hidden between shown lines (suppress with `-S`).
