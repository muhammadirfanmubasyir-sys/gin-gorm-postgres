import http from 'k6/http';
import { check, sleep } from 'k6';
import { Rate, Trend } from 'k6/metrics';

const BASE_URL = __ENV.BASE_URL || 'http://localhost:8080';

const errorRate = new Rate('errors');
const reqDuration = new Trend('req_duration', true);

export const options = {
  scenarios: {
    health: {
      executor: 'constant-vus',
      vus: 10,
      duration: '30s',
      exec: 'healthCheck',
    },
    crud: {
      executor: 'ramping-vus',
      startVUs: 0,
      stages: [
        { duration: '10s', target: 20 },
        { duration: '20s', target: 50 },
        { duration: '10s', target: 0 },
      ],
      exec: 'crudFlow',
    },
  },
  thresholds: {
    http_req_duration: ['p(95)<500', 'p(99)<1000'],
    errors: ['rate<0.1'],
  },
};

export function healthCheck() {
  const res = http.get(`${BASE_URL}/health`);
  check(res, {
    'health status 200': (r) => r.status === 200,
  });
  errorRate.add(res.status !== 200);
  reqDuration.add(res.timings.duration);
}

export function crudFlow() {
  // Create
  const createRes = http.post(
    `${BASE_URL}/api/v1/users`,
    JSON.stringify({
      name: `K6User${__VU}_${__ITER}`,
      email: `k6${__VU}${__ITER}@example.com`,
      password: 'k6pass123',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(createRes, { 'create status 201': (r) => r.status === 201 });
  errorRate.add(createRes.status !== 201);
  reqDuration.add(createRes.timings.duration);

  sleep(0.1);

  // List
  const listRes = http.get(`${BASE_URL}/api/v1/users`);
  check(listRes, { 'list status 200': (r) => r.status === 200 });
  errorRate.add(listRes.status !== 200);
  reqDuration.add(listRes.timings.duration);

  sleep(0.1);

  // Get
  const getRes = http.get(`${BASE_URL}/api/v1/users/1`);
  check(getRes, { 'get status 200': (r) => r.status === 200 });
  errorRate.add(getRes.status !== 200);
  reqDuration.add(getRes.timings.duration);

  sleep(0.1);

  // Update
  const updateRes = http.put(
    `${BASE_URL}/api/v1/users/1`,
    JSON.stringify({
      name: 'K6Updated',
      email: 'k6.updated@example.com',
      password: 'newk6pass123',
    }),
    { headers: { 'Content-Type': 'application/json' } }
  );
  check(updateRes, { 'update status 200': (r) => r.status === 200 });
  errorRate.add(updateRes.status !== 200);
  reqDuration.add(updateRes.timings.duration);

  sleep(0.1);

  // Delete
  const deleteRes = http.del(`${BASE_URL}/api/v1/users/1`);
  check(deleteRes, { 'delete status 200': (r) => r.status === 200 });
  errorRate.add(deleteRes.status !== 200);
  reqDuration.add(deleteRes.timings.duration);
}
