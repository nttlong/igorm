import { check } from 'k6';
import http from 'k6/http';

export let options = {
    stages: [
        { duration: '5s', target: 200 },  // Tăng đột ngột lên 200 VUs trong 5 giây
        { duration: '30s', target: 200 }, // Giữ nguyên 200 VUs trong 30 giây
        { duration: '5s', target: 0 },    // Giảm đột ngột về 0 VUs trong 5 giây
    ],
};

// Đọc file ở init stage
const binFile = open('./wx-test001.js', 'b');

export default function () {
    const data = {
        File: http.file(binFile, 'test.png', 'image/png'),
        description: 'This is a test upload',
    };

    const res = http.post(
        'http://localhost:8081/api/media/upload',
        data
    );

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}