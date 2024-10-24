import React from 'react';
import { PackageDto, SimDto } from '@/client/graphql/generated';
import CloseIcon from '@mui/icons-material/Close';
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
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';

interface TopUpProps {
  onCancel: () => void;
  isToPup: boolean;
  subscriberName: string;
  handleTopUp: (planId: string, simId: string) => void;
  packages: PackageDto[];
  loadingTopUp: boolean;
  sims: SimDto[];
}

interface FormValues {
  simIccid: string;
  planId: string;
}

// Styled components
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
  // Add consistent height for both fields
  '& .MuiInputBase-root': {
    height: '56px', // Standard Material-UI height
  },
  // Ensure Autocomplete internal input matches height
  '& .MuiAutocomplete-input': {
    height: '23px', // Adjust this value to match the Select input height
    padding: '7.5px 4px !important', // Add important to override Autocomplete's default styles
  },
}));

const StyledDialogActions = styled(DialogActions)(({ theme }) => ({
  padding: theme.spacing(1, 0),
}));

const validationSchema = Yup.object().shape({
  simIccid: Yup.string().required('SIM ICCID is required'),
  planId: Yup.string().required('Data plan is required'),
});

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
    planId: '',
  };

  const handleSubmit = (values: FormValues) => {
    const selectedSim = sims.find((sim) => sim.iccid === values.simIccid);
    if (selectedSim) {
      handleTopUp(values.planId, selectedSim.id);
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
          Add more data for subscriber. Note: new data plan will only come into
          effect after current data expires.
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
                      options={sims}
                      getOptionLabel={(option) => option.iccid || ''}
                      value={
                        sims.find((sim) => sim.iccid === field.value) || null
                      }
                      onChange={(_, newValue) => {
                        form.setFieldValue('simIccid', newValue?.iccid || '');
                      }}
                      renderInput={(params) => (
                        <TextField
                          {...params}
                          {...field}
                          label="SIM ICCID*"
                          error={touched.simIccid && Boolean(errors.simIccid)}
                          helperText={touched.simIccid && errors.simIccid}
                          InputLabelProps={{ shrink: true }}
                          fullWidth
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
                    No data plans available. Please contact support to set up
                    data plans.
                  </Typography>
                ) : (
                  <FormControl fullWidth sx={{ mt: 3 }}>
                    <InputLabel shrink required>
                      DATA PLAN
                    </InputLabel>
                    <Field name="planId">
                      {({ field }: any) => (
                        <Select
                          {...field}
                          input={<OutlinedInput notched label="DATA PLAN" />}
                          error={touched.planId && Boolean(errors.planId)}
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
                    !values.planId
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
