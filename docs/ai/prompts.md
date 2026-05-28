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
```go
defer os.RemoveAll(tempDir)```

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
```Bash
docker exec -it goboxd sh
python3 --version```

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
```Go
if string (output) == req.ExpectedOutput```
and structured JSON responses.
