import http from 'k6/http';
import { check, sleep } from 'k6';
//import { htmlReport } from "https://raw.githubusercontent.com/benc-uk/k6-reporter/main/dist/bundle.js";

export let options = {
  vus: 200,
  duration: '60s',
  thresholds: {
    'http_req_duration': ['p(95)<500'], // 95th percentile dưới 500ms
    'http_req_failed': ['rate<0.01'],   // Dưới 1% thất bại
  },
  // Thêm timeout để tránh chờ quá lâu
  http: {
    timeout: '10s',
  },
};

export default function () {
  const userIndex = __VU;
  const code = `xx06_${__VU}_${Date.now()}`;
  
  const rolename = `${code}${userIndex}`;
  
  //const AccessToken = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDk4MzQ4NzUsImlhdCI6MTc0OTgzMTI3NSwicm9sZSI6InVzZXIiLCJzY29wZSI6InJlYWQgd3JpdGUiLCJ1c2VySWQiOiIifQ.yMRg60VGd9Gg9utTPQDHFNmXxaSNJktUtE0ZMGwJ9Jk`; //goi 1 request khac de login
  const AccessToken =`eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE3NDk4OTI2NTQsImlhdCI6MTc0OTg4OTA1NCwicm9sZSI6InVzZXIiLCJzY29wZSI6InJlYWQgd3JpdGUiLCJ1c2VySWQiOiIifQ.Z1k1S3qQhrw7BqBzBhJB84g7jjIlmQvVhfjbTPtVtNU`
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'authorization': `Bearer ${AccessToken}`,
    },
  };

  const payload = JSON.stringify({
    args: { code, name: rolename }, // Thêm email và password
    language: "string",
    tenant: "string",
  });

  const res = http.post('http://localhost:8080/api/v1/invoke/create%40unvs.br.auth.roles', payload, params);

  const success = check(res, {
    'status is 200 or 409 or 401': (r) => r.status === 200 || r.status === 409 || r.status === 401,
    
  });

  

  sleep(0.5); // Giảm sleep để tăng thông lượng
}

// // Tùy chọn tạo báo cáo HTML
// export function handleSummary(data) {
//   return {
//     'summary.html': htmlReport(data),
//   };
// }