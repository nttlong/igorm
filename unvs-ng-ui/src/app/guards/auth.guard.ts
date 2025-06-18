import { Injectable } from '@angular/core';
import { CanActivateFn, Router, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { AuthService } from '../services/auth.service'; // Import AuthService
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';

// CanActivateFn là cách hiện đại (Angular 15+) để định nghĩa guards.
// Nó là một hàm, không phải một class.
export const authGuard: CanActivateFn = (
  route: ActivatedRouteSnapshot,
  state: RouterStateSnapshot
): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree => {

  const authService = inject(AuthService); // Inject AuthService
  const router = inject(Router);         // Inject Router

  return authService.isAuthenticated$.pipe(
    map(isAuthenticated => {
      if (isAuthenticated) {
        return true; // Cho phép truy cập route
      } else {
        // Chuyển hướng đến trang đăng nhập
        // Lấy tenantname từ URL hiện tại hoặc mặc định
        const tenantname = route.parent?.paramMap.get('tenantname') || 'default';
        return router.createUrlTree([`/${tenantname}/login`]); // Chuyển hướng
      }
    })
  );
};

// Cần inject hàm `inject` từ @angular/core
import { inject } from '@angular/core';
import { UrlTree } from '@angular/router'; // Import UrlTree
