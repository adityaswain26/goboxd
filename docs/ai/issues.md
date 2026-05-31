# Issues Log

## 2026-05-25 : Go build failed because entrypoint directory had no Go files

### Problem
Docker build failed with:
`no Go files in /src/cmd/goboxd`

### Cause
The Dockerfile expected a Go entrypoint at `cmd/goboxd/main.go`, but the directory existed without any `.go` files.

### Investigation
Used:
- `find . -name "*.go"`
- inspected Docker build logs
- checked expected build target path

### Resolution
Created a minimal `main.go` implementing:
- `/healthz`
- basic HTTP server

### Learning 
Understood how Docker build targets depend on project structure and how Go entrypoints are resolved during compilation.

---

## 2026-05-26 : Runtime container could not find pyton3

### Problem
Execution pipeline failed with:
`exec: "python": executable file not found in $PATH`

### Cause
Initially installed `python3` in the builder stage of a multi-stage Dockerfile instead of the runtime Stage.

### Investigation
Verified directly inside container using:
`docker exec -it goboxd sh`
and:
`python3 --version`

this confirmed pyton was missing from the final runtime image.

### Resolution
Added a proper runtime stage:
```dockerfile
From debian:${DEBIAN_VERSION}-slim As runtime```
and installed python3 inside that stage.

### Learning
Learned the difference between:
- builder stage
- runtime stage
and understood that dependencies required during execution must exist in the final runtime image, not only the build image.

---

## 2026-05-26 : Circular dependency error in Docker stages

### Problem
Docker build failed with:
circular depencdency detected on stage: builder

### Cause
Accidentally renamed the runtime stage to builder , causing Docker to attempt copying artifacts from the same stage into itself.

### Investigation
Inspected stage names and COPY instructions:
```dockerfile
COPY --from=builder ...```

### Resolution 
Restored separate runtime stage.

### Learning 
Learned how multi-stage Docker builds reference artifacts between stages using stage aliases.

---

## 2026-05-27 : Infinite loops incorrectly marked as accepted

### Problem
Programs like:
`Python
while True:
	pass`
were being terminated internally but still returned:
`JSON
{
  "status":"accepted"
}`

## Cause
Execution status logic checked output matching before checking timeout state. Since the infinite loop produced no stdout and expected output was an empty string, the request was incorrectly classified as accepted.

## Investigation 
Tested the backend using intentionally infinite Python programs and inspected execution logs:
`signal:killed`
This confirmed timeout termination was working internally.

## Resolution
Added timeout-state detection using:
`Go
ctx.Err() == context.DeadlineExceeded`
befor output comparison logic.

## Learning
Learned that execution classification order matters and process termination state must be evaluted before result comparison.

---

# 2026-05-28 : stdout and stderr mixed together
## Problem
Runtime errors like:
`Python
print(x)`
were being returned as 
`wrong answer`
because runtime errors and normal program output were merged together.

## Cause
The backend intially used:
`Go
CombinedOutput()`
which merged stdout and stderr into one stream.
## Ivestigation 
Observed that Python tracebacks appeared inside the stdout response field.
## Resolution
Replaced combined output handling with separate buffers:
`Go
cmd.Stdout = &stdoutBuf
cmd.Sterr = &stderrBuf`
## Learning 
Learned that process execution produces multiple output streams and proper execution classification depends on separating them.

---

# 2026-05-28 : Go variable shadowing caused sourcePath bug
## Problem
Build failed with:
` declared and not used: sourcePath`
## Cause
Inside:
`Go
if req.Language == "py3"`
I accidentally used:
`Go
:=`
instead of:
`Go
=`
which created a new local variable instead of updating the outer sourcePath.
## Investigation
Inspected Language-selection block and compared variable scope behavior.
## Resolution
Replaced:
`Go
sourcePath :=`
with:
`Go
sourcePath =`
## Learning 
Learned how Go variable shadowing works inside conditional blocks and how := creates new scoped variables.

---

# 2026-05-28 : Integrating compiled-language execution pipeline
## Problems
The backend initially assumed all languages could execute using:
`Go
python3 sourcePath`
which failed conceptually for compiled languages like C++.
## Cause
Execution architecture was still interpreter-specific.
## Invesigation 
Mapped execution flow differeces between:
- interpreted languages
- compiled languages
## Resolution
Implemented separate execution pipeline for C++:
`source -> compile -> binary -> execute`
using:
`Go
g++`
inside the runtime container.
## Learning 
Learned that execution systems require language-specific exection strategie rather than one universal execution path.
---
## 2026-05-29 : Language metadata duplicated across the codebase

### Problem

Language-specific information was spread across multiple sections of the code.

