import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 200,
    duration: '40s',
};

export default function () {

    let url = "http://localhost:8080//api/v1/example/media/file/tes-004.pdf"

    let res = http.get(url);

    check(res, {
        'status is 200': (r) => r.status === 200,
        //'response has access_token': (r) => r.body && r.body.includes('access_token'),
    });


}
