/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { SimPoolResDto } from '@/client/graphql/generated';
import { SubscriberDetailsType } from '@/types';
import styled from '@emotion/styled';
import CloseIcon from '@mui/icons-material/Close';
import { IconButton } from '@mui/material';
import {
  Autocomplete,
  Box,
  Button,
  DialogActions,
  DialogContent,
  DialogTitle,
  TextField,
  Typography,
} from '@mui/material';
import { Field, FormikProps } from 'formik';

const CloseButtonStyle = styled(IconButton)({
  position: 'absolute',
  right: 10,
  top: 14,
});

interface SimSelectionStepProps {
  formik: FormikProps<SubscriberDetailsType>;
  sims: SimPoolResDto[];
  onNext: () => void;
  onBack: () => void;
  onClose: () => void;
  onSimSelected: (sim: SimPoolResDto | null) => void;
}

/** Step 1 — subscriber selects or types a SIM ICCID. */
const SimSelectionStep: React.FC<SimSelectionStepProps> = ({
  formik,
  sims,
  onNext,
  onBack,
  onClose,
  onSimSelected,
}) => (
  <>
    <DialogTitle>
      Enter SIM ICCID
      <CloseButtonStyle aria-label="close" onClick={onClose}>
        <CloseIcon />
      </CloseButtonStyle>
    </DialogTitle>

    <DialogContent>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 3 }}>
        Please select the SIM you would like to assign to the subscriber, and
        enter the ICCID found on the back of the card. Please ensure the ICCID
        is correct, because this cannot be undone.
      </Typography>
      <Field name="simIccid">
        {({ field, form }: { field: { value: string }; form: FormikProps<SubscriberDetailsType> }) => (
          <Autocomplete
            freeSolo
            value={field.value}
            options={sims.map((s) => s.iccid)}
            sx={{
              '.MuiAutocomplete-inputRoot .MuiAutocomplete-input': {
                padding: '13px !important',
              },
            }}
            onInputChange={(_, newValue) => {
              form.setFieldValue('simIccid', newValue || undefined);
              onSimSelected(sims.find((s) => s.iccid === newValue) ?? null);
            }}
            renderInput={(params) => (
              <TextField
                {...params}
                {...field}
                label="SIM ICCID*"
                InputLabelProps={{ shrink: true }}
                sx={{ '.MuiOutlinedInput-root': { padding: 0 } }}
              />
            )}
            fullWidth
          />
        )}
      </Field>
    </DialogContent>

    <DialogActions sx={{ display: 'flex', justifyContent: 'space-between' }}>
      <Button sx={{ p: 0 }} onClick={onBack}>
        Back
      </Button>
      <Box sx={{ display: 'flex', gap: 1 }}>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={onNext}
          disabled={!formik.isValid || !formik.values.simIccid}
          variant="contained"
        >
          Next
        </Button>
      </Box>
    </DialogActions>
  </>
);

export default SimSelectionStep;
