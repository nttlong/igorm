import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 200, // 200 virtual users
  duration: '30s', // Test duration of 30 seconds
};

export default function () {
    
  
    const token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJlMzQ3M2UwMS1kYmFkLTQxYzktOTRjMS01N2Q2YjdmODE5MzAiLCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJhZG1pbkBhZG1pbi5jb20iLCJzdWIiOiJlMzQ3M2UwMS1kYmFkLTQxYzktOTRjMS01N2Q2YjdmODE5MzAiLCJleHAiOjE3NDk2Mzc5NzgsImlhdCI6MTc0OTU1MTU3OH0.eHMVbY_eup7EBoXN0E-33SKts7IX2HSq5EkGoYGzpNM`
  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Accept-Encoding': 'gzip, deflate, br, zstd',
      'Accept-Language': 'en-US,en;q=0.9',
      'Authorization': 'Bearer '+token, // Full token from headers
      'Connection': 'keep-alive',
    },
  };

  const payload = JSON.stringify({
    "action": "list",
  "language": "string",
  "params": {
    "sort":"username desc",
"page":10,
"size":20
  },
  "tenant": "test001",
  "viewId": "auth/users"
  });

  const res = http.post('http://localhost:8080/api/v1/callers/call', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.1); // 100ms delay between iterations to manage load
}