import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { TranslateService, TranslateModule } from '@ngx-translate/core';
import { Subscription } from 'rxjs';

@Component({
  selector: 'app-login-page',
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
  private langChangeSubscription: Subscription | undefined; // Thêm subscription cho langChange

  constructor(
    private router: Router,
    private activatedRoute: ActivatedRoute,
    private translate: TranslateService
  ) {
    console.log('LoginPageComponent constructor called.');
  }

  ngOnInit(): void {
    this.paramSubscription = this.activatedRoute.paramMap.subscribe(params => {
      this.tenantname = params.get('tenantname');
      console.log('Login Page - Tenantname:', this.tenantname);
    });

    // Thiết lập ngôn ngữ mặc định và sử dụng ngôn ngữ trình duyệt
    this.translate.setDefaultLang('en');
    const browserLang = this.translate.getBrowserLang();
    const initialLang = browserLang?.match(/en|vi/) ? browserLang : 'en';
    this.translate.use(initialLang); // Sử dụng ngôn ngữ ban đầu

    console.log(`TranslateService: Default lang set to ${this.translate.defaultLang}`);
    console.log(`TranslateService: Current lang set to ${this.translate.currentLang}`);

    // Lắng nghe sự kiện khi ngôn ngữ thay đổi (và khi bản dịch được tải)
    this.langChangeSubscription = this.translate.onLangChange.subscribe(() => {
      console.log(`TranslateService: Language changed to ${this.translate.currentLang}`);
      // Kiểm tra xem bản dịch cho 'login.title' đã có sẵn chưa
      const translatedTitle = this.translate.instant('login.title');
      console.log(`Translated 'login.title' (after lang change):`, translatedTitle);
      if (translatedTitle === 'login.title') {
        console.warn("Translation for 'login.title' still returns the key. Check your translation files/API response.");
      }
    });

    // Thử lấy bản dịch ngay sau khi thiết lập ngôn ngữ ban đầu
    const initialTranslatedTitle = this.translate.instant('login.title');
    console.log(`Translated 'login.title' (initial):`, initialTranslatedTitle);
    if (initialTranslatedTitle === 'login.title') {
      console.warn("Translation for 'login.title' initially returns the key. This might be because the translation file hasn't loaded yet.");
    }
  }

  ngOnDestroy(): void {
    if (this.paramSubscription) {
      this.paramSubscription.unsubscribe();
    }
    if (this.langChangeSubscription) { // Hủy bỏ subscription cho langChange
      this.langChangeSubscription.unsubscribe();
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
      await new Promise(resolve => setTimeout(resolve, 1500));

      if (this.username === 'admin' && this.password === '123') {
        console.log('Đăng nhập thành công!');
        localStorage.setItem('authToken', 'mock_angular_token_123');
        localStorage.setItem('username', this.username);
        localStorage.setItem('tenantname', this.tenantname);

        this.router.navigate([`/${this.tenantname}/dashboard`]);
      } else {
        this.error = this.translate.instant('login.error.authenticationFailed');
      }
    } catch (err) {
      console.error('Lỗi đăng nhập:', err);
      this.error = this.translate.instant('login.error.networkError');
    } finally {
      this.isLoading = false;
    }
  }
}
