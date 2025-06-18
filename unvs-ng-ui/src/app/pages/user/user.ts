import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule, DatePipe, JsonPipe } from '@angular/common'; // <-- Cần cho *ngIf, date pipe, json pipe
import { FormsModule } from '@angular/forms'; // <-- Cần cho [(ngModel)]
import { ActivatedRoute } from '@angular/router'; // Để lấy tenantname
import { Observable, Subscription } from 'rxjs'; // Thêm Subscription

import { ApiCallerService } from '../../services/apiCaller.service'; // <-- Import ApiCallerService
import {SearchInput} from '../../components/common/search-input/search-input'
// Giả định cấu trúc dữ liệu người dùng của bạn
interface UserData {
  username: string;
  userId: string;
  isLocked: boolean;
  email: string;
  createdAt: string; // Hoặc Date
  createdBy: string;
  isSupperUser: boolean;
  showDetails?: boolean; // Thêm trường này cho UI
}

// Giả định cấu trúc phản hồi API của bạn
interface ApiResponse<T> {
  results: T; // API trả về một mảng user
}

@Component({
  selector: 'app-user', // Selector của component này
  templateUrl: './user.html',
  styleUrls: ['./user.scss'], // Đảm bảo đúng styleUrl
  standalone: true, // <-- ĐÂY LÀ ĐIỂM QUAN TRỌNG: Định nghĩa là standalone component
  imports: [
    CommonModule,   // <-- CUNG CẤP *ngIf và các directive cơ bản khác
    FormsModule,    // Cần cho [(ngModel)]
    DatePipe,       // <-- CUNG CẤP pipe 'date'
    JsonPipe,        // <-- CUNG CẤP pipe 'json'
    SearchInput
  ],
  providers: [] // Không cần ApiCallerService ở đây nếu nó providedIn: 'root'
})
export class User implements OnInit, OnDestroy { // Implement OnInit, OnDestroy
  data: ApiResponse<UserData[]> | null = null;
  users: UserData[] = [];
  filteredUsers: UserData[] = [];
  searchTerm: string = '';
  
  private tenantname: string | null = null;
  private paramSubscription: Subscription | undefined;

  constructor(
    private apiCallerService: ApiCallerService, // <-- Inject ApiCallerService
    private activatedRoute: ActivatedRoute
  ) { }

  ngOnInit(): void {
    // Lấy tenantname từ URL params
    this.paramSubscription = this.activatedRoute.paramMap.subscribe(params => {
      this.tenantname = params.get('tenantname');
      console.log('Users Page - Tenantname:', this.tenantname);
      // Gọi API để lấy danh sách người dùng sau khi có tenantname
      this.getListOfUsers(); 
    });
  }

  ngOnDestroy(): void {
    if (this.paramSubscription) {
      this.paramSubscription.unsubscribe();
    }
  }
  private async getListOfUsersAsync() {
    await  this.apiCallerService.Api("list@unvs.br.auth.users").CallAsync<ApiResponse<UserData[]>>({})
  }
  private getListOfUsers(): void {
    this.apiCallerService.Api("list@unvs.br.auth.users")
      .Call<ApiResponse<any[]>>({})
      .subscribe({
        next: (response) => {
          this.data = response;
          if (Array.isArray(response.results)) {
            this.users = response.results.map(user => ({ ...user, showDetails: false }));
            this.filterUsers();
            console.log('User data loaded:', this.data);
          } else {
            console.warn('API response.results is not an array:', response.results);
            this.users = [];
          }
        },
        error: (error:any) => {
          console.error('Error fetching users:', error);
          this.data = null;
          this.users = [];
          this.filteredUsers = [];
        }
      });
  }
  onSearch(evt: Event): void {
    // this.searchTerm = searchTerm;
    this.filterUsers();
  }

  filterUsers(): void {
    if (!this.searchTerm) {
      this.filteredUsers = [...this.users];
    } else {
      const lowerCaseSearchTerm = this.searchTerm.toLowerCase();
      this.filteredUsers = this.users.filter(user =>
        user.username.toLowerCase().includes(lowerCaseSearchTerm) ||
        user.email.toLowerCase().includes(lowerCaseSearchTerm) ||
        user.userId.toLowerCase().includes(lowerCaseSearchTerm)
      );
    }
  }

  copyToClipboard(text: string): void {
    const el = document.createElement('textarea');
    el.value = text;
    document.body.appendChild(el);
    el.select();
    try {
      document.execCommand('copy');
      console.log('Copied to clipboard:', text);
      alert('Copied to clipboard: ' + text);
    } catch (err) {
      console.error('Failed to copy text: ', err);
    }
    document.body.removeChild(el);
  }

  toggleDetails(user: UserData): void {
    user.showDetails = !user.showDetails;
  }
}