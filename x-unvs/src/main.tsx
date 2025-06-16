// src/main.tsx

import React from 'react';
import ReactDOM from 'react-dom/client';
import App from './App.tsx';
import './index.css';
import { BrowserRouter, Routes, Route } from 'react-router-dom';
import Login from './pages/Login.tsx';

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <BrowserRouter>
      <Routes>
        {/* Route động cho trang Login với tenant-name */}
        <Route path="/:tenantName/login" element={<Login />} />
        <Route path="/" element={<App />} />
        {/* Các route khác của bạn */}
      </Routes>
    </BrowserRouter>
  </React.StrictMode>,
);