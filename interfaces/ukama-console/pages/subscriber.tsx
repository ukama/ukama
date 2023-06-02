import { commonData, snackbarMessage } from '@/app-recoil';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import {
  SubscribersResDto,
  useGetSubscribersByNetworkQuery,
} from '@/generated';
import {
  ContainerMax,
  PageContainer,
  VerticalContainer,
} from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSnackMessage } from '@/types';
import { DataTableWithOptions, LoadingWrapper } from '@/ui/components';
import PageContainerHeader from '@/ui/components/PageContainerHeader';
import { AlertColor } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { useRecoilValue, useSetRecoilState } from 'recoil';

const Page = async () => {
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

  const onTableMenuItem = (id: string, type: string) => {
    console.log(id, type);
  };

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
    <LoadingWrapper
      radius="small"
      width={'100%'}
      isLoading={loading}
      cstyle={{
        backgroundColor: loading ? colors.white : 'transparent',
      }}
    >
      <PageContainer>
        <PageContainerHeader
          title={'My subscribers'}
          subtitle={`${subscriber.subscribers.length}`}
          buttonTitle={'Add Subscriber'}
          handleButtonAction={() => console.log('Add subscriber')}
          onSearchChange={(e: string) => setSearch(e)}
          search={search}
        />

        <VerticalContainer>
          <ContainerMax mt={4.5}>
            <DataTableWithOptions
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
  );
};

export default Page;
