import React, { useState } from 'react';
import {
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Button,
  CircularProgress,
} from '@mui/material';
import { globalUseStyles } from '@/styles/global';
import { Formik, Form, Field, ErrorMessage } from 'formik';
import { UpdateSiteSchema } from '@/helpers/formValidators';

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
      <DialogTitle>Edit Site Name</DialogTitle>
      <DialogContent>
        <Formik
          initialValues={{ siteName: currentSiteName }}
          validationSchema={UpdateSiteSchema}
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
                      endAdornment: updateSiteLoading ? (
                        <CircularProgress size={20} />
                      ) : null,
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
                <Button
                  type="submit"
                  color="primary"
                  variant="contained"
                  disabled={updateSiteLoading}
                >
                  {updateSiteLoading ? <CircularProgress size={24} /> : 'Save'}
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
