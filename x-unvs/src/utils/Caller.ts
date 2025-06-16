// src/utils/Caller.ts

import axios, { type AxiosInstance, type AxiosResponse, type AxiosError } from 'axios';
// import type  AxiosInstance, AxiosResponse, AxiosError  from 'axios';
// Định nghĩa các kiểu dữ liệu cho phản hồi API
export interface ApiResult<T> {
  results?: T
}
export interface ApiResponse<T> {
  success: boolean;
  data?: ApiResult<T>;
  error?: string;
  statusCode?: number;
}

// Kiểu dữ liệu cho các tùy chọn cấu hình yêu cầu
interface RequestOptions {
  method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH';
  url: string;
  data?: any;
  params?: any;
  headers?: Record<string, string>;
}

// Class Caller để xây dựng và thực hiện các yêu cầu API
export class Caller {
  private axiosInstance: AxiosInstance;
  private _apiPath: string = '';
  
  private _method: 'GET' | 'POST' | 'PUT' | 'DELETE' | 'PATCH' = 'GET';
  private _data: any = null;
  private _params: any = null;
  private _headers: Record<string, string> = {};
  private _feature:string='';
  private _tenant:string='';
  private _lang:string='';


  constructor(instance: AxiosInstance) {
    this.axiosInstance = instance;
    this._headers['Content-Type'] = 'application/json'; // Mặc định
  }

  // Phương thức static để tạo một instance mới của Caller
  public static create(apiPath: string): Caller {
    // Đảm bảo axiosInstance được cấu hình ở một nơi tập trung
    // Ví dụ: import { apiInstance } from './axiosInstance';
    // Hoặc bạn có thể truyền thẳng base URL và Caller tự tạo axiosInstance
    // Để đơn giản, ở đây ta giả sử đã có một axios instance dùng chung
    if (!Caller._instance) {
      Caller._instance = new Caller(axios.create({
        baseURL: 'http://localhost:3000/api', // Thay thế bằng URL API thực tế của bạn
        timeout: 10000, // Timeout mặc định
      }));

      // Thêm interceptors nếu chưa có (chỉ làm một lần cho instance dùng chung)
      Caller._instance.axiosInstance.interceptors.request.use(config => {
        const token = localStorage.getItem('authToken');
        if (token) {
          config.headers.Authorization = `Bearer ${token}`;
        }
        return config;
      }, Promise.reject);

      Caller._instance.axiosInstance.interceptors.response.use(response => response, error => {
        if (error.response && error.response.status === 401) {
          console.error('Unauthorized: Token expired or invalid. Redirecting to login.');
          // Logic redirect ở đây, ví dụ: window.location.href = '/login';
        }
        return Promise.reject(error);
      });
    }
    // Reset trạng thái của builder cho mỗi lần gọi 'create'
    const newCaller = new Caller(Caller._instance.axiosInstance);
    newCaller._apiPath = apiPath;
    return newCaller;
  }

  // Dùng singleton pattern cho axios instance bên trong Caller nếu bạn muốn
  private static _instance: Caller;

  

  // Setter cho dữ liệu gửi đi (cho POST, PUT, PATCH)
  public withData(data: any): this {
    this._data = data;
    return this;
  }
  public withLanguage(lang:string): this {
    this._lang=lang;
    return this;
  }

  // Setter cho query parameters (cho GET)
  public withParams(params: any): this {
    this._params = params;
    return this;
  }
 public withTenant(tenant:string) {
    this._tenant=tenant;
    return this;
  
 }
  // Setter cho headers tùy chỉnh
  public withHeaders(headers: Record<string, string>): this {
    this._headers = { ...this._headers, ...headers };
    return this;
  }
  public withFeature(feature:string): this {
    this._feature = feature;
    return this;
  }

  // Các phương thức HTTP (trả về Promise của ApiResponse)
  public async getAsync<T>(): Promise<ApiResponse<T>> {
    this._method = 'GET';
    return this.executeRequest<T>();
  }

  public async postAsync<T>(): Promise<ApiResponse<T>> {
    
    this._method = 'POST';
    return this.executeRequest<T>();
  }

  public async putAsync<T>(): Promise<ApiResponse<T>> {
    this._method = 'PUT';
    return this.executeRequest<T>();
  }

  public async deleteAsync<T>(): Promise<ApiResponse<T>> {
    this._method = 'DELETE';
    return this.executeRequest<T>();
  }

  public async patchAsync<T>(): Promise<ApiResponse<T>> {
    this._method = 'PATCH';
    return this.executeRequest<T>();
  }
  // Phương thức nội bộ để thực hiện yêu cầu
  private async executeRequest<T>(): Promise<ApiResponse<T>> {
    // let url = this._apiPath;
    debugger;
    const module=this._apiPath.split('@')[1];
    const action:string=this._apiPath.split('@')[0];
    const lang=this._lang;
    const url=`http://localhost:8080/api/v1/invoke?module=${module}&action=${action}&feature=${this._feature}&tenant=${this._tenant}&lan=${lang}`;
    console.log(url);

    const requestOptions: RequestOptions = {
      method: this._method,
      url: url,
      headers: this._headers,
    };

    if (this._data) {
      requestOptions.data = this._data;
    }
    if (this._params) {
      requestOptions.params = this._params;
    }
    
    try {
      const response: AxiosResponse<T> = await this.axiosInstance.request(requestOptions);
     
      return { success: true, data: response.data, statusCode: response.status };
    } catch (error: any) {
      if (axios.isAxiosError(error)) {
        return {
          success: false,
          error: error.response?.data?.message || error.message || 'Request failed',
          statusCode: error.response?.status,
        };
      }
      return { success: false, error: 'An unexpected error occurred' };
    } finally {
        // Reset trạng thái sau mỗi yêu cầu để instance có thể tái sử dụng
        // Hoặc bạn có thể tạo instance mới mỗi lần gọi 'create' nếu không muốn reset
        
        this._data = null;
        this._params = null;
        this._method = 'GET'; // Reset về mặc định GET
        this._headers = { 'Content-Type': 'application/json' }; // Reset headers
    }
  }
}

// Export Caller class để có thể sử dụng static method `create`
export default Caller;