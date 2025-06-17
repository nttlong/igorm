import { Component, Input, Output, EventEmitter, OnInit, OnDestroy } from '@angular/core';
import { CommonModule, NgClass, NgIf, NgFor } from '@angular/common';
import { Router, RouterLink, RouterLinkActive, ActivatedRoute, NavigationEnd } from '@angular/router'; // <-- ĐÃ THÊM NavigationEnd
import { Subscription } from 'rxjs';
import { filter } from 'rxjs/operators'; // Đảm bảo filter cũng được import

// Định nghĩa interface cho MenuItem
interface MenuItem {
  id: string;
  label: string;
  path: string;
}

@Component({
  selector: 'app-sidebar',
  templateUrl: './sidebar.html',
  styleUrls: ['./sidebar.scss'],
  standalone: true,
  imports: [CommonModule, NgClass, NgIf, NgFor, RouterLink, RouterLinkActive]
})
export class Sidebar implements OnInit, OnDestroy {
  @Input() isSidebarCollapsed: boolean = false;
  @Input() menuItems: MenuItem[] = [];
  @Input() activeMenuItem: string = 'dashboard';

  @Output() setActiveMenuItem = new EventEmitter<string>();

  private tenantname: string | null = null;
  private paramSubscription: Subscription | undefined;

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute
  ) { }

  ngOnInit(): void {
    // Để lấy tenantname, chúng ta cần tìm route chứa param 'tenantname'.
    // `activatedRoute` sẽ là route tương ứng với Sidebar.
    // `root` là route cấp cao nhất.
    // Chúng ta duyệt qua cây route từ root đến child để tìm param.
    this.paramSubscription = this.router.events.pipe(
      filter(event => event instanceof NavigationEnd)
    ).subscribe(() => {
      let currentRoute: ActivatedRoute | null = this.activatedRoute.root;
      let tenantnameFound: string | null = null;

      // Duyệt qua các route con để tìm param 'tenantname'
      while (currentRoute) {
        if (currentRoute.snapshot.paramMap.has('tenantname')) {
          tenantnameFound = currentRoute.snapshot.paramMap.get('tenantname');
          break; // Đã tìm thấy, thoát vòng lặp
        }
        currentRoute = currentRoute.firstChild;
      }
      this.tenantname = tenantnameFound;
      console.log('Sidebar - Tenantname retrieved from URL (using Router events):', this.tenantname);
    });

    // Lấy tenantname lần đầu khi component load, trong trường hợp đã có trên URL
    let initialRoute: ActivatedRoute | null = this.activatedRoute.root;
    while (initialRoute) {
      if (initialRoute.snapshot.paramMap.has('tenantname')) {
        this.tenantname = initialRoute.snapshot.paramMap.get('tenantname');
        console.log('Sidebar - Initial Tenantname from URL (snapshot):', this.tenantname);
        break;
      }
      initialRoute = initialRoute.firstChild;
    }
  }

  ngOnDestroy(): void {
    if (this.paramSubscription) {
      this.paramSubscription.unsubscribe();
    }
  }

  handleMenuItemClick(id: string): void {
    alert(`Menu item with ID ${id} clicked!`);

    const selectedMenuItem = this.menuItems.find(item => item.id === id);

    if (selectedMenuItem && selectedMenuItem.path && this.tenantname) {
      this.router.navigate([`/${this.tenantname}/${selectedMenuItem.path}`]);
    } else if (selectedMenuItem && selectedMenuItem.path) {
      // Trường hợp không có tenantname, điều hướng đến path gốc (ví dụ: /dashboard nếu không có tenant)
      this.router.navigate([`/${selectedMenuItem.path}`]);
    }

    this.setActiveMenuItem.emit(id);
  }
}