//http://localhost:8081/api/media/list-files
// k6 run media-list-files.js
import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 200,
    duration: '1s',
};

export default function () {
    let url = 'http://localhost:8081/api/media/list-files';// 'http://localhost:8080/api/oauth/token';

    // let payload = 'grant_type=password&username=admin&password=123456';

    // let params = {
    //     headers: {
    //         'Content-Type': 'application/json',//'application/x-www-form-urlencoded',
    //     },
    // };

    let res = http.post(url);

    check(res, {
        'status is 200': (r) => r.status === 200,
        //'response has access_token': (r) => r.body && r.body.includes('access_token'),
    });


}
