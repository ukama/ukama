/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import AddSubscriberStepperDialog from '@/app/(main)/console/subscribers/_components/AddSubscriber';
import SubscriberDetails from '@/app/(main)/console/subscribers/_components/SubscriberDetails';
import { useSubscribersPage } from '@/app/(main)/console/subscribers/_hooks/useSubscribersPage';
import TopUpData from '@/app/(main)/manage/billing/_components/TopUpData';
import PlanCard from '@/app/(main)/manage/data-plans/_components/PlanCard';
import DataTableWithOptions from '@/components/ui/DataTableWithOptions';
import DeleteConfirmation from '@/components/ui/DeleteDialog';
import LoadingWrapper from '@/components/ui/LoadingWrapper';
import PageContainerHeader from '@/components/ui/PageContainerHeader';
import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import {
  CardWrapper,
  DataPlanEmptyView,
  NavigationButton,
  NavigationWrapper,
  ScrollableContent,
  ScrollContainer,
} from '@/styles/global';
import colors from '@/theme/colors';
import KeyboardArrowLeftIcon from '@mui/icons-material/KeyboardArrowLeft';
import KeyboardArrowRightIcon from '@mui/icons-material/KeyboardArrowRight';
import SubscriberIcon from '@mui/icons-material/PeopleAlt';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { Box, Paper, Stack, Typography } from '@mui/material';

const Page = () => {
  const vm = useSubscribersPage();

  return (
    <Stack
      mt={2}
      spacing={2}
      direction={'column'}
      sx={{ height: { xs: 'calc(100vh - 158px)', md: 'calc(100vh - 172px)' } }}
    >
      {vm.subscribersLoading ? (
        <LoadingWrapper
          radius="small"
          width={'100%'}
          isLoading={true}
          cstyle={{ height: '240px' }}
        >
          <br />
        </LoadingWrapper>
      ) : (
        <Paper elevation={1} sx={{ borderRadius: '10px', p: { xs: 2, md: 4 } }}>
          <Stack direction="column" spacing={{ xs: 0.5, md: 1.5 }}>
            <Box sx={{ display: 'flex', alignItems: 'center' }}>
              <Stack direction={'row'} spacing={1} alignItems={'center'}>
                <Typography variant="h6">Data plans</Typography>
                <Typography variant="subtitle2" sx={{ color: colors.black38 }}>
                  ({vm.packagesData?.getPackages.packages?.length ?? 0})
                </Typography>
              </Stack>

              {vm.packagesData && (vm.packagesData.getPackages.packages?.length ?? 0) > 4 && (
                <NavigationWrapper>
                  <NavigationButton onClick={() => vm.scroll('left')} disabled={false}>
                    <KeyboardArrowLeftIcon fontSize="small" />
                  </NavigationButton>
                  <NavigationButton onClick={() => vm.scroll('right')} disabled={false}>
                    <KeyboardArrowRightIcon fontSize="small" />
                  </NavigationButton>
                </NavigationWrapper>
              )}
            </Box>

            <ScrollContainer>
              <ScrollableContent ref={vm.scrollContainerRef}>
                {!(vm.packagesData?.getPackages.packages?.length) ? (
                  <DataPlanEmptyView>
                    <UpdateIcon sx={{ fontSize: 40, mb: 1 }} />
                    <Typography variant="body1">No data plan created yet!</Typography>
                  </DataPlanEmptyView>
                ) : (
                  vm.packagesData.getPackages.packages.map((pkg) => (
                    <CardWrapper key={pkg.uuid}>
                      <PlanCard
                        uuid={pkg.uuid}
                        name={pkg.name}
                        amount={String(pkg.amount)}
                        isOptions={false}
                        dataUnit={pkg.dataUnit}
                        duration={pkg.duration}
                        currency={vm.currencyData?.getCurrencySymbol.symbol ?? ''}
                        dataVolume={String(pkg.dataVolume)}
                      />
                    </CardWrapper>
                  ))
                )}
              </ScrollableContent>
            </ScrollContainer>
          </Stack>
        </Paper>
      )}

      {vm.subscribersLoading || vm.dataUsageLoading ? (
        <LoadingWrapper
          radius="small"
          width={'100%'}
          isLoading={true}
          cstyle={{ height: '100%' }}
        >
          <br />
        </LoadingWrapper>
      ) : (
        <Paper
          sx={{
            height: '100%',
            overflow: 'hidden',
            borderRadius: '10px',
            px: { xs: 2, md: 3 },
            py: { xs: 2, md: 4 },
          }}
        >
          <PageContainerHeader
            search={vm.search}
            title={'My subscribers'}
            buttonTitle={'Add Subscriber'}
            handleButtonAction={vm.handleAddSubscriberModal}
            onSearchChange={(e: string) => vm.setSearch(e)}
            subtitle={`${vm.subscriberCount}`}
          />
          <br />
          <DataTableWithOptions
            icon={SubscriberIcon}
            isRowClickable={false}
            columns={SUBSCRIBER_TABLE_COLUMNS}
            dataset={vm.tableRows}
            menuOptions={SUBSCRIBER_TABLE_MENU}
            onMenuItemClick={vm.handleTableMenuAction}
            emptyViewLabel={'No subscribers yet!'}
          />
        </Paper>
      )}

      <AddSubscriberStepperDialog
        isOpen={vm.openAddSubscriber}
        currencySymbol={vm.currencyData?.getCurrencySymbol.symbol ?? ''}
        handleCloseAction={vm.handleCloseAddSubscriber}
        handleAddSubscriber={vm.handleAddSubscriber}
        sims={vm.simPoolData?.getSimsFromPool.sims ?? []}
        packages={vm.packagesData?.getPackages.packages ?? []}
        isLoading={vm.addSubscriberLoading || vm.allocateSimLoading}
      />

      <DeleteConfirmation
        open={vm.isConfirmationOpen}
        onDelete={vm.handleDeleteSubscriber}
        onCancel={vm.handleCancel}
        itemName={vm.deletedSubscriber}
        loading={vm.deleteSubscriberLoading}
      />

      <SubscriberDetails
        ishowSubscriberDetails={vm.isSubscriberDetailsOpen}
        handleClose={vm.handleCloseSubscriberDetails}
        subscriberInfo={vm.subscriberDetails}
        handleSimActionOption={vm.handleSimAction}
        handleUpdateSubscriber={vm.handleUpdateSubscriber}
        loading={vm.updateSubscriberLoading}
        handleDeleteSubscriber={vm.handleSubscriberMenuAction}
        simStatusLoading={vm.toggleSimStatusLoading}
      />

      <TopUpData
        isToPup={vm.isTopupData}
        onCancel={vm.handleCloseTopUp}
        handleTopUp={vm.handleTopUp}
        loadingTopUp={vm.packagesLoading || vm.addPackagesToSimLoading}
        packages={vm.packagesData?.getPackages.packages ?? []}
        sims={(vm.subscriberSimList ?? []) as unknown as import('@/client/graphql/generated').SimDto[]}
        subscriberName={vm.topUpSubscriberName}
      />
    </Stack>
  );
};

export default Page;
