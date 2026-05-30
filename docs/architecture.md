# Architecture

## Overview

goboxd is a Go HTTP service that executes untrusted user code inside isolated temporary workspaces. The service accepts source code through a JSON API, prepares the execution environment, runs the program with resource limits, captures the result, and returns structured execution metadata.

The current implementation supports:

* Python 3
* JavaScript (Node.js)
* C++

## Request Lifecycle

A typical request follows this flow:

```text
POST /run
    ↓
Parse JSON request
    ↓
Validate language
    ↓
Create isolated temporary directory
    ↓
Write source file
    ↓
Build execution command
    ↓
Execute with timeout
    ↓
Capture stdout/stderr
    ↓
Classify result
    ↓
Cleanup temporary directory
    ↓
Return JSON response
```

## Language Registry

Languages are registered through a central registry.

Each language defines:

* Source filename
* Whether the language is compiled
* Runtime command

Examples:

* Python → `main.py` → `python3`
* JavaScript → `main.js` → `node`
* C++ → `main.cpp` → compile with `g++`

The registry allows execution behavior to be configured in one place rather than scattered across request handling code.

## Execution Model

### Interpreted Languages

Python and JavaScript follow:

```text
source
    ↓
interpreter
    ↓
stdout/stderr
```

### Compiled Languages

C++ follows:

```text
source
    ↓
compiler
    ↓
binary
    ↓
execution
    ↓
stdout/stderr
```

Compilation failures are reported separately from runtime failures.

## Isolation Strategy

Each request creates a unique temporary directory using Go filesystem APIs.

The temporary directory contains:

* Source files
* Compiled binaries
* Execution artifacts

The directory is removed after request completion using deferred cleanup.

This prevents collisions between concurrent requests and isolates execution artifacts.

## Result Classification

Execution results are classified into:

* accepted
* wrong answer
* runtime error
* compile error
* time limit exceeded

Status classification occurs after execution completes and is based on process state, output comparison, and timeout state.

## Timeout Handling

Program execution is wrapped with:

```go
context.WithTimeout(...)
```

to prevent infinite execution.

Programs that exceed the configured limit are terminated and classified as:

```text
time limit exceeded
```

## Output Handling

stdout and stderr are captured separately.

This allows the service to distinguish:

* normal program output
* compiler errors
* runtime errors

and return more accurate execution results.

## Future Work

Planned improvements include:

* YAML-based language registration
* nsjail integration
* readiness endpoint (`/readyz`)
* information endpoint (`/info`)
* bounded concurrency controls
* benchmarking and load testing
* per-language resource limits
* security hardening
