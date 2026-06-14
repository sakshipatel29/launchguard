# LaunchGuard

![CI](https://github.com/sakshipatel29/launchguard/actions/workflows/ci.yml/badge.svg)

LaunchGuard is a feature flag and progressive delivery platform built in Go.

It helps engineering teams safely release features using deterministic percentage-based rollouts, Redis-backed flag evaluation caching, PostgreSQL persistence, and Kafka-based evaluation event streaming.

## Why LaunchGuard?

Modern software teams need to release features safely without deploying new code for every configuration change. LaunchGuard provides a backend platform that allows developers to control feature rollouts, quickly disable risky features, and stream feature evaluation events for analytics, observability, and reliability workflows.

This project demonstrates backend engineering, distributed systems, caching, event-driven architecture, database persistence, Docker-based infrastructure, and CI automation.

## Features

* Create, update, list, and delete feature flags
* Deterministic percentage-based rollout evaluation
* Environment-based flag configuration
* User-specific rollout bucketing
* PostgreSQL-backed feature flag persistence
* Redis caching for fast flag evaluation
* Kafka event streaming for every flag evaluation
* Go-based REST API
* Docker Compose infrastructure for local development
* GitHub Actions CI for automated testing

## Tech Stack

* Go
* PostgreSQL
* Redis
* Apache Kafka
* Docker Compose
* GitHub Actions
* REST APIs

## Architecture

```text
Client / Application
        |
        v
LaunchGuard Go API
        |
        +--> Redis Cache
        |       |
        |       +--> Fast flag lookup for evaluation
        |
        +--> PostgreSQL
        |       |
        |       +--> Persistent feature flag storage
        |
        +--> Kafka
                |
                +--> feature_flag_evaluations topic
```

## System Flow

### Feature Flag Creation

```text
Client
  |
  v
POST /flags/
  |
  v
LaunchGuard API
  |
  +--> Store flag in PostgreSQL
  |
  +--> Cache flag in Redis
```

### Feature Flag Evaluation

```text
Client
  |
  v
POST /evaluate
  |
  v
LaunchGuard API
  |
  +--> Check Redis cache
  |
  +--> If cache miss, read from PostgreSQL
  |
  +--> Run deterministic rollout algorithm
  |
  +--> Publish evaluation event to Kafka
  |
  v
Return evaluation response
```

## API Endpoints

### Health Check

```http
GET /health
```

Example response:

```json
{
  "status": "ok",
  "service": "launchguard-api",
  "version": "v0.1.0"
}
```

## Feature Flag APIs

### Create Feature Flag

```http
POST /flags/
```

Example request:

```json
{
  "name": "Payment Retry",
  "key": "payment_retry_v2",
  "description": "Controls rollout of payment retry workflow",
  "enabled": true,
  "rollout_percentage": 35,
  "environment": "dev"
}
```

Example response:

```json
{
  "id": "0f663bcd-7368-401f-a3da-256de7936906",
  "name": "Payment Retry",
  "key": "payment_retry_v2",
  "description": "Controls rollout of payment retry workflow",
  "enabled": true,
  "rollout_percentage": 35,
  "environment": "dev",
  "created_at": "2026-06-14T14:56:51Z",
  "updated_at": "2026-06-14T14:56:51Z"
}
```

### List Feature Flags

```http
GET /flags/
```

### Get Feature Flag by ID

```http
GET /flags/{id}
```

### Update Feature Flag

```http
PUT /flags/{id}
```

Example request:

```json
{
  "name": "Payment Retry",
  "description": "Updated rollout percentage for payment retry workflow",
  "enabled": true,
  "rollout_percentage": 50,
  "environment": "dev"
}
```

### Delete Feature Flag

```http
DELETE /flags/{id}
```

## Evaluation API

### Evaluate Feature Flag

```http
POST /evaluate
```

Example request:

```json
{
  "flag_key": "payment_retry_v2",
  "user_id": "user_123",
  "environment": "dev"
}
```

Example response:

```json
{
  "flag_key": "payment_retry_v2",
  "user_id": "user_123",
  "environment": "dev",
  "enabled": true,
  "rollout_percentage": 35,
  "bucket": 18,
  "reason": "user included in rollout"
}
```

## Rollout Logic

LaunchGuard uses deterministic hashing to assign each user to a rollout bucket between 1 and 100.

The same user and same feature flag always produce the same bucket. This makes feature rollout behavior stable across repeated requests.

Example:

```text
flag_key = payment_retry_v2
user_id  = user_123
bucket   = 18
```

If the rollout percentage is 35, the user is included because bucket 18 is within the first 35 percent.

## Kafka Event Streaming

Every feature flag evaluation publishes an event to Kafka.

Kafka topic:

```text
feature_flag_evaluations
```

Example event:

```json
{
  "event_type": "flag_evaluated",
  "flag_key": "payment_retry_v2",
  "user_id": "user_123",
  "environment": "dev",
  "enabled": true,
  "rollout_percentage": 35,
  "bucket": 18,
  "reason": "user included in rollout",
  "timestamp": "2026-06-14T20:00:00Z"
}
```

These events can be used for analytics, audit trails, experimentation dashboards, or automated rollback workflows.

## Local Development Setup

### Prerequisites

Make sure you have the following installed:

* Go
* Docker
* Docker Compose
* Git

### Clone the Repository

```bash
git clone https://github.com/sakshipatel29/launchguard.git
cd launchguard
```

### Start Infrastructure

```bash
docker compose up -d
```

This starts:

* PostgreSQL on port `5433`
* Redis on port `6380`
* Kafka on port `9092`

### Run the API

```bash
go run cmd/api/main.go
```

The API will start on:

```text
http://localhost:8080
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Example Usage

### Create a Feature Flag

```bash
curl -X POST http://localhost:8080/flags/ \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Payment Retry",
    "key": "payment_retry_v2",
    "description": "Controls rollout of payment retry workflow",
    "enabled": true,
    "rollout_percentage": 35,
    "environment": "dev"
  }'
