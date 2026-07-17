
# Performance Test Report

This document summarizes the performance test results and includes the benchmark screenshots captured during the evaluation.

## Overview

The tests were executed to assess the behavior of the broker under sustained load and to compare observed throughput and stability across different runs.

## Benchmark Screenshots

### 1. Initial observation

![Performance test screenshot 1](images/2026-07-15%20210413.png)

### 2. Runtime behavior

![Performance test screenshot 2](images/2026-07-15%20210222.png)

### 3. Mid-test metrics

![Performance test screenshot 3](images/2026-07-15%20205955.png)

### 4. Stability view

![Performance test screenshot 4](images/2026-07-15%20210413.png)

### 5. Additional throughput view

![Performance test screenshot 5](images/2026-07-15%20210222.png)

### 6. Additional runtime capture

![Performance test screenshot 6](images/2026-07-15%20205955.png)

### 7. Additional benchmark capture

![Performance test screenshot 7](images/2026-07-15%20204726.png)

### 8. Final snapshot

![Performance test screenshot 8](images/2026-07-15%20204505.png)

## Notes

- The images above are included directly from the performance test artifacts.
## Analysis of captures

Summary of observations from the added captures (Grafana + Redis UI):

- Throughput / handled per second: the Grafana panel labeled "Handled per second sum" shows an increase from approximately 200,000 to 260,000. This indicates the test pushed ~200k–260k messages/sec in the observed window.
- Latency (percentiles): the percentiles panel shows P99 around ~10 milliseconds while P95 and P50 remain low and stable. This indicates low and stable latencies under the applied load.
- CPU: low CPU utilization (~1% → 4%).
- Memory: memory usage increases from ~36% to ~42% during the test; memory growth correlates with increased throughput.
- Availability: the availability panel remains near 100% in the observed window.
- Redis stream / consumer group: the Redis UI captures show a `consumer-group` with 3 consumers and `Pending = 0`, plus several entries with IDs and timestamps. These captures correspond to the asynchronous pipeline (SQS → Redis) and show that the Redis queue used by the SQS pipeline had no backlog during the observed window.

Interpretation for the tested scenario (bidirectional gRPC + SQS/Redis):

- Bidirectional (gRPC fan-in / fan-out): the bidirectional flow exercised by `dominus.BrokerAPI.BidirectionalStream` is a pure gRPC fan-in/fan-out pattern and does not use Redis. The performance observed in Grafana (throughput and latencies) reflects the gRPC processing behavior for this pattern.
- SQS / Redis (async): the Redis UI captures correspond to the asynchronous pipeline (SQS producer → Redis consumer group) and show 3 consumers with `Pending = 0`, i.e. no accumulation of pending messages in Redis during the window.
- The "Handled per second" metric suggests the aggregated system (all load sources) reached ~250k msgs/s. It is necessary to confirm whether that value is aggregated across both tests (bidirectional + SQS) or corresponds to a single path.
- Low P99 latencies (~10 ms) indicate the critical paths (gRPC processing and/or enqueue/ack) responded with good headroom during this run.

## Evidence links

- Redis stream - consumer groups and entries (SQS pipeline): images/2026-07-17 153124.png, images/2026-07-17 153117.png
- Grafana dashboard snapshot (CPU / Memory / Handled/sec / Percentiles / Availability): images/2026-07-17 150918.png

## All captures

Embedded below are all images present in `doc/images/` (for convenience and to prevent broken links):

![2026-07-15 204505](images/2026-07-15%20204505.png)
![2026-07-15 204726](images/2026-07-15%20204726.png)
![2026-07-15 205955](images/2026-07-15%20205955.png)
![2026-07-15 210222](images/2026-07-15%20210222.png)
![2026-07-15 210413](images/2026-07-15%20210413.png)
![2026-07-16 220735](images/2026-07-16%20220735.png)
![2026-07-16 220752](images/2026-07-16%20220752.png)
![2026-07-16 235655](images/2026-07-16%20235655.png)
![2026-07-16 235755](images/2026-07-16%20235755.png)
![2026-07-17 000124](images/2026-07-17%20000124.png)
![2026-07-17 150611](images/2026-07-17%20150611.png)
![2026-07-17 150653](images/2026-07-17%20150653.png)
![2026-07-17 150739](images/2026-07-17%20150739.png)
![2026-07-17 150807](images/2026-07-17%20150807.png)
![2026-07-17 150912](images/2026-07-17%20150912.png)
![2026-07-17 150918](images/2026-07-17%20150918.png)
![2026-07-17 153117](images/2026-07-17%20153117.png)
![2026-07-17 153124](images/2026-07-17%20153124.png)

