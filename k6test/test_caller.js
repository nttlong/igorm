import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 1, // Number of virtual users
  duration: '30s', // Test duration
};

export default function () {
  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const payload = JSON.stringify({
    action: 'login',
    params: {
      username: 'admin',
      password: '123456',
    },
    tenant: 'test',
    viewId: 'auth',
  });

  const res = http.post('http://localhost:8080/api/v1/callers/call', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.1); // Add a 0.1-second delay between iterations to manage load
}