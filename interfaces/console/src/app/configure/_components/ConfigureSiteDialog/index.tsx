/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { NetworkDto } from '@/client/graphql/generated';
import { useAppContext } from '@/context';
import { TSiteForm } from '@/types';
import { useFetchAddress } from '@/utils/useFetchAddress';
import CloseIcon from '@mui/icons-material/Close';
import {
  Dialog,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  Step,
  StepLabel,
  Stepper,
} from '@mui/material';
import React, { useEffect, useState } from 'react';
import ComponentSelectionStep, { Component } from './ComponentSelectionStep';
import SiteDetailsStep from './SiteDetailsStep';

interface IConfigureSiteDialog {
  open: boolean;
  site: TSiteForm;
  onClose: () => void;
  addSiteLoading: boolean;
  networks: NetworkDto[];
  components: Component[];
  handleSiteConfiguration: (values: TSiteForm) => void;
}

const STEPS = [
  'Select Switch, Power, Backhaul, and spectrum band',
  'Enter your site details',
];

const EMPTY_FORM: TSiteForm = {
  switch: '',
  power: '',
  backhaul: '',
  access: '',
  spectrum: '',
  siteName: '',
  network: '',
  latitude: 0,
  longitude: 0,
  address: '',
};

/**
 * Two-step dialog for configuring a new site installation.
 * Step orchestration lives here; per-step form logic is in
 * ComponentSelectionStep and SiteDetailsStep respectively.
 */
const ConfigureSiteDialog: React.FC<IConfigureSiteDialog> = ({
  open,
  site,
  onClose,
  networks,
  components,
  addSiteLoading = false,
  handleSiteConfiguration,
}) => {
  const { setSnackbarMessage } = useAppContext();
  const [activeStep, setActiveStep] = useState(0);
  const [formValues, setFormValues] = useState<TSiteForm>(site);
  const [currentAddress, setCurrentAddress] = useState('');
  const { address, isLoading, error, fetchAddress } = useFetchAddress();

  // Reset whenever the dialog opens
  useEffect(() => {
    if (open) {
      setFormValues(EMPTY_FORM);
      setActiveStep(0);
      setCurrentAddress('');
    }
  }, [open]);

  useEffect(() => {
    if (address) setCurrentAddress(address);
  }, [address]);

  useEffect(() => {
    if (error) {
      setSnackbarMessage({
        id: 'error-fetching-address',
        type: 'error',
        show: true,
        message: 'Error fetching address from coordinates',
      });
    }
  }, [error, setSnackbarMessage]);

  const handleClose = () => {
    setFormValues(EMPTY_FORM);
    setActiveStep(0);
    setCurrentAddress('');
    onClose();
  };

  const handleStep0Complete = (values: Partial<TSiteForm>) => {
    setFormValues((prev) => ({ ...prev, ...values }));
    setActiveStep(1);
  };

  const handleStep1Complete = (values: TSiteForm) => {
    if (address === 'Location not found') {
      setSnackbarMessage({
        id: 'error-fetching-address',
        type: 'error',
        show: true,
        message: 'Error fetching address from coordinates',
      });
    }
    handleSiteConfiguration({ ...values, address });
    handleClose();
  };

  const handleFetchAddress = async (lat: number, lng: number) => {
    setSnackbarMessage({
      id: 'fetching-address',
      type: 'success',
      show: true,
      message: 'Fetching address with coordinates',
    });
    await fetchAddress(lat.toString(), lng.toString());
  };

  const switchComponents = components.filter((c) => c.category === 'SWITCH');
  const powerComponents = components.filter((c) => c.category === 'POWER');
  const backhaulComponents = components.filter((c) => c.category === 'BACKHAUL');
  const accessComponents = components.filter((c) => c.category === 'ACCESS');

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      sx={{ '& .MuiDialog-paper': { width: '60%', maxWidth: '40%' } }}
    >
      <DialogTitle>
        Configure site installation ({activeStep + 1}/2)
      </DialogTitle>
      <IconButton
        aria-label="close"
        onClick={handleClose}
        sx={{
          position: 'absolute',
          right: 8,
          top: 8,
          color: (theme) => theme.palette.grey[500],
        }}
      >
        <CloseIcon />
      </IconButton>

      <DialogContent>
        <DialogContentText>
          {`You have successfully installed your site, and need to configure
              it. Please note that if your power or backhaul choice is "other",
              it can't be monitored within Ukama's Console.`}
        </DialogContentText>
      </DialogContent>

      <DialogContent>
        <Stepper activeStep={activeStep} sx={{ mb: 4 }}>
          {STEPS.map((label) => (
            <Step key={label}>
              <StepLabel>{label}</StepLabel>
            </Step>
          ))}
        </Stepper>

        {activeStep === 0 && (
          <ComponentSelectionStep
            initialValues={formValues}
            switchComponents={switchComponents}
            powerComponents={powerComponents}
            backhaulComponents={backhaulComponents}
            accessComponents={accessComponents}
            onComplete={handleStep0Complete}
            onCancel={handleClose}
          />
        )}

        {activeStep === 1 && (
          <SiteDetailsStep
            initialValues={formValues}
            networks={networks}
            address={address}
            currentAddress={currentAddress}
            isAddressLoading={isLoading}
            addressError={error}
            addSiteLoading={addSiteLoading}
            onBack={() => setActiveStep(0)}
            onCancel={handleClose}
            onComplete={handleStep1Complete}
            onFetchAddress={handleFetchAddress}
            onFieldChange={(field, value) =>
              setFormValues((prev) => ({ ...prev, [field]: value }))
            }
          />
        )}
      </DialogContent>
    </Dialog>
  );
};

export default ConfigureSiteDialog;
