# ADR-001 : Use isolated temporary directories per execution request

## Status 
Accepted

## Context
The execution service needs to write user-submitted source code to disk before execution.

An early implementation idea was using a shared file path like:
`text
/tmp/main.py`

However, this creates conflicts when multiple requests excute concurrently because submissions can overwrite each other's files.

## Decision
Each execution request will create its own isolated temporary workspace using:
`Go
os.MkdirTemp("","goboxd-*")`
Source files, compiled binaries, and temporary execution artifacts are stored inside this request-specific directory.

## Consequences

### Positive
- avoids file collisions
- supports concurrent execution
- simplifies cleanup
- isolates execution artifacts

### Negative
- additional filesystem operations per request
- temporary directory cleanup becomes necessary

---

# ADR-002 : Use context-based execution timeouts

## Status
Accepted

## Context
User-submitted programs can execute indefinitely:
`python
while True:
	pass`
Without execution limits, the service can hang permanently and exhaust server resources.

## Decision
Use:
`Go
context.WithTimeout(...)`
with:
`Go
exec.CommandContext(...)`
to automatically terminate long-running processes.

## Consequences
### Positive
- prevents infinite execution
- improves service staility
- introduces basic execution resource control
### Negative
- timeout duration currently hardcoded
- long-running legitimate programs may terminate early

# ADR-003 : Separate stdout and stderr handling 

## Status
Accepted 

## Context
Initially the backend used:
`Go
CombinedOutput()`
which merged:
- standard output
- runtime error
into a single stream.

This prevented accurate execution result classification.

Decision
Use separate stdout/stder buffers:
`Go
cmd.Stdout = &stdoutBuf
cmd.Stderr = &stderrBuf`
instead of combined output captue.

## Consequences

### Positive
- enables runtime error detection
- improves execution classification
- produces clearer API responses

### Nagative 
- slightly more execution-handling complexity

# ADR-004 · Centralize language metadata in a registry

## Status

Accepted

## Context

Language-specific information such as source filenames and execution behavior was being determined through conditional branches distributed throughout the codebase.

As support for additional languages increased, this approach became harder to maintain and required modifying multiple code locations whenever a new language was introduced.

## Decision

Introduce a central language registry:

```go
var languages = map[string]LanguageConfig
```

Each language entry stores execution-related metadata in a single location.

## Consequences

### Positive

* single source of truth for language configuration
* easier onboarding of new languages
* reduced duplication
* foundation for future YAML-based language registration

### Negative

* requires additional abstraction compared to direct conditional logic
* language configuration structure will likely evolve as new requirements are added

# ADR-005 · Separate HTTP handling from execution construction

## Status

Accepted

## Context

The request handler was responsible for:

* parsing requests
* selecting language behavior
* building execution commands
* handling compilation
* executing programs

This caused the request handler to accumulate multiple responsibilities.

## Decision

Extract command creation into a dedicated function:

```go
buildCommand(...)
```

The request handler now delegates language-specific execution preparation to this function.

## Consequences

### Positive

* clearer separation of concerns
* simpler request handler
* easier testing of execution behavior
* improved maintainability

### Negative

* introduces an additional abstraction layer
* future language-specific features may require expanding the interface

# ADR-006 · Derive service metadata from the language registry

## Status

Accepted

## Context

The `/info` endpoint initially maintained a manually defined list of supported languages.

This created a risk that service metadata could become inconsistent with the actual language registry.

## Decision

Generate supported-language information directly from the registry at request time.

The `/info` endpoint now derives language information from the same configuration source used by execution logic.

## Consequences

### Positive

* eliminates duplicated configuration
* keeps service metadata synchronized automatically
* supports future language additions with fewer code changes

### Negative

* introduces a small runtime lookup step
* response ordering depends on registry iteration unless explicitly sorted
