// src/components/MainLayout.tsx (Hoặc Layout.tsx, tùy bạn đặt tên)

import React, { useState, useEffect } from "react";
import { FiMenu, FiHome, FiSettings, FiUsers, FiFolder, FiBell, FiSun, FiMoon } from "react-icons/fi";
import { FaUserCircle } from "react-icons/fa";
import { useTranslation } from 'react-i18next'; // <-- Thêm dòng này
import { Outlet, Link } from 'react-router-dom'; // <-- Thêm Outlet và Link

const MainLayout = () => {
  const { t, i18n } = useTranslation(); // <-- Khởi tạo useTranslation
  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [activeMenuItem, setActiveMenuItem] = useState("dashboard");
  const [isDarkMode, setIsDarkMode] = useState(false);
  const [isMobileMenuOpen, setIsMobileMenuOpen] = useState(false); // Chưa dùng nhưng giữ lại

  const menuItems = [
    // Sử dụng t() để dịch các label
    { id: "dashboard", label: t("dashboard"), icon: <FiHome size={20} />, path: "dashboard" },
    { id: "users", label: t("users"), icon: <FiUsers size={20} />, path: "users" },
    { id: "projects", label: t("projects"), icon: <FiFolder size={20} />, path: "projects" },
    { id: "settings", label: t("settings"), icon: <FiSettings size={20} />, path: "settings" }
  ];

  useEffect(() => {
    const handleResize = () => {
      // Chỉ collapse sidebar nếu chiều rộng màn hình nhỏ hơn hoặc bằng 768px
      if (window.innerWidth <= 768) {
        setIsSidebarCollapsed(true);
      } else {
        setIsSidebarCollapsed(false); // Mở sidebar khi màn hình lớn hơn
      }
    };
    window.addEventListener("resize", handleResize);
    handleResize(); // Gọi một lần khi component mount
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const toggleDarkMode = () => {
    setIsDarkMode(!isDarkMode);
    document.documentElement.classList.toggle("dark");
  };

  const changeLanguage = (lng: string) => { // <-- Hàm đổi ngôn ngữ
    i18n.changeLanguage(lng);
  };

  return (
    <div className={`min-h-screen ${isDarkMode ? "dark" : ""}`}>
      <div className="flex h-screen bg-gray-100 dark:bg-gray-900">
        {/* Toolbar */}
        <header className="fixed top-0 z-50 w-full bg-white dark:bg-gray-800 shadow-md">
          <div className="flex items-center justify-between px-4 py-3">
            <div className="flex items-center space-x-4">
              <button
                onClick={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
                className="p-2 rounded-lg hover:bg-gray-200 dark:hover:bg-gray-700"
                aria-label="Toggle sidebar"
              >
                <FiMenu size={24} className="text-gray-600 dark:text-gray-300" />
              </button>
              <h1 className="text-xl font-bold text-gray-800 dark:text-white">AppName</h1>
            </div>
            <div className="flex items-center space-x-4">
              {/* Language Selector */}
              <div className="flex items-center space-x-2 text-gray-600 dark:text-gray-300">
                <button
                  onClick={() => changeLanguage('en')}
                  className={`px-2 py-1 rounded-md text-sm ${i18n.language === 'en' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}
                >
                  EN
                </button>
                <button
                  onClick={() => changeLanguage('vi')}
                  className={`px-2 py-1 rounded-md text-sm ${i18n.language === 'vi' ? 'bg-blue-600 text-white' : 'hover:bg-gray-200 dark:hover:bg-gray-700'}`}
                >
                  VI
                </button>
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
              <div className="relative">
                <button
                  className="flex items-center space-x-2"
                  aria-label="User profile"
                >
                  <FaUserCircle size={32} className="text-gray-600 dark:text-gray-300" />
                </button>
              </div>
            </div>
          </div>
        </header>

        {/* Sidebar */}
        <aside
          className={`fixed left-0 top-16 h-[calc(100vh-4rem)] bg-white dark:bg-gray-800 shadow-lg transition-all duration-300 ${
            isSidebarCollapsed ? "w-16" : "w-64"
          } overflow-hidden z-40`}
        >
          <nav className="mt-4">
            {menuItems.map((item) => (
              // Sử dụng Link từ react-router-dom
              <Link
                to={item.path} // Đường dẫn tương đối với route cha
                key={item.id}
                onClick={() => setActiveMenuItem(item.id)}
                className={`w-full flex items-center px-4 py-3 transition-colors ${
                  activeMenuItem === item.id
                    ? "bg-blue-100 dark:bg-blue-900 text-blue-600 dark:text-blue-300"
                    : "text-gray-600 dark:text-gray-300 hover:bg-gray-100 dark:hover:bg-gray-700"
                }`}
              >
                <span className="flex-shrink-0">{item.icon}</span>
                <span
                  className={`ml-4 transition-opacity duration-300 ${
                    isSidebarCollapsed ? "opacity-0" : "opacity-100"
                  }`}
                >
                  {item.label}
                </span>
              </Link>
            ))}
          </nav>
        </aside>

        {/* Main Content (đặt Outlet ở đây) */}
        <main
          className={`flex-1 mt-16 p-6 transition-all duration-300 ${
            isSidebarCollapsed ? "ml-16" : "ml-64"
          }`}
        >
          {/* Outlet sẽ render các component con của route */}
          <Outlet />
        </main>
      </div>
    </div>
  );
};

export default MainLayout;