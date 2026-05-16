# Regex Engine in Go

A regex engine built from scratch in Go — no `regexp` package. Parses a pattern into a token tree, builds an NFA using Thompson's construction, then matches input strings by traversing the automaton recursively.

Made this to get a concrete feel for how automata theory works in practice.

---

## How it works

**1. Parsing**
Scans the pattern left to right into a token tree. Each operator gets its own token type with relevant metadata — a repeat token carries `min`/`max` bounds, a bracket token carries the expanded character set, and so on.

**2. NFA construction**
Converts the token tree into an NFA using Thompson's construction. Each token maps to a small state fragment connected by ε-transitions. Fragments are chained together into a single NFA with one start and one terminal state.

**3. Matching**
Traverses the NFA with recursive depth-first search. At each step it tries character transitions (consume one character, advance) and ε-transitions (move for free). Accepts if any path reaches a terminal state after consuming the full input.

---

## Supported syntax

| Syntax | Description |
|---|---|
| `abc` | Literal match |
| `[a-zA-Z0-9]` | Character class with ranges |
| `\|` | Alternation |
| `( )` | Grouping |
| `*` | Zero or more |
| `+` | One or more |
| `?` | Zero or one |
| `{m,n}` | Between m and n repetitions |
| `{m}` | Exactly m repetitions |
| `{m,}` | At least m repetitions |

---

## Getting started

**Requirements:** Go 1.18+

```bash
git clone https://github.com/your-username/regex-engine
cd regex-engine
go test ./...
```

Usage in Go:

```go
ctx := parse("[a-zA-Z][a-zA-Z0-9_.]+@[a-zA-Z0-9]+.[a-zA-Z]{2,}")
nfa := toNfa(ctx)

nfa.check("hello@example.com", -1) // true
nfa.check("not-an-email", -1)      // false
```

---

## Known limitations

- **No `.` wildcard** — treated as a literal
- **No `^` / `$` anchors**
- **Single transition per state** — only the first character transition is followed, which can miss valid matches in some cases
- **No visited-state tracking** — ε-cycles could cause infinite recursion
- **No Unicode** — raw bytes only

---

## What's next

- **Error handling** — the parser panics on bad input right now; replacing that with proper error messages for things like unclosed brackets or invalid ranges like `[z-a]`
- **CLI** — a simple `regex-engine <pattern> <input>` interface so you can use it without writing Go code

---

## References

- Thompson, K. (1968). *Programming Techniques: Regular expression search algorithm.* Communications of the ACM.
