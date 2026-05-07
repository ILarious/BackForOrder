.PHONY: test load

test:
	go test ./...

load:
	docker run --rm \
		--network host \
		-e BASE_URL=$${BASE_URL:-http://localhost:8080} \
		-e VUS=$${VUS:-10} \
		-e DURATION=$${DURATION:-30s} \
		-v "$$(pwd)/load:/scripts" \
		grafana/k6:0.49.0 run /scripts/k6/orders.js
