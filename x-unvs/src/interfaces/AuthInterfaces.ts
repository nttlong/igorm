// src/interfaces/AuthInterfaces.ts


  
  // Interface chính cho toàn bộ phản hồi đăng nhập
  export interface LoginResponse {
    access_token: string;
    token_type: string;
    expires_in: number;
    scope: string;
    refresh_token: string;
    message: string;
    roleId: string;
    userId: string;
    username: string;
    email: string;
  }