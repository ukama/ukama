import { useAppContext } from '@/context';
import {
  useGetMemberByUserIdLazyQuery,
  useGetUserLazyQuery,
} from '@/generated';
import { doesHttpOnlyCookieExist } from '@/utils';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

const whoami = async () => {
  return await fetch(`${process.env.NEXT_PUBLIC_API_GW}/get-user`, {
    method: 'GET',
    // cache: 'reload',
    
    credentials: 'include',
  })
    .then((response) => response.text())
    .then((data) => JSON.parse(data))
    .catch((err) => err);
};
interface IAuthWrapper {
  children: React.ReactNode;
}

const AuthWrapper = ({ children }: IAuthWrapper) => {
  const route = useRouter();
  const {
    user,
    setUser,
    skeltonLoading,
    setSkeltonLoading,
    isValidSession,
    setIsValidSession,
  } = useAppContext();

  const [getMember] = useGetMemberByUserIdLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      setUser({
        ...user,
        role: data.getMemberByUserId.role,
      });
    },
  });

  const [getUser, { data: userData, loading: userLoading }] =
    useGetUserLazyQuery({
      fetchPolicy: 'cache-and-network',
      onCompleted: (data) => {
        setUser({
          role: '',
          id: data.getUser.uuid,
          name: data.getUser.name,
          email: data.getUser.email,
        });

        getMember({
          variables: {
            userId: data.getUser.uuid,
          },
        });
      },
    });

  useEffect(() => {
    if (!isValidSession && doesHttpOnlyCookieExist('ukama_session')) {
      whoami()
        .then((res) => {
          console.log(res);
          // setUser({
          //   ...user,
          //   id: res.uuid,
          // });
          setIsValidSession(true);
        })
        .catch(() => handleLogoutAction());
    } else route.push('/unauthorized');
  }, []);

  useEffect(() => {
    if (isValidSession && route.pathname === '/home' && user.id === '') {
      if (userData) {
        // if (
        //   doesHttpOnlyCookieExist('ukama_session') &&
        //   doesHttpOnlyCookieExist('user_session')
        // )
        //   getUser({
        //     variables: {
        //       userId: userId,
        //     },
        //   });
        // if (!orgId && !orgName) {
        //   route.replace('/onboarding', undefined, { shallow: true });
        // } else {
        //   route.replace(route.pathname, undefined, { shallow: true });
        // }
      } else {
        //failed to get userId
        // handleLogoutAction();
        console.log('GO to login');
      }
    }
  }, [isValidSession]);

  const handleGoToLogin = () => {
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_AUTH_APP_URL || '');
  };

  const handleLogoutAction = () => {
    typeof window !== 'undefined' &&
      window.location.replace(`${process.env.NEXT_PUBLIC_AUTH_APP_URL}/logout`);
  };

  return <div>{children}</div>;
};

export default AuthWrapper;
