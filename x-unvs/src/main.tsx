// src/main.tsx

import React from 'react';
import ReactDOM from 'react-dom/client';
import { BrowserRouter } from 'react-router-dom';
import App from './App';
import { AuthProvider } from './context/AuthContext';
import './index.css'; // Giữ file CSS của bạn

// Import đối tượng i18n đã cấu hình của bạn
import i18n from './i18n/i18n'; // <-- Đảm bảo đường dẫn này đúng
import { I18nextProvider } from 'react-i18next'; // <-- Import I18nextProvider

// KHÔNG CẦN import LanguageProvider tùy chỉnh nữa nếu bạn dùng i18next cho việc dịch
// import { LanguageProvider } from './context/LanguageContext';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <AuthProvider>
        {/* Bọc I18nextProvider ở đây, cung cấp đối tượng i18n đã cấu hình */}
        <I18nextProvider i18n={i18n}> 
          <App />
        </I18nextProvider>
      </AuthProvider>
    </BrowserRouter>
  </React.StrictMode>
);