'use client';
import '@/styles/console.css';
import { CenterContainer } from '@/styles/global';
import { Stack, Typography } from '@mui/material';
import Link from 'next/link';

const Page = () => {
  return (
    <CenterContainer>
      <Stack spacing={0.5} alignItems={'center'}>
        <Typography variant="body1">
          {"Sorry, You don't have permissions to view this page"}
        </Typography>
        <Link href={`${process.env.NEXT_PUBLIC_AUTH_APP_URL}/user/logout`}>
          Log me out
        </Link>
      </Stack>
    </CenterContainer>
  );
};

export default Page;
