// src/components/Layout.tsx

import React from 'react';
import { useTranslation } from 'react-i18next'; // Import useTranslation

interface LayoutProps {
  children: React.ReactNode;
}

const Layout: React.FC<LayoutProps> = ({ children }) => {
  const { t, i18n } = useTranslation();

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };

  return (
    <div className="flex flex-col min-h-screen">
      {/* Header hoặc navigation bar */}
      <header className="bg-gray-800 text-white p-4 flex justify-between items-center">
        <h1 className="text-xl font-semibold">My Multi-Tenant App</h1>
        
        {/* Language Selector */}
        <div className="flex items-center space-x-2">
          <span>{t('language_selector')}:</span> {/* Dịch "Select Language" */}
          <button 
            onClick={() => changeLanguage('en')} 
            className={`px-3 py-1 rounded ${i18n.language === 'en' ? 'bg-blue-600' : 'bg-gray-700'} hover:bg-blue-500`}
          >
            English
          </button>
          <button 
            onClick={() => changeLanguage('vi')} 
            className={`px-3 py-1 rounded ${i18n.language === 'vi' ? 'bg-blue-600' : 'bg-gray-700'} hover:bg-blue-500`}
          >
            Tiếng Việt
          </button>
        </div>
      </header>

      {/* Main content area */}
      <main className="flex-grow">
        {children} {/* Đây là nơi các route content sẽ được render */}
      </main>

      {/* Footer (tùy chọn) */}
      <footer className="bg-gray-800 text-white p-4 text-center">
        &copy; {new Date().getFullYear()} My Multi-Tenant App
      </footer>
    </div>
  );
};

export default Layout;