/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { commonData, snackbarMessage } from '@/app-recoil';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import {
  SubscribersResDto,
  useGetPackagesQuery,
  useGetSubscribersByNetworkQuery,
} from '@/generated';
import {
  ContainerMax,
  PageContainer,
  VerticalContainer,
} from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSnackMessage } from '@/types';
import DataTableWithOptions from '@/ui/molecules/DataTableWithOptions';
import EmptyView from '@/ui/molecules/EmptyView';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import PlanCard from '@/ui/molecules/PlanCard';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import { AlertColor, Grid, Stack, Typography } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const Page = () => {
  const [search, setSearch] = useState<string>('');
  const _commonData = useRecoilValue<TCommonData>(commonData);
  const setSnackbarMessage = useSetRecoilState<TSnackMessage>(snackbarMessage);
  const [subscriber, setSubscriber] = useState<SubscribersResDto>({
    subscribers: [],
  });

  const { loading, data } = useGetSubscribersByNetworkQuery({
    variables: { networkId: _commonData.networkId },
    fetchPolicy: 'cache-first',
    onCompleted: (data) => {
      if (data.getSubscribersByNetwork.subscribers.length > 0) {
        setSubscriber(() => ({
          subscribers: [...data.getSubscribersByNetwork.subscribers],
        }));
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'subscriber-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const { data: dataPlanData, loading: dataPlanLoading } = useGetPackagesQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'data-plan-err-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  useEffect(() => {
    if (search.length > 3) {
      const subscribers = data?.getSubscribersByNetwork.subscribers.filter(
        (subscriber) => {
          const s = search.toLowerCase();
          if (
            subscriber.firstName.toLowerCase().includes(s) ||
            subscriber.lastName.toLowerCase().includes(s)
          )
            return subscriber;
        },
      );
      setSubscriber(() => ({ subscribers: subscribers ?? [] }));
    } else if (search.length === 0) {
      setSubscriber(() => ({
        subscribers: data?.getSubscribersByNetwork.subscribers ?? [],
      }));
    }
  }, [search]);

  const onTableMenuItem = (id: string, type: string) => {};

  const structureData = useCallback(
    (data: SubscribersResDto) =>
      data.subscribers.map((subscriber) => ({
        id: subscriber.uuid,
        email: subscriber.email,
        name: `${subscriber.firstName} ${subscriber.lastName}`,
        dataUsage: '',
        dataPlan: '',
        actions: '',
      })),
    [],
  );

  return (
    <Stack direction={'column'}>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={loading}
        cstyle={{
          backgroundColor: loading ? colors.white : 'transparent',
        }}
      >
        <PageContainer
          sx={{ height: 'fit-content', maxHeight: 'calc(100vh - 400px)' }}
        >
          <PageContainerHeader
            search={search}
            title={'My subscribers'}
            buttonTitle={'Add Subscriber'}
            handleButtonAction={() => {}}
            subtitle={`${subscriber.subscribers.length}`}
            onSearchChange={(e: string) => setSearch(e)}
          />
          <VerticalContainer>
            <ContainerMax mt={4.5}>
              <DataTableWithOptions
                icon={SubscriberIcon}
                isRowClickable={false}
                columns={SUBSCRIBER_TABLE_COLUMNS}
                dataset={structureData(subscriber)}
                menuOptions={SUBSCRIBER_TABLE_MENU}
                onMenuItemClick={onTableMenuItem}
                emptyViewLabel={'No subscribers yet!'}
              />
            </ContainerMax>
          </VerticalContainer>
        </PageContainer>
      </LoadingWrapper>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={dataPlanLoading}
        cstyle={{
          backgroundColor: dataPlanLoading ? colors.white : 'transparent',
        }}
      >
        <PageContainer
          sx={{ height: 'fit-content', maxHeight: 'calc(100vh - 550px)' }}
        >
          <Stack direction={'row'} alignItems={'center'}>
            <Typography variant="h6" mr={1}>
              Data plans
            </Typography>
            <Typography variant="subtitle2">
              <i>(view only)</i>
            </Typography>
          </Stack>
          <Stack my={4}>
            {dataPlanData?.getPackages &&
            dataPlanData?.getPackages?.packages?.length > 0 ? (
              <Grid container rowSpacing={2} columnSpacing={2}>
                {dataPlanData?.getPackages?.packages.map(
                  ({
                    uuid,
                    name,
                    duration,
                    users,
                    currency,
                    dataVolume,
                    dataUnit,
                    amount,
                  }: any) => (
                    <Grid item xs={12} sm={6} md={4} key={uuid}>
                      <PlanCard
                        uuid={uuid}
                        name={name}
                        users={users}
                        amount={amount}
                        dataUnit={dataUnit}
                        duration={duration}
                        currency={currency}
                        dataVolume={dataVolume}
                        isOptions={false}
                      />
                    </Grid>
                  ),
                )}
              </Grid>
            ) : (
              <EmptyView
                size="medium"
                title={
                  'No data plans yet! Go to “Manage data plans” in your organization settings to add one'
                }
                icon={SubscriberIcon}
              />
            )}
          </Stack>
        </PageContainer>
      </LoadingWrapper>
    </Stack>
  );
};

export default Page;
