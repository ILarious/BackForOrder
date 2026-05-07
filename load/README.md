# Load tests

The load test uses k6 and targets the running HTTP API.

Start services first:

```bash
docker network create back-for-order-net
docker compose -f docker-compose.kafka.yml up -d
docker compose up --build -d
```

Run the default load test:

```bash
k6 run load/k6/orders.js
```

Or through Docker:

```bash
make load
```

Override load profile:

```bash
VUS=50 DURATION=2m BASE_URL=http://localhost:8080 k6 run load/k6/orders.js
```

With Docker:

```bash
VUS=50 DURATION=2m BASE_URL=http://localhost:8080 make load
```
