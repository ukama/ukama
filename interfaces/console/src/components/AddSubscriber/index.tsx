import React, { useState } from 'react';
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
  name: Yup.string()
    .required('Name is required')
    .min(2, 'Name must be at least 2 characters')
    .max(50, 'Name must not exceed 50 characters'),
  simIccid: Yup.string()
    .required('SIM ICCID is required')
    .matches(/^\d{19,20}$/, 'Invalid ICCID format'),
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
  handleAddSubscriber: (subscriber: SubscriberDetailsType) => void;
  packages: PackageDto[];
  sims: SimDto[];
  data: AllocateSimApiDto;
}

const SubscriberForm: React.FC<SubscriberFormProps> = ({
  isOpen,
  handleCloseAction,
  handleAddSubscriber,
  packages,
  sims,
  data,
}) => {
  const [activeStep, setActiveStep] = useState(0);
  const [showQrCode, setShowQrCode] = useState(false);
  const gclasses = globalUseStyles();
  const classes = useStyles();

  const initialValues: SubscriberDetailsType = {
    name: '',
    simIccid: '',
    plan: '',
  };

  const handleClose = () => {
    setActiveStep(0);
    handleCloseAction();
  };

  const handleSubmit = (values: SubscriberDetailsType) => {
    handleAddSubscriber(values);
    setActiveStep(2);
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

  const renderStepContent = (formik: any) => {
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
                Authorize subscriber to use your network. Please ensure the
                ICCID is correct, because it cannot be undone once assigned.
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
                <Field name="simIccid">
                  {({ field, form, meta }: any) => (
                    <Autocomplete
                      options={sims}
                      getOptionLabel={(option) => option.iccid || ''}
                      value={sims.find((sim) => sim.id === field.value) || null}
                      onChange={(_, newValue) => {
                        form.setFieldValue('simIccid', newValue?.iccid || '');
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
                      fullWidth
                    />
                  )}
                </Field>
              </Stack>
            </DialogContent>
            <DialogActions>
              <Button onClick={handleClose}>Cancel</Button>
              <Button
                onClick={() => handleNext(formik.values, formik.setValues)}
                disabled={!formik.isValid || !formik.dirty}
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
              Add Subscriber: [{formik.values.name}]
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
                Select the purchased data plan.
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
                      sx={{
                        '& .MuiOutlinedInput-notchedOutline': {
                          textAlign: 'left',
                        },
                        '& .MuiInputLabel-outlined': {
                          transform: 'translate(14px, -6px) scale(0.75)',
                        },
                        '& .MuiInputLabel-outlined.MuiInputLabel-shrink': {
                          transform: 'translate(14px, -6px) scale(0.75)',
                        },
                        '& .MuiSelect-select': {
                          paddingTop: '16px',
                          paddingBottom: '16px',
                        },
                      }}
                      fullWidth
                      required
                      label="DATA PLAN"
                      labelId="outlined-plan-label"
                      MenuProps={{
                        disablePortal: false,
                        PaperProps: {
                          sx: {
                            boxShadow:
                              '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
                            borderRadius: '4px',
                          },
                        },
                      }}
                      className={classes.selectStyle}
                    >
                      {packages.map((plan) => (
                        <MenuItem key={plan.uuid} value={plan.uuid}>
                          <Typography variant="body1">
                            {`${plan.name} - ${plan.currency} ${plan.amount}/${plan.dataVolume} ${plan.dataUnit}`}
                          </Typography>
                        </MenuItem>
                      ))}
                    </Select>
                    {meta.touched && meta.error && (
                      <FormHelperText>{meta.error}</FormHelperText>
                    )}
                  </FormControl>
                )}
              </Field>
            </DialogContent>
            <DialogActions>
              <Box sx={{ flexGrow: 1 }}>
                <Button onClick={handleBack}>Back</Button>
              </Box>
              <Button onClick={handleClose}>Cancel</Button>
              <Button
                onClick={formik.handleSubmit}
                disabled={!formik.isValid || !formik.dirty}
                variant="contained"
              >
                Add Subscriber
              </Button>
            </DialogActions>
          </>
        );

      case 2:
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
              {data?.is_physical ? (
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
                      UID: {data.subscriber_id}
                    </Typography>
                    <Typography fontFamily="monospace">
                      SIM ICCID: {data.iccid}
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
                        value={data?.iccid || ''}
                        style={{ height: 164, width: 164 }}
                      />
                    </AccordionDetails>
                  </Accordion>
                </>
              )}
              <Box sx={{ display: 'flex', gap: 1, mb: 2 }}>
                <TextField
                  fullWidth
                  label="Email"
                  variant="outlined"
                  size="small"
                />
                <Button variant="contained">SHARE </Button>
              </Box>
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
            activeStep === 0 ? subscriberDetailsSchema : subscriberDetailsSchema
          }
          onSubmit={handleSubmit}
          validateOnMount
          // onSubmit={(values) => {
          //   handleAddSubscriber(values);
          //   setActiveStep(2);
          // }}
        >
          {(formik) => <Form>{renderStepContent(formik)}</Form>}
        </Formik>
      </Box>
    </Dialog>
  );
};

export default SubscriberForm;
