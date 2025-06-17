import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần thiết cho *ngIf trong template

@Component({
  selector: 'app-header', // Selector của component này, sẽ được sử dụng trong template cha
  templateUrl: './header.html', // Trỏ đến file HTML template
  styleUrls: ['./header.scss'], // Trỏ đến file CSS style (có thể để trống)
  standalone: true, // <-- Đánh dấu component này là standalone
  imports: [CommonModule] // <-- Import CommonModule để sử dụng directives như *ngIf
})
export class Header { // <-- Tên lớp của component
  // @Input() để nhận dữ liệu từ component cha
  @Input() isDarkMode: boolean = false; // Nhận trạng thái dark mode từ bên ngoài

  // @Output() để gửi sự kiện lên component cha
  // Khi nút toggle sidebar được click, sự kiện này sẽ được emit
  @Output() onSidebarToggle = new EventEmitter<void>();
  // Khi nút toggle dark mode được click, sự kiện này sẽ được emit
  @Output() onToggleDarkMode = new EventEmitter<void>();

  constructor() {

    
   }
}
