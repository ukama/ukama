'use client';

import AppSnackbar from '@/components/AppSnackbar/page';
import BackButton from '@/components/BackButton';
import { useAppContext } from '@/context';
import { MANAGE_MENU_LIST } from '@/routes';
import '@/styles/console.css';
import colors from '@/theme/colors';
import { Container, Divider, Paper, Stack, Typography } from '@mui/material';
import Link from 'next/link';
import { usePathname } from 'next/navigation';

const ManageMenu = () => {
  const pathname = usePathname();
  const { isDarkMode } = useAppContext();
  return (
    <Paper
      sx={{
        width: '436px',
        height: 'fit-content',
        display: 'flex',
        borderRadius: '10px',
        maxWidth: 'fit-content',
      }}
    >
      <Stack px={2} py={3} spacing={1.5} direction="column">
        {MANAGE_MENU_LIST.map(({ id, icon: Icon, name, path }) => (
          <Link
            key={id}
            href={path}
            prefetch={id === 'manage-members' ? true : false}
            style={{
              borderRadius: 4,
              textDecoration: 'none',
              backgroundColor:
                pathname === path ? colors.white38 : 'transparent',
            }}
          >
            <Stack
              pl={2}
              pr={4}
              py={1}
              spacing={2}
              direction={'row'}
              alignItems={'flex-start'}
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
        ))}
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
    <Container maxWidth={'xl'} sx={{ my: 8 }}>
      <Stack direction={'column'} spacing={4}>
        <BackButton title="BACK TO CONSOLE" />
        <Divider />
        <Stack direction={'row'} spacing={4}>
          <ManageMenu />
          {children}
          <AppSnackbar />
        </Stack>
      </Stack>
    </Container>
  );
};

export default ManageLayout;
