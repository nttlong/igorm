import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 200, // Number of virtual users
  duration: '30s', // Test duration
};

export default function () {
  const userIndex = __VU; // Hoặc dùng __ITER nếu muốn tăng theo lần lặp
  const code="xx05"
  const email = `${code}${__VU}_${Date.now()}@example.com`;
  const username = `${code}${userIndex}`;
  const password = `123456`;
  const payload = JSON.stringify({
    "args": ["root","root","2025-01-10T17:00:00.000Z"],
    "language": "string",
    "tenant": "string"
  });

  const params = {
    headers: {
      'Content-Type': 'application/json',
    },
  };

  const res = http.post('http://localhost:8080/api/v1/invoke/login%40unvs.br.auth.users', payload, params);
  check(res, {
    'status is 200': (r) => r.status === 200,
    // SỬA DÒNG NÀY: Thay đổi chuỗi kiểm tra để khớp với tiếng Việt
    //'response message is correct': (r) => r.json().message === 'Đăng nhập thành công', 
  });
//   check(res, {
//     'status is 200': (r) => r.status === 200,
//     'response contains success': (r) => r.json().message === 'Login successful',
//   });

  //sleep(1); // Add a 1-second delay between iterations
}