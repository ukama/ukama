/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import {
  AllocateSimApiDto,
  PackageDto,
  SimPoolResDto,
} from '@/client/graphql/generated';
import { useUIContext } from '@/context';
import { SubscriberDetailsType } from '@/types';
import {
  Box,
  CircularProgress,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  Button,
  Typography,
} from '@mui/material';
import { Form, Formik } from 'formik';
import React, { useState } from 'react';
import * as Yup from 'yup';
import PlanSelectionStep from './PlanSelectionStep';
import SimSelectionStep from './SimSelectionStep';
import SubscriberInfoStep from './SubscriberInfoStep';
import SuccessStep from './SuccessStep';

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

const INITIAL_VALUES: SubscriberDetailsType = {
  name: '',
  email: '',
  simIccid: '',
  plan: '',
};

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

/**
 * Multi-step dialog for adding a new subscriber.
 * Step orchestration and shared Formik context live here; each step's UI is
 * in its own focused component (SubscriberInfoStep, SimSelectionStep,
 * PlanSelectionStep, SuccessStep).
 */
const AddSubscriberStepperDialog: React.FC<SubscriberFormProps> = ({
  isOpen,
  handleCloseAction,
  handleAddSubscriber,
  packages,
  sims,
  isLoading,
  currencySymbol,
}) => {
  const { setSnackbarMessage } = useUIContext();
  const [activeStep, setActiveStep] = useState(0);
  const [error, setError] = useState<string | null>(null);
  const [submissionData, setSubmissionData] = useState<AllocateSimApiDto | null>(null);
  const [selectedSim, setSelectedSim] = useState<SimPoolResDto | null>(null);

  const handleClose = () => {
    setActiveStep(0);
    setError(null);
    handleCloseAction();
  };

  const handleSubmit = async (values: SubscriberDetailsType) => {
    if (!values.plan) {
      setSnackbarMessage({ show: true, type: 'error', id: 'plane-not-found', message: 'Please select a plan' });
      return;
    }
    if (!values.simIccid) {
      setSnackbarMessage({ show: true, type: 'error', id: 'sim-not-found', message: 'Please provide sim iccid.' });
      return;
    }
    try {
      setError(null);
      const response = await handleAddSubscriber(values);
      setSubmissionData(response);
      setActiveStep(3);
    } catch (err) {
      setError(err instanceof Error ? err.message : 'Failed to process request');
    }
  };

  return (
    <Dialog open={isOpen} onClose={handleClose} maxWidth="sm" fullWidth>
      <Box sx={{ py: 1 }}>
        <Formik
          initialValues={INITIAL_VALUES}
          validationSchema={activeStep === 0 ? stepZeroSchema : subscriberDetailsSchema}
          onSubmit={handleSubmit}
          validateOnMount
        >
          {(formik) => (
            <Form>
              {/* Global loading overlay */}
              {isLoading && (
                <DialogContent>
                  <Box display="flex" flexDirection="column" alignItems="center" gap={2} py={4}>
                    <CircularProgress />
                    <Typography>Processing your request...</Typography>
                  </Box>
                </DialogContent>
              )}

              {/* Error screen */}
              {!isLoading && error && (
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
              )}

              {/* Step screens */}
              {!isLoading && !error && activeStep === 0 && (
                <SubscriberInfoStep
                  formik={formik}
                  onNext={() => setActiveStep(1)}
                  onClose={handleClose}
                />
              )}
              {!isLoading && !error && activeStep === 1 && (
                <SimSelectionStep
                  formik={formik}
                  sims={sims}
                  onNext={() => setActiveStep(2)}
                  onBack={() => setActiveStep(0)}
                  onClose={handleClose}
                  onSimSelected={setSelectedSim}
                />
              )}
              {!isLoading && !error && activeStep === 2 && (
                <PlanSelectionStep
                  formik={formik}
                  packages={packages}
                  currencySymbol={currencySymbol}
                  onBack={() => setActiveStep(1)}
                  onClose={handleClose}
                  onSubmit={() => formik.handleSubmit()}
                />
              )}
              {!isLoading && !error && activeStep === 3 && submissionData && (
                <SuccessStep
                  subscriberName={formik.values.name}
                  submissionData={submissionData}
                  selectedSim={selectedSim}
                  onClose={handleClose}
                />
              )}
            </Form>
          )}
        </Formik>
      </Box>
    </Dialog>
  );
};

export default AddSubscriberStepperDialog;
