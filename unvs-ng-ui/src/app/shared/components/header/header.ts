import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common';
import { AuthService } from '../../../services/auth.service'; // <-- Import AuthService

@Component({
  selector: 'app-header',
  templateUrl: './header.html',
  styleUrls: ['./header.scss'],
  standalone: true,
  imports: [CommonModule]
})
export class Header {
  @Input() isDarkMode: boolean = false;
  @Input() isProfileDropdownOpen: boolean = false; // Thêm input cho dropdown
  
  @Output() onSidebarToggle = new EventEmitter<void>();
  @Output() onToggleDarkMode = new EventEmitter<void>();
  @Output() setIsProfileDropdownOpen = new EventEmitter<boolean>(); // Output cho dropdown
  @Output() onLogout = new EventEmitter<void>(); // Output cho logout

  constructor(private authService: AuthService) { } // Inject AuthService

  // Phương thức để gọi logout từ AuthService
  logout(): void {
    this.authService.logout();
    this.onLogout.emit(); // Emit sự kiện logout lên component cha
  }
}
