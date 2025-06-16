import { Component, OnInit, OnDestroy, HostBinding } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần thiết cho *ngIf, *ngFor, ngClass

@Component({
  selector: 'app-simple-dashboard',
  templateUrl: './simple-dashboard.component.html',
  styleUrls: ['./simple-dashboard.component.css'],
  standalone: true, // <-- Đây là một standalone component
  imports: [CommonModule] // Import CommonModule
})
export class SimpleDashboardComponent implements OnInit, OnDestroy {
  isSidebarCollapsed: boolean = false;
  isDarkMode: boolean = false;
  activeMenuItem: string = "dashboard";

  // Dữ liệu cho các mục menu
  menuItems = [
    { id: "dashboard", label: "Dashboard" },
    { id: "users", label: "Users" },
    { id: "projects", label: "Projects" },
    { id: "settings", label: "Settings" }
  ];

  // HostBinding để thêm class 'dark' vào thẻ <html> khi isDarkMode là true
  @HostBinding('class.dark') get themeClass() {
    return this.isDarkMode;
  }

  constructor() { }

  ngOnInit(): void {
    this.handleResize();
    window.addEventListener('resize', this.handleResize);

    if (localStorage.getItem('theme') === 'dark') {
      this.isDarkMode = true;
      document.documentElement.classList.add('dark');
    }
  }

  ngOnDestroy(): void {
    window.removeEventListener('resize', this.handleResize);
  }

  toggleSidebar(): void {
    this.isSidebarCollapsed = !this.isSidebarCollapsed;
  }

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

  setActiveMenuItem(id: string): void {
    this.activeMenuItem = id;
  }

  getActiveMenuLabel(): string {
    const item = this.menuItems.find(item => item.id === this.activeMenuItem);
    return item ? item.label : 'Nội dung';
  }

  private handleResize = () => {
    if (window.innerWidth <= 768) {
      this.isSidebarCollapsed = true;
    } else {
      this.isSidebarCollapsed = false;
    }
  };
}