import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Button,
} from '@mui/material';
import { globalUseStyles } from '@/styles/global';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import * as Yup from 'yup';

const validationSchema = Yup.object().shape({
  siteName: Yup.string()
    .required('Site name is required')
    .matches(
      /^[a-z0-9-]*$/,
      'Site name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
    ),
});

interface EditSiteDialogProps {
  open: boolean;
  siteId: string;
  currentSiteName: string;
  onClose: () => void;
  onSave: (siteId: string, newSiteName: string) => void;
}

const EditSiteDialog: React.FC<EditSiteDialogProps> = ({
  open,
  siteId,
  currentSiteName,
  onClose,
  onSave,
}) => {
  const handleSubmit = (values: { siteName: string }) => {
    onSave(siteId, values.siteName);
    onClose();
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
      <DialogTitle>Edit Site Name</DialogTitle>
      <DialogContent>
        <Formik
          initialValues={{ siteName: currentSiteName }}
          validationSchema={validationSchema}
          onSubmit={handleSubmit}
        >
          {({ touched, errors }) => (
            <Form>
              <Field name="siteName">
                {({ field }: { field: any }) => (
                  <TextField
                    {...field}
                    autoFocus
                    margin="dense"
                    label="Site Name"
                    fullWidth
                    variant="outlined"
                    InputLabelProps={{ shrink: true }}
                    InputProps={{
                      classes: { input: globalUseStyles().inputFieldStyle },
                    }}
                    error={touched.siteName && !!errors.siteName}
                    helperText={<ErrorMessage name="siteName" />}
                  />
                )}
              </Field>
              <DialogActions>
                <Button type="button" onClick={onClose} color="secondary">
                  Cancel
                </Button>
                <Button type="submit" color="primary">
                  Save
                </Button>
              </DialogActions>
            </Form>
          )}
        </Formik>
      </DialogContent>
    </Dialog>
  );
};

export default EditSiteDialog;
