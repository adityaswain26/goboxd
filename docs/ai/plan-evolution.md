## 2026-05-25 · Initial understanding of the system

Right now my understanding is that goboxd is basically a backend service that receives code through an HTTP request, runs it inside a restricted sandbox, captures the output, and returns execution results as JSON.

At the moment I understand the HTTP/API side much more than the sandboxing side. I still do not fully understand how namespaces, cgroups, and nsjail work internally, but I understand that they are used to isolate untrusted code and limit resource usage.

Initially I thought the hackathon was mainly about writing Go code, but after reading the spec more carefully I realized the bigger focus is on system design, security, concurrency, and safe execution under load.
