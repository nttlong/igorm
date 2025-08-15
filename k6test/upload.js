import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 200,
    duration: '10s',
};

// Đọc file ở init stage
const binFile = open('./tes-004.pdf', 'b');

export default function () {
    const data = {
        File: http.file(binFile, 'test.png', 'image/png'),
        description: 'This is a test upload',
    };

    const res = http.post(
        'http://localhost:8080/api/v1/example/media/upload2/xxxx',
        data
    );

    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
