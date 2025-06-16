// src/services/AuthService.ts

// Không cần API_BASE_URL ở đây nữa vì nó đã nằm trong ApiService
// const API_BASE_URL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080';

import { callApi } from './ApiService'; // <--- Import callApi từ ApiService

interface LoginResponse {
  token: string;
  user: {
    id: string;
    username: string;
    // ... thêm các trường thông tin người dùng khác
  };
  message?: string;
}

// Định nghĩa cấu trúc dữ liệu mà bạn muốn gửi trong trường 'data' của API call
interface LoginRequestPayload {
  username: string;
  password: string;
  // Các trường khác mà backend của bạn mong đợi trong body request
  // Ví dụ:
  // deviceId?: string;
  // clientType?: string;
  // lan được thêm vào đây, nếu backend đọc nó từ body JSON
  lan?: string;
}

export const loginUser = async (
  username: string,
  password: string,
  tenantName?: string,
  language?: string // Tham số ngôn ngữ
): Promise<LoginResponse> => {
  try {
    console.log('Login API called via ApiService with:', { username, password, tenantName, language });

    // Chuẩn bị payload cho trường 'data' của yêu cầu API
    const requestPayload: LoginRequestPayload = {
      username: username,
      password: password,
    };

    // Thêm ngôn ngữ vào payload nếu có
    if (language) {
      requestPayload.lan = language; // Gắn lan vào payload
    }

    // Chuyển payload thành chuỗi JSON
    const jsonDataString = JSON.stringify(requestPayload);

    // Gọi API bằng cách sử dụng service chung ApiService
    const responseData: LoginResponse = await callApi({
      queryParams: {
        feature: 'login', // Cố định theo yêu cầu Swagger
        action: 'login',  // Cố định theo yêu cầu Swagger
        module: 'unvs.br.auth.users', // Cố định theo yêu cầu Swagger
        tenant: tenantName || 'default', // Sử dụng tenantName hoặc 'default'
        lan: language || 'en', // Sử dụng ngôn ngữ hoặc mặc định 'en'
      },
      method: 'POST', // Theo yêu cầu Swagger cho login
      data: jsonDataString, // Dữ liệu đăng nhập đã JSON stringify
      // Bạn có thể thêm các headers tùy chỉnh nếu cần thiết cho AuthService
      // Ví dụ: headers: { 'X-Custom-Auth': 'SomeValue' }
    });

    // Xử lý phản hồi thành công
    if (responseData && responseData.token) {
      localStorage.setItem('authToken', responseData.token);
    }

    return responseData;
  } catch (error) {
    console.error('Error during login API call via ApiService:', error);
    throw error;
  }
};