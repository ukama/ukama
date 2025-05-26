/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { Sim_Types, useUploadSimsMutation } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import { fileToBase64, setQueryParam } from '@/utils';
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined';
import {
  AlertColor,
  Box,
  Button,
  IconButton,
  Stack,
  Typography,
} from '@mui/material';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';
import { useDropzone } from 'react-dropzone';
import LoadingSkeleton from './skelton';

const Sims = () => {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [loading, setLoading] = useState(false);
  const { env, setSnackbarMessage } = useAppContext();
  const [file, setFile] = useState<any>();
  const { acceptedFiles, getRootProps, getInputProps } = useDropzone({
    accept: {
      'text/html': ['.csv'],
    },
    maxFiles: 1,
  });

  const [uploadSimPool, { loading: uploadSimsLoading }] = useUploadSimsMutation(
    {
      onCompleted: () => {
        setSnackbarMessage({
          id: 'sim-pool-uploaded',
          message: 'Sims uploaded successfully',
          type: 'success' as AlertColor,
          show: true,
        });
        const p = setQueryParam(
          'pool',
          'true',
          searchParams.toString(),
          pathname,
        );
        router.push(`/configure/complete?${p.toString()}`);
      },
      onError: (error) => {
        setLoading(false);
        setSnackbarMessage({
          id: 'sim-pool-error',
          message: error.message,
          type: 'error' as AlertColor,
          show: true,
        });
      },
    },
  );

  useEffect(() => {
    if (uploadSimsLoading) {
      setLoading(true);
    }
  }, [uploadSimsLoading]);

  useEffect(() => {
    if (acceptedFiles.length > 0) {
      setFile(acceptedFiles[0]);
    }
  }, [acceptedFiles]);

  const handleSkip = () => {
    const p = setQueryParam('pool', 'false', searchParams.toString(), pathname);
    router.push(`/configure/complete?${p.toString()}`);
  };

  const handleNext = () => {
    if (
      acceptedFiles &&
      acceptedFiles.length > 0 &&
      env.SIM_TYPE !== Sim_Types.Unknown
    ) {
      const file: any = acceptedFiles[0];
      fileToBase64(file)
        .then((base64String) => {
          handleUploadSimsAction('success', base64String);
        })
        .catch((error) => {
          handleUploadSimsAction('error', error);
        });
    } else {
      setSnackbarMessage({
        id: 'sim-pool-error',
        message: 'Please add/drop file to upload',
        type: 'error' as AlertColor,
        show: true,
      });
    }
  };

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

  if (loading) return <LoadingSkeleton />;

  return (
    <Stack>
      <Typography variant="h4" fontWeight={500} mb={2}>
        Upload Sims
      </Typography>

      <Stack direction={'column'} spacing={4}>
        <Typography variant={'body1'} fontWeight={400}>
          Upload the SIMs belonging to your organization, so that you can later
          authorize subscribers to start using your network.
        </Typography>

        {file ? (
          <Stack direction={'row'} spacing={2} alignItems={'center'}>
            <Typography variant="body1">{acceptedFiles[0].name}</Typography>
            <IconButton
              onClick={() => {
                setFile(null);
                // acceptedFiles?.pop();
              }}
              size="small"
            >
              <DeleteOutlineOutlinedIcon fontSize="small" />
            </IconButton>
          </Stack>
        ) : (
          <Box
            sx={{
              py: 8,
              width: '100%',
              height: '94px',
              display: 'flex',
              cursor: 'pointer',
              borderRadius: '4px',
              justifyContent: 'center',
              border: '1px dashed grey',
              backgroundColor: colors.primaryMain02,
              ':hover': {
                border: `1px dashed ${colors.primaryMain}`,
              },
            }}
          >
            <div
              id="csv-file-input-onboarding"
              {...getRootProps({ className: 'dropzone' })}
              style={{
                width: '100%',
                height: '100%',
                display: 'flex',
                alignItems: 'center',
                justifyContent: 'center',
              }}
            >
              <input {...getInputProps()} />
              <Typography variant="body2" sx={{ cursor: 'inherit' }}>
                Drag & Drop Or Choose file to upload.
              </Typography>
            </div>
          </Box>
        )}
      </Stack>

      <Stack
        mt={{ xs: 4, md: 6 }}
        spacing={2}
        direction={'row'}
        justifyContent={'space-between'}
      >
        <Button
          variant="text"
          onClick={handleSkip}
          sx={{ color: colors.black70, p: 0 }}
        >
          Skip
        </Button>
        <Button variant="contained" onClick={handleNext}>
          Upload sims
        </Button>
      </Stack>
    </Stack>
  );
};

export default Sims;
