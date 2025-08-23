//k6 run create-user-go.js
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
    vus: 10,         // số VU cùng lúc
    duration: "30s",
    // scenarios: {
    //     shock_load: {
    //         executor: "constant-vus",
    //         vus: 200,         // số VU cùng lúc
    //         duration: "30s",  // chạy trong 30 giây
    //     },
    // },
};
export default function () {
    const url = "http://localhost:8080/api/users/createt-user";
    //let url = "http://localhost:8080/api/media/hello"
    let payload = JSON.stringify({
        "user_name": "admin",
        "password": "12344566"
    });
    let res = http.post(url, payload, {
        headers: {
            "Content-Type": "application/json",
        },
    })

    check(res, {
        "status is 200": (r) => r.status === 200,

    });

    sleep(1); // nghỉ 1 giây giữa các request
}
