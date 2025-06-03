import {
  AllocateSimApiDto,
  PackageDto,
  SimPoolResDto,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import colors from '@/theme/colors';
import { SubscriberDetailsType } from '@/types';
import styled from '@emotion/styled';
import CloseIcon from '@mui/icons-material/Close';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Autocomplete,
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  FormHelperText,
  IconButton,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  Stack,
  TextField,
  Typography,
} from '@mui/material';
import { Field, Form, Formik } from 'formik';
import QRCode from 'qrcode.react';
import React, { useState } from 'react';
import * as Yup from 'yup';

const subscriberDetailsSchema = Yup.object().shape({
  simIccid: Yup.string()
    .required('SIM ICCID is required')
    .matches(/^\d{19,20}$/, 'Invalid ICCID format'),
});
const stepZeroSchema = Yup.object().shape({
  name: Yup.string()
    .required('Name is required')
    .min(2, 'Name must be at least 2 characters')
    .max(50, 'Name must not exceed 50 characters'),
  email: Yup.string()
    .required('Email is required')
    .matches(
      /^(([^<>()[\]\\.,;:\s@"]+(\.[^<>()[\]\\.,;:\s@"]+)*)|(".+"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$/,
      'Invalid email format',
    ),
});

const SelectStyle = styled(Select)({
  width: '100%',
  height: '48px',
});

const CloseButtonStyle = styled(IconButton)({
  position: 'absolute',
  right: 10,
  top: 14,
});

interface SubscriberFormProps {
  isOpen: boolean;
  currencySymbol: string;
  handleCloseAction: () => void;
  handleAddSubscriber: (
    subscriber: SubscriberDetailsType,
  ) => Promise<AllocateSimApiDto>;
  packages: PackageDto[];
  sims: SimPoolResDto[];
  isLoading: boolean;
}

const AddSubscriberStepperDialog: React.FC<SubscriberFormProps> = ({
  isOpen,
  handleCloseAction,
  handleAddSubscriber,
  packages,
  sims,
  isLoading,
  currencySymbol,
}) => {
  const { setSnackbarMessage } = useAppContext();
  const [activeStep, setActiveStep] = useState(0);
  const [showQrCode, setShowQrCode] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submissionData, setSubmissionData] =
    useState<AllocateSimApiDto | null>(null);
  const [selectedSim, setSelectedSim] = useState<SimPoolResDto | null>(null);

  const initialValues: SubscriberDetailsType = {
    name: '',
    email: '',
    simIccid: '',
    plan: '',
  };

  const handleClose = () => {
    setActiveStep(0);
    setError(null);
    handleCloseAction();
  };

  const handleSubmit = async (values: SubscriberDetailsType) => {
    try {
      if (values.plan === '') {
        setSnackbarMessage({
          show: true,
          type: 'error',
          id: 'plane-not-found',
          message: 'Please select a plan',
        });
        return;
      }
      if (values.simIccid === '') {
        setSnackbarMessage({
          show: true,
          type: 'error',
          id: 'sim-not-found',
          message: 'Please provide sim iccid.',
        });
        return;
      }
      setError(null);
      const response = await handleAddSubscriber(values);
      setSubmissionData(response);
      setActiveStep(3);
    } catch (err) {
      setError(
        err instanceof Error ? err.message : 'Failed to process request',
      );
    }
  };

  const handleNext = (
    values: SubscriberDetailsType,
    setValues: (values: SubscriberDetailsType) => void,
  ) => {
    setValues(values);
    setActiveStep((prev) => prev + 1);
  };

  const handleBack = () => {
    setActiveStep((prev) => prev - 1);
  };
  const NoItemMessage = ({ message }: { message: string }) => (
    <MenuItem
      disabled
      value={''}
      sx={{
        m: 0,
        p: '6px 16px',
      }}
    >
      <Typography variant="body1">{message}</Typography>
    </MenuItem>
  );

  const renderStepContent = (formik: any) => {
    if (isLoading) {
      return (
        <DialogContent>
          <Box
            display="flex"
            flexDirection="column"
            alignItems="center"
            gap={2}
            py={4}
          >
            <CircularProgress />
            <Typography>Processing your request...</Typography>
          </Box>
        </DialogContent>
      );
    }

    if (error) {
      return (
        <>
          <DialogTitle>Error</DialogTitle>
          <DialogContent>
            <Typography color="error">{error}</Typography>
          </DialogContent>
          <DialogActions>
            <Button onClick={() => setError(null)}>Try Again</Button>
            <Button onClick={handleCloseAction}>Close</Button>
          </DialogActions>
        </>
      );
    }
    switch (activeStep) {
      case 0:
        return (
          <>
            <DialogTitle sx={{ color: colors.black }}>
              Add Subscriber
              <CloseButtonStyle aria-label="close" onClick={handleClose}>
                <CloseIcon />
              </CloseButtonStyle>
            </DialogTitle>
            <DialogContent>
              <Typography
                variant="subtitle1"
                color="text.secondary"
                sx={{ mb: 3 }}
              >
                Enter basic information about the subscriber, so that they can
                be authorized to use the network.{' '}
              </Typography>
              <Stack direction="column" spacing={2}>
                <Field name="name">
                  {({ field, meta }: any) => (
                    <TextField
                      {...field}
                      required
                      fullWidth
                      label="Name"
                      sx={{
                        height: '48px',
                        '& .MuiInputBase-root': {
                          height: '100%',
                        },
                      }}
                      error={meta.touched && Boolean(meta.error)}
                      helperText={meta.touched && meta.error}
                      slotProps={{
                        inputLabel: {
                          shrink: true,
                        },
                      }}
                    />
                  )}
                </Field>
                <Field name="email">
                  {({ field, meta }: any) => (
                    <TextField
                      {...field}
                      required
                      fullWidth
                      label="Email"
                      sx={{
                        height: '48px',
                        '& .MuiInputBase-root': {
                          height: '100%',
                        },
                      }}
                      error={meta.touched && Boolean(meta.error)}
                      helperText={meta.touched && meta.error}
                      slotProps={{
                        inputLabel: {
                          shrink: true,
                        },
                      }}
                    />
                  )}
                </Field>
              </Stack>
            </DialogContent>
            <DialogActions>
              <Button onClick={handleClose}>Cancel</Button>
              <Button
                onClick={() => handleNext(formik.values, formik.setValues)}
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

      case 1:
        return (
          <>
            <DialogTitle>
              Enter SIM ICCID
              <CloseButtonStyle aria-label="close" onClick={handleClose}>
                <CloseIcon />
              </CloseButtonStyle>
            </DialogTitle>
            <DialogContent>
              <Typography
                variant="subtitle1"
                color="text.secondary"
                sx={{ mb: 3 }}
              >
                Please select the SIM you would like to assign to the
                subscriber, and enter the ICCID found on the back of the card.
                Please ensure the ICCID is correct, because this cannot be
                undone.{' '}
              </Typography>
              <Field name="simIccid">
                {({ field, form }: any) => (
                  <Autocomplete
                    freeSolo
                    value={field.value}
                    options={sims.map((option) => option.iccid)}
                    sx={{
                      '.MuiAutocomplete-inputRoot .MuiAutocomplete-input': {
                        padding: '13px !important',
                      },
                    }}
                    onInputChange={(_, newValue) => {
                      form.setFieldValue('simIccid', newValue || undefined);
                      setSelectedSim(
                        sims.find((sim) => sim.iccid === newValue) || null,
                      );
                    }}
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        {...field}
                        label="SIM ICCID*"
                        InputLabelProps={{ shrink: true }}
                        sx={{
                          '.MuiOutlinedInput-root': { padding: 0 },
                        }}
                      />
                    )}
                    fullWidth
                  />
                )}
              </Field>
            </DialogContent>
            <DialogActions
              sx={{ display: 'flex', justifyContent: 'space-between' }}
            >
              <Button sx={{ p: 0 }} onClick={handleBack}>
                Back
              </Button>

              <Box sx={{ display: 'flex', gap: 1 }}>
                <Button onClick={handleClose}>Cancel</Button>
                <Button
                  onClick={() => handleNext(formik.values, formik.setValues)}
                  disabled={!formik.isValid || !formik.values.simIccid}
                  variant="contained"
                >
                  Next
                </Button>
              </Box>
            </DialogActions>
          </>
        );
      case 2: // Plan Selection
        return (
          <>
            <DialogTitle>
              Select data plan
              <CloseButtonStyle aria-label="close" onClick={handleClose}>
                <CloseIcon />
              </CloseButtonStyle>
            </DialogTitle>
            <DialogContent>
              <Typography
                variant="subtitle1"
                color="text.secondary"
                sx={{ mb: 2 }}
              >
                Select the purchased data plan
              </Typography>
              <Field name="plan" id="add-subscriber-plan-select">
                {({ field, meta }: any) => (
                  <FormControl
                    fullWidth
                    error={meta.touched && Boolean(meta.error)}
                  >
                    <InputLabel htmlFor="outlined-plan" shrink>
                      DATA PLAN
                    </InputLabel>
                    <SelectStyle
                      {...field}
                      label="DATA PLAN"
                      input={
                        <OutlinedInput
                          notched
                          label="DATA PLAN"
                          id="outlined-plan"
                        />
                      }
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
            <DialogActions
              sx={{ display: 'flex', justifyContent: 'space-between' }}
            >
              <Button sx={{ p: 0 }} onClick={handleBack}>
                Back
              </Button>
              <Box sx={{ display: 'flex', gap: 1 }}>
                <Button onClick={handleClose}>Cancel</Button>
                <Button
                  onClick={() => formik.handleSubmit()}
                  disabled={!formik.isValid || !formik.values.plan}
                  variant="contained"
                >
                  ADD SUBSCRIBER
                </Button>
              </Box>
            </DialogActions>
          </>
        );
      case 3:
        return (
          <>
            <DialogTitle>
              Successfully added [{formik.values.name}]
              <CloseButtonStyle aria-label="close" onClick={handleClose}>
                <CloseIcon />
              </CloseButtonStyle>
            </DialogTitle>
            <DialogContent>
              {submissionData?.is_physical ? (
                <>
                  <Typography
                    variant="subtitle1"
                    color="text.secondary"
                    sx={{ mb: 3 }}
                  >
                    You have successfully added {formik.values.name} as a
                    subscriber to your network, and a unique ID has been
                    generated for them, which must be used to create a Ukama
                    subscriber app.
                  </Typography>
                  <Box
                    sx={{
                      bgcolor: 'grey.50',
                      p: 2,
                      borderRadius: 1,
                      mb: 3,
                    }}
                  >
                    <Typography fontFamily="monospace" sx={{ mb: 1 }}>
                      UID: {submissionData.subscriber_id}
                    </Typography>
                    <Typography fontFamily="monospace">
                      SIM ICCID: {submissionData.iccid}
                    </Typography>
                  </Box>
                </>
              ) : (
                <>
                  <Typography
                    variant="subtitle1"
                    color="text.secondary"
                    sx={{ mb: 3 }}
                  >
                    You have successfully added {formik.values.name} as a
                    subscriber to your network, and an eSIM installation
                    invitation has been sent out to them. If they would rather
                    install their eSIM now, have them scan the QR code below.
                  </Typography>
                  <Accordion
                    sx={{ boxShadow: 'none', background: 'transparent' }}
                    onChange={(_, isExpanded: boolean) => {
                      setShowQrCode(isExpanded);
                    }}
                  >
                    <AccordionSummary
                      expandIcon={<ExpandMoreIcon color="primary" />}
                      sx={{
                        p: 0,
                        m: 0,
                        justifyContent: 'flex-start',
                        '& .MuiAccordionSummary-content': {
                          flexGrow: 0.02,
                        },
                      }}
                    >
                      <Typography
                        fontWeight={500}
                        variant="caption"
                        color={colors.primaryMain}
                      >
                        {showQrCode ? 'HIDE QR CODE' : 'SHOW QR CODE'}
                      </Typography>
                    </AccordionSummary>
                    <AccordionDetails
                      sx={{ p: 2, display: 'flex', justifyContent: 'center' }}
                    >
                      <QRCode
                        id="qrCodeId"
                        value={selectedSim?.qrCode || ''}
                        style={{ height: 180, width: 180 }}
                      />
                    </AccordionDetails>
                  </Accordion>
                </>
              )}
            </DialogContent>
            <DialogActions>
              <Button onClick={handleClose} variant="contained">
                Close
              </Button>
            </DialogActions>
          </>
        );

      default:
        return null;
    }
  };

  return (
    <Dialog open={isOpen} onClose={handleClose} maxWidth="sm" fullWidth>
      <Box sx={{ py: 1 }}>
        <Formik
          initialValues={initialValues}
          validationSchema={
            activeStep === 0 ? stepZeroSchema : subscriberDetailsSchema
          }
          onSubmit={handleSubmit}
          validateOnMount
        >
          {(formik) => <Form>{renderStepContent(formik)}</Form>}
        </Formik>
      </Box>
    </Dialog>
  );
};

export default AddSubscriberStepperDialog;
