import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 50, // 200 virtual users
  duration: '30s', // Test duration of 30 seconds
};

export default function () {
    
  
    const token = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6bnVsbCwiZXhwIjoxNzUwMzg3MjUzLCJpYXQiOjE3NTAzODM2NTMsInJvbGUiOiI0MTJmNTIwZC1lYzM0LTQyZGItYTRiOC1lMzBiYTIwZTZjM2QiLCJzY29wZSI6InJlYWQgd3JpdGUiLCJ1c2VySWQiOiI0MTJmNTIwZC1lYzM0LTQyZGItYTRiOC1lMzBiYTIwZTZjM2QiLCJ1c2VybmFtZSI6InJvb3QifQ.hq0Z-69tlew0mzPeMmUc2Us35hKXzJx4HvKataRWeWs`
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
    "pageIndex": 0,
    "pageSize": 50
});

  const res = http.post('http://localhost:8080/api/v1/invoke?feature=unvs.br.auth.users&action=list&module=unvs.br.auth.users&tenant=default&lan=en', payload,params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.1); // 100ms delay between iterations to manage load
}