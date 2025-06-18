import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter, Routes, Router, ActivatedRoute } from '@angular/router';
import { HttpClient, provideHttpClient, withInterceptorsFromDi } from '@angular/common/http'; // <-- Thêm withInterceptorsFromDi
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
// import { TranslateHttpLoader } from '@ngx-translate/http-loader';

// Import AppDashboard component (layout component)
import { AppDashboard } from './shared/components/app-dashboard/app-dashboard';

// Import các trang component thực tế
import { Dashboard } from './pages/dashboard/dashboard';
import { User } from './pages/user/user';
import { Projects } from './pages/projects/projects';
import { Settings } from './pages/settings/settings';
import { Login } from './pages/login/login'; // Import Login component

// Import UnvsTranslateLoader của bạn
import { UnvsTranslateLoader } from './services/unvs-translate-loader.service'; // <-- Import custom loader

// Import AuthInterceptor
import { AuthInterceptor } from './interceptors/auth.interceptor'; // <-- Import AuthInterceptor
import { HTTP_INTERCEPTORS } from '@angular/common/http'; // <-- Import HTTP_INTERCEPTORS

// Import AuthGuard
import { authGuard } from './guards/auth.guard'; // <-- Import authGuard

// Hàm factory để tạo instance của UnvsTranslateLoader
export function createUnvsTranslateLoader(http: HttpClient, router: Router) {
  return new UnvsTranslateLoader(http, router);
}

const routes: Routes = [
  {
    path: ':tenantname/login',
    component: Login
  },
  {
    path: 'login',
    component: Login
  },
  {
    path: ':tenantname',
    component: AppDashboard,
    canActivate: [authGuard], // <-- Áp dụng AuthGuard vào route cha này
    children: [
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' },
      {
        path: 'dashboard',
        loadComponent: () => import('./pages/dashboard/dashboard').then(m => m.Dashboard)
      },
      {
        path: 'users',
        loadComponent: () => import('./pages/user/user').then(m => m.User)
      },
      {
        path: 'projects',
        loadComponent: () => import('./pages/projects/projects').then(m => m.Projects)
      },
      {
        path: 'settings',
        loadComponent: () => import('./pages/settings/settings').then(m => m.Settings)
      },
      { path: '**', redirectTo: 'dashboard' }
    ]
  },
  {
    path: '',
    redirectTo: 'login',
    pathMatch: 'full'
  },
  { path: '**', redirectTo: 'login' }
];

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes),
    // Cấu hình HttpClient để sử dụng interceptors từ DI
    provideHttpClient(withInterceptorsFromDi()), // <-- Thêm withInterceptorsFromDi
    {
      provide: HTTP_INTERCEPTORS, // Cung cấp interceptor
      useClass: AuthInterceptor,
      multi: true // Quan trọng: cho phép nhiều interceptor
    },
    importProvidersFrom(
      TranslateModule.forRoot({
        loader: {
          provide: TranslateLoader,
          useFactory: (http: HttpClient, router: Router) => createUnvsTranslateLoader(http, router),
          deps: [HttpClient, Router]
        }
      })
    )
  ]
};
