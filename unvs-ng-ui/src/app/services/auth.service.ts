import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, BehaviorSubject, of, throwError } from 'rxjs';
import { tap, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment'; // Giả sử bạn có environment.apiUrl
import {UnvsTranslateLoader} from './unvs-translate-loader.service'
import { TranslateService } from '@ngx-translate/core'; 
// Định nghĩa interface cho phản hồi đăng nhập từ API của bạn
interface LoginResponse {
  access_token: string;
  refresh_token?: string; // Tùy chọn
  username?: string;
  userId?: string;
  tenantname?: string;
  // ... các trường khác từ API của bạn
}

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  // Observable để theo dõi trạng thái xác thực
  private _isAuthenticated = new BehaviorSubject<boolean>(this.hasAuthToken());
  isAuthenticated$ = this._isAuthenticated.asObservable(); // Public Observable

  // Observable để theo dõi thông tin người dùng (ví dụ: username, tenantname)
  private _currentUser = new BehaviorSubject<any | null>(null); // Bạn có thể định nghĩa interface cụ thể hơn
  currentUser$ = this._currentUser.asObservable();

  constructor(
    private http: HttpClient,
    private router: Router,
    private translate: TranslateService
  ) {
    
    this.loadAuthState();
  }

  // Kiểm tra xem đã có token trong localStorage chưa
  private hasAuthToken(): boolean {
    return !!localStorage.getItem('authToken');
  }

  // Tải trạng thái xác thực và thông tin người dùng từ localStorage
  private loadAuthState(): void {
    if (this.hasAuthToken()) {
      // Đọc token và các thông tin khác từ localStorage
      const token = localStorage.getItem('authToken');
      const username = localStorage.getItem('username');
      const tenantname = localStorage.getItem('tenantname');
      // ... (thêm các trường khác nếu bạn lưu)

      // Cập nhật BehaviorSubject
      this._isAuthenticated.next(true);
      this._currentUser.next({ username, tenantname, token }); // Cần một object với các thuộc tính này
      console.log('AuthService: User is authenticated, loaded from localStorage.');
    } else {
      this._isAuthenticated.next(false);
      this._currentUser.next(null);
      console.log('AuthService: User is not authenticated.');
    }
  }

  // Phương thức đăng nhập
  // (Giả sử bạn đã có Caller service hoặc bạn sẽ gọi trực tiếp HttpClient)
  login(username: string, password: string, tenantname: string): Observable<LoginResponse> {
    alert(this.translate.currentLang)
    //http://localhost:8080/api/v1/invoke?feature=common&action=login&module=unvs.br.auth.users&tenant=default&lan=vi
    const loginUrl = `${environment.apiUrl}/oauth/token`; // Hoặc endpoint login cụ thể của bạn
    // Chuẩn bị dữ liệu form-urlencoded cho OAuth2 Password Flow
    const body = new URLSearchParams();
    body.set('grant_type', 'password');
    body.set('username', `${username}@${tenantname}`);
    body.set('password', password);
    // Nếu API của bạn yêu cầu client_id/client_secret trong body hoặc header, hãy thêm vào đây
    // body.set('client_id', 'your-client-id');

    // Headers cho form-urlencoded
    const headers = { 'Content-Type': 'application/x-www-form-urlencoded' };

    // Gửi yêu cầu POST đến API đăng nhập
    return this.http.post<LoginResponse>(loginUrl, body.toString(), { headers }).pipe(
      tap(response => {
        // Lưu token và thông tin người dùng vào localStorage
        localStorage.setItem('authToken', response.access_token);
        localStorage.setItem('username', response.username || username); // Lưu username từ response hoặc từ input
        localStorage.setItem('tenantname', tenantname); // Lưu tenantname

        // Cập nhật trạng thái xác thực
        this._isAuthenticated.next(true);
        this._currentUser.next({ username: response.username || username, tenantname, token: response.access_token });
        console.log('AuthService: Login successful, state updated.');
      }),
      catchError(error => {
        // Xử lý lỗi đăng nhập
        console.error('AuthService: Login failed', error);
        this.logout(); // Đảm bảo trạng thái không xác thực nếu có lỗi
        return throwError(() => new Error('Login failed')); // Ném lỗi để component xử lý
      })
    );
  }

  // Phương thức đăng xuất
  logout(): void {
    // Xóa token và thông tin người dùng từ localStorage
    localStorage.removeItem('authToken');
    localStorage.removeItem('username');
    localStorage.removeItem('tenantname');

    // Cập nhật trạng thái xác thực
    this._isAuthenticated.next(false);
    this._currentUser.next(null);
    console.log('AuthService: User logged out, state updated.');

    // Chuyển hướng về trang đăng nhập của tenant hiện tại hoặc mặc định
    const currentTenantname = this.router.url.split('/')[1] || 'default';
    this.router.navigate([`/${currentTenantname}/login`]);
  }

  // Phương thức để lấy token
  getToken(): string | null {
    return localStorage.getItem('authToken');
  }

  // Lấy tenantname từ localStorage
  getTenantname(): string | null {
    return localStorage.getItem('tenantname');
  }
}