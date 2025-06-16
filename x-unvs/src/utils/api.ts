// src/utils/api.ts
import { useTranslation } from 'react-i18next'; 
interface ApiResponse<T> {
    success: boolean;
    data?: T;
    error?: string;
    statusCode?: number;
  }
  
  const API_BASE_URL = 'http://localhost:8080/api/v1'; // Thay đổi thành URL API thực tế của bạn
  
  // Hàm trợ giúp để xử lý phản hồi API
  async function handleResponse<T>(response: Response): Promise<ApiResponse<T>> {
    if (!response.ok) {
      const errorData = await response.json().catch(() => ({ message: 'Something went wrong!' }));
      return {
        success: false,
        error: errorData.message || `HTTP error! Status: ${response.status}`,
        statusCode: response.status,
      };
    }
    const data: T = await response.json();
    return { success: true, data };
  }
  
  // Hàm GET request
  export async function get<T>(endpoint: string): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}`, {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
          // Thêm Authorization header nếu cần
          // 'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
        },
      });
      return handleResponse(response);
    } catch (error: any) {
      return { success: false, error: error.message || 'Network error' };
    }
  }
  
  // Hàm POST request
  export async function post<T>(endpoint: string, data: any): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          // 'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
        },
        body: JSON.stringify(data),
      });
      return handleResponse(response);
    } catch (error: any) {
      return { success: false, error: error.message || 'Network error' };
    }
  }
  
  // Hàm PUT request
  export async function put<T>(endpoint: string, data: any): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}`, {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          // 'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
        },
        body: JSON.stringify(data),
      });
      return handleResponse(response);
    } catch (error: any) {
      return { success: false, error: error.message || 'Network error' };
    }
  }
  
  // Hàm DELETE request
  export async function del<T>(endpoint: string): Promise<ApiResponse<T>> {
    try {
      const response = await fetch(`${API_BASE_URL}/${endpoint}`, {
        method: 'DELETE',
        headers: {
          'Content-Type': 'application/json',
          // 'Authorization': `Bearer ${localStorage.getItem('authToken')}`,
        },
      });
      return handleResponse(response);
    } catch (error: any) {
      return { success: false, error: error.message || 'Network error' };
    }
  }
  