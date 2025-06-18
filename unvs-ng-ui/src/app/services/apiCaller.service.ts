import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Router } from '@angular/router';
import { Observable, BehaviorSubject, of, throwError } from 'rxjs';
import { tap, catchError } from 'rxjs/operators';
import { environment } from '../../environments/environment'; 
import { TranslateService, TranslateModule } from '@ngx-translate/core';
export class ApiCallerAction {
    
    private _action: string;
    private _module: string;
    private _tenant?: string|null;
    private http?: HttpClient;
    private _language?: string|null;
    constructor(public subEndpoint:string,public tenant?:string|null,http?:HttpClient,language?: string|null ) { 
        if (!subEndpoint.includes('@')) {
            throw "Invalid subEndpoint format. It should be in format 'action@module'"
        }
        this._module = subEndpoint.split('@')[1];
        this._action = subEndpoint.split('@')[0];
        this._tenant = tenant;
        this._language = language;
        this.http=http;


    }
    public Call<T>(data : any): Observable<T> {
        //http://localhost:8080/api/v1/invoke?feature=common&action=login&module=unvs.br.auth.users&tenant=default&lan=vi
        let url = environment.apiUrl + `/invoke?feature=${this._module}&action=${this._action}&module=${this._module}&tenant=${this._tenant}&lan=${this._language}`;
        return this.http!.post<T>(url,data);
    }
    public CallAsync<T>(data : any): Promise<T> {

        return new Promise<T>((resolve, reject) => {
            this.Call<T>(data).subscribe(
                (response: T) => {
                    resolve(response);
                },
                (error: any) => {
                    reject(error);
                }
            );
        });
    }
       
}
@Injectable({providedIn: 'root'})
export class ApiCallerService {
    public TenantName?: string|null;
    private _language?: string|null;
constructor(
    private http: HttpClient,
    private router: Router,
    
    private translate: TranslateService
    
    ) {
        this.TenantName = localStorage.getItem('tenantname')
        this._language = translate.currentLang
    }
    public Api(api:string):ApiCallerAction {
        return new ApiCallerAction(api,this.TenantName,this.http,this._language);
    }
    
}
