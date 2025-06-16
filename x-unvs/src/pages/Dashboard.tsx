
import { useAuth } from '../context/AuthContext';
import { useNavigate, useParams } from 'react-router-dom';

const DashboardPage = () => {
  const { user, logout } = useAuth();
  const { tenantname } = useParams<{ tenantname: string }>();
  const navigate = useNavigate();

  const handleLogout = () => {
    logout();
    // Sau khi đăng xuất, quay về trang login của tenant
    navigate(`/${tenantname}/login`);
  };

  return (
    <div className="flex flex-col items-center justify-center min-h-screen bg-gray-50">
      <div className="p-10 text-center bg-white rounded-xl shadow-md">
        <h1 className="text-4xl font-bold">Dashboard của {tenantname}</h1>
        {user && (
            <p className="mt-4 text-xl">
                Đăng nhập với tư cách: <span className="font-semibold">{user.name}</span>
            </p>
        )}
        <button 
          onClick={handleLogout} 
          className="w-full px-4 py-3 mt-8 font-semibold text-white bg-red-500 rounded-lg"
        >
          Đăng Xuất
        </button>
      </div>
    </div>
  );
};

export default DashboardPage;