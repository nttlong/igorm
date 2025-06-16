// src/components/Header.tsx

import React from 'react';
import { FiMenu, FiBell, FiSun, FiMoon, FiLogOut } from "react-icons/fi";
import { FaUserCircle } from "react-icons/fa";
import { useTranslation } from 'react-i18next'; // <-- Vẫn giữ import này
import { useNavigate } from 'react-router-dom'; // <-- Thêm useNavigate để điều hướng

// KHÔNG CẦN import useLanguage từ LanguageContext của bạn ở đây nếu dùng react-i18next
// import { useLanguage } from '../context/LanguageContext'; 

// Khai báo interface cho props của Header
export interface HeaderProps {
  onSidebarToggle: () => void;
  isDarkMode: boolean;
  toggleDarkMode: () => void;
  // XÓA CÁC DÒNG NÀY NẾU DÙNG react-i18next
  // changeLanguage: (lng: string) => void;
  // currentLanguage: string; 
  onLogout: () => void; // onLogout này có thể được truyền từ App để clear auth context
  isProfileDropdownOpen: boolean;
  isSidebarCollapsed: boolean;
  setIsProfileDropdownOpen: (isOpen: boolean) => void;
}

const Header: React.FC<HeaderProps> = ({
  onSidebarToggle,
  
  isDarkMode,
  isProfileDropdownOpen,
  onLogout, // Giữ onLogout nếu bạn muốn App quản lý nó
  setIsProfileDropdownOpen,
  toggleDarkMode,
  // XÓA CÁC PROPS NÀY NẾU DÙNG react-i18next
  // changeLanguage, 
  // currentLanguage,
}) => {
  const { t, i18n } = useTranslation(); // <-- Lấy t và i18n từ useTranslation
  const navigate = useNavigate(); // <-- Sử dụng useNavigate

  // Logic handleLogout nội bộ của Header
  const handleInternalLogout = () => {
    // Gọi prop onLogout để cha (App) xử lý logout từ AuthContext
    onLogout(); 
    navigate('/login'); // Điều hướng mượt mà hơn
    setIsProfileDropdownOpen(false);
  };

  return (
    <header className="fixed top-0 z-50 w-full bg-white dark:bg-gray-800 shadow-md">
      <div className="flex items-center justify-between px-4 py-3">
        <div className="flex items-center space-x-4">
          <button
            onClick={onSidebarToggle}
            className="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"
            aria-label="Toggle sidebar"
          >
            <FiMenu size={24} className="text-gray-600 dark:text-gray-300" />
          </button>
          <h1 className="text-xl font-bold text-gray-800 dark:text-white">{t('appName')}</h1>
        </div>
        <div className="flex items-center space-x-4">
          {/* Language Selector */}
          <div className="flex items-center space-x-2 text-gray-600 dark:text-gray-300">
            <button
              onClick={() => i18n.changeLanguage('en')} // <-- Dùng i18n.changeLanguage
              className={`px-2 py-1 rounded-md text-sm ${i18n.language === 'en' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}
            >
              EN
            </button>
            <button
              onClick={() => i18n.changeLanguage('vi')} // <-- Dùng i18n.changeLanguage
              className={`px-2 py-1 rounded-md text-sm ${i18n.language === 'vi' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}
            >
              VI
            </button>
            {/* Bạn có thể thêm các ngôn ngữ khác nếu muốn */}
          </div>

          {/* Dark Mode Toggle */}
          <button
            onClick={toggleDarkMode}
            className="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"
            aria-label="Toggle dark mode"
          >
            {isDarkMode ? (
              <FiSun size={24} className="text-gray-300" />
            ) : (
              <FiMoon size={24} className="text-gray-600" />
            )}
          </button>
          <button
            className="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"
            aria-label="Notifications"
          >
            <FiBell size={24} className="text-gray-600 dark:text-gray-300" />
          </button>
          <div className="relative profile-dropdown">
            <button
              onClick={() => setIsProfileDropdownOpen(!isProfileDropdownOpen)}
              className="flex items-center space-x-2"
              aria-label="User profile"
            >
              <FaUserCircle size={32} className="text-gray-600 dark:text-gray-300" />
            </button>
            {isProfileDropdownOpen && (
              <div className="absolute right-0 mt-2 w-48 bg-white dark:bg-gray-700 rounded-md shadow-lg py-1 z-50">
                <button
                  onClick={handleInternalLogout} // <-- Gọi hàm logout nội bộ
                  className="w-full text-left px-4 py-2 text-sm text-gray-700 dark:text-gray-200 hover:bg-gray-100 dark:hover:bg-gray-600 flex items-center"
                >
                  <FiLogOut size={18} className="mr-2" /> {t('logout')}
                </button>
              </div>
            )}
          </div>
        </div>
      </div>
    </header>
  );
};

export default Header;