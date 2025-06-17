// import { ApplicationConfig } from '@angular/core';
// import { provideRouter, Routes } from '@angular/router';

// // Import SimpleDashboardComponent của bạn
// // import { SimpleDashboardComponent } from './simple-dashboard-delete/simple-dashboard.component';

// // Định nghĩa các Routes cho ứng dụng
// const routes: Routes = [
//   {
//     // Đây là route chính để hiển thị SimpleDashboardComponent ngay lập tức.
//     // Khi truy cập '/' (đường dẫn gốc) sẽ hiển thị SimpleDashboardComponent.
//     path: '',
//     component: SimpleDashboardComponent,
//   },
//   {
//     // Route 404 cho bất kỳ đường dẫn nào khác không khớp
//     path: '**',
//     component: SimpleDashboardComponent // Sử dụng SimpleDashboardComponent làm placeholder cho 404
//     // Hoặc tạo một NotFoundPageComponent nếu bạn muốn trang 404 riêng biệt:
//     // component: NotFoundPageComponent
//   }
// ];

// export const appConfig: ApplicationConfig = {
//   providers: [
//     provideRouter(routes) // Cung cấp Router với các route đã định nghĩa
//   ]
// };
import { ApplicationConfig } from '@angular/core';
import { provideRouter, Routes } from '@angular/router';

// Import AppDashboard component (đảm bảo đường dẫn này đúng)
// Đã điều chỉnh đường dẫn file để không có hậu tố .component
import { AppDashboard } from './shared/components/app-dashboard/app-dashboard';

const routes: Routes = [
  {
    // Đây là route chính để hiển thị AppDashboard ngay lập tức.
    // Khi truy cập '/' (đường dẫn gốc) sẽ hiển thị AppDashboard.
    path: '',
    component: AppDashboard, // <-- Sử dụng AppDashboard ở đây
  },
  {
    // Route 404 cho bất kỳ đường dẫn nào khác không khớp
    path: '**',
    component: AppDashboard // <-- Sử dụng AppDashboard làm placeholder cho 404
    // Hoặc nếu bạn đã tạo NotFoundPageComponent, hãy sử dụng nó:
    // component: NotFoundPageComponent
  }
];

export const appConfig: ApplicationConfig = {
  providers: [
    provideRouter(routes) // Cung cấp Router với các route đã định nghĩa
  ]
};