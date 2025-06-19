import React, { forwardRef, useImperativeHandle, useRef, type PropsWithChildren } from 'react';
import { useNavigate, useLocation, useParams } from 'react-router-dom';

import { useTranslation } from 'react-i18next';
const apiBaseUrl = import.meta.env.VITE_API_BASE_URL;

export type IBaseComponent = {
  callAPIAsync: (apiPath:string,data:any) => any;
  setFeatureId: (id: string) => void;
  getFeatureId: () => any;
  getTenantName: () => string|undefined;
  getLanguage: () => string;
};

type BaseComponentProps = PropsWithChildren<{}>;

const BaseComponent = forwardRef<IBaseComponent, BaseComponentProps>((props, ref) => {
  const { children } = props;
  const { i18n } = useTranslation();
  const hiddenInputRef = useRef<HTMLInputElement>(null);
  const { tenantname } = useParams<{ tenantname: string }>();
  const currentLang = i18n.language;
  const getLanguage = ():string => {
    return currentLang;
  };
  const callAPIAsync = async (path:string,data:any) =>  {
    if (!path.includes("@")){
      console.error("BaseComponent: Invalid API path, it should contain '@' symbol")
    }
    let action=path.split("@")[0]
    let module=path.split("@")[1]
    let lang=i18n.language
    //http://localhost:8080/api/v1/invoke?feature=common&action=list&module=unvs.br.auth.users&tenant=default&lan=vi
    let url = `${apiBaseUrl}/invoke?feature=${getFeatureId()}&action=${action}&module=${module}&tenant=${tenantname}&lan=${lang}`
    const token = localStorage.getItem('token'); 
    console.log(url);
    try {
      const response = await fetch(url, {
        method: 'POST', // hoặc GET nếu cần
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${token}`, // Thêm Bearer Token ở đây
        },
        body: JSON.stringify(data), // data gửi đi
      });
  
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`);
      }
  
      const result = await response.json();
      console.log('API Response:', result);
      return result;
    } catch (error) {
      console.error('API Call failed:', error);
    }
  };
  const getTenantName = ():string|undefined => {
    return tenantname;
  };
  const setFeatureId = (id: string) => {
    if (hiddenInputRef.current) {
      hiddenInputRef.current.value = id;
      console.log('Hidden input set to:', id);
    }
  };
  const getFeatureId = ():string|undefined => {
    return hiddenInputRef.current?.value;
  };

  const initData:IBaseComponent= {
    callAPIAsync: callAPIAsync,
    getFeatureId: getFeatureId,
    setFeatureId: setFeatureId,
    getTenantName: getTenantName,
    getLanguage: getLanguage,
  }
  useImperativeHandle(ref, () => initData);

  return (
    
     
      <div className='dock-full'>
      <input type='hidden' ref={hiddenInputRef} />
      {children}
      </div>
    
  );
});

export default BaseComponent;
