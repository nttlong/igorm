import {Renderer2, Component, OnInit, OnDestroy, HostBinding, ViewChild, ElementRef, AfterViewInit } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần thiết cho *ngIf, *ngFor, ngClass
import { RouterOutlet, Router, NavigationEnd, ActivatedRoute } from '@angular/router'; // Import Router, NavigationEnd, ActivatedRoute
import { filter } from 'rxjs/operators'; // Import filter
import { Subscription } from 'rxjs'; // Import Subscription

// Import Header và Sidebar components
// Đảm bảo đường dẫn này CHÍNH XÁC với vị trí file của bạn
// Ví dụ: './../header/header' nếu Header nằm trong `src/app/components/header/header.ts`
// và AppDashboard nằm trong `src/app/shared/components/app-dashboard/app-dashboard.ts`
import { Header } from '../../components/header/header'; // <-- Đảm bảo đường dẫn này đúng
import { Sidebar } from '../../components/sidebar/sidebar'; // <-- Đảm bảo đường dẫn này đúng

// Định nghĩa interface cho MenuItem
interface MenuItem {
  id: string;
  label: string;
  path: string;
}

@Component({
  selector: 'app-dashboard',
  templateUrl: './app-dashboard.html',
  // template:`<div style="border: 14px solid red; padding: 10px;"><router-outlet></router-outlet></div>`,
  styleUrls: ['./app-dashboard.scss'],
  standalone: true, // Đây là một standalone component
  imports: [
    CommonModule,
    RouterOutlet, // <-- Đảm bảo RouterOutlet được import
    Header,       // <-- Đảm bảo Header được import
    Sidebar       // <-- Đảm bảo Sidebar được import
  ]
})
export class AppDashboard implements OnInit, OnDestroy, AfterViewInit {
  isSidebarCollapsed: boolean = false;
  isDarkMode: boolean = false;
  activeMenuItem: string = "dashboard";
  currentTenantname: string | null = null; // Để lưu tenantname từ URL

  menuItems: MenuItem[] = [
    { id: "dashboard", label: "Dashboard", path: "dashboard" },
    { id: "users", label: "Users", path: "users" },
    { id: "projects", label: "Projects", path: "projects" },
    { id: "settings", label: "Settings", path: "settings" }
  ];
  @ViewChild('mainContent') mainContentElementRef!: ElementRef; // Sử dụng `!` để đảm bảo nó sẽ được gán
  @HostBinding('class.dark') get themeClass() {
    return this.isDarkMode;
  }
  
  private routerSubscription: Subscription | undefined;
  private tenantSubscription: Subscription | undefined;

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute, // Cần ActivatedRoute để lấy tenantname
    private renderer: Renderer2
  ) { }

  ngOnInit(): void {
    // Xử lý responsive sidebar
    this.handleResize();
    window.addEventListener('resize', this.handleResize);

    // Load dark mode preference
    if (localStorage.getItem('theme') === 'dark') {
      this.isDarkMode = true;
      document.documentElement.classList.add('dark');
    }

    // Lắng nghe params của route để lấy tenantname
    // tenantname là param của route cha của AppDashboard
    this.tenantSubscription = this.activatedRoute.paramMap.subscribe(params => {
      this.currentTenantname = params.get('tenantname');
      console.log('AppDashboard - Current Tenantname:', this.currentTenantname);
    });


    // Lắng nghe sự kiện NavigationEnd để cập nhật activeMenuItem
    this.routerSubscription = this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe(() => {
      // Lấy segment cuối cùng của URL
      const currentPathSegment = this.router.url.split('/').pop();
      // Nếu URL kết thúc bằng tenantname (ví dụ: /acme), thì active là dashboard
      if (currentPathSegment === this.currentTenantname) {
        this.activeMenuItem = 'dashboard';
      } else {
        // Ngược lại, sử dụng segment cuối cùng làm activeMenuItem
        this.activeMenuItem = currentPathSegment || 'dashboard';
      }
    });

    // Cập nhật activeMenuItem ban đầu
    const initialPathSegment = this.router.url.split('/').pop();
    if (initialPathSegment === this.currentTenantname) {
      this.activeMenuItem = 'dashboard';
    } else {
      this.activeMenuItem = initialPathSegment || 'dashboard';
    }
  }
  
  // Phương thức để lấy và log chiều cao của phần tử <main>
  logMainContentHeight(): void {
   
    if (this.mainContentElementRef) {
      const mainElement: HTMLElement = this.mainContentElementRef.nativeElement;
      this.setMainContentHeight(window.innerHeight-70);
      // offsetHeight: bao gồm padding, border, và chiều cao nội dung (đã làm tròn)
      // clientHeight: bao gồm padding và chiều cao nội dung (không bao gồm border, scrollbar)
      // getBoundingClientRect().height: chiều cao chính xác của phần tử (bao gồm padding)
      console.log('Chiều cao của Main Content (offsetHeight):', mainElement.offsetHeight + 'px');
      console.log('Chiều cao của Main Content (clientHeight):', mainElement.clientHeight + 'px');
      console.log('Chiều cao của Main Content (getBoundingClientRect().height):', mainElement.getBoundingClientRect().height + 'px');
    }
  }
  /**
   * Phương thức để đặt chiều cao của phần tử main content.
   * @param height Chiều cao mong muốn tính bằng pixel (số).
   */
  setMainContentHeight(height: number): void {
    if (this.mainContentElementRef) {
      const mainElement: HTMLElement = this.mainContentElementRef.nativeElement;
      // Cách 1: Sử dụng Renderer2 (được khuyến nghị)
      this.renderer.setStyle(mainElement, 'height', `${height}px`);
      console.log(`Chiều cao của mainContentElementRef được đặt thành ${height}px bằng Renderer2.`);

      
    } else {
      console.warn('Không thể đặt chiều cao: mainContentElementRef không tồn tại.');
    }
  }
  private handleResize = () => {
    if (window.innerWidth <= 768) {
      this.isSidebarCollapsed = true;
    } else {
      this.isSidebarCollapsed = false;
    }
    // Gọi lại logMainContentHeight khi kích thước cửa sổ thay đổi
    if (this.mainContentElementRef) {
      this.logMainContentHeight();
    }
    
  };

// --- Lifecycle Hook: ngAfterViewInit ---
  // Được gọi sau khi Angular đã khởi tạo các view của component và các view con.
  // Đây là nơi an toàn để truy cập các phần tử DOM đã được render.
  ngAfterViewInit(): void {
    // Đảm bảo phần tử tồn tại trước khi truy cập
    if (this.mainContentElementRef) {
      this.logMainContentHeight();
      // Bạn cũng có thể lắng nghe sự kiện resize để cập nhật chiều cao động
      window.addEventListener('resize', () => this.logMainContentHeight());
    }
  }

  ngOnDestroy(): void {
    window.removeEventListener('resize', this.handleResize);
    if (this.routerSubscription) {
      this.routerSubscription.unsubscribe();
    }
    if (this.tenantSubscription) {
      this.tenantSubscription.unsubscribe();
    }
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
    const item = this.menuItems.find(menuItem => menuItem.id === id);
    if (item && this.currentTenantname) {
      this.router.navigate([`/${this.currentTenantname}/${item.path}`]);
    } else if (item) {
      // Fallback nếu không có tenantname (dù lý thuyết không xảy ra với cấu hình route này)
      this.router.navigate([`/${item.path}`]);
    }
  }

  getActiveMenuLabel(): string {
    const item = this.menuItems.find(item => item.id === this.activeMenuItem);
    return item ? item.label : 'Nội dung';
  }

  
  

}