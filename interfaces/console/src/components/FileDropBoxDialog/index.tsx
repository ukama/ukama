/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { fileToBase64 } from '@/utils';
import CloseIcon from '@mui/icons-material/Close';
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  Typography,
} from '@mui/material';
import { useEffect, useState } from 'react';
import { useDropzone } from 'react-dropzone';

type FileDropBoxDialogProps = {
  title: string;
  note: string;
  isOpen: boolean;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction: any;
  labelNegativeBtn?: string;
};

const FileDropBoxDialog = ({
  note,
  title,
  isOpen,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
}: FileDropBoxDialogProps) => {
  const [file, setFile] = useState<any>();
  const { acceptedFiles, getRootProps, getInputProps } = useDropzone({
    accept: {
      'text/html': ['.csv'],
    },
    maxFiles: 1,
  });

  useEffect(() => {
    if (acceptedFiles.length > 0) {
      setFile(acceptedFiles[0]);
    }
  }, [acceptedFiles]);

  const handleUploadAction = () => {
    if (acceptedFiles && acceptedFiles.length > 0) {
      const file: any = acceptedFiles[0];
      fileToBase64(file)
        .then((base64String) => {
          handleSuccessAction('success', base64String);
          handleCloseAction();
        })
        .catch((error) => {
          handleSuccessAction('error', error);
          handleCloseAction();
        });
    }
  };

  return (
    <Dialog
      fullWidth
      open={isOpen}
      maxWidth="sm"
      onClose={handleCloseAction}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <DialogTitle>{title}</DialogTitle>
        <IconButton
          onClick={handleCloseAction}
          sx={{ position: 'relative', right: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </Stack>

      <DialogContent>
        <Stack
          height={'100%'}
          direction="column"
          alignItems="flex-start"
          justifyContent="center"
          spacing={2}
        >
          <Typography variant="body1" fontWeight={400}>
            Upload the SIM CSV file you received so that you can digitally
            assign SIMs to your subscribers, and authorize them to use your
            network.
          </Typography>

          {file ? (
            <Stack direction={'row'} spacing={2} alignItems={'center'}>
              <Typography variant="body1">{acceptedFiles[0].name}</Typography>
              <IconButton
                onClick={() => {
                  setFile(null);
                  // acceptedFiles.pop();
                }}
                size="small"
              >
                <DeleteOutlineOutlinedIcon fontSize="small" />
              </IconButton>
            </Stack>
          ) : (
            <Box
              sx={{
                width: '100%',
                height: '196px',
                cursor: 'pointer',
                borderRadius: '4px',
                border: '1px dashed grey',
                backgroundColor: colors.white38,
                ':hover': {
                  border: '1px dashed black',
                },
              }}
            >
              <div
                id="csv-file-input"
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
                  {note}
                </Typography>
              </div>
            </Box>
          )}
        </Stack>
      </DialogContent>

      <DialogActions>
        <Stack direction={'row'} alignItems="center" spacing={2}>
          {labelNegativeBtn && (
            <Button
              variant="text"
              color={'primary'}
              onClick={handleCloseAction}
            >
              {labelNegativeBtn}
            </Button>
          )}
          {labelSuccessBtn && (
            <Button variant="contained" onClick={handleUploadAction}>
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default FileDropBoxDialog;
