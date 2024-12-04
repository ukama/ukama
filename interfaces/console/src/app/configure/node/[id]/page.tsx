/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import SiteMapComponent from '@/components/SiteMapComponent';
import { LField } from '@/components/Welcome';
import colors from '@/theme/colors';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { Button, Paper, Skeleton, Stack, Typography } from '@mui/material';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';

interface INodeConfigure {
  params: {
    id: string;
  };
}

const NodeConfigure: React.FC<INodeConfigure> = ({ params }) => {
  const { id } = params;
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const qpLat = searchParams.get('lat') ?? '';
  const qpLng = searchParams.get('lng') ?? '';
  const flow = searchParams.get('flow') ?? 'onb';
  const step = parseInt(searchParams.get('step') ?? '1');
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

  const setQueryParam = (key: string, value: string) => {
    const p = new URLSearchParams(searchParams.toString());
    p.set(key, value);
    window.history.replaceState({}, '', `${pathname}?${p.toString()}`);
    return p;
  };

  const handleBack = () => {
    setQueryParam('step', (step - 1).toString());
    router.back();
  };

  const handleNext = () => {
    if (address) {
      const p = setQueryParam('step', (step + 1).toString());
      router.push(`/configure/node/${id}/site/name?${p.toString()}`);
    }
  };

  return (
    <Paper elevation={0} sx={{ px: { xs: 2, md: 4 }, py: { xs: 1, md: 2 } }}>
      <Stack direction={'row'}>
        <Typography variant="h6">{'Install site'}</Typography>
        <Typography
          variant="h6"
          fontWeight={400}
          sx={{
            color: colors.black70,
          }}
        >
          {flow === 'onb' && <i>&nbsp;- optional</i>}&nbsp;({step}/
          {flow === 'onb' ? 6 : 4})
        </Typography>
      </Stack>

      <Stack mt={3} mb={3} direction={'column'} spacing={2}>
        <Typography variant={'body1'} fontWeight={400}>
          You have successfully installed your site, and it is online now!
          Please check to make sure you are satisfied with the location details.
        </Typography>

        {isLoading || addressLoading ? (
          <Skeleton variant="rounded" width={'100%'} height={128} />
        ) : (
          <SiteMapComponent
            posix={[latlng[0], latlng[1]]}
            address={address}
            height={'128px'}
          />
        )}

        <LField label="Node Id" value={id} />
        <LField
          label="SITE LOCATION"
          value={`${address} [${qpLat}, ${qpLng}]`}
        />
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
