import http from 'k6/http';
import { check } from 'k6';

export const options = {
  vus: 50,
  duration: '10s',
};

export default function () {
  const payload = JSON.stringify({
    name: "Học .NET Core",
    isComplete: false,
  });

  const res = http.post('http://localhost:5278/api/todo', payload, {
    headers: { 'Content-Type': 'application/json' },
  });

  console.log(`Status: ${res.status} - Body: ${res.body}`);
  check(res, {
    'status is 201': (r) => r.status === 201,
    'response has id': (r) => JSON.parse(r.body).id !== undefined,
  });
}
