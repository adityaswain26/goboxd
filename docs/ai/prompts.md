## 2026-05-25 : Fixing Docker build failure for Go entrypoint

**Prompt:**
The Docker build failed with "no Go files in /src/cmd/goboxd". I asked for help understanding why the container build was failing and how the Go project structure should look.

**Response summary:**
The issue was that the Dockerfile expected a Go entrypoint in cmd/goboxd/main.go, but the directory was empty. Suggested creating a minimal HTTP server with /healthz.

**What we used / didn't use:**
Used the suggested project structure and minimal HTTP server implementation. Did not copy additional architecture suggestions yet because we only needed the build to succeed first.

## 2026-05-25 · Implementing the /run endpoint

**Prompt:**
Asked for step-by-step guidance on creating a minimal POST /run endpoint in Go that accepts JSON and returns a JSON response.

**Response summary:**
Suggested creating a RunRequest struct, decoding request JSON using encoding/json, and registering the handler with net/http.

**What we used / didn't use:**
Used the JSON decoding and route registration approach. Kept the response intentionally minimal for Day 1 before adding actual code execution.

## 2026-05-25 · Understanding sandbox execution flow

**Prompt:**
Asked for a simplified explanation of how an online judge style sandbox backend works, including what happens between receiving a POST /run request and returning the execution result.

**Response summary:**
Explained the flow as: receive request → validate JSON → create temporary workspace → execute code → capture output → cleanup → return response. Also explained namespaces, cgroups, and nsjail conceptually.

**What we used / didn't use:**
Used the simplified request-flow model to understand the architecture before implementing real execution. Did not yet implement sandboxing or resource isolation.

