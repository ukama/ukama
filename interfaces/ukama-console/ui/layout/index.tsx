'use client';
import { HorizontalContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import { useMediaQuery } from '@mui/material';
import Box from '@mui/material/Box';
import { useTheme } from '@mui/material/styles';
import { useRouter } from 'next/router';
import * as React from 'react';
import { useEffect } from 'react';
import { LoadingWrapper } from '../components';
import Header from './Header';
import Sidebar from './Sidebar';

interface ILayoutProps {
  page: string;
  isLoading: boolean;
  isDarkMode: boolean;
  handlePageChange: Function;
  children: React.ReactNode;
}

const Layout = ({
  page,
  children,
  isLoading,
  isDarkMode,
  handlePageChange,
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
    handlePageChange(name);
    router.push(path);
  };

  return (
    <Box sx={{ display: 'flex', overflow: 'hidden' }}>
      <Header
        isOpen={open}
        isLoading={isLoading}
        onNavigate={onNavigate}
        isDarkMode={isDarkMode}
      />
      <HorizontalContainer>
        <Sidebar
          page={page}
          isOpen={open}
          isLoading={isLoading}
          onNavigate={onNavigate}
          isDarkMode={isDarkMode}
        />

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
              p: {
                xs: '18px 18px 0px 18px !important',
                md: '32px 32px 0px 32px !important',
              },
              m: {
                xs: `44px 0px 44px 62px !important`,
                md: `60px 0px 60px 218px !important`,
              },
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
              {children}
            </LoadingWrapper>
          </Box>
        </Box>
      </HorizontalContainer>
    </Box>
  );
};

export default Layout;
