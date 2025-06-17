import { ApplicationConfig, importProvidersFrom } from '@angular/core';
import { provideRouter, Routes, ActivatedRoute, Router } from '@angular/router'; // <-- Vẫn import Router ở đây, nhưng không dùng ActivatedRoute trong factory
import { HttpClient, provideHttpClient } from '@angular/common/http';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';

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

// Hàm factory để tạo instance của UnvsTranslateLoader
// Nó sẽ nhận HttpClient và Router
export function createUnvsTranslateLoader(http: HttpClient, router: Router) { // <-- Nhận Router thay vì ActivatedRoute
  return new UnvsTranslateLoader(http, router); // <-- Truyền Router vào loader
}

// Hàm factory để tạo instance của TranslateHttpLoader
// export function HttpLoaderFactory(http: HttpClient) {
//   return new TranslateHttpLoader(http, '/assets/i18n/', '.json');
// }

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
    provideHttpClient(),
    importProvidersFrom(
      TranslateModule.forRoot({
        loader: {
          provide: TranslateLoader,
          // Sử dụng useFactory để inject HttpClient và Router
          useFactory: (http: HttpClient, router: Router) => createUnvsTranslateLoader(http, router), // <-- Truyền Router
          deps: [HttpClient, Router] // <-- Khai báo Router là một dependency
        }
      })
    )
  ]
};
