// src/components/AsideBar.tsx

import React from 'react'; // Cần import React để sử dụng JSX (cho icon) và React.ReactNode
import { Link } from 'react-router-dom'; // Import Link từ react-router-dom

// Định nghĩa kiểu dữ liệu cho một MenuItem (giữ nguyên từ file constants/menuItems.ts)
export interface MenuItem {
  id: string;
  label: string; // Đã thống nhất dùng 'caption'
  path: string;
  icon?: React.ReactNode;
  // Thêm các thuộc tính khác nếu cần
}

// Định nghĩa kiểu dữ liệu cho Props của AsideBar
export interface AsideBarProps {
  menuItems: MenuItem[]; // Mảng các mục menu
  isSidebarCollapsed: boolean; // Trạng thái collapse của sidebar
  activeMenuItem: string; // Mục menu đang active
  setActiveMenuItem: (id: string) => void; // Hàm để set mục active
}

const AsideBar: React.FC<AsideBarProps> = ({
  menuItems,
  isSidebarCollapsed,
  activeMenuItem,
  setActiveMenuItem,
}) => {
  return (
    <aside
    className={`fixed left-0 top-16 h-[calc(100vh-4rem)] bg-white shadow-lg transition-all duration-300 ${
      isSidebarCollapsed ? "w-16" : "w-64"
    } overflow-hidden z-40`}
    >
      <nav className="mt-4">
        {menuItems.map((item) => (
          <Link // Sử dụng Link thay vì button để điều hướng
            to={item.path}
            key={item.id}
            onClick={() => setActiveMenuItem(item.id)} // Cập nhật trạng thái active
            className={`w-full flex items-center px-4 py-3 transition-colors ${
              activeMenuItem === item.id
                ? "bg-blue-100 text-blue-600"
                : "text-gray-600 hover:bg-gray-100"
            }`}
          >
            <span className="flex-shrink-0">{item.icon}</span>
            <span
              className={`ml-4 whitespace-nowrap transition-opacity duration-300 ${
                isSidebarCollapsed ? "opacity-0" : "opacity-100"
              }`}
            >
              {item.label} {/* Sử dụng item.caption */}
            </span>
          </Link>
        ))}
      </nav>
    </aside>
  );
};

export default AsideBar;