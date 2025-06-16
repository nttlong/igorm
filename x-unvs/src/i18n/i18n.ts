// src/i18n/i18n.ts (hoặc .js)
import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';
import LanguageDetector from 'i18next-browser-languagedetector'; // Tự động phát hiện ngôn ngữ trình duyệt

// Bạn có thể chia nhỏ các file dịch ra (ví dụ: en.json, vi.json)
// và import chúng ở đây. Đối với ví dụ này, tôi sẽ định nghĩa trực tiếp.
const resources = {
  en: {
    translation: {
      "appName": "My App",
      "logout": "Logout",
      "welcome": "Welcome",
      "loginTitle": "Login",
      // ... thêm các key dịch khác
    },
  },
  vi: {
    translation: {
      "appName": "Ứng dụng của tôi",
      "logout": "Đăng xuất",
      "welcome": "Chào mừng",
      "loginTitle": "Đăng nhập",
      // ... thêm các key dịch khác
    },
  },
};

i18n
  .use(LanguageDetector) // Sử dụng trình phát hiện ngôn ngữ của trình duyệt
  .use(initReactI18next) // Gắn i18n vào React
  .init({
    resources, // Các file dịch của bạn
    fallbackLng: "en", // Ngôn ngữ mặc định nếu không tìm thấy ngôn ngữ của người dùng
    debug: true, // Bật debug để xem log trong console (có thể tắt trong production)

    interpolation: {
      escapeValue: false, // React đã tự động escape XSS
    },
    detection: {
      order: ['queryString', 'cookie', 'localStorage', 'navigator'], // Thứ tự ưu tiên phát hiện ngôn ngữ
      caches: ['localStorage'], // Lưu ngôn ngữ đã chọn vào localStorage
    },
  });

export default i18n;