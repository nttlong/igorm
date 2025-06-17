import { Component, OnInit, OnDestroy, HostBinding } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần cho *ngIf
import { FormsModule } from '@angular/forms'; // Cần cho [(ngModel)]
import { ActivatedRoute, Router } from '@angular/router'; // Cần cho Router và ActivatedRoute
import { Subscription } from 'rxjs'; // Cần để quản lý subscriptions
// import { AuthService } from '../../services/auth.service'; // Giả định bạn sẽ có một AuthService
// import { ApiService } from '../../services/api.service'; // Giả định bạn sẽ có một ApiService

@Component({
  selector: 'app-login', // Selector của component
  templateUrl: './login.html', // Trỏ đến file HTML
  styleUrls: ['./login.scss'], // Trỏ đến file SCSS
  standalone: true, // Đây là một standalone component
  imports: [
    CommonModule,
    FormsModule // Import FormsModule để sử dụng [(ngModel)]
  ]
})
export class Login implements OnInit, OnDestroy { // Tên lớp là Login
  username = '';
  password = '';
  error: string | null = null;
  isLoading = false;
  tenantname: string | null = null;
  private paramSubscription: Subscription | undefined;

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute,
    // private authService: AuthService, // Kích hoạt nếu bạn có AuthService
    // private apiService: ApiService // Kích hoạt nếu bạn có ApiService
  ) { }

  ngOnInit(): void {
    // Lấy tenantname từ URL params (được định nghĩa trong app.config.ts)
    this.paramSubscription = this.activatedRoute.paramMap.subscribe(params => {
      this.tenantname = params.get('tenantname');
      console.log('Login Page - Tenantname:', this.tenantname);
      // Bạn có thể thiết lập base URL cho API service ở đây nếu cần,
      // tương tự cách setBaseApiUrl trong React/Vue app
      // if (this.tenantname) {
      //   this.apiService.setBaseUrl(`http://localhost:8080/api/v1/${this.tenantname}`);
      // } else {
      //   this.apiService.setBaseUrl(`http://localhost:8080/api/v1`);
      // }
    });
  }

  ngOnDestroy(): void {
    if (this.paramSubscription) {
      this.paramSubscription.unsubscribe(); // Hủy subscription để tránh rò rỉ bộ nhớ
    }
  }

  async handleLogin(event: Event): Promise<void> {
    event.preventDefault(); // Ngăn chặn hành vi submit mặc định của form
    this.error = null;
    this.isLoading = true;

    if (!this.tenantname) {
      this.error = 'Không tìm thấy thông tin tenant.'; // Tạm dịch
      this.isLoading = false;
      return;
    }

    try {
      // --- LOGIC GỌI API ĐĂNG NHẬP CỦA BẠN (SỬ DỤNG AJAX/FETCH HOẶC API SERVICE) ---
      // Ví dụ mô phỏng:
      await new Promise(resolve => setTimeout(resolve, 1500)); // Simulate API call

      if (this.username === 'admin' && this.password === '123') {
        console.log('Đăng nhập thành công!');
        // Lưu thông tin xác thực vào localStorage hoặc AuthService
        localStorage.setItem('authToken', 'mock_angular_token_123');
        localStorage.setItem('username', this.username);
        localStorage.setItem('tenantname', this.tenantname);

        // Sử dụng AuthService để login nếu bạn có
        // this.authService.login(response.access_token, response.username, this.tenantname);

        // Chuyển hướng đến dashboard của tenant
        this.router.navigate([`/${this.tenantname}/dashboard`]);
      } else {
        this.error = 'Tên đăng nhập hoặc mật khẩu không đúng.'; // Tạm dịch
      }
    } catch (err) {
      console.error('Lỗi đăng nhập:', err);
      this.error = 'Đăng nhập thất bại. Vui lòng kiểm tra lại tên đăng nhập hoặc mật khẩu.'; // Tạm dịch
    } finally {
      this.isLoading = false;
    }
  }
}