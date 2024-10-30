'use client';

import { CenterContainer } from '@/styles/global';
import { CircularProgress } from '@mui/material';
import { useEffect, useRef } from 'react';

export default function LogoutAction({ deleteTokens }: { deleteTokens: any }) {
  const deleteTokensRef = useRef(deleteTokens);

  useEffect(() => {
    deleteTokensRef.current = deleteTokens;
  });

  useEffect(() => {
    deleteTokensRef.current();
    setTimeout(() => {
      window.location.href = `${process.env.NEXT_PUBLIC_AUTH_APP_URL}/user/logout`;
    }, 5000);
  }, []);

  return (
    <CenterContainer>
      <CircularProgress />
    </CenterContainer>
  );
}
