/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import AppSnackbar from '@/components/AppSnackbar/page';
import BackButton from '@/components/BackButton';
import { useAppContext } from '@/context';
import { MANAGE_MENU_LIST, MANAGE_MENU_LIST_SMALL } from '@/routes';
import '@/styles/console.css';
import colors from '@/theme/colors';
import {
  Container,
  Divider,
  Paper,
  Stack,
  Typography,
  useMediaQuery,
  useTheme,
} from '@mui/material';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const ManageMenu = () => {
  const theme = useTheme();
  const matches = useMediaQuery(theme.breakpoints.down('md'));
  const pathname = usePathname();
  const { isDarkMode } = useAppContext();
  return (
    <Paper
      sx={{
        display: 'flex',
        borderRadius: '10px',
        width: 'max-content',
        height: 'fit-content',
      }}
    >
      <Stack
        px={{ xs: 1.4, md: 1 }}
        py={{ xs: 1.5, md: 3 }}
        spacing={{ xs: 0.5, md: 1.5 }}
        direction={{ xs: 'row', md: 'column' }}
      >
        {(matches ? MANAGE_MENU_LIST_SMALL : MANAGE_MENU_LIST).map(
          ({ id, icon: Icon, name, path }) => (
            <Link
              key={id}
              href={path}
              prefetch={id === 'manage-members'}
              style={{
                borderRadius: 4,
                textDecoration: 'none',
                backgroundColor:
                  pathname === path ? colors.white38 : 'transparent',
              }}
            >
              <Stack
                pl={{ xs: 1, md: 2 }}
                pr={{ xs: 2, md: 4 }}
                py={1}
                spacing={{ xs: 1, md: 2 }}
                alignItems={{ xs: 'center', md: 'flex-start' }}
                direction={{ xs: 'column', md: 'row' }}
              >
                <Icon
                  sx={{
                    color: isDarkMode ? colors.white70 : colors.vulcan100,
                  }}
                />
                <Typography
                  variant="body1"
                  fontWeight={400}
                  color={colors.vulcan100}
                >
                  {name}
                </Typography>
              </Stack>
            </Link>
          ),
        )}
      </Stack>
    </Paper>
  );
};

const ManageLayout = ({
  children,
}: Readonly<{
  children: React.ReactNode;
}>) => {
  return (
    <Container maxWidth={'xl'} sx={{ my: { xs: 2, md: 8 } }}>
      <Stack direction={'column'} spacing={{ xs: 2, md: 4 }}>
        <BackButton title="BACK TO CONSOLE" />
        <Divider />
        <Stack
          direction={{ xs: 'column', md: 'row' }}
          spacing={{ xs: 2, md: 4 }}
        >
          <ManageMenu />
          {children}
          <AppSnackbar />
        </Stack>
      </Stack>
    </Container>
  );
};

export default ManageLayout;
