import { Component, OnInit, OnDestroy, HostBinding } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần thiết cho *ngIf, *ngFor, ngClass

// Import Header và Sidebar components
// Đã điều chỉnh đường dẫn file để không có hậu tố .component
import { Header } from '../header/header'; // <-- ĐÃ SỬA ĐƯỜNG DẪN
import { Sidebar } from '../sidebar/sidebar'; // <-- ĐÃ SỬA ĐƯỜNG DẪN

// Định nghĩa interface cho MenuItem
interface MenuItem {
  id: string;
  label: string;
  // Bạn có thể thêm icon: string; nếu muốn truyền tên icon dưới dạng string và xử lý nó trong template con
}

@Component({
  selector: 'app-dashboard', // Selector của component này
  templateUrl: './app-dashboard.html', // Trỏ đến file HTML template (đã sửa từ .component.html)
  styleUrls: ['./app-dashboard.scss'], // Trỏ đến file SCSS style
  standalone: true, // <-- Đây là một standalone component
  imports: [
    CommonModule,
    Header,   // <-- Thêm Header vào imports
    Sidebar   // <-- Thêm Sidebar vào imports
  ]
})
export class AppDashboard implements OnInit, OnDestroy {
  // Thêm các thuộc tính cần thiết để template hoạt động
  isSidebarCollapsed: boolean = false;
  isDarkMode: boolean = false;
  activeMenuItem: string = "dashboard"; // Giá trị mặc định

  // Định nghĩa menuItems để truyền xuống Sidebar
  menuItems: MenuItem[] = [
    { id: "dashboard", label: "Dashboard" },
    { id: "users", label: "Users" },
    { id: "projects", label: "Projects" },
    { id: "settings", label: "Settings" }
  ];

  // HostBinding để thêm/xóa class 'dark' vào thẻ <html>
  @HostBinding('class.dark') get themeClass() {
    return this.isDarkMode;
  }

  constructor() { }

  ngOnInit(): void {
    // Logic khởi tạo và quản lý responsive
    this.handleResize();
    window.addEventListener('resize', this.handleResize);

    // Load dark mode preference từ localStorage
    if (localStorage.getItem('theme') === 'dark') {
      this.isDarkMode = true;
      document.documentElement.classList.add('dark');
    }
  }

  ngOnDestroy(): void {
    // Gỡ bỏ event listener khi component bị hủy
    window.removeEventListener('resize', this.handleResize);
  }

  // Phương thức để chuyển đổi trạng thái thu gọn/mở rộng của sidebar
  toggleSidebar(): void {
    this.isSidebarCollapsed = !this.isSidebarCollapsed;
  }

  // Phương thức để chuyển đổi chế độ sáng/tối
  toggleDarkMode(): void {
    this.isDarkMode = !this.isDarkMode;
    if (this.isDarkMode) {
      document.documentElement.classList.add('dark');
      localStorage.setItem('theme', 'dark');
    } else {
      document.documentElement.classList.remove('dark');
      localStorage.setItem('theme', 'light');
    }
  }

  // Phương thức để đặt mục menu đang hoạt động
  setActiveMenuItem(id: string): void {
    this.activeMenuItem = id;
  }

  // Phương thức để lấy nhãn (label) của mục menu đang hoạt động
  getActiveMenuLabel(): string {
    const item = this.menuItems.find(item => item.id === this.activeMenuItem);
    return item ? item.label : 'Nội dung'; // Trả về nhãn hoặc một giá trị mặc định
  }

  // Hàm xử lý khi kích thước cửa sổ thay đổi để điều chỉnh sidebar
  private handleResize = () => {
    if (window.innerWidth <= 768) { // Màn hình nhỏ hơn hoặc bằng 768px (mobile/tablet)
      this.isSidebarCollapsed = true;
    } else {
      this.isSidebarCollapsed = false;
    }
  };
}