## SLI / SLO (proposal based on observed data)

SLI definitions (what to measure):

- SLI: Effective throughput (messages processed / s), measured as the aggregated value shown in the `Handled per second` panel.
- SLI: Message processing latency (P50, P95, P99) measured end-to-end from ingestion to ack.
- SLI: Availability (percentage of requests or operations processed without error) during the observed window.
- SLI: Consumer lag / pending (number of messages pending in Redis stream / SQS queue). Note: Redis applies to the asynchronous SQS pipeline; the bidirectional gRPC flow does not use Redis.

Proposed SLOs (targets and thresholds, informed by observations):

- Throughput SLO: maintain >= 200,000 messages/second (aggregated) over 1-minute windows under the target load. Compliance target: 95% of windows over a 1-hour evaluation.
- Latency SLO (P99): P99 < 20 ms. Compliance target: 99% of samples in a 5-minute window.
- Latency SLO (P95): P95 < 10 ms. Compliance target: 99.5% of samples.
- Availability SLO: 99.95% successful operations in the evaluation window (alert if below 99.9%).
- Consumer lag SLO: `Pending` in Redis stream = 0 under normal conditions; tolerance: < 5 pending messages per consumer. Compliance target: 99.9% of time.
- SQS-specific SLO (if SQS is in the pipeline): end-to-end processing P99 < 100 ms; P95 < 50 ms. Compliance target: 99% for P99.

Notes on SLOs:

- Thresholds include margin over current observations (e.g. observed P99 ~10 ms → proposed SLO P99 < 20 ms) to allow for normal variation.
- Before publishing SLOs to production, validate them with additional runs (different windows, longer peaks, and fault injection) to avoid noisy alerts.

## Recommendations

- Add dashboards/alerts for: P99 latency, throughput per instance and aggregated, pending messages (Redis stream for SQS), and memory growth trends (alert at 75–80% to prevent OOM).
- Monitor consumer-group size and ack latency to detect backpressure.
- Run additional variations: increase consumer count, simulate consumer failures, and run sustained-load tests to detect memory leaks or other resource trends.
- Confirm whether the ~200k–260k value is aggregated across producers/consumers or per-instance — adjust SLOs to topology.

## Test configuration (ghz)

The tests were executed with `ghz` using the configuration files included in the repository. Key parameters are summarized below with direct links to the YAML files for easy re-run and audit:

- Bidirectional stream test: [tests/integration/workload_test/broker_bidirectionalstream.yml](tests/integration/workload_test/broker_bidirectionalstream.yml)
  - `call`: dominus.BrokerAPI.BidirectionalStream
  - `proto` and `import-paths`: ../dominus-proto-definition/proto/dominus.proto
  - `host`: localhost:5000
  - `data.subscribers`: 3 endpoints (192.168.1.3:5001, :5002, :5003)
  - `rps`: 680
  - Concurrency schedule: step from 20 to 1000, step 20 (step duration 5s)
  - Load schedule: step (start 10 → end 10000)
  - `duration`: 5m, `timeout`: 60s
  - Stream specifics: `stream-call-duration`: 500ms, `stream-call-count`: 20, `stream-interval`: 100ms
  - `cpus`: 2

- SQS producer test: [tests/integration/workload_test/broker_sqs_producer.yml](tests/integration/workload_test/broker_sqs_producer.yml)
  - `call`: dominus.SqsAPI.Producer
  - `proto` and `import-paths`: ../dominus-proto-definition/proto/dominus.proto
  - `host`: localhost:5000
  - `data.payload`: random string template (`{{randomString 8}}`)
  - `rps`: 20000
  - `concurrency`: 1, `total`: 100
  - Load schedule: step (start 20, step 20, end 10000)
## Evidence links

- Redis stream - consumer groups and entries (SQS pipeline): images/2026-07-17 153124.png, images/2026-07-17 153117.png
- Grafana dashboard snapshot (CPU / Memory / Handled/sec / Percentiles / Availability): images/2026-07-17 150918.png
- To reproduce tests locally run `ghz --config <file.yml>` targeting the correct `host` and ensure the proto file is served from the relative path specified.



## Next suggested steps
