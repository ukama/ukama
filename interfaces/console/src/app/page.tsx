'use client';
import AppSnackbar from '@/components/AppSnackbar/page';
import LayoutSkelton from '@/components/Layout/skelton';
import { Box } from '@mui/material';

const Page = () => {
  return (
    <Box>
      <LayoutSkelton />
      <AppSnackbar />
    </Box>
  );
};

export default Page;
