/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import {
  SimPoolResDto,
  Sim_Status,
  Sim_Types,
  useGetSimsFromPoolQuery,
  useUploadSimsMutation,
} from '@/client/graphql/generated';
import EmptyView from '@/components/EmptyView';
import FileDropBoxDialog from '@/components/FileDropBoxDialog';
import LoadingWrapper from '@/components/LoadingWrapper';
import PageContainerHeader from '@/components/PageContainerHeader';
import SimpleDataTable from '@/components/SimpleDataTable';
import { MANAGE_SIM_POOL_COLUMN } from '@/constants';
import { useAppContext } from '@/context';
import SimCardIcon from '@mui/icons-material/SimCard';
import { AlertColor, Box, Paper } from '@mui/material';
import { useState } from 'react';

const Page = () => {
  const [data, setData] = useState<SimPoolResDto[]>([]);
  const { setSnackbarMessage, env } = useAppContext();
  const [isUploadSims, setIsUploadSims] = useState<boolean>(false);

  const { loading: simsLoading, refetch: refetchSims } =
    useGetSimsFromPoolQuery({
      fetchPolicy: 'cache-and-network',
      skip: false,
      variables: {
        data: {
          status: Sim_Status.All,
          type: env.SIM_TYPE as Sim_Types,
        },
      },
      onCompleted: (data) => {
        setData(data?.getSimsFromPool?.sims ?? []);
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'sim-pool',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    });

  const [uploadSimPool, { loading: uploadSimsLoading }] = useUploadSimsMutation(
    {
      onCompleted: () => {
        refetchSims();
        setSnackbarMessage({
          id: 'sim-pool-uploaded',
          message: 'Sims uploaded successfully',
          type: 'success' as AlertColor,
          show: true,
        });
        setIsUploadSims(false);
      },
      onError: (error) => {
        setSnackbarMessage({
          id: 'sim-pool-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    },
  );

  const handleUploadSimsAction = (action: string, value: string) => {
    if (action === 'error') {
      setSnackbarMessage({
        id: 'sim-pool-parsing-error',
        message: value,
        type: 'error' as AlertColor,
        show: true,
      });
    } else if (action === 'success') {
      uploadSimPool({
        variables: {
          data: {
            data: value,
            simType: env.SIM_TYPE as Sim_Types,
          },
        },
      });
    }
  };

  return (
    <LoadingWrapper
      width={'100%'}
      radius="medium"
      height={'calc(100vh - 244px)'}
      isLoading={uploadSimsLoading ?? simsLoading}
    >
      <Paper
        sx={{
          py: { xs: 1.5, md: 3 },
          px: { xs: 2, md: 4 },
          overflow: 'hidden',
          borderRadius: '10px',
          height: '100%',
        }}
      >
        <Box sx={{ width: '100%', height: '100%' }}>
          <PageContainerHeader
            showSearch={false}
            title={'My SIM pool'}
            buttonTitle={'CLAIM SIMS'}
            subtitle={data.length.toString() ?? '0'}
            handleButtonAction={() => setIsUploadSims(true)}
          />
          <br />
          {data.length === 0 ? (
            <EmptyView icon={SimCardIcon} title="No sims in sim pool!" />
          ) : (
            <SimpleDataTable
              dataset={data}
              isIdHyperlink={true}
              columns={MANAGE_SIM_POOL_COLUMN}
              hyperlinkPrefix="/console/subscribers?"
              height="calc(100vh - 320px)"
            />
          )}
        </Box>
        {isUploadSims && (
          <FileDropBoxDialog
            isOpen={isUploadSims}
            labelSuccessBtn={'Claim'}
            labelNegativeBtn={'Cancel'}
            title={'Upload SIMs'}
            note={'Drag & Drop Or Choose file to upload. (format:*.csv)'}
            handleSuccessAction={handleUploadSimsAction}
            handleCloseAction={() => setIsUploadSims(false)}
          />
        )}
      </Paper>
    </LoadingWrapper>
  );
};

export default Page;
