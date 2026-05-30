## 2026-05-25 : Initial understanding of the system

Right now my understanding is that goboxd is basically a backend service that receives code through an HTTP request, runs it inside a restricted sandbox, captures the output, and returns execution results as JSON.

At the moment I understand the HTTP/API side much more than the sandboxing side. I still do not fully understand how namespaces, cgroups, and nsjail work internally, but I understand that they are used to isolate untrusted code and limit resource usage.

Initially I thought the hackathon was mainly about writing Go code, but after reading the spec more carefully I realized the bigger focus is on system design, security, concurrency, and safe execution under load.

---

## 2026-05-26 : Initial execution model

At the start of Day 2, my understanding of the system was still mostly request/response oriented. I thought the main task was simply reciving code through `/run`, executing it somehow, and returning output.

`text
request -> run code -> return response`

I had not yet thought carefully about:
- isolated execution workspaces
- runtime environments inside containers
- concurrenct execution conflicts
- differences between builder/runtime Docker stages

---

## 2026-05-26 : Understanding isolated execution workspaces
While implementing the execution pipeline, I initially considered writing all submission into a shared path like:
`/tmp/main.py`
After thinking through concurrent requests, I realized multiple executions could overwrite each other's files or interfere with execution state.

This changed my understanding form:
`single execution file`
to 
` per-request isolated execution workspace`
I implemented:
`Go
os.MkdirTemp("","goboxd-*")`
tp create temporary directories for each request.

This was first point where I started thinking about the system as a concurrent execution service instead of a single local script runner.

---

## 2026-05-26 : Understanding runtime Vs builder containers
A major misunderstanding today was assuming that intalling Python in the Docker builder stage would automatically make it available during runtime execution.

After debugging:
`exec: "python3": executable file not found in $PATH`
I realized:
- build stages
- runtime stages
are isolated in multi-stage Docker builds.

This changed my understanding of the architecture significantly. I started viewing the runtime image as the actual executio environment instead of just the final output of the build process.

---

## 2026-05-26 : Evolving understanding of goboxd architecture
At the start of the project, I mostly viewed goboxd as:
`an API that runs code`
After implementing the execution pipeline, I now understand it more as 
`a request-driven execution orchestration system`
where each request involves:
- isolated workspace creation
- source file generation 
- execution environment setup
- proces execution 
- output capture
- cleanup
- structure response generation

I still do not fully understand:
- nsjail internals
- cgroups
- namespace isolation
- secure sandboxing
but I now understand wher those pieces fit into the overall execution lifecycle.

---

## 2026-05-27 : Execution model becoming language-aware

At the start of Day 3, the backend still behaved like a Python-specific execution service. Even though the `/run` API accepted a `language` field, the actual execution pipeline was still hardcoded around:

`text
python3 main.py`

My mental model was still:
`single execution flow`
rather than:
`different execution strategies depending on language type`

### Understanding execution safety 
Today I implemented timeout-based execution control using:
`Go
context.WithTimeout(...)`
Initially, I only thought about whether code execution worked or not.After testing:
`Pyhton
while True:
	pass`
I realized the backend also needs execution limits and process termination logic.

This changed my understanding from:
`"execute submitted code"`
to:
`"execute submitted code safely under constraints"`
I also started understanding that online judges are not only execution systems, but resource-management systems.

### Separating stdout and stderr
Initially I used:
`Go
CombinedOutput()`
which merged:
- normal program output
- runtime errors
into a single stream.

While testing runtime failures like:
`Python
print(x)`
I realized the backend could not properly distinguish:
- wrong answer
- runtime error
because stderr and stdout were mixed together.

This changed my understanding of execution handling signigicantly. I started thinking about execution results as structured process metadata rather than just raw terminal output.

I replaced:
`Go
CombinedOutput()`
with separate stdout/stderr buffers.

### Understanding execution classification
At the start of the project, I thought result handling was mostly:
`output matches expected output`
Today the execution model became more structured.

The backend now classifies:
- accepted
- wrong answer
- runtime error
- compile error
- time limit exceeded
This made me understand that online judges are effectively state classification systems built around process execution.

### Interpreter Vs compiler execution pipelines
The biggest architectural shift today came from adding C++ support.

Before today, I viewed all  execution as: 
`source -> execute`
After implementing C++ support, I understood there are fundamentally different execution models.
Python execution:
`source -> interpreter -> output`
C++ execution:
`source -> compiler -> binary -> execution -> output`
This is changed how I think about the backend architecture. I no longer see the system as a Python runner, but as a multi-language execution system where each language may require its own execution pipeline.
 
