//k6 run upload-file-wx.js
import { check, sleep } from "k6";
import http from "k6/http";

export const options = {
    vus: 50,           // số lượng virtual users
    duration: "10s",   // chạy trong 10 giây
    // scenarios: {
    //     shock_load: {
    //         executor: "constant-vus",
    //         vus: 200,
    //         duration: "30s",
    //     },
    // },
};
const bigFile20MB = new ArrayBuffer(20 * 1024 * 1024); // 20MB zero buffer
const binFile = open("D:/code/go/news2/igorm/comparr-framework/wx_hello/bm_test.go", "b");

export default function () {
    // Đọc file từ local để gửi
    //const binFile = open("D:/code/go/news2/igorm/comparr-framework/wx_hello/bm_test.go", "b");

    const data = {
        File: http.file(bigFile20MB, "bm_test.go"), // field "File" phải khớp với server handler
    };
    const res = http.post("http://localhost:5000/api/test-api/upload", data);

    check(res, {
        "status is 200": (r) => r.status === 200,
    });

    sleep(1);
}
