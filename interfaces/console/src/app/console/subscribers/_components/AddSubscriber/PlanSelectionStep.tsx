/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { PackageDto } from '@/client/graphql/generated';
import { SubscriberDetailsType } from '@/types';
import styled from '@emotion/styled';
import CloseIcon from '@mui/icons-material/Close';
import { IconButton, Select } from '@mui/material';
import {
  Box,
  Button,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  FormHelperText,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Typography,
} from '@mui/material';
import { Field, FormikProps } from 'formik';

const SelectStyle = styled(Select)({ width: '100%', height: '48px' });
const CloseButtonStyle = styled(IconButton)({ position: 'absolute', right: 10, top: 14 });

const NoItemMessage = ({ message }: { message: string }) => (
  <MenuItem disabled value="" sx={{ m: 0, p: '6px 16px' }}>
    <Typography variant="body1">{message}</Typography>
  </MenuItem>
);

interface PlanSelectionStepProps {
  formik: FormikProps<SubscriberDetailsType>;
  packages: PackageDto[];
  currencySymbol: string;
  onBack: () => void;
  onClose: () => void;
  onSubmit: () => void;
}

/** Step 2 — subscriber selects a data plan. */
const PlanSelectionStep: React.FC<PlanSelectionStepProps> = ({
  formik,
  packages,
  currencySymbol,
  onBack,
  onClose,
  onSubmit,
}) => (
  <>
    <DialogTitle>
      Select data plan
      <CloseButtonStyle aria-label="close" onClick={onClose}>
        <CloseIcon />
      </CloseButtonStyle>
    </DialogTitle>

    <DialogContent>
      <Typography variant="subtitle1" color="text.secondary" sx={{ mb: 2 }}>
        Select the purchased data plan
      </Typography>
      <Field name="plan" id="add-subscriber-plan-select">
        {({ field, meta }: { field: object; meta: { touched: boolean; error?: string } }) => (
          <FormControl fullWidth error={meta.touched && Boolean(meta.error)}>
            <InputLabel htmlFor="outlined-plan" shrink>
              DATA PLAN
            </InputLabel>
            <SelectStyle
              {...field}
              label="DATA PLAN"
              input={<OutlinedInput notched label="DATA PLAN" id="outlined-plan" />}
            >
              {packages.length === 0 ? (
                <NoItemMessage message="No packages available. Please add packages first." />
              ) : (
                packages.map((plan) => (
                  <MenuItem key={plan.uuid} value={plan.uuid}>
                    <Typography variant="body1">
                      {`${plan.name} - ${currencySymbol} ${plan.amount}/${plan.dataVolume} ${plan.dataUnit}`}
                    </Typography>
                  </MenuItem>
                ))
              )}
            </SelectStyle>
            {meta.touched && meta.error && (
              <FormHelperText>{meta.error}</FormHelperText>
            )}
          </FormControl>
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
          onClick={onSubmit}
          disabled={!formik.isValid || !formik.values.plan}
          variant="contained"
        >
          ADD SUBSCRIBER
        </Button>
      </Box>
    </DialogActions>
  </>
);

export default PlanSelectionStep;
