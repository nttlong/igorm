import React, { useState, useContext, createContext, useEffect } from 'react';
import type { LoginResponse } from '../interfaces/AuthInterfaces';
// Định nghĩa kiểu dữ liệu cho user và context
interface AuthContextType {
  user: { name: string } | null;
  // Cập nhật hàm login để nhận tenant
  login: (userData: LoginResponse, tenant: string) => Promise<void>;
  logout: () => void;
  loading: boolean;
}

const AuthContext = createContext<AuthContextType | null>(null);

export const AuthProvider = ({ children }: { children: React.ReactNode }) => {
  const [user, setUser] = useState<{ name: string } | null>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    try {
      const storedUser = localStorage.getItem('user');
      if (storedUser) {
        setUser(JSON.parse(storedUser));
      }
    } catch (error) {
      console.error("Failed to parse user from localStorage", error);
    } finally {
        setLoading(false);
    }
  }, []);

  // Hàm đăng nhập - giờ sẽ nhận thêm tenant
  const login = async (userData: { username: string }, tenant: string) => {
    console.log(`Logging in user ${userData.username} for tenant ${tenant}`);
    // Trong thực tế, bạn sẽ gửi cả thông tin user và tenant lên server để xác thực
    const loggedInUser = { name: userData.username };
    localStorage.setItem('user', JSON.stringify(loggedInUser));
    setUser(loggedInUser);
  };

  // Hàm đăng xuất
  const logout = () => {
    localStorage.removeItem('user');
    setUser(null);
  };

  const value = { user, login, logout, loading };

  return (
    <AuthContext.Provider value={value}>
      {children}
    </AuthContext.Provider>
  );
};

// Hook tùy chỉnh để dễ dàng sử dụng context
export const useAuth = () => {
  const context = useContext(AuthContext);
  if (!context) {
    throw new Error('useAuth must be used within an AuthProvider');
  }
  return context;
};
