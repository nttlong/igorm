// src/pages/UsersPage.tsx

import { useTranslation } from 'react-i18next';

const UsersPage = () => {
  const { t } = useTranslation();
  return (
    <div className="bg-white dark:bg-gray-800 rounded-lg shadow-md p-6">
      <h2 className="text-2xl font-bold text-gray-800 dark:text-white mb-4">
        {t('users')}
      </h2>
      <p className="text-gray-600 dark:text-gray-300">
        {t('users_page_content')}
      </p>
    </div>
  );
};

export default UsersPage;