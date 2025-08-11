import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 200,
    duration: '1s',
};

export default function () {
    //let url = 'http://localhost:8080/api/v1/example/media/file';// 'http://localhost:8080/api/oauth/token';
    // let url = "http://localhost:8080//api/v1/example/media/tes-004.pdf"
    //let url = "http://localhost:8080/api/v1/example/media/test.mp4.mp4"
    let url = "http://127.0.0.1:8081/download/tes-004.pdf"
    let payload = 'grant_type=password&username=admin&password=123456';

    let params = {
        headers: {
            'Content-Type': 'application/json',//'application/x-www-form-urlencoded',
        },
    };

    let res = http.get(url, "{}", params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        //'response has access_token': (r) => r.body && r.body.includes('access_token'),
    });


}