```

### Evaluate a Feature Flag

```bash
curl -X POST http://localhost:8080/evaluate \
  -H "Content-Type: application/json" \
  -d '{
    "flag_key": "payment_retry_v2",
    "user_id": "user_123",
    "environment": "dev"
  }'
```

### Verify Redis Cache

```bash
docker exec -it launchguard-redis redis-cli
```

Inside Redis CLI:

```redis
KEYS *
```

Expected example:

```text
flag:dev:payment_retry_v2
```

Exit Redis:

```redis
exit
```

### Verify Kafka Events

```bash
docker exec -it launchguard-kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic feature_flag_evaluations \
  --from-beginning \
  --timeout-ms 10000
```

Expected output:

```json
{"event_type":"flag_evaluated","flag_key":"payment_retry_v2","user_id":"user_123","environment":"dev","enabled":true,"rollout_percentage":35,"bucket":18,"reason":"user included in rollout","timestamp":"2026-06-14T20:00:00Z"}
```

## Running Tests

```bash
go test ./...
```

## CI/CD

This project uses GitHub Actions to automatically run Go tests on every push and pull request to the `main` branch.

CI workflow:

```text
Checkout repository
Set up Go
Download dependencies
Verify formatting
Run tests
```

## Project Structure

```text
launchguard
├── cmd
│   └── api
│       └── main.go
├── internal
│   ├── cache
│   │   └── redis.go
│   ├── db
│   │   └── postgres.go
│   ├── evaluator
│   │   ├── evaluator.go
│   │   └── evaluator_test.go
│   ├── events
│   │   └── kafka_publisher.go
│   ├── handlers
│   │   ├── flags.go
│   │   └── health.go
│   ├── models
│   │   └── flag.go
│   └── store
│       ├── cached_flag_store.go
│       ├── flag_store.go
│       └── postgres_flag_store.go
├── .github
│   └── workflows
│       └── ci.yml
├── docker-compose.yml
├── go.mod
├── go.sum
└── README.md
```

## Engineering Highlights

* Designed a Go REST API for feature flag management and runtime evaluation
* Implemented deterministic rollout evaluation using hashing-based user bucketing
* Added PostgreSQL persistence for durable feature flag storage
* Integrated Redis caching to reduce repeated database lookups during flag evaluation
* Published evaluation events to Kafka for event-driven analytics and audit workflows
* Added automated unit tests for rollout evaluation logic
* Configured GitHub Actions CI to run tests on every push and pull request

## Future Improvements

* Add JWT authentication
* Add role-based access control
* Add audit log API
* Add React dashboard for managing flags
* Add OpenTelemetry distributed tracing
* Add Prometheus metrics
* Add Grafana dashboard
* Add automated rollback rules based on error rate and latency
* Add Kubernetes deployment manifests
* Add SDKs for Go and Python applications
* Add support for multi-variant experiments
* Add organization and project-level flag grouping
