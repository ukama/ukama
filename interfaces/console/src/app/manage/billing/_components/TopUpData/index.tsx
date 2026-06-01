/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { PackageDto, SimDto } from '@/client/graphql/generated';
import colors from '@/theme/colors';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
import CloseIcon from '@mui/icons-material/Close';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import {
  Autocomplete,
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  IconButton,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  TextField,
  Typography,
  styled,
} from '@mui/material';
import { Field, FieldArray, Form, Formik, FormikErrors } from 'formik';
import React from 'react';
import * as Yup from 'yup';

interface TopUpProps {
  onCancel: () => void;
  isToPup: boolean;
  subscriberName: string;
  handleTopUp: (simId: string, planIds: string[]) => void;
  packages: PackageDto[];
  loadingTopUp: boolean;
  sims: SimDto[];
}

interface PlanEntry {
  planId: string;
}

interface FormValues {
  simIccid: string;
  plans: PlanEntry[];
}

const StyledDialogContent = styled(DialogContent)(({ theme }) => ({
  padding: theme.spacing(3),
}));

const StyledDialogTitle = styled(DialogTitle)(({ theme }) => ({
  padding: theme.spacing(2, 3),
}));

const FormContainer = styled(Box)(() => ({
  '& .MuiAutocomplete-input': {
    height: '23px',
    padding: '7.5px 4px !important',
  },
  '& .MuiSelect-select': {
    height: '53px',
    display: 'flex',
    alignItems: 'center',
    padding: '14.5px 14px',
  },
  '& .MuiOutlinedInput-notchedOutline': {
    borderRadius: '4px',
  },
  '& .MuiTextField-root': {
    '& .MuiInputBase-root': {
      height: '53px',
    },
  },
}));
const StyledDialogActions = styled(DialogActions)(({ theme }) => ({
  padding: theme.spacing(1, 0),
}));

const PlanContainer = styled(Box)(({ theme }) => ({
  position: 'relative',
  marginTop: theme.spacing(2),
  '&:not(:last-child)': {
    marginBottom: theme.spacing(2),
  },
}));

const validationSchema = Yup.object().shape({
  simIccid: Yup.string().required('SIM ICCID is required'),
  plans: Yup.array()
    .of(
      Yup.object().shape({
        planId: Yup.string().required('Data plan is required'),
      }),
    )
    .min(1, 'At least one plan is required'),
});

const getOptionLabel = (option: SimDto) => {
  if (!option) return '';
  return `${option.iccid} - ${option.isPhysical ? 'pSIM' : 'eSIM'}`;
};

