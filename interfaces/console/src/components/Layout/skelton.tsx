import { Skeleton, Stack } from '@mui/material';

const LayoutSkelton = () => {
  return (
    <Stack direction={'column'} height={'100vh'} overflow={'hidden'}>
      <Skeleton
        height={69}
        width={'100%'}
        animation="pulse"
        variant="rectangular"
      />
      <Stack direction="row" height="100%" spacing={4}>
        <Skeleton
          width={232}
          height={'auto'}
          animation="pulse"
          variant="rectangular"
        />

        <Stack
          py={3}
          pr={3}
          spacing={3}
          width="100%"
          height="100%"
          direction="column"
        >
          <Skeleton width={200} height={40} animation="pulse" variant="text" />
          <Skeleton
            width={'100%'}
            height={'calc(100% - 10%)'}
            animation="pulse"
            variant="rectangular"
            sx={{
              borderRadius: '10px',
            }}
          />
        </Stack>
      </Stack>
    </Stack>
  );
};

export default LayoutSkelton;
