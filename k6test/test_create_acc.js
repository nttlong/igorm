import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 10, // Number of virtual users
  duration: '30s', // Test duration
};

export default function () {
  const payload = JSON.stringify({
    email: 'test@example.com',
    password: 'testpassword',
    username: 'testuser',
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post('http://localhost:8080/api/v1/accounts/create', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
    'response contains success': (r) => r.json().message === 'Account created',
  });

  sleep(1); // Add a 1-second delay between iterations
}
//k6 run test_create_acc.js