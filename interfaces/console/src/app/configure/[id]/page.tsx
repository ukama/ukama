'use client';

import SiteMapComponent from '@/components/SiteMapComponent';
import { CenterContainer } from '@/styles/global';
import colors from '@/theme/colors';
import { useFetchAddress } from '@/utils/useFetchAddress';
import GradientWrapper from '@/wrappers/gradiantWrapper';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import { Button, Paper, Skeleton, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';
import { useEffect, useState } from 'react';

interface INodeFound {
  params: {
    id: string;
  };
}

const NodeFound: React.FC<INodeFound> = ({ params }) => {
  const { id } = params;
  const router = useRouter();
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [latlng, setLatlng] = useState<[number, number]>([0, 0]);
  const {
    address,
    isLoading: addressLoading,
    error,
    fetchAddress,
  } = useFetchAddress();

  //TODO: GET NODE BY ID

  useEffect(() => {
    setTimeout(() => {
      setLatlng([-4.322447, 15.307045]);
    }, 2000);
  }, []);

  useEffect(() => {
    if (latlng[0] !== 0 && latlng[1] !== 0) handleFetchAddress();
  }, [latlng]);

  useEffect(() => {
    if (address) setIsLoading(false);
  }, [address]);

  const handleFetchAddress = async () => {
    await fetchAddress(latlng[0], latlng[1]);
  };

  const handleBack = () => {
    router.back();
  };

  const handleNext = () => {};

  return (
    <CenterContainer>
      <GradientWrapper>
        <Paper elevation={0} sx={{ px: 4, py: 2 }}>
          <Typography variant="h6" fontWeight={500}>
            Install site -{' '}
            <span style={{ color: colors.black70, fontWeight: 400 }}>
              <i>optional</i> (2/6)
            </span>
          </Typography>
          <Stack mt={3} mb={3} direction={'column'} spacing={2}>
            {isLoading ? (
              <Stack direction="row" alignItems={'center'} spacing={1}>
                <Skeleton variant="circular" width={24} height={24} />
                <Skeleton variant="text" width={200} height={20} />
              </Stack>
            ) : (
              <Stack direction="row" alignItems={'center'} spacing={1}>
                <CheckCircleIcon sx={{ color: colors.green }} />
                <Typography variant={'body1'} sx={{ fontWeight: 700 }}>
                  Your site is online
                </Typography>
              </Stack>
            )}

            {isLoading || addressLoading ? (
              <Skeleton variant="rounded" width={'100%'} height={88} />
            ) : (
              <SiteMapComponent
                posix={[latlng[0], latlng[1]]}
                address={address}
                height={'88px'}
              />
            )}
          </Stack>
          <Stack mb={1} direction={'row'} justifyContent={'space-between'}>
            <Button
              variant="text"
              onClick={handleBack}
              sx={{ color: colors.black70, p: 0 }}
            >
              Back
            </Button>
            <Button variant="contained" onClick={handleNext}>
              Next
            </Button>
          </Stack>
        </Paper>
      </GradientWrapper>
    </CenterContainer>
  );
};

export default NodeFound;
