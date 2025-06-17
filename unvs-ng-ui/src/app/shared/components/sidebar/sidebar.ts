import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule, NgClass, NgIf, NgFor } from '@angular/common'; // Cần cho *ngIf, *ngFor, ngClass

// Định nghĩa interface cho MenuItem để có kiểu dữ liệu rõ ràng hơn
interface MenuItem {
  id: string;
  label: string;
  // Bạn có thể thêm icon: string; nếu muốn truyền tên icon dưới dạng string và xử lý nó trong template
}

@Component({
  selector: 'app-sidebar', // Selector của component này
  templateUrl: './sidebar.html', // <-- ĐÃ THAY ĐỔI TỪ './sidebar.component.html' SANG './sidebar.html'
  styleUrls: ['./sidebar.scss'], // Có thể để trống nếu chỉ dùng Tailwind
  standalone: true, // <-- Đánh dấu component này là standalone
  imports: [CommonModule, NgClass, NgIf, NgFor] // <-- Import CommonModule và các directives cần thiết
})
export class Sidebar { // <-- Tên lớp của component
  // @Input() để nhận dữ liệu từ component cha
  @Input() isSidebarCollapsed: boolean = false;
  @Input() menuItems: MenuItem[] = []; // Mảng các mục menu
  @Input() activeMenuItem: string = 'dashboard'; // Mục menu đang active

  // @Output() để gửi sự kiện lên component cha
  // Khi một mục menu được click, sự kiện này sẽ được emit với id của mục đó
  @Output() setActiveMenuItem = new EventEmitter<string>();

  constructor() {
    

   }
}
