import http from "k6/http";
import { check, sleep } from "k6";

export const options = {
  scenarios: {
    smoke: {
      executor: "constant-vus",
      vus: Number(__ENV.VUS || 10),
      duration: __ENV.DURATION || "30s",
    },
  },
  thresholds: {
    http_req_failed: ["rate<0.01"],
    http_req_duration: ["p(95)<500"],
  },
};

const baseURL = __ENV.BASE_URL || "http://localhost:8080";

export default function () {
  const username = `load_user_${__VU}_${__ITER}`;

  const createResponse = http.post(
    `${baseURL}/orders`,
    JSON.stringify({ username }),
    { headers: { "Content-Type": "application/json" } },
  );
  check(createResponse, {
    "create order status is 201": (response) => response.status === 201,
  });

  const listResponse = http.get(`${baseURL}/orders`);
  check(listResponse, {
    "list orders status is 200": (response) => response.status === 200,
  });

  const healthResponse = http.get(`${baseURL}/health`);
  check(healthResponse, {
    "health status is 200": (response) => response.status === 200,
  });

  sleep(1);
}