const TopUpData: React.FC<TopUpProps> = ({
  handleTopUp,
  onCancel,
  isToPup,
  subscriberName,
  packages,
  sims,
  loadingTopUp = false,
}) => {
  const defaultSim = sims.length > 0 ? sims[0] : null;

  const initialValues: FormValues = {
    simIccid: defaultSim ? defaultSim.iccid : '',
    plans: [{ planId: '' }],
  };

  const handleSubmit = (values: FormValues) => {
    const selectedSim = sims.find((sim) => sim.iccid === values.simIccid);
    if (selectedSim) {
      const planIdsToSubmit = values.plans.map((plan) => plan.planId);
      handleTopUp(selectedSim.id, planIdsToSubmit);
    }
  };

  return (
    <Dialog
      open={isToPup}
      onClose={onCancel}
      maxWidth="sm"
      fullWidth
      PaperProps={{
        sx: {
          borderRadius: 1,
        },
      }}
    >
      <StyledDialogTitle>
        <Box
          sx={{
            display: 'flex',
            alignItems: 'center',
            justifyContent: 'space-between',
          }}
        >
          <Typography variant="h6">Top up data</Typography>
          <IconButton aria-label="close" onClick={onCancel} size="small">
            <CloseIcon />
          </IconButton>
        </Box>
      </StyledDialogTitle>

      <StyledDialogContent>
        <Typography sx={{ mb: 3 }} variant="body1">
          Add one or more data plans for subscriber. Note: new data plans will
          only come into effect after current data expires.
        </Typography>

        <Box sx={{ mb: 4 }}>
          <InputLabel
            shrink
            htmlFor="email"
            sx={{
              transition: 'all 0.2s',
              zIndex: 1,
            }}
          >
            NAME
          </InputLabel>

          <Box sx={{ display: 'flex', alignItems: 'center' }}>
            <Typography
              variant="body1"
              sx={{
                flexGrow: 1,
                color: subscriberName ? 'inherit' : 'text.secondary',
              }}
            >
              {subscriberName}
            </Typography>
          </Box>
        </Box>

        {sims.length === 0 ? (
          <Typography color="error" variant="body2">
            No SIMs available. Please add SIMs to the SIM pool.
          </Typography>
        ) : (
          <Formik
            initialValues={initialValues}
            validationSchema={validationSchema}
            onSubmit={handleSubmit}
          >
            {({ errors, touched, values }) => (
              <Form>
                <FormContainer>
                  <Field name="simIccid">
                    {({ field, form }: any) => (
                      <Autocomplete
                        className="w-full"
                        options={sims}
                        getOptionLabel={getOptionLabel}
                        value={
                          sims.find((sim) => sim.iccid === field.value) || null
                        }
                        disabled={true}
                        onChange={(_) => {
                          form.setFieldValue(
                            'simIccid',
                            defaultSim?.iccid || '',
                          );
                        }}
                        renderInput={(params) => (
                          <TextField
                            {...params}
                            required
                            label="SIM ICCID"
                            placeholder="Select a SIM"
                            className="w-full"
                            InputProps={{
                              ...params.InputProps,
                              readOnly: true,
                            }}
                          />
                        )}
                        noOptionsText={
                          sims.length === 0
                            ? 'SIM pool is empty. Please upload SIMs to SIM pool first.'
                            : 'No matching SIMs found'
                        }
                      />
                    )}
                  </Field>

                  {packages.length === 0 ? (
                    <Typography color="error" variant="body2" sx={{ mt: 2 }}>
                      No data plans available. Please upload data plans.
                    </Typography>
                  ) : (
                    <FieldArray name="plans">
                      {({ push, remove }) => (
                        <Box>
                          {values.plans.map((_, index) => (
                            <PlanContainer key={index}>
                              <FormControl fullWidth>
                                <InputLabel shrink required>
                                  DATA PLAN {index + 1}
                                </InputLabel>
                                <Field name={`plans.${index}.planId`}>
                                  {({ field }: any) => (
                                    <Select
                                      {...field}
                                      input={
                                        <OutlinedInput
                                          notched
                                          label={`DATA PLAN ${index + 1}`}
                                        />
                                      }
                                      error={
                                        touched.plans?.[index]?.planId &&
                                        Boolean(
                                          (
                                            errors.plans as FormikErrors<PlanEntry>[]
                                          )?.[index]?.planId,
                                        )
                                      }
                                      fullWidth
                                    >
                                      {packages.map((pkg) => (
                                        <MenuItem
                                          key={pkg.uuid}
                                          value={pkg.uuid}
                                        >
                                          {`${pkg.name} - $${pkg.amount}/${Number(pkg.dataVolume) / 1024} GB`}
                                        </MenuItem>
                                      ))}
                                    </Select>
                                  )}
                                </Field>
                              </FormControl>
                              {values.plans.length > 1 && (
                                <IconButton
                                  size="small"
                                  onClick={() => remove(index)}
                                  sx={{
                                    position: 'absolute',
                                    right: 20,
                                    top: 8,
                                  }}
                                >
                                  <DeleteOutlineIcon />
                                </IconButton>
                              )}
                            </PlanContainer>
                          ))}

                          {values.plans.length < 5 && (
                            <Button
                              startIcon={<AddCircleOutlineIcon />}
                              onClick={() => push({ planId: '' })}
                              sx={{ mt: 1, color: colors.primaryMain }}
                            >
                              ADD ANOTHER DATA PLAN
                            </Button>
                          )}

                          {values.plans.length === 5 && (
                            <Typography
                              color="info"
                              variant="body2"
                              sx={{ mt: 1, color: colors.green }}
                            >
                              Maximum of 5 data plans reached
                            </Typography>
                          )}
                        </Box>
                      )}
                    </FieldArray>
                  )}
                </FormContainer>

                <StyledDialogActions>
                  <Button
                    onClick={onCancel}
                    color="primary"
                    disabled={loadingTopUp}
                  >
                    Cancel
                  </Button>
                  <Button
                    variant="contained"
                    type="submit"
                    disabled={
                      loadingTopUp ||
                      packages.length === 0 ||
                      sims.length === 0 ||
                      !values.simIccid ||
                      values.plans.some((plan) => !plan.planId)
                    }
                  >
                    TOP UP
                  </Button>
                </StyledDialogActions>
              </Form>
            )}
          </Formik>
        )}
      </StyledDialogContent>
    </Dialog>
  );
};

export default TopUpData;
