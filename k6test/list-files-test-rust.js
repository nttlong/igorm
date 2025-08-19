//k6 run list-files-test-rust.js
import { check, sleep } from "k6";
import http from "k6/http";

export const options = {
  vus: 200,           // 10 virtual users
  duration: "30s",   // chạy trong 30 giây
};

export default function () {
  const url = "http://127.0.0.1:8082/api/media/list-files";
  //let url = "http://localhost:8080/api/media/hello"

  let res = http.get(url);

  check(res, {
    "status is 200": (r) => r.status === 200,

  });

  sleep(1); // nghỉ 1 giây giữa các request
}
