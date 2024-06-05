/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

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
import { NotificationRes } from '@/generated/metrics';

interface ILayoutProps {
  page: string;
  isLoading: boolean;
  placeholder: string;
  isDarkMode: boolean;
  isFullScreen: boolean;
  handlePageChange: Function;
  networks: NetworkDto[];
  children: React.ReactNode;
  handleAddNetwork: Function;
  handleNetworkChange: Function;
  alerts: NotificationRes[] | undefined,
  setAlerts: Function;
}

const Layout = ({
  page,
  children,
  networks,
  isLoading,
  isDarkMode,
  placeholder,
  isFullScreen,
  handlePageChange,
  handleAddNetwork,
  handleNetworkChange,
  alerts,
  setAlerts,
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
      {!isFullScreen && (
        <Header
          isOpen={open}
          isLoading={isLoading}
          onNavigate={onNavigate}
          isDarkMode={isDarkMode}
          alerts={alerts}
          setAlerts={setAlerts}
        />
      )}
      <HorizontalContainer>
        {!isFullScreen && (
          <Sidebar
            page={page}
            isOpen={open}
            isLoading={isLoading}
            onNavigate={onNavigate}
            isDarkMode={isDarkMode}
            placeholder={placeholder}
            networks={networks || []}
            handleAddNetwork={handleAddNetwork}
            handleNetworkChange={handleNetworkChange}
          />
        )}

        <Box
          sx={{
            width: '100%',
            height: '100%',
            overflow: 'hidden',
            background: (theme) =>
              theme.palette.mode === 'light'
                ? colors.black40
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
                      md: `60px 0px 0px 218px !important`,
                    },
                  }),
              height: '100%',
              // backgroundColor: (theme) =>
              //   theme.palette.mode === 'light'
              //     ? colors.black10
              //     : colors.nightGrey,
            }}
          >
            <LoadingWrapper
              radius="small"
              width={'100%'}
              isLoading={isLoading}
              height={isLoading ? '90vh' : '100%'}
            >
              <Stack height={'100%'} direction={'column'}>
                {page !== '404' &&
                  page !== 'Unauthorized' &&
                  page != 'OnBoarding' && (
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
