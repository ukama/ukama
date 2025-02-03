/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Button,
  CircularProgress,
  Typography,
} from '@mui/material';
import { globalUseStyles } from '@/styles/global';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';

interface EditSiteDialogProps {
  open: boolean;
  siteId: string;
  currentSiteName: string;
  onClose: () => void;
  onSave: (siteId: string, newSiteName: string) => void;
  updateSiteLoading: boolean;
}

const UpdateSiteSchema = Yup.object().shape({
  siteName: Yup.string()
    .required('Site name is required')
    .max(253, 'Site name must be less than 253 characters')
    .matches(
      /^[a-z0-9][a-z0-9-]*[a-z0-9]$/,
      'Site name must consist of lowercase letters, numbers, and hyphens only. It must start and end with an alphanumeric character',
    )
    .test(
      'no-consecutive-hyphens',
      'Site name cannot contain consecutive hyphens',
      (value) => !value || !value.includes('--'),
    )
    .test(
      'length-per-label',
      'Each part between hyphens must be 63 characters or less',
      (value) =>
        !value || value.split('-').every((label) => label.length <= 63),
    ),
});

const EditSiteDialog: React.FC<EditSiteDialogProps> = ({
  open,
  siteId,
  currentSiteName,
  onClose,
  onSave,
  updateSiteLoading,
}) => {
  const handleSubmit = async (
    values: { siteName: string },
    { setSubmitting }: any,
  ) => {
    try {
      await onSave(siteId, values.siteName.trim());
      onClose();
    } catch (error) {
      console.error('Error updating site name:', error);
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <Dialog
      open={open}
      onClose={onClose}
      PaperProps={{
        sx: {
          width: '600px',
          maxWidth: '90vw',
        },
      }}
    >
      <DialogTitle>Edit Site Name</DialogTitle>
      <Formik
        initialValues={{ siteName: currentSiteName }}
        validationSchema={UpdateSiteSchema}
        onSubmit={handleSubmit}
        validateOnMount
      >
        {({ isValid, dirty, touched, errors, isSubmitting }) => (
          <Form>
            <DialogContent>
              <Field name="siteName">
                {({ field }: { field: any }) => (
                  <>
                    <TextField
                      {...field}
                      autoFocus
                      margin="dense"
                      label="Site Name"
                      fullWidth
                      variant="outlined"
                      InputLabelProps={{ shrink: true }}
                      error={touched.siteName && !!errors.siteName}
                      helperText={touched.siteName && errors.siteName}
                      disabled={updateSiteLoading || isSubmitting}
                      sx={{
                        '& .MuiOutlinedInput-root': {
                          '&.Mui-error': {
                            '& .MuiOutlinedInput-notchedOutline': {
                              borderColor: 'error.main',
                            },
                          },
                        },
                      }}
                      InputProps={{
                        classes: { input: globalUseStyles().inputFieldStyle },
                        endAdornment: (updateSiteLoading || isSubmitting) && (
                          <CircularProgress size={20} color="inherit" />
                        ),
                      }}
                    />
                    <Typography
                      variant="caption"
                      sx={{
                        display: 'block',
                        mt: 1,
                        ml: 1,
                      }}
                    >
                      Site name must be lowercase, may contain hyphens, and must
                      start and end with a letter or number.
                    </Typography>
                  </>
                )}
              </Field>
            </DialogContent>
            <DialogActions sx={{ padding: 2, justifyContent: 'flex-end' }}>
              <Button
                onClick={onClose}
                color="secondary"
                disabled={updateSiteLoading || isSubmitting}
              >
                Cancel
              </Button>
              <Button
                type="submit"
                variant="contained"
                disabled={
                  !isValid || !dirty || updateSiteLoading || isSubmitting
                }
                sx={{ ml: 1 }}
              >
                {updateSiteLoading || isSubmitting ? (
                  <CircularProgress size={24} color="inherit" />
                ) : (
                  'Save'
                )}
              </Button>
            </DialogActions>
          </Form>
        )}
      </Formik>
    </Dialog>
  );
};

export default EditSiteDialog;