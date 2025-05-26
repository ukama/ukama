/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogContent,
  DialogTitle,
  FormControlLabel,
  IconButton,
  Stack,
  Switch,
  TextField,
  Typography,
} from '@mui/material';

import { Formik } from 'formik';
import * as Yup from 'yup';

type AddNetworkDialogProps = {
  title: string;
  isOpen: boolean;
  loading: boolean;
  description: string;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction?: any;
  labelNegativeBtn?: string;
};

interface AddNetworkForm {
  name: string;
  isDefault: boolean;
  budget: number;
  countries: { name: string; code: string }[];
  networks: { id: string; name: string; isDefault: boolean }[];
}

const validationSchema = Yup.object({
  networks: Yup.array().optional().default([]),
  isDefault: Yup.boolean().default(false),
  countries: Yup.array().optional().default([]),
  name: Yup.string()
    .required('Network name is required')
    .matches(
      /^[a-z0-9-]*$/,
      'Network name must be lowercase alphanumeric and should not contain spaces, "-" are allowed.',
    ),
  budget: Yup.number().default(0),
});

const initialValues: AddNetworkForm = {
  name: '',
  budget: 0,
  isDefault: false,
  countries: [],
  networks: [],
};

const AddNetworkDialog = ({
  title,
  isOpen,
  loading,
  description,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
}: AddNetworkDialogProps) => {
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
          <CloseIcon fontSize="small" />
        </IconButton>
      </Stack>

      <DialogContent sx={{ maxHeight: '400px', overflow: 'auto' }}>
        <Formik
          initialValues={initialValues}
          validationSchema={validationSchema}
          onSubmit={(values) => {
            handleSuccessAction(values);
          }}
        >
          {({
            values,
            errors,
            touched,
            handleChange,
            handleSubmit,
            handleBlur,
            setFieldValue,
          }) => (
            <form onSubmit={handleSubmit}>
              <Stack spacing={2} direction={'column'} alignItems="start">
                {description && (
                  <Typography variant="body1">{description}</Typography>
                )}
                <TextField
                  fullWidth
                  name={'name'}
                  size="medium"
                  placeholder="network-name"
                  label={'Network name'}
                  InputLabelProps={{
                    shrink: true,
                  }}
                  onBlur={handleBlur}
                  onChange={handleChange}
                  value={values.name}
                  helperText={touched.name && errors.name}
                  error={touched.name && Boolean(errors.name)}
                  id={'name'}
                />
                <FormControlLabel
                  sx={{ display: 'none' }}
                  control={
                    <Switch
                      defaultChecked={false}
                      value={values.isDefault}
                      checked={values.isDefault}
                      onChange={() =>
                        setFieldValue('isDefault', !values.isDefault)
                      }
                    />
                  }
                  label="Make this network default"
                />
                <Stack
                  width="100%"
                  spacing={2}
                  direction={'row'}
                  alignItems="center"
                  justifyContent={'flex-end'}
                >
                  <Button
                    variant="text"
                    color={'primary'}
                    onClick={handleCloseAction}
                  >
                    {labelNegativeBtn}
                  </Button>

                  <Button type="submit" variant="contained" disabled={loading}>
                    {labelSuccessBtn}
                  </Button>
                </Stack>
              </Stack>
            </form>
          )}
        </Formik>
      </DialogContent>
    </Dialog>
  );
};

export default AddNetworkDialog;
