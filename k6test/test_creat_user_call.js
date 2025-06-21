import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
  vus: 200, // 200 virtual users
  duration: '30s', // Test duration of 30 seconds
};

export default function () {
    const description = `Các bước tiếp theo (Dựa trên phân tích trước đó của bạn):

    Bạn đã xác định được bcrypt.GenerateFromPassword took 68 ms là một yếu tố quan trọng. Đây chắc chắn là một đóng góp lớn vào độ trễ.
    
    Theo dõi tài nguyên server (CPU, RAM, Network I/O, Disk I/O): Trong khi chạy bài test k6, bạn cần theo dõi chặt chẽ các chỉ số này trên máy chủ ứng dụng Go và máy chủ cơ sở dữ liệu.
    CPU: Rất có thể CPU của máy chủ Go của bạn đang bị quá tải (gần 100% sử dụng) do bcrypt.
    DB: Kiểm tra CPU, RAM, Disk I/O của DB server. Xem có truy vấn nào bị chậm trễ kéo dài không.
    Đánh giá lại cost factor của Bcrypt: 68ms là trong phạm vi chấp nhận được, nhưng nếu CPU của bạn đang bị giới hạn, bạn có thể cân nhắc giảm cost factor xuống một chút để giảm thời gian băm mật khẩu (ví dụ: xuống 40-50ms) và xem xét tác động đến hiệu suất tổng thể. Tuy nhiên, đừng làm giảm bảo mật quá nhiều.
    Tối ưu hóa DB:
    Kiểm tra các chỉ mục (indexes) trên bảng User. Đảm bảo các cột được sử dụng trong WHERE clauses (nếu có trong quá trình tạo user, ví dụ kiểm tra email tồn tại) và các khóa ngoại được index phù hợp.
    Kiểm tra cấu hình kết nối DB (connection pooling) trong ứng dụng Go của bạn.
    Kiểm tra tình trạng lock trên DB.
    Mở rộng (Scaling):
    Nếu CPU là điểm nghẽn chính, giải pháp hiệu quả nhất là scale out (thêm nhiều instance của ứng dụng Go của bạn). Điều này sẽ phân tán gánh nặng tính toán bcrypt trên nhiều CPU core/máy chủ khác nhau.
    Hoặc scale up (nâng cấp máy chủ hiện tại lên CPU mạnh hơn/nhiều core hơn).
    Phân tích code: Sử dụng các công cụ profiling của Go (pprof) để xác định chính xác phần nào trong code của bạn (ngoài bcrypt) đang tiêu tốn nhiều thời gian nhất.
    Tóm lại, hệ thống của bạn đang hoạt động nhưng đang bị quá tải. Bạn cần tập trung vào việc xác định và tối ưu hóa các điểm nghẽn, mà khả năng cao nhất là tài nguyên CPU do bcrypt và/hoặc `;


    const userIndex = __VU; // Hoặc dùng __ITER nếu muốn tăng theo lần lặp
  const code=`uer-${Date.now()}`
  const email = `${code}${__VU}_${Date.now()}@example.com`;
  const username = `${code}${userIndex}`;
  const password = `123456`;
  
    const token1 = `eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiJlMzQ3M2UwMS1kYmFkLTQxYzktOTRjMS01N2Q2YjdmODE5MzAiLCJ1c2VybmFtZSI6ImFkbWluIiwiZW1haWwiOiJhZG1pbkBhZG1pbi5jb20iLCJzdWIiOiJlMzQ3M2UwMS1kYmFkLTQxYzktOTRjMS01N2Q2YjdmODE5MzAiLCJleHAiOjE3NDk2Mzc5NzgsImlhdCI6MTc0OTU1MTU3OH0.eHMVbY_eup7EBoXN0E-33SKts7IX2HSq5EkGoYGzpNM`;
  const token="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJlbWFpbCI6bnVsbCwiZXhwIjoxNzUwNDg3OTc2LCJpYXQiOjE3NTA0ODQzNzYsInJvbGUiOiI0MTJmNTIwZC1lYzM0LTQyZGItYTRiOC1lMzBiYTIwZTZjM2QiLCJzY29wZSI6InJlYWQgd3JpdGUiLCJ1c2VySWQiOiI0MTJmNTIwZC1lYzM0LTQyZGItYTRiOC1lMzBiYTIwZTZjM2QiLCJ1c2VybmFtZSI6InJvb3QifQ.7tmaC9YUqQL2LY-lciEcXJHpgPOF4i9UQuBP2tL4Jv4"
    const params = {
    headers: {
      'Content-Type': 'application/json',
      'Accept-Encoding': 'gzip, deflate, br, zstd',
      'Accept-Language': 'en-US,en;q=0.9',
      'Authorization': 'Bearer '+token, // Full token from headers
      'Connection': 'keep-alive',
    },
  };

  const payload = JSON.stringify({
    username: username,
      password: password,
      email: email,
      description: description
  });

  const res = http.post('http://localhost:8080/api/v1/invoke?feature=common&action=create&module=unvs.br.auth.users&tenant=default&lan=vi', payload, params);

  check(res, {
    'status is 200': (r) => r.status === 200,
  });

  sleep(0.1); // 100ms delay between iterations to manage load
}