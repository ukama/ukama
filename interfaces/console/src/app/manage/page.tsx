'use client';
import { Role_Type } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import { Skeleton } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect } from 'react';

const Page = () => {
  const router = useRouter();
  const { token, user } = useAppContext();
  useEffect(() => {
    if (
      user.role === Role_Type.RoleOwner ||
      user.role === Role_Type.RoleAdmin
    ) {
      router.push('/manage/members');
    }
  }, []);
  return (
    <Skeleton
      variant="rectangular"
      sx={{
        width: '100%',
        borderRadius: '10px',
        height: 'calc(100vh - 400px)',
      }}
    />
  );
};

export default Page;
