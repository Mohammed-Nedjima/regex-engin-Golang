# Regex Engine in Go

A simplified regular expression engine built from scratch in Go, without using the standard library's `regexp` package. Implements the classical pipeline: lexing and parsing a regex pattern into a token tree, converting that tree into a Nondeterministic Finite Automaton (NFA) using Thompson's construction, and traversing the NFA recursively to determine whether a string matches.

Built as a learning project to make automata theory concrete.

---

## How it works

The engine runs in three stages:

**1. Parsing**
The pattern string is scanned left to right and broken into a token tree. Each operator — character class, quantifier, group, alternation — is parsed into a typed token with its relevant metadata (e.g. a repeat token carries its `min` and `max` bounds, a bracket token carries the expanded character set).

**2. NFA construction (Thompson's construction)**
The token tree is recursively converted into an NFA. Each token type maps to a small fragment of states connected by epsilon (ε) transitions. Fragments are composed by chaining their start and end states with ε-transitions, producing a single NFA with one start state and one terminal state.

**3. Matching**
The NFA is traversed using recursive depth-first search. At each step, the algorithm tries character transitions (consuming one input character) and ε-transitions (consuming no input). The string is accepted if any path through the NFA reaches a terminal state after consuming the full input.

---

## Supported syntax

| Syntax | Description |
|---|---|
| `abc` | Literal character match |
| `[a-zA-Z0-9]` | Character class with ranges |
| `\|` | Alternation (either left or right) |
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
```

Run the tests:

```bash
go test ./...
```

Use the engine in your own Go code:

```go
// Parse the pattern into a token tree
ctx := parse("[a-zA-Z][a-zA-Z0-9_.]+@[a-zA-Z0-9]+.[a-zA-Z]{2,}")

// Convert the token tree into an NFA
nfa := toNfa(ctx)

// Check if a string matches
nfa.check("hello@example.com", -1) // true
nfa.check("not-an-email", -1)      // false
```

---

## Known limitations

This is a learning implementation, not a production engine. Current gaps worth noting:

- **No `.` wildcard** — the dot is treated as a literal character
- **No `^` / `$` anchors** — start-of-text and end-of-text boundaries are not supported
- **Single transition followed per state** — the matcher follows only the first character transition available at each state, which can cause it to miss valid matches in patterns with multiple overlapping transitions
- **No visited-state tracking** — certain epsilon cycles could cause infinite recursion; the current test patterns don't trigger this, but it is a known gap
- **No Unicode support** — operates on raw bytes only

The natural next step would be converting the NFA to a DFA using subset construction, which eliminates backtracking entirely and gives O(n) matching time.

---

## Project structure

```
regex-engine/
├── main.go          # Parser, NFA builder, and matching logic
├── email_test.go    # Test cases using an email validation pattern
└── go.mod
```

---

## References
- Thompson, K. (1968). *Programming Techniques: Regular expression search algorithm.* Communications of the ACM.
