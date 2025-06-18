import { Component, Input, Output, EventEmitter } from '@angular/core';
import { CommonModule } from '@angular/common'; // Cần cho ngClass, ngIf nếu dùng
import { FormsModule } from '@angular/forms'; // Cần cho [(ngModel)]



@Component({
  selector: 'unvs-search-input', // Selector của component này
  templateUrl: './search-input.html',
  styleUrls: ['./search-input.scss'], // Có thể trống nếu chỉ dùng Tailwind
  standalone: true, // Đánh dấu là standalone component
  imports: [
    CommonModule,
    FormsModule // Quan trọng: cần FormsModule để sử dụng [(ngModel)]
  ]
})
export class SearchInput {
  // @Input() để nhận placeholder text từ component cha
  @Input() placeholder: string = 'Search...';

  // @Output() để gửi sự kiện tìm kiếm và giá trị input lên component cha
  @Output() onSearch = new EventEmitter<string>();

  // Biến để lưu giá trị nhập vào của input
  searchTerm: string = '';

  constructor() { }

  /**
   * Phương thức được gọi khi giá trị input thay đổi.
   * Emit giá trị hiện tại của searchTerm lên component cha.
   */
  onInputChange(): void {
    this.onSearch.emit(this.searchTerm);
  }
}
