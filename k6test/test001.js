import { check } from 'k6';

export let options = {
    vus: 100,
    duration: '5s',
};

export default function () {
    let url = 'https://vnresource.vn/vnresource-hrm-pro-v10-tai-jaccs-chinh-thuc-van-hanh-cu-hich-manh-me-cho-quan-tri-nhan-su-so/?fbclid=IwY2xjawMBEzRleHRuA2FlbQIxMQBicmlkETE4REhmOFcwOWFnOUdXY1hFAR4rk4N96Zi5zvD3PsrZI54Awby6A3NytHPErP4Mmz8-tgw4R-yWI9DiNL0g6A_aem_ysheTchpnhGzpPxWuY0ZGA';

    // let payload = 'grant_type=password&username=admin&password=123456';

    // let params = {
    //     headers: {
    //         'Content-Type': 'application/x-www-form-urlencoded',
    //     },
    // };

    // let res = http.post(url, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        // 'response has access_token': (r) => r.body && r.body.includes('access_token'),
    });


}
