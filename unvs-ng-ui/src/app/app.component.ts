import { Component } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần nếu sử dụng các directives chung
import { RouterOutlet } from '@angular/router'; // <-- Quan trọng: Import RouterOutlet

@Component({
  selector: 'app-root', // Selector gốc của ứng dụng Angular
  templateUrl: './app.component.html', // Trỏ đến file HTML (chỉ có <router-outlet>)
  styleUrls: ['./app.component.css'], // Trỏ đến file CSS (có thể trống)
  standalone: true, // <-- Đánh dấu đây là một standalone component
  imports: [
    CommonModule,   // Cung cấp các directives như *ngIf, *ngFor (nếu có)
    RouterOutlet    // <-- Cần thiết để sử dụng <router-outlet> trong template
  ]
})
export class AppComponent {
  // AppComponent không cần logic phức tạp vì nó chỉ là host cho router.
  // Mọi logic của ứng dụng (layout, trang, v.v.) sẽ nằm trong các component được load bởi router.
  constructor() {
    // console.log("AppComponent loaded!"); // Có thể thêm log để debug nếu cần
  }
}
