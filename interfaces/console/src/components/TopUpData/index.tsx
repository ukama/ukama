import React from 'react';
import { PackageDto, SimDto } from '@/client/graphql/generated';
import CloseIcon from '@mui/icons-material/Close';
import DeleteOutlineIcon from '@mui/icons-material/DeleteOutline';
import AddCircleOutlineIcon from '@mui/icons-material/AddCircleOutline';
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
import { Formik, Form, Field, FieldArray, FormikErrors } from 'formik';
import * as Yup from 'yup';
import colors from '@/theme/colors';

interface TopUpProps {
  onCancel: () => void;
  isToPup: boolean;
  subscriberName: string;
  handleTopUp: (plans: { planId: string; simId: string }[]) => void;
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

// Styled components remain the same
const StyledDialogContent = styled(DialogContent)(({ theme }) => ({
  padding: theme.spacing(3),
}));

const StyledDialogTitle = styled(DialogTitle)(({ theme }) => ({
  padding: theme.spacing(2, 3),
}));

const NameLabel = styled(Typography)(({ theme }) => ({
  color: theme.palette.text.secondary,
  marginBottom: theme.spacing(0.5),
  fontSize: '0.875rem',
  fontWeight: 500,
}));

const NameValue = styled(Typography)(({ theme }) => ({
  marginBottom: theme.spacing(3),
}));

const FormContainer = styled(Box)(({ theme }) => ({
  '& .MuiFormControl-root': {
    marginBottom: theme.spacing(1),
  },
  '& .MuiInputBase-root': {
    height: '56px',
  },
  '& .MuiAutocomplete-input': {
    height: '23px',
    padding: '7.5px 4px !important',
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
  const initialValues: FormValues = {
    simIccid: '',
    plans: [{ planId: '' }],
  };

  const handleSubmit = (values: FormValues) => {
    const selectedSim = sims.find((sim) => sim.iccid === values.simIccid);
    if (selectedSim) {
      const plansToSubmit = values.plans.map((plan) => ({
        planId: plan.planId,
        simId: selectedSim.id,
      }));
      handleTopUp(plansToSubmit);
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
        <Typography sx={{ color: 'text.secondary', mb: 3 }}>
          Add one or more data plans for subscriber. Note: new data plans will
          only come into effect after current data expires.
        </Typography>

        <NameLabel>NAME</NameLabel>
        <NameValue>{subscriberName}</NameValue>

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
                      onChange={(_, newValue) => {
                        form.setFieldValue('simIccid', newValue?.iccid || '');
                      }}
                      renderInput={(params) => (
                        <TextField
                          {...params}
                          label="SIM ICCID"
                          placeholder="Select a SIM"
                          className="w-full"
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
                                      <MenuItem key={pkg.uuid} value={pkg.uuid}>
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
                                sx={{ position: 'absolute', right: 20, top: 8 }}
                              >
                                <DeleteOutlineIcon />
                              </IconButton>
                            )}
                          </PlanContainer>
                        ))}
                        <Button
                          startIcon={<AddCircleOutlineIcon />}
                          onClick={() => push({ planId: '' })}
                          sx={{ mt: 1, color: colors.primaryMain }}
                        >
                          ADD ANOTHER DATA PLAN
                        </Button>
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
      </StyledDialogContent>
    </Dialog>
  );
};

export default TopUpData;
