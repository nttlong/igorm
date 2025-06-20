import { useRef, useEffect,useState } from 'react';
import React from 'react';
import { useTranslation } from 'react-i18next';
import SearchBox from '../components/SearchBox';
import BaseComponent, { type IBaseComponent } from '../components/BaseComponent';
import { VirtualScroller } from 'primereact/virtualscroller';
const UsersPage = () => {
  const { t } = useTranslation();
  const baseRef = useRef<IBaseComponent>(null);
  const [users, setUsers] = useState([]);
  const getUsers= async ()=>{
    debugger
    const res = await baseRef.current?.callAPIAsync("list@unvs.br.auth.users", {
      pageIndex: 0,
      pageSize: 90
    });
    setUsers(res.results || []); // nếu API kiểu { results: [...] }
  }
  const userItem=(user:any, index:number)=>{
    return <div key={user.userId || index} className="p-4 bg-white rounded shadow">
    <p><strong>{t('Username')}:</strong> {user.username}</p>
    <p><strong>{t('Email')}:</strong> {user.email || 'N/A'}</p>
    <p><strong>{t('Create by')}:</strong> {user.createdBy}</p>
  </div>
  }
  

  useEffect(() => {
    // Gọi hàm callAPI() trong BaseComponent nếu cần
    if (baseRef.current) {
      baseRef.current.setFeatureId('users-manager');
      getUsers();
     
    }
  }, []);

  return (
    <BaseComponent ref={baseRef}>
      <div className='dock-full bg-white'>
      <h1 className=''>
        <SearchBox />
      </h1>
      <div className='dock-full  overflow-y-scroll'>
      <div className='grid grid-cols-3 gap-2'>
        
          {users.map((user:any, index) => (
            userItem(user, index)
          ))}
        </div>
      </div>

    </div>
    </BaseComponent>
  );
};

export default UsersPage;
