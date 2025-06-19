import React, { useState } from 'react';
import { useNavigate, useLocation, useParams } from 'react-router-dom';
import { useRef, useEffect } from 'react';
import { useAuth } from '../context/AuthContext';
import './Login.css';
import {Caller} from '../utils/Caller'
import { useTranslation } from 'react-i18next';
import type {LoginResponse} from '../interfaces/AuthInterfaces';
import {setBaseApiUrl} from '../utils/Caller'; 
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;

const LoginPage = () => {
  setBaseApiUrl(apiBaseUrl);
  const { t, i18n } = useTranslation();
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();
  const currentLang = i18n.language; 
  // Lấy tenantname từ URL params
  const { tenantname } = useParams<{ tenantname: string }>();

  const [username, setUsername] = useState('');
  const [password, setPassword] = useState('');
  const [error, setError] = useState('');

  const from = location.state?.from?.pathname || `/${tenantname}/dashboard`;

  const handleLogin = async (e: React.FormEvent) => {
    e.preventDefault();
    setError('');

    if (!tenantname) {
      setError(t('login.error.tenantname.required'));
      return;
    }
    try {
      // Gửi cả tenantname khi gọi hàm login
      debugger;
      const ret= await Caller.create(`login@unvs.br.auth.users`)
      .withTenant(tenantname)
      .withLanguage(currentLang)
      .withFeature('login')
      .withData({ username, password }).postAsync<LoginResponse>();
      console.log(ret);
      //save token to local storage
      if (ret.data == null || ret.data.results == null){
        setError(t('Login fail, please check your username or password and try again.'));
        return;
      }
      localStorage.setItem('token', ret.data.results.access_token);
      localStorage.setItem('refresh_token', ret.data.results.refresh_token);
      localStorage.setItem('username', username);
      localStorage.setItem('tenantname', tenantname);
      await login(ret.data.results, tenantname);
      navigate(from, { replace: true });
    } catch (err) {
      setError(t('Login fail, please check your username or password and try again.'));
    }
  };

  return (
    <div className="flex items-center justify-center min-h-screen bg-gray-100">
      <div className="w-full max-w-md p-8 space-y-6 bg-white rounded-xl shadow-lg">
        <h2 className="text-3xl font-bold text-center text-gray-800">
        {t('Login')}
        </h2>
        <form onSubmit={handleLogin} className="space-y-6">
          <div>
          <label htmlFor="username">{t('Username')}</label>
            <input id="username" type="text" value={username} onChange={(e) => setUsername(e.target.value)} required placeholder="Nhập 'admin'" className="w-full px-4 py-2 mt-2 border rounded-md"/>
          </div>
          <div>
            <label htmlFor="password">{t('password')}</label>
            <input id="password" type="password" value={password} onChange={(e) => setPassword(e.target.value)} required placeholder="Nhập '123'" className="w-full px-4 py-2 mt-2 border rounded-md"/>
          </div>
          {error && <p className="text-sm text-center text-red-500">{error}</p>}
          <button type="submit" className="w-full py-3 font-semibold text-white bg-blue-600 rounded-lg">
          {t('Login')}
          </button>
        </form>
      </div>
    </div>
  );
};

export default LoginPage;