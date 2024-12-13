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
import { ONBOARDING_FLOW } from '@/constants';
import colors from '@/theme/colors';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { Button, Stack, Typography } from '@mui/material';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import LoadingSkelton from './skelton';

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
  const flow = searchParams.get('flow') ?? ONBOARDING_FLOW;
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
    if (latlng[0] === 0 && latlng[1] === 0)
      router.push(`/configure/check?flow=${flow}`);
  }, []);

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
    router.back();
  };

  const handleNext = () => {
    if (address) {
      router.push(`/configure/node/${id}/site/name?${searchParams.toString()}`);
    }
  };

  if (isLoading || addressLoading) <LoadingSkelton />;

  return (
    <Stack spacing={2}>
      <Stack direction={'row'}>
        <Typography variant="h4" fontWeight={500}>
          Site installed
        </Typography>
      </Stack>

      <Stack mt={3} mb={3} direction={'column'} spacing={3}>
        <Typography variant={'body1'} fontWeight={400}>
          You have successfully installed your site, and it is online now!
          Please check to make sure you are satisfied with the location details.
          <br />
          <br />
        </Typography>

        <SiteMapComponent
          posix={[latlng[0], latlng[1]]}
          address={address}
          height={'128px'}
        />

        <LField label="Node Id" value={id} />
        <LField
          label="SITE LOCATION"
          value={`${address} [${latlng[0]}, ${latlng[1]}]`}
        />
      </Stack>

      <Stack
        direction={'row'}
        pt={{ xs: 4, md: 6 }}
        justifyContent={'space-between'}
      >
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
    </Stack>
  );
};

export default NodeConfigure;
