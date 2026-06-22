there was tcp server
producer sent logs to the server

now fanout is added to process the logs, used bufio instead of conn.Write to handle lines instead of buffer
Total Execution Time (Latency)1 Worker: Took 31.760 seconds total.5 Workers: Took 8.233 seconds total.Improvement: A ~74% reduction in total execution time (nearly a $4\times$ speedup).

how it works:
there are 5 workers waiting for job
when readloop sends job any worker that is free takes that job and does it

when all 5 has job, the channel is blocked, the moment a worker is ready it takes the next job

| Workers | Batches needed    | Expected time | Observed      |
| ------- | ----------------- | ------------- | ------------- |
| 10      | `ceil(10/10) = 1` | ~3 s          | **3.01 s**   |
| 5       | `ceil(10/5) = 2`  | ~6 s          | **6.00 s**   |
| 1       | `ceil(10/1) = 10` | ~30 s         | **30.02 s**  |
