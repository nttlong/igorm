import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 200,
    duration: '40s',
};

export default function () {

    //let url = "http://localhost:8080//api/v1/example/media/file/tes-004.pdf"
    let url = "http://localhost:8081/api/media/files/2025/08/17/0009.png"
    let res = http.get(url);

    check(res, {
        'status is 200': (r) => r.status === 200,
        //'response has access_token': (r) => r.body && r.body.includes('access_token'),
    });


}
// k6 run download-pdf-vapi.js