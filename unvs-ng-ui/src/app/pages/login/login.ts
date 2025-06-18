import { Component, OnInit, OnDestroy } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';
import { Subscription } from 'rxjs';
import { TranslateService, TranslateModule } from '@ngx-translate/core';
import { AuthService } from '../../services/auth.service'; // <-- Import AuthService

@Component({
  selector: 'app-login',
  templateUrl: './login.html',
  styleUrls: ['./login.scss'],
  standalone: true,
  imports: [
    CommonModule,
    FormsModule,
    TranslateModule
  ]
})
export class Login implements OnInit, OnDestroy {
  username = '';
  password = '';
  error: string | null = null;
  isLoading = false;
  tenantname: string | null = null;
  private paramSubscription: Subscription | undefined;

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private translate: TranslateService,
    private authService: AuthService // <-- Inject AuthService
  ) { }

  ngOnInit(): void {
    this.paramSubscription = this.activatedRoute.paramMap.subscribe(params => {
      this.tenantname = params.get('tenantname');
      console.log('Login Page - Tenantname:', this.tenantname);
    });

    this.translate.setDefaultLang('en');
    const browserLang = this.translate.getBrowserLang();
    this.translate.use(browserLang?.match(/en|vi/) ? browserLang : 'en');
  }

  ngOnDestroy(): void {
    if (this.paramSubscription) {
      this.paramSubscription.unsubscribe();
    }
  }

  async handleLogin(event: Event): Promise<void> {
    event.preventDefault();
    this.error = null;
    this.isLoading = true;

    if (!this.tenantname) {
      this.error = this.translate.instant('login.error.tenantnameRequired');
      this.isLoading = false;
      return;
    }

    try {
      // Gọi phương thức login từ AuthService
      this.authService.login(this.username, this.password, this.tenantname).subscribe({
        next: (response) => {
          console.log('Đăng nhập thành công:', response);
          // AuthService đã tự lưu token và chuyển hướng
          // Bạn có thể chuyển hướng thêm ở đây nếu muốn logic khác
          this.router.navigate([`/${this.tenantname}/dashboard`]);
        },
        error: (err) => {
          console.error('Lỗi đăng nhập:', err);
          this.error = this.translate.instant('login.error.authenticationFailed');
        }
      }).add(() => {
        this.isLoading = false;
      });
    } catch (err) {
      console.error('Lỗi xảy ra ngoài Observable:', err);
      this.error = this.translate.instant('login.error.networkError');
      this.isLoading = false;
    }
  }
}
