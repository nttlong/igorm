// src/App.tsx

import React, { useState, useEffect } from "react";
import { FiHome, FiSettings, FiUsers, FiFolder } from "react-icons/fi"; // Đảm bảo import đầy đủ icon nếu dùng ở đây

// BỎ BrowserRouter khỏi đây. Chỉ giữ lại Routes, Route, Outlet
import { Routes, Route, Outlet } from 'react-router-dom'; 
import { useTranslation } from 'react-i18next'; 

import Header from "./components/Header";
import AsideBar from "./components/AsideBar";
import type { MenuItem } from "./components/AsideBar"; 

// Import các trang component thực tế
import DashboardPage from './pages/Dashboard'; 
import UsersPage from './pages/UsersPage';
import ProjectsPage from './pages/ProjectsPage';
import SettingsPage from './pages/SettingsPage';
import LoginPage from './pages/Login'; 
import ProtectedRoute from './components/ProtectedRoute'; 

// Component Layout chính của bạn
const MainLayout = () => {
  const { t, i18n } = useTranslation();

  const [isSidebarCollapsed, setIsSidebarCollapsed] = useState(false);
  const [activeMenuItem, setActiveMenuItem] = useState("dashboard"); 
  const [isDarkMode, setIsDarkMode] = useState(false); 
  const [isProfileDropdownOpen, setIsProfileDropdownOpen] = useState(false); 

  // Định nghĩa menuItems. Sử dụng 'label' thay vì 'caption' để khớp với interface MenuItem
  const rawMenuItems: MenuItem[] = [
    { id: "dashboard", label: "Dashboard", path: "dashboard", icon: <FiHome size={20} /> },
    { id: "users", label: "Users", path: "users", icon: <FiUsers size={20} /> },
    { id: "projects", label: "Projects", path: "projects", icon: <FiFolder size={20} /> },
    { id: "settings", label: "Settings", path: "settings", icon: <FiSettings size={20} /> }
  ];

  const menuItems = rawMenuItems.map(item => ({
    ...item,
    label: t(item.id) // Sử dụng i18n để dịch label
  }));

  const handleLogout = () => {
    console.log("Logging out...");
  };

  const toggleDarkMode = () => {
    setIsDarkMode(!isDarkMode);
    document.documentElement.classList.toggle("dark"); 
  };

  const changeLanguage = (lng: string) => {
    i18n.changeLanguage(lng);
  };

  useEffect(() => {
    const handleResize = () => {
      if (window.innerWidth <= 768) {
        setIsSidebarCollapsed(true);
      } else {
        setIsSidebarCollapsed(false); 
      }
    };
    window.addEventListener("resize", handleResize);
    handleResize(); 
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  useEffect(() => {
    const handleClickOutside = (event: MouseEvent) => {
      const profileDropdownElement = document.querySelector(".profile-dropdown");
      if (profileDropdownElement && !profileDropdownElement.contains(event.target as Node)) {
        setIsProfileDropdownOpen(false);
      }
    };
    document.addEventListener("click", handleClickOutside);
    return () => document.removeEventListener("click", handleClickOutside);
  }, []);

  return (
    <div className={`main-layout min-h-screen ${isDarkMode ? "dark" : ""}`}>
      <div className="flex h-screen bg-gray-100 dark:bg-gray-900">
        <Header 
          onSidebarToggle={() => setIsSidebarCollapsed(!isSidebarCollapsed)}
          isDarkMode={isDarkMode}
          toggleDarkMode={toggleDarkMode}

          onLogout={handleLogout}
          isProfileDropdownOpen={isProfileDropdownOpen}
          setIsProfileDropdownOpen={setIsProfileDropdownOpen}
          isSidebarCollapsed={isSidebarCollapsed}
        />

        <AsideBar
          menuItems={menuItems}
          isSidebarCollapsed={isSidebarCollapsed}
          activeMenuItem={activeMenuItem}
          setActiveMenuItem={setActiveMenuItem}
        />

        <main
          className={`flex-1  transition-all duration-300 ${
            isSidebarCollapsed ? "ml-16" : "ml-64"
          }`}
        >
          <Outlet /> 
        </main>
      </div>
    </div>
  );
};

// Component chính của ứng dụng nơi bạn định nghĩa Routes
// BỎ BrowserRouter bao bọc Routes ở đây
function App() {
  return (
    <Routes>
      <Route path="/:tenantname/login" element={<LoginPage />} />

      <Route path="/:tenantname" element={
        <ProtectedRoute>
          <MainLayout />
        </ProtectedRoute>
      }>
        <Route index element={<DashboardPage />} /> 
        <Route path="dashboard" element={<DashboardPage />} />
        <Route path="users" element={<UsersPage />} />
        <Route path="projects" element={<ProjectsPage />} />
        <Route path="settings" element={<SettingsPage />} />
        
        <Route path="*" element={
          <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
            <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-4">404 Not Found within Tenant</h2>
            <p className="text-gray-600 dark:text-gray-300">The page you are looking for does not exist within this tenant's context.</p>
          </div>
        } />
      </Route>

      <Route path="/" element={
        <div className="flex flex-col items-center justify-center h-screen bg-gray-100 dark:bg-gray-900 text-gray-800 dark:text-white">
          <h1 className="text-3xl font-bold">Welcome to Multi-Tenant App</h1>
          <p className="mt-4 text-lg">Please access your tenant's login page.</p>
          <p className="mt-2 text-gray-500">For example: <a href="/acme/login" className="text-blue-500 hover:underline">/acme/login</a></p>
        </div>
      } />
      
      <Route path="*" element={
        <div className="flex flex-col items-center justify-center h-screen bg-gray-100 dark:bg-gray-900 text-gray-800 dark:text-white">
          <h1 className="text-3xl font-bold">404 Not Found</h1>
          <p className="mt-4 text-lg">The requested page does not exist.</p>
        </div>
      } />
    </Routes>
  );
}

export default App;