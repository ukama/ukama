import React, { useEffect, useState } from 'react';
import { Formik, Form, Field } from 'formik';
import * as Yup from 'yup';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  Typography,
  TextField,
  Autocomplete,
  Select,
  MenuItem,
  Box,
  Stack,
  OutlinedInput,
  FormControl,
  InputLabel,
  FormHelperText,
  IconButton,
  AccordionSummary,
  Accordion,
  AccordionDetails,
  CircularProgress,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import { globalUseStyles } from '@/styles/global';
import colors from '@/theme/colors';
import { makeStyles } from '@mui/styles';
import {
  PackageDto,
  SimDto,
  AllocateSimApiDto,
} from '@/client/graphql/generated';
import QRCode from 'qrcode.react';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { SubscriberDetailsType } from '@/types';

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
const useStyles = makeStyles(() => ({
  selectStyle: {
    width: '100%',
    height: '48px',
  },
  formControl: {
    width: '100%',
    height: '48px',
  },
  closeButton: {
    position: 'absolute',
    right: 8,
    top: 8,
  },
}));

interface SubscriberFormProps {
  isOpen: boolean;
  handleCloseAction: () => void;
  handleAddSubscriber: (
    subscriber: SubscriberDetailsType,
  ) => Promise<AllocateSimApiDto>;
  packages: PackageDto[];
  sims: SimDto[];
  isLoading: boolean;
}

const AddSubscriberStepperDialog: React.FC<SubscriberFormProps> = ({
  isOpen,
  handleCloseAction,
  handleAddSubscriber,
  packages,
  sims,
  isLoading,
}) => {
  const [activeStep, setActiveStep] = useState(0);
  const [showQrCode, setShowQrCode] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [submissionData, setSubmissionData] =
    useState<AllocateSimApiDto | null>(null);
  const [selectedSim, setSelectedSim] = useState<SimDto | null>(null);

  const gclasses = globalUseStyles();
  const classes = useStyles();

  const initialValues: SubscriberDetailsType = {
    name: '',
    email: '',
    simIccid: '',
    plan: '',
  };

  useEffect(() => {
    console.log('HELLO', selectedSim);
  }, [selectedSim]);

  const handleClose = () => {
    setActiveStep(0);
    handleCloseAction();
  };

  const handleSubmit = async (values: SubscriberDetailsType) => {
    try {
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
              <IconButton
                aria-label="close"
                onClick={handleClose}
                className={classes.closeButton}
              >
                <CloseIcon />
              </IconButton>
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
                      error={meta.touched && Boolean(meta.error)}
                      helperText={meta.touched && meta.error}
                      InputLabelProps={{ shrink: true }}
                      InputProps={{
                        classes: { input: gclasses.inputFieldStyle },
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
                      error={meta.touched && Boolean(meta.error)}
                      helperText={meta.touched && meta.error}
                      InputLabelProps={{ shrink: true }}
                      InputProps={{
                        classes: { input: gclasses.inputFieldStyle },
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
              <IconButton
                aria-label="close"
                onClick={handleClose}
                className={classes.closeButton}
              >
                <CloseIcon />
              </IconButton>
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
                {({ field, form, meta }: any) => (
                  <Autocomplete
                    options={sims}
                    getOptionLabel={(option) => option.iccid || ''}
                    value={sims.find((sim) => sim.id === field.value) || null}
                    onChange={(_, newValue) => {
                      form.setFieldValue('simIccid', newValue?.iccid || '');
                      setSelectedSim(newValue);
                    }}
                    renderInput={(params) => (
                      <TextField
                        {...params}
                        {...field}
                        label="SIM ICCID*"
                        error={meta.touched && Boolean(meta.error)}
                        helperText={meta.touched && meta.error}
                        InputLabelProps={{ shrink: true }}
                      />
                    )}
                    noOptionsText={
                      sims.length === 0
                        ? 'SIM pool is empty. Please upload SIMs to SIM pool first.'
                        : 'No matching SIMs found'
                    }
                    fullWidth
                  />
                )}
              </Field>
            </DialogContent>
            <DialogActions
              sx={{ display: 'flex', justifyContent: 'space-between' }}
            >
              <Box>
                <Button onClick={handleBack}>Back</Button>
              </Box>
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
              <IconButton
                aria-label="close"
                onClick={handleClose}
                className={classes.closeButton}
              >
                <CloseIcon />
              </IconButton>
            </DialogTitle>
            <DialogContent>
              <Typography
                variant="subtitle1"
                color="text.secondary"
                sx={{ mb: 2 }}
              >
                Select the purchased data plan
              </Typography>
              <Field name="plan">
                {({ field, meta }: any) => (
                  <FormControl
                    fullWidth
                    error={meta.touched && Boolean(meta.error)}
                  >
                    <InputLabel htmlFor="outlined-plan" shrink>
                      DATA PLAN
                    </InputLabel>
                    <Select
                      {...field}
                      input={
                        <OutlinedInput
                          notched
                          label="DATA PLAN"
                          id="outlined-plan"
                        />
                      }
                      className={classes.selectStyle}
                    >
                      {packages.length === 0 ? (
                        <NoItemMessage message="No packages available. Please add packages first." />
                      ) : (
                        packages.map((plan) => (
                          <MenuItem key={plan.uuid} value={plan.uuid}>
                            <Typography variant="body1">
                              {`${plan.name} - ${plan.currency} ${plan.amount}/${plan.dataVolume} ${plan.dataUnit}`}
                            </Typography>
                          </MenuItem>
                        ))
                      )}
                    </Select>
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
              <Box>
                <Button onClick={handleBack}>Back</Button>
              </Box>
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
              <IconButton
                aria-label="close"
                onClick={handleClose}
                className={classes.closeButton}
              >
                <CloseIcon />
              </IconButton>
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
