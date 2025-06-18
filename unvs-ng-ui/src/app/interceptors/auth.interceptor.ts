import { Injectable } from '@angular/core';
import {
  HttpRequest,
  HttpHandler,
  HttpEvent,
  HttpInterceptor,
  HttpErrorResponse
} from '@angular/common/http';
import { Observable, throwError } from 'rxjs';
import { catchError } from 'rxjs/operators';
import { AuthService } from '../services/auth.service'; // Import AuthService
import { Router } from '@angular/router';

@Injectable()
export class AuthInterceptor implements HttpInterceptor {

  constructor(private authService: AuthService, private router: Router) {}

  intercept(request: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const authToken = this.authService.getToken();

    // Nếu có token, thêm nó vào header Authorization
    if (authToken) {
      request = request.clone({
        setHeaders: {
          Authorization: `Bearer ${authToken}`
        }
      });
    }

    // Tiếp tục request và bắt lỗi
    return next.handle(request).pipe(
      catchError((error: HttpErrorResponse) => {
        // Xử lý lỗi 401 (Unauthorized) hoặc 403 (Forbidden)
        if (error.status === 401 || error.status === 403) {
          console.warn('AuthInterceptor: Unauthorized or Forbidden request. Logging out...');
          this.authService.logout(); // Đăng xuất người dùng
          // Chuyển hướng về trang đăng nhập sẽ được xử lý bởi AuthService.logout()
        }
        return throwError(() => error); // Ném lỗi để các subscriber khác xử lý
      })
    );
  }
}