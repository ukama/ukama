import { NetworkDto } from '@/generated';
import { HorizontalContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { Divider, Stack, Typography, useMediaQuery } from '@mui/material';
import Box from '@mui/material/Box';
import { useTheme } from '@mui/material/styles';
import { useRouter } from 'next/router';
import * as React from 'react';
import { useEffect } from 'react';
import BackButton from '../molecules/BackButton';
import LoadingWrapper from '../molecules/LoadingWrapper';
import Header from './Header';
import Sidebar from './Sidebar';

interface ILayoutProps {
  page: string;
  networkId: string;
  isLoading: boolean;
  placeholder: string;
  isDarkMode: boolean;
  isFullScreen: boolean;
  handlePageChange: Function;
  networks: NetworkDto[] | undefined;
  children: React.ReactNode;
  handleNetworkChange: Function;
}

const Layout = ({
  page,
  children,
  networks,
  isLoading,
  networkId,
  isDarkMode,
  placeholder,
  isFullScreen,
  handlePageChange,
  handleNetworkChange,
}: ILayoutProps) => {
  const theme = useTheme();
  const router = useRouter();
  const [open, setOpen] = React.useState(true);
  const matches = useMediaQuery(theme.breakpoints.down('md'));

  useEffect(() => {
    if (matches) {
      setOpen(false);
    } else {
      setOpen(true);
    }
  }, [matches]);

  const onNavigate = (name: string, path: string) => {
    console.log(name, path);
    handlePageChange(name);
    router.push(path);
  };

  return (
    <Box sx={{ display: 'flex', overflow: 'hidden' }}>
      {!isFullScreen && (
        <Header
          isOpen={open}
          isLoading={isLoading}
          onNavigate={onNavigate}
          isDarkMode={isDarkMode}
        />
      )}
      <HorizontalContainer>
        {!isFullScreen && (
          <Sidebar
            page={page}
            isOpen={open}
            networkId={networkId}
            isLoading={isLoading}
            onNavigate={onNavigate}
            isDarkMode={isDarkMode}
            placeholder={placeholder}
            networks={networks || []}
            handleNetworkChange={handleNetworkChange}
          />
        )}

        <Box
          sx={{
            width: '100%',
            height: '100vh',
            overflow: 'auto',
            background: (theme) =>
              theme.palette.mode === 'light'
                ? colors.black10
                : colors.nightGrey,
          }}
        >
          <Box
            sx={{
              ...(isFullScreen
                ? {
                    m: {
                      xs: `38px 18px !important`,
                      md: `52px 84px !important`,
                    },
                  }
                : {
                    p: {
                      xs: '8px 18px 0px 18px !important',
                      md: '16px 32px 0px 32px !important',
                    },
                    m: {
                      xs: `44px 0px 44px 62px !important`,
                      md: `60px 0px 60px 218px !important`,
                    },
                  }),
              backgroundColor: (theme) =>
                theme.palette.mode === 'light'
                  ? colors.black10
                  : colors.nightGrey,
            }}
          >
            <LoadingWrapper
              radius="small"
              width={'100%'}
              isLoading={isLoading}
              height={isLoading ? '100vh' : '100%'}
              cstyle={{ background: isLoading ? colors.white : 'inherit' }}
            >
              <Stack direction={'column'}>
                {page !== '404' && (
                  <Box>
                    <Stack
                      direction={'row'}
                      alignItems={'center'}
                      spacing={{ xs: 4, md: 10.5 }}
                    >
                      {isFullScreen && <BackButton title="BACK TO CONSOLE" />}
                      <Typography variant="h5">{page}</Typography>
                    </Stack>
                    <Divider sx={{ my: 1 }} />
                  </Box>
                )}

                {children}
              </Stack>
            </LoadingWrapper>
          </Box>
        </Box>
      </HorizontalContainer>
    </Box>
  );
};

export default Layout;
