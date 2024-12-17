/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  PackageDto,
  useAddPackageMutation,
  useDeletePackageMutation,
  useGetCurrencySymbolLazyQuery,
  useGetNetworkQuery,
  useGetPackagesQuery,
  useUpdatePacakgeMutation,
} from '@/client/graphql/generated';
import DataPlanDialog from '@/components/DataPlanDialog';
import EmptyView from '@/components/EmptyView';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import PlanCard from '@/components/PlanCard';
import { useAppContext } from '@/context';
import { colors } from '@/theme';
import { CreatePlanType, DataUnitType } from '@/types';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { AlertColor, Box, Grid, Paper } from '@mui/material';
import { useState } from 'react';

const INIT_DATAPLAN = {
  id: '',
  name: '',
  dataVolume: undefined,
  dataUnit: DataUnitType.GigaBytes,
  amount: undefined,
  duration: 0,
  currency: '',
  country: '',
};

const Page = () => {
  const [data, setData] = useState<any>([]);
  const { network, user, setSnackbarMessage } = useAppContext();
  const [dataplan, setDataplan] = useState<CreatePlanType>(INIT_DATAPLAN);
  const [isDataPlan, setIsDataPlan] = useState<boolean>(false);

  const [getCurrencySymbol, { data: currencyData }] =
    useGetCurrencySymbolLazyQuery({
      fetchPolicy: 'cache-and-network',
      onError: (error) => {
        setSnackbarMessage({
          id: 'currency-info-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const { data: networkData } = useGetNetworkQuery({
    variables: {
      networkId: network.id,
    },
    onCompleted: (data) => {
      getCurrencySymbol({
        variables: {
          code: user.currency,
        },
      });
    },
  });

  const {
    data: packagesData,
    loading: packagesLoading,
    refetch: getDataPlans,
  } = useGetPackagesQuery({
    fetchPolicy: 'cache-and-network',
    onCompleted: (data) => {
      setData(data?.getPackages.packages ?? []);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'packages',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [addDataPlan, { loading: dataPlanLoading }] = useAddPackageMutation({
    onCompleted: () => {
      getDataPlans();
      setSnackbarMessage({
        id: 'add-data-plan',
        message: 'Data plan added successfully',
        type: 'success' as AlertColor,
        show: true,
      });
      setIsDataPlan(false);
      setDataplan(INIT_DATAPLAN);
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'data-plan-error',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [deletePackage, { loading: deletePkgLoading }] =
    useDeletePackageMutation({
      onCompleted: () => {
        getDataPlans().then((res) => {
          setData(res?.data?.getPackages.packages ?? []);
        });
        setSnackbarMessage({
          id: 'delete-data-plan',
          message: 'Data plan deleted successfully',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'data-plan-delete-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [updatePackage, { loading: updatePkgLoading }] =
    useUpdatePacakgeMutation({
      onCompleted: () => {
        getDataPlans();
        setSnackbarMessage({
          id: 'update-data-plan',
          message: 'Data plan updated successfully',
          type: 'success' as AlertColor,
          show: true,
        });
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'data-plan-update-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const handleAddDataPlanAction = () => {
    if (user.currency) {
      setDataplan(INIT_DATAPLAN);
      setIsDataPlan(true);
    } else {
      setSnackbarMessage({
        id: 'network-not-selected',
        message: 'Something went wrong, please try again',
        type: 'warning' as AlertColor,
        show: true,
      });
    }
  };

  const handleDataPlanAction = (action: string) => {
    if (action === 'add' && dataplan.amount && dataplan.dataVolume) {
      addDataPlan({
        variables: {
          data: {
            name: dataplan.name,
            amount: dataplan.amount,
            dataUnit: dataplan.dataUnit,
            dataVolume: dataplan.dataVolume,
            duration: dataplan.duration,
            country: user.country ?? '',
            currency: user.currency ?? '',
          },
        },
      });
    } else if (action === 'update') {
      updatePackage({
        variables: {
          packageId: dataplan.id,
          data: {
            name: dataplan.name,
            active: true,
          },
        },
      });
      setIsDataPlan(false);
    }
  };

  const handleOptionMenuItemAction = (id: string, action: string) => {
    if (action === 'delete') {
      deletePackage({
        variables: {
          packageId: id,
        },
      });
    } else if (action === 'edit') {
      const d: PackageDto | undefined = packagesData?.getPackages.packages.find(
        (pkg: PackageDto) => pkg.uuid === id,
      );
      setDataplan({
        id: id,
        name: d?.name ?? '',
        duration: d?.duration ?? 0,
        dataVolume: d?.dataVolume ?? 0,
        country: d?.country ?? '',
        currency: d?.currency ?? '',
        amount: typeof d?.rate.amount === 'number' ? d.rate.amount : 0,
        dataUnit: d?.dataUnit as DataUnitType,
      });
      setIsDataPlan(true);
    }
  };
  return (
    <LoadingWrapper
      width={'100%'}
      radius="medium"
      isLoading={
        packagesLoading ??
        dataPlanLoading ??
        updatePkgLoading ??
        deletePkgLoading
      }
      height={'calc(100vh - 244px)'}
    >
      <Paper
        sx={{
          py: { xs: 1.5, md: 3 },
          px: { xs: 2, md: 4 },
          overflow: 'scroll',
          borderRadius: '10px',
          bgcolor: colors.white,
          height: '100%',
        }}
      >
        <Box sx={{ width: '100%', height: '100%' }}>
          <PageContainerHeader
            showSearch={false}
            title={'Data plans'}
            buttonTitle={'CREATE DATA PLAN'}
            handleButtonAction={handleAddDataPlanAction}
          />
          <br />
          {data.length === 0 ? (
            <EmptyView icon={UpdateIcon} title="No data plan created yet!" />
          ) : (
            <Grid container rowSpacing={2} columnSpacing={2}>
              {data.map(
                ({
                  uuid,
                  name,
                  duration,
                  currency,
                  dataVolume,
                  dataUnit,
                  amount,
                }: any) => (
                  <Grid item xs={12} sm={6} md={4} key={uuid}>
                    <PlanCard
                      uuid={uuid}
                      name={name}
                      amount={amount}
                      dataUnit={dataUnit}
                      duration={duration}
                      currency={currency}
                      dataVolume={dataVolume}
                      handleOptionMenuItemAction={handleOptionMenuItemAction}
                    />
                  </Grid>
                ),
              )}
            </Grid>
          )}
        </Box>
        {isDataPlan && (
          <DataPlanDialog
            data={dataplan}
            isOpen={isDataPlan}
            setData={setDataplan}
            currencySymbol={currencyData?.getCurrencySymbol.symbol ?? ''}
            title={'Create data plan'}
            labelNegativeBtn={'Cancel'}
            action={dataplan.id ? 'update' : 'add'}
            handleSuccessAction={handleDataPlanAction}
            handleCloseAction={() => setIsDataPlan(false)}
            labelSuccessBtn={
              dataplan.id ? 'Update Data Plan' : 'Save Data Plan'
            }
          />
        )}
      </Paper>
    </LoadingWrapper>
  );
};

export default Page;
