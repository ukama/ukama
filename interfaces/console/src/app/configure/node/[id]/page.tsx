/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { useGetNetworksQuery } from '@/client/graphql/generated';
import SiteMapComponent from '@/components/SiteMapComponent';
import { LField } from '@/components/Welcome';
import { INSTALLATION_FLOW, ONBOARDING_FLOW } from '@/constants';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import { setQueryParam } from '@/utils';
import { useFetchAddress } from '@/utils/useFetchAddress';
import { Button, Skeleton, Stack, Typography } from '@mui/material';
import dynamic from 'next/dynamic';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import LoadingSkelton from './skelton';

const BasicDropdown = dynamic(() => import('@/components/BasicDropdown'), {
  ssr: false,
  loading: () => <Skeleton variant="rectangular" width={'100%'} height={29} />,
});

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
  const networkId = searchParams.get('networkid') ?? '';
  const flow = searchParams.get('flow') ?? ONBOARDING_FLOW;
  const [networkSelected, setNetworkSelected] = useState<string>(networkId);
  const [isLoading, setIsLoading] = useState<boolean>(true);
  const { setSnackbarMessage, setNetwork } = useAppContext();
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
      setQueryParam('address', address, searchParams.toString(), pathname);
      setIsLoading(false);
    }
  }, [address]);

  const { data: networksData } = useGetNetworksQuery({
    fetchPolicy: 'cache-first',
  });

  const getDropDownData = () =>
    networksData?.getNetworks.networks.map((network) => ({
      id: network.id,
      label: network.name,
      value: network.id,
    }));

  const handleFetchAddress = async () => {
    await fetchAddress(latlng[0], latlng[1]);
  };

  const handleBack = () => {
    router.back();
  };

  const handleNext = () => {
    const n = searchParams.get('networkid') ?? '';
    if (!n || !networkSelected) {
      setSnackbarMessage({
        id: 'network-msg',
        message: 'Please select network',
        type: 'error',
        show: true,
      });
      return;
    }

    if (address && networkSelected === n) {
      setNetwork({
        id: networkSelected,
        name:
          networksData?.getNetworks.networks.find(
            (n) => n.id === networkSelected,
          )?.name ?? '',
      });
      router.push(`/configure/node/${id}/site/name?${searchParams.toString()}`);
    }
  };

  const handleNetworkChange = (id: string) => {
    if (id) {
      const filterNetwork = networksData?.getNetworks.networks.find(
        (n) => n.id === id,
      );
      setNetworkSelected(id);
      setQueryParam(
        'networkid',
        filterNetwork?.id ?? '',
        searchParams.toString(),
        pathname,
      );
    }
  };

  if (isLoading || addressLoading) <LoadingSkelton />;

  return (
    <Stack spacing={2} overflow={'scroll'}>
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

        <Stack
          width={'fit-content'}
          direction={'column'}
          display={flow === INSTALLATION_FLOW || !networkId ? 'flex' : 'none'}
        >
          <LField label="Network" value={''} />
          <BasicDropdown
            id={'network-dropdown'}
            value={networkSelected}
            isShowAddOption={false}
            placeholder={'Select Network'}
            list={getDropDownData() || []}
            handleOnChange={handleNetworkChange}
            handleAddNetwork={() => {}}
          />
        </Stack>

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
        <Button id="next-button" variant="contained" onClick={handleNext}>
          Next
        </Button>
      </Stack>
    </Stack>
  );
};

export default NodeConfigure;
