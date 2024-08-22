'use client';
import SiteMapComponent from '@/components/SiteMapComponent';
import colors from '@/theme/colors';
import { useFetchAddress } from '@/utils/useFetchAddress';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import { Button, Paper, Skeleton, Stack, Typography } from '@mui/material';
import { useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';

interface INodeConfigure {
  params: {
    id: string;
  };
}

const NodeConfigure: React.FC<INodeConfigure> = ({ params }) => {
  const { id } = params;
  const router = useRouter();
  const searchParams = useSearchParams();
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const stepTracker = searchParams.get('step') ?? '1';
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const [latlng] = useState<[number, number]>([
    parseFloat(qpLat),
    parseFloat(qpLng),
  ]);
  const {
    address,
    isLoading: addressLoading,
    fetchAddress,
  } = useFetchAddress();

  useEffect(() => {
    if (latlng[0] !== 0 && latlng[1] !== 0) handleFetchAddress();
  }, [latlng]);

  useEffect(() => {
    if (address) {
      setQueryParam('address', address);
      setIsLoading(false);
    }
  }, [address]);

  const handleFetchAddress = async () => {
    await fetchAddress(latlng[0], latlng[1]);
  };

  const handleBack = () => router.back();

  const setQueryParam = (key: string, value: string) => {
    const params = new URLSearchParams(searchParams.toString());
    params.set(key, value);
    window.history.pushState(null, '', `?${params.toString()}`);
  };

  const handleNext = () => {
    if (address) {
      router.push(`/configure/node/${id}/site?${searchParams.toString()}`);
    }
  };

  return (
    <Paper elevation={0} sx={{ px: 4, py: 2 }}>
      <Stack direction={'row'}>
        <Typography variant="h6">{'Install site'}</Typography>
        <Typography
          variant="h6"
          fontWeight={400}
          sx={{
            color: colors.black70,
            display: stepTracker !== '1' ? 'none' : 'flex',
          }}
        >
          <i>&nbsp;- optional</i>&nbsp;(3/6)
        </Typography>
      </Stack>

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
          <Skeleton variant="rounded" width={'100%'} height={128} />
        ) : (
          <SiteMapComponent
            posix={[latlng[0], latlng[1]]}
            address={address}
            height={'128px'}
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
  );
};

export default NodeConfigure;
