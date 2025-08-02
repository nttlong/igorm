import { check, sleep } from 'k6';
import http from 'k6/http';

// 1. Cấu hình các tùy chọn cho bài kiểm thử
// Trong ví dụ này, chúng ta sẽ chạy 10 người dùng ảo (VUs) trong 30 giây.
export const options = {
    vus: 250,
    duration: '45s',
};

// Dữ liệu JSON bạn muốn gửi trong mỗi request POST
const payload = JSON.stringify({
    Code: 'A001',
    Name: 'Test',
});

// Các header cần thiết cho request JSON
const params = {
    headers: {
        'Content-Type': 'application/json',
        "authorization": `Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpc3MiOiJteS1hdXRoLXNlcnZpY2UiLCJzdWIiOiJ1c2VyLTEyMyIsImV4cCI6MTc1NDE1MTI5MCwiaWF0IjoxNzU0MTQ3NjkwLCJ1c2VyX2lkIjoidXNlci0xMjMiLCJ1c2VybmFtZSI6ImFkbWluIiwicm9sZXMiOlsiYWRtaW4iLCJtZW1iZXIiXX0.LAFjGuo9nOH8IoEVZb4bVXfAegWUombEhLPvDR1bwmQ`
    },
};

// 2. Hàm mặc định (default function)
// K6 sẽ chạy hàm này cho mỗi người dùng ảo trong suốt thời gian của bài kiểm thử.
export default function () {
    // Gửi request POST đến API
    const res = http.post('http://localhost:8080/api/v1/main/test/test', payload, params);

    // 3. Kiểm tra phản hồi (response)
    // k6 sẽ ghi lại các kết quả kiểm tra này.
    check(res, {
        // Kiểm tra mã trạng thái HTTP phải là 200 (OK)
        'status is 200': (r) => r.status === 200,
        // Kiểm tra body của phản hồi phải là chuỗi "OK"
        'body is "OK"': (r) => r.body === '{"Code":"A001","Name":"Test"}',
    });

    // Tạm dừng một chút giữa các request để mô phỏng hành vi của người dùng
    // sleep(1) tạm dừng 1 giây.
    sleep(1);
}