## 2026-05-25 : Fixing Docker build failure for Go entrypoint

**Prompt:**
The Docker build failed with "no Go files in /src/cmd/goboxd". I asked for help understanding why the container build was failing and how the Go project structure should look.

**Response summary:**
The issue was that the Dockerfile expected a Go entrypoint in cmd/goboxd/main.go, but the directory was empty. Suggested creating a minimal HTTP server with /healthz.

**What we used / didn't use:**
Used the suggested project structure and minimal HTTP server implementation. Did not copy additional architecture suggestions yet because we only needed the build to succeed first.

---

## 2026-05-25 : Implementing the /run endpoint

**Prompt:**
Asked for step-by-step guidance on creating a minimal POST /run endpoint in Go that accepts JSON and returns a JSON response.

**Response summary:**
Suggested creating a RunRequest struct, decoding request JSON using encoding/json, and registering the handler with net/http.

**What we used / didn't use:**
Used the JSON decoding and route registration approach. Kept the response intentionally minimal for Day 1 before adding actual code execution.

---

## 2026-05-25 : Understanding sandbox execution flow

**Prompt:**
Asked for a simplified explanation of how an online judge style sandbox backend works, including what happens between receiving a POST /run request and returning the execution result.

**Response summary:**
Explained the flow as: receive request → validate JSON → create temporary workspace → execute code → capture output → cleanup → return response. Also explained namespaces, cgroups, and nsjail conceptually.

**What we used / didn't use:**
Used the simplified request-flow model to understand the architecture before implementing real execution. Did not yet implement sandboxing or resource isolation.

---

## 2026-05-26 : Building minimal execution pipeline

**Promts:**
Asked for step-by-step guidance to implement a minimal `/run` excution pipeline in Go.

**Response summary:**
Guidance included:
- parsing JSON requests
- creating temporary execution directories
- writing submitted source code into files
- executing python using `exec.Command`
- capturing stdout
- returning JSON responses

**What we used / didn't use:**
Implemented:
- `os.MkdirTemp`
- `os.WriteFile`
- `exec.Command`
- `CombinedOutput`
- JSON response encoding

Adjusted implementation incrementally during degging and integrated cleanup logic separately using:
`go
defer os.RemoveAll(tempDir)`

---

## 2026-05-26 : Debugging runtime container dependency issues

**Prompt:**
Asked why:
`exec: "python3": executable file not found in $PATH`
occurred even after installing Python in Docker.

**Response summary:**
Receiced explanation of:
- multi-stage Docker builds
- difference between builder and runtime stages
- runtime dependency isolation
Guidance included verifying directly inside the running container using:
`Bash
docker exec -it goboxd sh
python3 --version

**What we used / didn't use:**
Used container inspection and corrected runtime stage structure.

---

## 2026-05-26: Output comparison implementation

**Prompt:**
Asked how to compare execution output with expected output and return status responses.

**Response Summary:**
Suggested:
- extending request schema
- comparing stdout with expected output
- returing accepted/wrong answer statuses

**What we used / didn't use:**
Implemented basic string comparison:
`Go
if string (output) == req.ExpectedOutput`
and structured JSON responses.

---

## 2026-05-27 : Adding execution timeouts

**Prompt:**
Asked for step-by-step guidance to prevent user programs from running forever and hanging the execution service.

**Response Summary:**
Suggested using:
`go
context.WithTimeout(...)`
together with:
`Gp
exec.CommandContext(...)`
to automatically terminate long-running processes.

**What we used / didn't use:**
Added a 2-second timeout around program execution and tested using:
`Python
while True:
	pass`

---

## 2026-05-28 : Separating stdout and stderr
**Prompt:**
Asked why runtime errors were being returned as wrong answers and how execution results should be handled.
**Response Summary:**
Explained that:
`Go
CombinedOutput()`
merges stdout and stderr, making it difficult to distinguish normal output from output from runtime failures.
Suggested using separate output buffers.

**What we used / didn't use:**
Replaced:
`Go
CombinedOutput()`
with:
`Go
cmd.Stdout = &stdoutBuf
cmd.Stderr = &stderrBuf`

Runtime errors are now separated from normarl program output.

---

## 2026-05-28 : Runtime error classification
**Prompt**
Asked how execution status should be determined after separating stdout and stderr.

**Response Summary:**
Suggested prioritizing execution states:
`time limit exceeded
-> runtime error
-> accepted
-> wrong answer`
instead of relying only on output comparison.

**What we used / didn't use:**
Added runtime error detection using:
`Go
if err != nil`
and timeout detection using:
`Go
ctx.Err() == context.DeadlineExceeded`

---

## 2026-05-28 : Adding C++ execution support
**Prompt:**
Asked how to extend the execution service beyond Python and support compiled languages.
**Response Summary:**
Explained the difference between:
`interpreter execution`
and 
`compile -> execute`
pipelines.

Suggested:
`write source
-> compile
-> executed binary`
for C++.

**What we used / didn't use:**
Added:
`Go
g++`
to the runtime container and implemented:
`main.cpp
-> g++
-> binary
-> execution`
including compile error handling.
