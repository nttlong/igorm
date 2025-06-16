// src/services/ApiService.ts

// Định nghĩa base URL của API từ biến môi trường
// Đảm bảo bạn đã định nghĩa VITE_API_BASE_URL trong file .env (ví dụ: http://localhost:8080)
const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

/**
 * Định nghĩa interface cho các tham số query chung của API Swagger.
 * Các trường này là bắt buộc theo Swagger bạn cung cấp.
 */
interface ApiQueryParams {
  feature: string; // The specific id of feature. e.g., login
  action: string;   // The specific action to invoke (e.g., login, register, logout)
  module: string;   // The specific module to invoke (e.g., unvs.br.auth.users)
  tenant: string;   // The specific tenant to invoke (e.g., default, name)
  lan: string;      // The specific language to invoke (e.g., en, pt)
}

/**
 * Định nghĩa interface cho các tùy chọn gửi request API.
 * Bao gồm các query params và data (dạng JSON stringify).
 */
interface ApiRequestOptions {
  queryParams: ApiQueryParams;
  data?: string; // JSON stringify from browser
  method?: 'GET' | 'POST' | 'PUT' | 'DELETE'; // Mặc định là POST hoặc GET tùy ngữ cảnh
  headers?: HeadersInit; // Cho phép thêm các header tùy chỉnh
  // Bạn có thể thêm các tùy chọn khác như responseType, timeout, v.v.
}

/**
 * Hàm tiện ích để build URL với các query parameters.
 * @param baseUrl Base URL của API (ví dụ: http://localhost:8080)
 * @param queryParams Đối tượng chứa các tham số query
 * @returns Chuỗi URL đã được build
 */
function buildUrlWithQueryParams(baseUrl: string, queryParams: ApiQueryParams): string {
  const url = new URL(baseUrl+"/api/v1/invoke");
  url.searchParams.append('feature', queryParams.feature);
  url.searchParams.append('action', queryParams.action);
  url.searchParams.append('module', queryParams.module);
  url.searchParams.append('tenant', queryParams.tenant);
  url.searchParams.append('lan', queryParams.lan);
  return url.toString();
}

/**
 * Dịch vụ chung để gửi các yêu cầu API dựa trên cấu trúc Swagger của bạn.
 * @param options Các tùy chọn cho yêu cầu API bao gồm queryParams và data.
 * @returns Promise chứa phản hồi từ API.
 */
export async function callApi(options: ApiRequestOptions): Promise<any> {
  const { queryParams, data, method = 'POST', headers } = options; // Mặc định method là POST

  // Xây dựng URL hoàn chỉnh với query parameters
  const url = buildUrlWithQueryParams(API_BASE_URL, queryParams);

  const defaultHeaders: HeadersInit = {
    'Content-Type': 'application/json',
    // Thêm các header mặc định khác nếu cần, ví dụ: Authorization token
    // 'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
  };

  const requestOptions: RequestInit = {
    method: method,
    headers: {
      ...defaultHeaders,
      ...headers, // Ghi đè hoặc thêm các header tùy chỉnh
    },
  };

  // Thêm body nếu method không phải GET hoặc HEAD và có data
  if (method !== 'GET' && method !== 'HEAD' && data) {
    // Swagger của bạn hiển thị `data` là `string (formData)`,
    // nhưng giá trị example là JSON stringify.
    // Nếu nó thực sự là form data (application/x-www-form-urlencoded hoặc multipart/form-data),
    // bạn cần xử lý khác.
    // Nếu nó là JSON stringify trong body, bạn cần đảm bảo 'Content-Type': 'application/json'.
    // Ở đây, tôi giả định `data` là một chuỗi JSON đã được stringify.
    requestOptions.body = data;
  } else if (method !== 'GET' && method !== 'HEAD' && !data && requestOptions.method === 'POST') {
     // Đảm bảo body là rỗng nếu không có data nhưng là POST
     requestOptions.body = JSON.stringify({});
  }


  console.log(`Calling API: ${method} ${url}`);
  console.log('Request Headers:', requestOptions.headers);
  if (requestOptions.body) {
    console.log('Request Body:', requestOptions.body);
  }

  try {
    const response = await fetch(url, requestOptions);

    if (!response.ok) {
      // Xử lý lỗi HTTP status (4xx, 5xx)
      const errorText = await response.text();
      console.error(`API Error: ${response.status} - ${errorText}`);
      throw new Error(`API request failed with status ${response.status}: ${errorText}`);
    }

    // Cố gắng parse JSON, nếu không được thì trả về text
    try {
      const jsonResponse = await response.json();
      return jsonResponse;
    } catch (e) {
      // Nếu không phải JSON, trả về text hoặc throw lỗi
      const textResponse = await response.text();
      console.warn("API response was not JSON, returning raw text:", textResponse);
      return textResponse;
    }

  } catch (error) {
    console.error('Network or unexpected API error:', error);
    throw error; // Ném lỗi để component gọi có thể bắt
  }
}