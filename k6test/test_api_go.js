import http from 'k6/http';
import { check } from 'k6';

export default function () {
  const payload = JSON.stringify({
    name: "Học Go với Echo",
    isComplete: false
  });

  const params = {
    headers: { 'Content-Type': 'application/json' },
  };

  const res = http.post('http://localhost:8080/api/todo', payload, params);

  console.log(`Status: ${res.status} - Body: ${res.body}`);

  check(res, {
    'status is 201': (r) => r.status === 201,
    'response has id': (r) => {
      if (!r.body) return false;
      try {
        const data = JSON.parse(r.body);
        return data.id !== undefined;
      } catch (e) {
        console.log("JSON parse error", r.body);
        return false;
      }
    },
  });
}
