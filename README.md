## 目前決定先以庫存時間夠的話思考如何選位

### Server Architecture
1. etcd -> grpc -> http gateway

- **etcd**: Service registration and discovery, load balancing.
- **grpc**: High-performance internal service communication.
- **http gateway**: Handles external HTTP requests and routes them to gRPC services.

- **Resilience**: Implement circuit breaking and degradation strategies.(closed register or other)
- **Rate Limiting**: Use Token Bucket or Sliding Window algorithms to protect services.

### Database & Storage
1. **Redis**: Primary stock management and high-concurrency handling.
2. **PostgreSQL**: Persistent storage for orders and users.
3. **Elasticsearch**: Log aggregation and analysis.
4. **Message Queue (Kafka/RabbitMQ)**: Asynchronous processing for peak shaving (Traffic Shaping).

### SRE & Observability
1. **Prometheus**: Metrics collection.
2. **Grafana**: Visualization dashboards.
3. **Jaeger**: Distributed tracing.

### Deployment
1. Docker
2. Kubernetes (k8s)

## Work flow

1. **Request Phase**: 
   - User initiates a seckill request.
2. **Redis Layer (Atomic Operation via Lua Script)**:
   - **Duplicate Check**: Use `SISMEMBER` (Redis Set) in Lua to atomically check if `user_id` exists in the buyer list. This prevents duplicate purchases more reliably than a Bloom filter (no false positives).
   - **Stock Check**: Verify if stock > 0.
   - **Execution**: If valid, perform `DECR` (stock) and `SADD` (add user to buyer list) atomically.
   - **Result**: 
       - If failed (duplicate or no stock) -> Return Error Page immediately.
       - If success -> Proceed to async processing.
3. **Async Processing (Traffic Shaping)**:
   - Send a message containing `user_id`, `product_id`, and `order_info` to the **Message Queue**.
   - This decouples the high-throughput Redis layer from the database, protecting the DB from traffic spikes.
4. **Database Layer (Consumer)**:
   - Consumer service reads from MQ and writes the order to PostgreSQL.
   - **Consistency Guarantee**: Use PostgreSQL `UNIQUE CONSTRAINT` on `(user_id, activity_id)` as the final safety net to ensure data consistency and idempotency.

### Work flow for log

- **App -> Kafka -> Logstash (或自研 Consumer) -> Elasticsearch**

## Roadmap

### Version 1: Baseline & Correctness
1. Implement basic consistency between Redis and PostgreSQL.
2. Set up the development environment and project structure.
3. **Goal**: Ensure stock is reduced correctly and orders are created without errors, prioritizing correctness over high performance.

### Version 2: Performance Optimization
1. **Lua Scripting**: Move stock deduction and duplicate purchase checks into Redis Lua scripts to ensure atomicity and reduce network round-trips.
2. **Message Queue Integration**: Introduce MQ to handle high concurrency and smooth out traffic spikes (Peak Shaving).
3. **gRPC Optimization**: Tune internal communication performance.

### Version 3: High Availability & Observability
1. **Service Governance**: Use etcd for advanced service discovery and load balancing to enhance gRPC server reliability.
2. **Observability**: Deploy Prometheus and Grafana for real-time metrics (QPS, Latency, Stock levels) and Jaeger for distributed tracing to identify bottlenecks.
