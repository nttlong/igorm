import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 100,
    duration: '5s',
};

export default function () {
    let url = 'http://localhost:8012/oauth/token';

    let payload = 'grant_type=password&username=admin&password=123456';

    let params = {
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
    };

    let res = http.post(url, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response has access_token': (r) => r.body && r.body.includes('access_token'),
    });


}