Adding a new language required updating several different conditional blocks.

### Cause

Language metadata such as filenames and execution behavior was embedded directly inside request-processing logic.

### Investigation

While preparing to add JavaScript support, I reviewed the code path for language handling and noticed that language-specific logic existed in multiple places.

### Resolution

Introduced a central language registry:

```go
var languages = map[string]LanguageConfig
```

and moved language metadata into a single location.

### Learning

Learned that configuration duplication creates maintenance overhead and makes language additions more error-prone.

---

## 2026-05-29 : Request handler accumulating too many responsibilities

### Problem

The request handler was responsible for:

* request parsing
* source file preparation
* language selection
* compilation
* execution
* response generation

This made the function increasingly difficult to extend.

### Cause

Execution behavior was implemented directly inside the HTTP handling layer.

### Investigation

While adding JavaScript support, the execution logic became more complex and highlighted the growing responsibility of the request handler.

### Resolution

Extracted execution-command creation into:

```go
buildCommand(...)
```

and delegated language-specific execution setup to this function.

### Learning

Learned that separating HTTP concerns from execution concerns improves readability and maintainability.

---

## 2026-05-29 : Compile failure caused by malformed struct literal

### Problem

The project failed to build with errors such as:

```text
unexpected newline in composite literal
```

and:

```text
syntax error: unexpected ) at end of statement
```

### Cause

While implementing compile-error handling, the response struct contained:

* a missing comma
* a misspelled field name (`Stder`)
* incorrect brace placement

### Investigation

Inspected the compiler output and reviewed the reported source lines.

### Resolution

Corrected the struct literal syntax and matched the response field names with the RunResponse definition.

### Learning

Learned that Go compiler line references are usually very precise and should be inspected before making broader changes.

---

## 2026-05-30 : Service metadata could become inconsistent with supported languages

### Problem

The `/info` endpoint initially contained a manually maintained list of supported languages.

Future language additions could update the registry without updating the endpoint.

### Cause

The supported-language list was duplicated instead of being derived from the execution configuration.

### Investigation

While reviewing the info endpoint, I noticed that language information existed both in the registry and inside the endpoint implementation.

### Resolution

Generated supported-language information directly from the language registry.

### Learning

Learned the value of maintaining a single source of truth for configuration data.

---

# Issues Log

## 2026-05-29 · Unbounded request bodies

### Problem

The service accepted request bodies without any size restrictions.

A client could submit extremely large payloads, causing excessive memory usage before validation or execution began.

### Cause

Request bodies were passed directly to the JSON decoder without any maximum size enforcement.

### Investigation

While reviewing the project specification, I identified request size limits as one of the documented security concerns.

I reviewed the request-processing flow and confirmed that no upper bound existed.

### Resolution

Added request size limiting using:

```go
http.MaxBytesReader(...)
```

before JSON decoding.

### Learning

Learned that input validation includes not only checking content correctness but also controlling resource consumption.

---

## 2026-05-31 : Unbounded child process output

### Problem

The service captured stdout and stderr into memory without any output limits.

Programs that continuously printed output could consume excessive memory and potentially destabilize the service.

### Cause

Output was collected using in-memory buffers with no maximum capacity.

### Investigation

After reviewing the specification's security section, I examined how process output was captured and noticed that output growth was unrestricted.

### Resolution

Implemented bounded output buffers with truncation support.

Output is now capped at a fixed size and marked as truncated when limits are exceeded.

### Learning

Learned that resource limits must apply not only to execution time but also to program output.

---

## 2026-05-31 : Prioritization mismatch between features and evaluation criteria

### Problem

My initial instinct for Day 5 was to continue adding more functionality and languages.

### Cause

I was focusing primarily on visible features rather than reviewing the evaluation criteria and judging priorities.

### Investigation

After rereading the specification, I noticed that security, documentation, testing, and software engineering practices are heavily weighted alongside functionality.

### Resolution

Shifted focus from feature expansion to security hardening and project completeness.

### Learning

Learned that successful project delivery depends on aligning implementation work with evaluation criteria rather than continuously adding features.

---

## 2026-05-31 : Security fixes needed to be documented, not only implemented

### Problem

Several security improvements had already been introduced throughout the project, but they were not being tracked in a structured way.

### Cause

The focus was initially on making the system work rather than documenting why certain design choices were made.

### Investigation

While preparing for submission requirements, I reviewed the list of security issues and compared them against existing implementation decisions.

### Resolution

Started documenting security-related decisions, fixes, and architectural changes through ADRs and project logs.

### Learning

Learned that maintainability includes explaining security decisions so reviewers and future contributors can understand them quickly.

