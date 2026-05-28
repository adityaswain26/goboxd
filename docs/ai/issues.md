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
