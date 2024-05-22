import { pageName, user } from '@/app-recoil';
import { TUser } from '@/types';
import { doesHttpOnlyCookieExist } from '@/utils';
import { useEffect } from 'react';
import { useRecoilState, useResetRecoilState } from 'recoil';

interface IAuthWrapper {
  children: React.ReactNode;
}

const AuthWrapper = ({ children }: IAuthWrapper) => {
  const resetData = useResetRecoilState(user);
  const resetPageName = useResetRecoilState(pageName);
  const [_user, _setUser] = useRecoilState<TUser>(user);

  useEffect(() => {
    const interval = setInterval(() => {
      if (!doesHttpOnlyCookieExist('ukama_session')) {
        handleGoToLogin();
      }
    }, 3000);
    return () => clearInterval(interval);
  }, []);

  const handleGoToLogin = () => {
    resetData();
    resetPageName();
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_AUTH_APP_URL || '');
  };

  return <div>{children}</div>;
};

export default AuthWrapper;
