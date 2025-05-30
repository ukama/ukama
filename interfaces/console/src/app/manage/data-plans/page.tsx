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
import { Box, Grid, Paper } from '@mui/material';
import { SetStateAction, useState } from 'react';

interface DialogState {
  isOpen: boolean;
  mode: 'create' | 'edit';
  data: CreatePlanType;
}

const INIT_DATAPLAN: CreatePlanType = {
  id: '',
  name: '',
  country: '',
  duration: '',
  currency: '',
  amount: undefined,
  dataVolume: undefined,
  dataUnit: DataUnitType.GigaBytes,
};

const useDataPlans = () => {
  const { network, user, setSnackbarMessage } = useAppContext();

  const [getCurrencySymbol, { data: currencyData }] =
    useGetCurrencySymbolLazyQuery({
      fetchPolicy: 'cache-first',
      onError: (error) => {
        setSnackbarMessage({
          id: 'currency-info-error',
          message: error.message,
          type: 'error',
          show: true,
        });
      },
    });

  useGetNetworkQuery({
    fetchPolicy: 'cache-first',
    skip: !network.id,
    variables: { networkId: network.id },
    onCompleted: () => {
      getCurrencySymbol({ variables: { code: user.currency } });
    },
  });

  const {
    data: packagesData,
    loading: packagesLoading,
    refetch: getDataPlans,
  } = useGetPackagesQuery({
    fetchPolicy: 'network-only',
    onError: (error) => {
      setSnackbarMessage({
        id: 'packages',
        message: error.message,
        type: 'error',
        show: true,
      });
    },
  });

  const [addDataPlan, { loading: addLoading }] = useAddPackageMutation({
    onCompleted: () => {
      getDataPlans();
      setSnackbarMessage({
        id: 'add-data-plan',
        message: 'Data plan added successfully',
        type: 'success',
        show: true,
      });
    },
  });

  const [updatePackage, { loading: updateLoading }] = useUpdatePacakgeMutation({
    onCompleted: () => {
      getDataPlans();
      setSnackbarMessage({
        id: 'update-data-plan',
        message: 'Data plan updated successfully',
        type: 'success',
        show: true,
      });
    },
  });

  return {
    packages: packagesData?.getPackages.packages ?? [],
    currencySymbol: currencyData?.getCurrencySymbol.symbol ?? '',
    loading: packagesLoading || addLoading || updateLoading,
    addDataPlan,
    updatePackage,
    user,
  };
};

const DataPlansPage = () => {
  const {
    packages,
    currencySymbol,
    loading,
    addDataPlan,
    updatePackage,
    user,
  } = useDataPlans();
  const { setSnackbarMessage } = useAppContext();
  const [dialogState, setDialogState] = useState<DialogState>({
    isOpen: false,
    mode: 'create',
    data: INIT_DATAPLAN,
  });

  const handleAddDataPlanAction = () => {
    if (!user.currency) {
      setSnackbarMessage({
        id: 'network-not-selected',
        message: 'Something went wrong, please try again',
        type: 'warning',
        show: true,
      });
      return;
    }

    setDialogState({
      isOpen: true,
      mode: 'create',
      data: INIT_DATAPLAN,
    });
  };

  const handleDataPlanAction = (action: string, values: CreatePlanType) => {
    if (action === 'add' && values.amount && values.dataVolume) {
      addDataPlan({
        variables: {
          data: {
            name: values.name,
            amount: values.amount,
            dataUnit: values.dataUnit,
            country: user.country ?? '',
            currency: user.currency ?? '',
            dataVolume: values.dataVolume,
            duration: parseInt(values.duration),
          },
        },
      });
    } else if (action === 'update') {
      updatePackage({
        variables: {
          packageId: values.id,
          data: {
            active: true,
            name: values.name,
          },
        },
      });
    }
    setDialogState((prev) => ({ ...prev, isOpen: false }));
  };

  const handleEdit = (id: string) => {
    const packageToEdit = packages.find((pkg: PackageDto) => pkg.uuid === id);
    if (!packageToEdit) return;

    setDialogState({
      isOpen: true,
      mode: 'edit',
      data: {
        id,
        name: packageToEdit.name,
        duration: packageToEdit.duration.toString(),
        dataVolume: packageToEdit.dataVolume,
        country: packageToEdit.country,
        currency: packageToEdit.currency,
        amount:
          typeof packageToEdit.rate.amount === 'number'
            ? packageToEdit.rate.amount
            : 0,
        dataUnit: packageToEdit.dataUnit as DataUnitType,
      },
    });
  };

  return (
    <LoadingWrapper
      width={'100%'}
      radius="medium"
      isLoading={loading}
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
            buttonId="create-data-plan-btn"
            buttonTitle={'CREATE DATA PLAN'}
            handleButtonAction={handleAddDataPlanAction}
          />
          <br />
          {packages.length === 0 ? (
            <EmptyView icon={UpdateIcon} title="No data plan created yet!" />
          ) : (
            <Grid container rowSpacing={2} columnSpacing={2}>
              {packages.map((pkg: PackageDto) => (
                <Grid item xs={12} sm={6} md={4} key={pkg.uuid}>
                  <PlanCard
                    {...pkg}
                    currency={currencySymbol}
                    amount={pkg.amount.toString()}
                    dataVolume={pkg.dataVolume.toString()}
                    handleOptionMenuItemAction={(type: string) => {
                      if (type === 'edit') handleEdit(pkg.uuid);
                    }}
                  />
                </Grid>
              ))}
            </Grid>
          )}
        </Box>
        {dialogState.isOpen && (
          <DataPlanDialog
            data={dialogState.data}
            isOpen={dialogState.isOpen}
            setData={(newData: SetStateAction<CreatePlanType>) =>
              setDialogState((prev) => ({
                ...prev,
                data:
                  typeof newData === 'function' ? newData(prev.data) : newData,
              }))
            }
            currencySymbol={currencySymbol}
            title={
              dialogState.mode === 'create'
                ? 'Create data plan'
                : 'Edit data plan'
            }
            labelNegativeBtn={'Cancel'}
            action={dialogState.mode === 'create' ? 'add' : 'update'}
            handleSuccessAction={handleDataPlanAction}
            handleCloseAction={() =>
              setDialogState((prev) => ({ ...prev, isOpen: false }))
            }
            labelSuccessBtn={
              dialogState.mode === 'create'
                ? 'Save Data Plan'
                : 'Update Data Plan'
            }
          />
        )}
      </Paper>
    </LoadingWrapper>
  );
};

export default DataPlansPage;
