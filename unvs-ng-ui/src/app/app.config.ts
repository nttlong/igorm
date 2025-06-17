import { ApplicationConfig } from '@angular/core';
import { provideRouter, Routes } from '@angular/router';

// Import AppDashboard component (layout component)
import { AppDashboard } from './shared/components/app-dashboard/app-dashboard';

// Import các trang component thực tế
import { Dashboard } from './pages/dashboard/dashboard';
import { User } from './pages/user/user';
import { Projects } from './pages/projects/projects';
import { Settings } from './pages/settings/settings';
import { Login } from './pages/login/login'; // Import Login component

const routes: Routes = [
  {
    // Route cho trang Login (có thể có tenantname hoặc không)
    // Ví dụ: /acme/login hoặc /login
    path: ':tenantname/login', // Route cho login với tenantname
    component: Login // <-- Đây là route độc lập cho Login
  },
  {
    path: 'login', // Route login không có tenantname (fallback hoặc mặc định)
    component: Login // <-- Đây cũng là route độc lập cho Login
  },
  {
    // Route cha cho toàn bộ ứng dụng sau khi đăng nhập
    // :tenantname là một route parameter, nó sẽ khớp với bất kỳ giá trị nào ở vị trí đó
    path: ':tenantname', // <-- ĐÃ THÊM :tenantname VÀO PATH CHÍNH
    component: AppDashboard, // AppDashboard sẽ là layout chính cho tenant này
    children: [ // Các route con này sẽ được hiển thị trong <router-outlet> của AppDashboard
      { path: '', redirectTo: 'dashboard', pathMatch: 'full' }, // Mặc định chuyển hướng đến dashboard của tenant
      {
        path: 'dashboard',
        loadComponent: () => import('./pages/dashboard/dashboard').then(m => m.Dashboard) // Lazy load Dashboard
      },
      {
        path: 'users',
        loadComponent: () => import('./pages/user/user').then(m => m.User) // Lazy load User
      },
      {
        path: 'projects',
        loadComponent: () => import('./pages/projects/projects').then(m => m.Projects) // Lazy load Projects
      },
      {
        path: 'settings',
        loadComponent: () => import('./pages/settings/settings').then(m => m.Settings) // Lazy load Settings
      },
      // Route 404 cho các đường dẫn không khớp trong layout của tenant
      { path: '**', redirectTo: 'dashboard' } // Bất kỳ đường dẫn con nào không khớp sẽ về dashboard của tenant
    ]
  },
  {
    // Route gốc mặc định, có thể chuyển hướng đến một trang chào mừng hoặc login
    path: '',
    redirectTo: 'login', // Chuyển hướng về trang login mặc định
    pathMatch: 'full'
  },
  // Route 404 tổng quát cho các đường dẫn không khớp ở cấp độ cao nhất
  // Sẽ chuyển hướng đến trang login nếu không tìm thấy route nào
  { path: '**', redirectTo: 'login' }
];

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes) // Cung cấp Router với các route đã định nghĩa
  ]
};
