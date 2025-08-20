//k6 run list-files-test-dot-net.js
import { check, sleep } from "k6";
import http from "k6/http";

// export const options = {
//   vus: 200,           // 10 virtual users
//   duration: "30s",   // chạy trong 30 giây
// };
// export let options = {
//   stages: [
//     { duration: '30s', target: 200 }, // tăng dần tới 200 VU
//   ],
// };
export let options = {
  scenarios: {
    shock_load: {
      executor: "constant-vus",
      vus: 200,         // số VU cùng lúc
      duration: "30s",  // chạy trong 30 giây
    },
  },
};
export default function () {
  const url = "http://localhost:5027/api/Media/list-files";
  //let url = "http://localhost:5027/api/Media/Hello"


  let res = http.get(url);

  check(res, {
    "status is 200": (r) => r.status === 200,

  });

  sleep(1); // nghỉ 1 giây giữa các request
}
