import React from 'react';
import { Navigate, useLocation } from 'react-router-dom';
import { useAuth } from '../context/AuthContext';

const ProtectedRoute = ({ children }: { children: JSX.Element }) => {
  const { user, loading } = useAuth();
  const location = useLocation();

  if (loading) {
    return <div>Loading...</div>; // Hoặc một component Spinner đẹp hơn
  }

  if (!user) {
    // Lấy tenant từ URL hiện tại (ví dụ: /acme/dashboard -> 'acme')
    const tenant = location.pathname.split('/')[1];
    // Nếu có tenant, chuyển hướng về trang login của tenant đó.
    // Nếu không, chuyển về trang gốc (có thể là trang chọn tenant).
    const loginPath = tenant ? `/${tenant}/login` : '/';

    return <Navigate to={loginPath} state={{ from: location }} replace />;
  }
  
  return children;
};

export default ProtectedRoute;