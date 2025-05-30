/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { UpdateSiteSchema } from '@/helpers/formValidators';
import { GlobalInput } from '@/styles/global';
import CloseIcon from '@mui/icons-material/Close';
import {
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogContent,
  DialogTitle,
  IconButton,
} from '@mui/material';
import { ErrorMessage, Field, Form, Formik } from 'formik';
import React from 'react';

interface EditSiteDialogProps {
  open: boolean;
  siteId: string;
  currentSiteName: string;
  onClose: () => void;
  onSave: (siteId: string, newSiteName: string) => void;
  updateSiteLoading: boolean;
}

const EditSiteDialog: React.FC<EditSiteDialogProps> = ({
  open,
  siteId,
  currentSiteName,
  onClose,
  onSave,
  updateSiteLoading,
}) => {
  const handleSubmit = async (values: { siteName: string }) => {
    try {
      await onSave(siteId, values.siteName);
    } finally {
      onClose();
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      PaperProps={{
        sx: {
          width: '600px',
        },
      }}
    >
      <DialogTitle sx={{ m: 0, p: 2, position: 'relative' }}>
        Edit Site Name
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{
            position: 'absolute',
            right: 10,
            top: 8,
            color: (theme) => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <Formik
        initialValues={{ siteName: currentSiteName }}
        validationSchema={UpdateSiteSchema}
        onSubmit={handleSubmit}
      >
        {({ touched, errors }) => (
          <Form>
            <DialogContent>
              <Field name="siteName">
                {({ field }: { field: any }) => (
                  <GlobalInput
                    {...field}
                    autoFocus
                    margin="dense"
                    label="Site Name"
                    fullWidth
                    variant="outlined"
                    slotProps={{
                      inputLabel: {
                        shrink: true,
                      },
                      input: {
                        endAdornment: updateSiteLoading ? (
                          <CircularProgress size={20} />
                        ) : null,
                      },
                    }}
                    error={touched.siteName && !!errors.siteName}
                    helperText={<ErrorMessage name="siteName" />}
                  />
                )}
              </Field>
            </DialogContent>

            <Box
              sx={{
                display: 'flex',
                justifyContent: 'flex-end',
                p: 3,
                pt: 0,
              }}
            >
              <Button
                type="button"
                onClick={onClose}
                color="secondary"
                sx={{ mr: 1 }}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                color="primary"
                variant="contained"
                disabled={updateSiteLoading}
              >
                {updateSiteLoading ? <CircularProgress size={24} /> : 'Save'}
              </Button>
            </Box>
          </Form>
        )}
      </Formik>
    </Dialog>
  );
};

export default EditSiteDialog;
