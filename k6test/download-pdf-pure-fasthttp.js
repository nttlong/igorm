import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 200,
    duration: '10s',
};

export default function () {
    let url = "http://127.0.0.1:8083/file";
    let res = http.get(url);
    check(res, {
        'status is 200': (r) => r.status === 200,
    });
}
