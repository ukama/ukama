import { useAppContext } from '@/context';
import { Role_Type, useGetTokenLazyQuery } from '@/generated';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import { doesHttpOnlyCookieExist } from '@/utils';
import { Box } from '@mui/material';
import { useRouter } from 'next/router';
import { useEffect } from 'react';

// const whoami = async () => {
//   return await fetch(`${process.env.NEXT_PUBLIC_API_GW}/get-user`, {
//     method: 'GET',
//     // cache: 'reload',

//     credentials: 'include',
//   })
//     .then((response) => response.text())
//     .then((data) => JSON.parse(data))
//     .catch((err) => err);
// };
interface IAuthWrapper {
  children: React.ReactNode;
}

const AuthWrapper = ({ children }: IAuthWrapper) => {
  const route = useRouter();
  const { token, setToken, user, setUser, isValidSession, setIsValidSession } =
    useAppContext();

  const [getToken, { loading: getTokenLoading }] = useGetTokenLazyQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setUser({
        ...user,
        id: data.getToken.userId,
        email: data.getToken.email,
        name: data.getToken.name,
        role: data.getToken.role,
        orgId: data.getToken.orgId,
        orgName: data.getToken.orgName,
      });
      setToken(data.getToken.token);
      setIsValidSession(true);
      if (!data.getToken.token) {
        route.push('/unauthorized');
      } else if (data.getToken.role === Role_Type.None) {
        route.push('/onboarding');
      } else {
        route.push(`${route.pathname ? route.pathname : '/home'}`);
      }
    },
    onError: (err) => {
      handleLogoutAction();
    },
  });

  useEffect(() => {
    if (doesHttpOnlyCookieExist('ukama_session')) {
      getToken();
    }
  }, []);

  const handleGoToLogin = () => {
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_AUTH_APP_URL || '');
  };

  const handleLogoutAction = () => {
    setIsValidSession(false);
    typeof window !== 'undefined' &&
      window.location.replace(`${process.env.NEXT_PUBLIC_AUTH_APP_URL}/logout`);
  };

  if (getTokenLoading)
    return (
      <Box sx={{ display: 'flex', alignItems: 'center', margin: 12 }}>
        <LoadingWrapper
          radius="small"
          isLoading={true}
          cstyle={{
            overflow: 'hidden',
            height: 'calc(100vh - 200px)',
            backgroundColor: colors.silver,
          }}
        >
          <p></p>
        </LoadingWrapper>
      </Box>
    );

  return <div>{children}</div>;
};

export default AuthWrapper;
