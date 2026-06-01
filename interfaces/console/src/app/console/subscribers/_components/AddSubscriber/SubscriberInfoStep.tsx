/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import colors from '@/theme/colors';
import { SubscriberDetailsType } from '@/types';
import styled from '@emotion/styled';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { Field, FormikProps } from 'formik';

const CloseButtonStyle = styled(IconButton)({
  position: 'absolute',
  right: 10,
  top: 14,
});

interface SubscriberInfoStepProps {
  formik: FormikProps<SubscriberDetailsType>;
  onNext: () => void;
  onClose: () => void;
}

/** Step 0 — collects subscriber name and email. */
const SubscriberInfoStep: React.FC<SubscriberInfoStepProps> = ({
  formik,
  onNext,
  onClose,
}) => (
  <>
    <DialogTitle sx={{ color: colors.black }}>
      Add Subscriber
      <CloseButtonStyle aria-label="close" onClick={onClose}>
        <CloseIcon />
      </CloseButtonStyle>
    </DialogTitle>

    <DialogContent>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 3 }}>
        Enter basic information about the subscriber, so that they can be
        authorized to use the network.
      </Typography>
      <Stack direction="column" spacing={2}>
        <Field name="name">
          {({
            field,
            meta,
          }: {
            field: object;
            meta: { touched: boolean; error?: string };
          }) => (
            <TextField
              {...field}
              required
              fullWidth
              label="Name"
              sx={{
                height: '48px',
                '& .MuiInputBase-root': { height: '100%' },
              }}
              error={meta.touched && Boolean(meta.error)}
              helperText={meta.touched && meta.error}
              slotProps={{ inputLabel: { shrink: true } }}
            />
          )}
        </Field>
        <Field name="email">
          {({
            field,
            meta,
          }: {
            field: object;
            meta: { touched: boolean; error?: string };
          }) => (
            <TextField
              {...field}
              required
              fullWidth
              label="Email"
              sx={{
                height: '48px',
                '& .MuiInputBase-root': { height: '100%' },
              }}
              error={meta.touched && Boolean(meta.error)}
              helperText={meta.touched && meta.error}
              slotProps={{ inputLabel: { shrink: true } }}
            />
          )}
        </Field>
      </Stack>
    </DialogContent>

    <DialogActions>
      <Button onClick={onClose}>Cancel</Button>
      <Button
        onClick={onNext}
        disabled={
          !formik.values.name || !formik.values.email || !formik.isValid
        }
        variant="contained"
      >
        Next
      </Button>
    </DialogActions>
  </>
);

export default SubscriberInfoStep;
