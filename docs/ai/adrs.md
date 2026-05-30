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
