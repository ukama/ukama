import { commonData, snackbarMessage } from '@/app-recoil';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import {
  SubscribersResDto,
  useGetSubscribersByNetworkQuery,
} from '@/generated';
import {
  ContainerMax,
  HorizontalContainerJustify,
  PageContainer,
  VerticalContainer,
} from '@/styles/global';
import { colors } from '@/styles/theme';
import { TCommonData, TSnackMessage } from '@/types';
import { DataTableWithOptions } from '@/ui/components';
import { Search } from '@mui/icons-material';
import { AlertColor, Button, Grid, TextField, Typography } from '@mui/material';
import { useEffect, useState } from 'react';
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
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      if (data.getSubscribersByNetwork.subscribers.length > 0) {
        setSubscriber((prev) => ({
          subscribers: [
            ...prev.subscribers,
            ...data.getSubscribersByNetwork.subscribers,
          ],
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
      setSubscriber({ subscribers: subscribers ?? [] });
    } else if (search.length === 0) {
      setSubscriber({
        subscribers: data?.getSubscribersByNetwork.subscribers ?? [],
      });
    }
  }, [search]);

  const onTableMenuItem = (id: string, type: string) => {
    console.log(id, type);
  };

  return (
    <PageContainer>
      <HorizontalContainerJustify>
        <Grid container justifyContent={'space-between'} spacing={1}>
          <Grid container item xs={12} md="auto" alignItems={'center'}>
            <Grid item xs={'auto'}>
              <Typography variant="h6" mr={1}>
                My subscribers
              </Typography>
            </Grid>
            <Grid item xs={'auto'}>
              <Typography
                variant="subtitle2"
                mr={1.4}
              >{`(${subscriber.subscribers.length})`}</Typography>
            </Grid>
            <Grid item xs={12} md={'auto'}>
              <TextField
                id="subscriber-search"
                label="Search"
                variant="outlined"
                size="small"
                value={search}
                onChange={(e) => setSearch(e.target.value)}
                sx={{ width: { xs: '100%', lg: '250px' } }}
                InputProps={{
                  endAdornment: <Search htmlColor={colors.black54} />,
                }}
              />
            </Grid>
          </Grid>
          <Grid item xs={12} md={'auto'}>
            <Button
              variant="contained"
              color="primary"
              size="medium"
              sx={{ width: { xs: '100%', md: '250px' } }}
            >
              Add Subscriber
            </Button>
          </Grid>
        </Grid>
      </HorizontalContainerJustify>
      <VerticalContainer>
        <ContainerMax mt={4.5}>
          <DataTableWithOptions
            columns={SUBSCRIBER_TABLE_COLUMNS}
            dataset={[
              {
                name: 'John Doe',
                network: 'Globe',
                dataUsage: '1.2 GB',
                dataPlan: '1.5 GB',
                actions: 'actions',
              },
              {
                name: 'John Do',
                network: 'Earth',
                dataUsage: '1.1 GB',
                dataPlan: '1.9 GB',
                actions: 'actions',
              },
            ]}
            menuOptions={SUBSCRIBER_TABLE_MENU}
            onMenuItemClick={onTableMenuItem}
            emptyViewLabel={'No subscribers yet!'}
          />
        </ContainerMax>
      </VerticalContainer>
    </PageContainer>
  );
};

export default Page;
