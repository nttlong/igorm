import { Component } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần thiết cho các directives như *ngIf, *ngFor
import { RouterOutlet } from '@angular/router'; // Cần thiết cho <router-outlet>

@Component({
  selector: 'app-root', // Selector chính của ứng dụng
  templateUrl: './app.component.html', // Trỏ đến file HTML template
  styleUrls: ['./app.component.css'], // Trỏ đến file CSS style
  standalone: true, // Đây là một standalone component
  imports: [
    CommonModule,   // Cung cấp các directives chung của Angular
    RouterOutlet    // Cho phép sử dụng <router-outlet>
  ]
})
export class AppComponent {
  // Không cần thêm code nào ở đây, vì logic UI chính được đặt trong SimpleDashboardComponent
}