
import { CgOverflow } from 'react-icons/cg';
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
    <div className='dock-full bg-white'>
      <h1 className='debug'>Dashboard</h1>
      <div className='dock-full debug overflow-y-scroll bg-red-900'>
        content
      </div>

    </div>
  );
};

export default DashboardPage;