import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 250, // Number of virtual users
  duration: '60s', // Test duration
};

export default function () {
  // Sử dụng __VU (Virtual User index) hoặc __ITER (Iteration index) để tạo chỉ số duy nhất
  const userIndex = __VU; // Hoặc dùng __ITER nếu muốn tăng theo lần lặp
  const code="xx06"
  const email = `${code}${__VU}_${Date.now()}@example.com`;
  const username = `${code}${userIndex}`;
  const password = `123456`;

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'authorization': 'Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJlMzQ3M2UwMS1kYmFkLTQxYzktOTRjMS01N2Q2YjdmODE5MzAiLCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJhZG1pbkBhZG1pbi5jb20iLCJzdWIiOiJlMzQ3M2UwMS1kYmFkLTQxYzktOTRjMS01N2Q2YjdmODE5MzAiLCJleHAiOjE3NDk0NzQyMzksImlhdCI6MTc0OTM4NzgzOX0.80_IU0G4rpdcFCz9gKPLULR6ipvnZm0Sc90NOFLhWF8',
    },
  };

  const payload = JSON.stringify({
    email: email,
    password: password,
    username: username,
  });

  const res = http.post('http://localhost:8080/api/v1/accounts/create', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200 || r.status === 409 || r.status === 401,
    // 'response contains success': (r) => r.json().message === 'Account created' || r.json().code === 'USERNAME_ALREADY_USED',
  });

  sleep(2); // Add a 1-second delay between iterations
}