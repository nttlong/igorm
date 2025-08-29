// k6 run post-user-gin.js
import { check } from "k6";
import http from "k6/http";

export const options = {
  scenarios: {
    shock_load: {
      executor: "constant-vus",
      vus: 200,
      duration: "30s",
    },
  },
};

export default function () {
  const payload = JSON.stringify({
    name: "Nguyen Van A",
    age: 30,
    email: "a@example.com",
    phones: ["0909xxx111", "0912xxx222"],
    address: { city: "Hanoi", district: "Ba Dinh", street: "Kim Ma" },
  });

  const headers = { "Content-Type": "application/json" };

  let res = http.post("http://localhost:8081/user", payload, { headers });
  //console.log(`Status: ${res.status}`)
  check(res, { "status is 200": (r) => r.status === 200 });
}
