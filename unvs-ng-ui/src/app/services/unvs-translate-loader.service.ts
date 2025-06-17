import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { TranslateLoader } from '@ngx-translate/core';
import { Observable, of } from 'rxjs';
import { map, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment';
import { Router } from '@angular/router'; // <-- Import Router

@Injectable({
  providedIn: 'root'
})
export class UnvsTranslateLoader implements TranslateLoader {
  private baseApiUrl: string = environment.apiUrl; // Lấy base API URL từ environment

  // Hàm khởi tạo (constructor) nhận cả HttpClient và Router
  constructor(private http: HttpClient, private router: Router) {
    console.log('UnvsTranslateLoader initialized.');
  }

  /**
   * Lấy bản dịch cho ngôn ngữ đã cho từ API backend.
   * Phương thức này sẽ động lấy tenantname từ URL hiện tại.
   * @param lang Ngôn ngữ cần lấy bản dịch (ví dụ: 'en', 'vi')
   * @returns Observable của object bản dịch (key-value)
   */
  getTranslation(lang: string): Observable<any> {
    // Lấy tenantname từ URL hiện tại
    // Ví dụ: /acme/dashboard -> 'acme'
    // hoặc /acme/login -> 'acme'
    const pathSegments = this.router.url.split('/').filter(segment => segment !== '');
    let tenantname: string = 'default'; // Mặc định là 'default'

    // Logic này cố gắng trích xuất tenantname từ URL
    // Cần đảm bảo logic này khớp với cấu trúc route của bạn: /:tenantname/...
    // Ví dụ: Nếu URL là /acme/dashboard, pathSegments[0] là 'acme'.
    // Nếu URL là /login, pathSegments[0] là 'login', tenantname vẫn là 'default'.
    if (pathSegments.length > 0 && pathSegments[0] !== 'login' && pathSegments[0] !== 'api') {
      tenantname = pathSegments[0];
    } else if (pathSegments.length > 1 && pathSegments[1] === 'login') {
      // Trường hợp URL là /:tenantname/login
      tenantname = pathSegments[0];
    }
    // Các trường hợp khác (ví dụ: /login, /), tenantname vẫn là 'default'
    //http://localhost:8080/api/v1/get/default/unvs.common/dictionary/vi?feature=testft&lan=vi
    // Xây dựng baseTranslationsUrl dựa trên tenantname động
    // URL sẽ có dạng: http://localhost:8080/api/v1/get/TENANTNAME/unvs.common/dictionary
    const baseTranslationsUrl = `${this.baseApiUrl}/get/${tenantname}/unvs.common/dictionary`;
    
    // URL đầy đủ sẽ là: http://localhost:8080/api/v1/get/TENANTNAME/unvs.common/dictionary/LANG_CODE?feature=translation&lan=LANG_CODE
    const url = `${baseTranslationsUrl}/${lang}?feature=translation&lan=${lang}`;
    
    console.log(`UnvsTranslateLoader: Fetching translations for ${lang} from: ${url} (Tenant: ${tenantname})`);
    
    return this.http.get(url).pipe(
      map((translation: any) => {
        console.log(`UnvsTranslateLoader: Translations received for ${lang}:`, translation);
        return translation;
      }),
      catchError(error => {
        console.error(`UnvsTranslateLoader: Error loading translations for ${lang} from ${url}:`, error);
        return of({});
      })
    );
  }
}
