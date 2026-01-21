# Introduction

A proof-of-concept CLI frontend for [Gogent](https://github.com/BrandonIrizarry/gogent).

Gogent is an API for frontends to provide agentic services to
users. This is one such frontend.


## Installation

`go install github.com/BrandonIrizarry/gogent_repl`

## Flags

### dir

The path to the project directory (the *working directory*) the LLM
should provide assistance for. If omitted, a TUI selection widget
appears prompting the user for a recent project directory.

### log

The log level used by the Gogent backend. One of `debug`, `info`,
`warn`, or `error`. Defaults to `error`.
