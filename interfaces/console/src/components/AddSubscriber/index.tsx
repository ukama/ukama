/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { PackageDto, SimDto } from '@/client/graphql/generated';
import { TAddSubscriberData } from '@/types';
import CloseIcon from '@mui/icons-material/Close';
import {
  Backdrop,
  Box,
  Button,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  Typography,
} from '@mui/material';
import React, { useEffect, useState } from 'react';
import Step1 from './Step1';
import Step2 from './Step2';
import Step4 from './Step4';

interface SubscriberDialogProps {
  open: boolean;
  onClose: () => void;
  submitButtonState: boolean;
  pkgList: PackageDto[];
  loading: boolean;
  sims: SimDto[];
  pSimCount: number | undefined;
  eSimCount: number | undefined;
  onSubmit: Function;
  onSuccess: boolean;
  qrCode: string;
}

const INIT_ADD_SUBSCRIBER = {
  name: '',
  email: '',
  phone: '',
  simType: 'eSim',
  roamingStatus: false,
  iccid: '',
  plan: '',
};

const AddSubscriberDialog: React.FC<SubscriberDialogProps> = ({
  open,
  onClose,
  onSubmit,
  qrCode,
  pkgList,
  loading = false,
  sims,
  pSimCount,
  eSimCount,
  onSuccess = false,
}) => {
  const [activeStep, setActiveStep] = useState(1);
  const [formData, setFormData] =
    useState<TAddSubscriberData>(INIT_ADD_SUBSCRIBER);

  const handleStep1Submit = async (values: TAddSubscriberData) => {
    setActiveStep((prevStep) => prevStep + 1);
    setFormData(values);
  };

  const handleDialogClose = () => {
    setActiveStep(1);
    setFormData(INIT_ADD_SUBSCRIBER);
    onSuccess = false;
    onClose();
  };
  const getSubscriberForm = (step: number) => {
    const commonProps = {
      onClose: () => {
        handleDialogClose();
      },
      handleStep1Submit,
      pSimCount,
      eSimCount,
      formData,
    };
    switch (step) {
      case 1:
        return <Step1 {...commonProps} />;
      case 2:
        return (
          <Step2
            sims={
              formData.simType === 'eSim'
                ? sims.filter((sim) => sim.isPhysical === 'false')
                : sims.filter((sim) => sim.isPhysical === 'true')
            }
            {...commonProps}
            packages={pkgList}
            formData={formData}
            setFormData={setFormData}
            handleSubmitButton={() => {
              onSubmit(formData);
              handleDialogClose();
            }}
            goBack={() => setActiveStep((prevStep) => prevStep - 1)}
          />
        );
      case 3:
        return <Step4 qrCode={qrCode} simType={formData.simType} />;
      default:
        return <Step1 {...commonProps} />;
    }
  };

  useEffect(() => {
    return () => {
      setActiveStep(1);
    };
  }, []);
  const getSubTitle = (step: number) => {
    switch (step) {
      case 1:
        return 'Add subscribers to your network.';
      case 2:
        return 'Enter the ICCID for the SIM you have assigned to the subscriber, and select their data plan. Please ensure the ICCID is correct, because it cannot be undone once assigned.';
      case 3:
        return formData.simType == 'eSim '
          ? `You have successfully added ${name} as a subscriber to your network, and an ${formData.simType} installation invitation has been sent out to them. If they would rather install their eSIM now, have them scan the QR code below.`
          : `You have successfully added ${name} as a subscriber to your network, and ${formData.simType} installation instructions have been sent out to them. `;
      default:
        return 'Add subscribers to your network.';
    }
  };

  return (
    <Dialog open={open} onClose={handleDialogClose} fullWidth maxWidth="sm">
      <DialogTitle>
        {onSuccess
          ? `Successfully added ${name}`
          : activeStep > 1
            ? `Add subscriber ${name}`
            : 'Add Subscriber'}
      </DialogTitle>
      <IconButton
        aria-label="close"
        onClick={handleDialogClose}
        sx={{
          position: 'absolute',
          right: 8,
          top: 8,
        }}
      >
        <CloseIcon />
      </IconButton>
      <DialogContent>{getSubTitle(onSuccess ? 3 : activeStep)}</DialogContent>
      <Backdrop
        sx={{ color: '#fff', zIndex: (theme) => theme.zIndex.drawer + 1 }}
        open={loading}
      >
        <Stack direction="row" spacing={2} alignItems="center">
          <CircularProgress color="inherit" />
          <Typography variant="body1" color="initial">
            Creating a subscriber...
          </Typography>
        </Stack>
      </Backdrop>
      <Box sx={{ px: 3, pt: 1, pb: 2 }}>
        {getSubscriberForm(onSuccess ? 3 : activeStep)}
      </Box>
      {onSuccess && (
        <DialogActions sx={{ justifyContent: 'flex-end' }}>
          <Button variant="contained" onClick={handleDialogClose}>
            Close
          </Button>
        </DialogActions>
      )}
    </Dialog>
  );
};

export default AddSubscriberDialog;
