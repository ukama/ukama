import React, { useEffect, useState } from 'react';
import {
  Dialog,
  Box,
  DialogTitle,
  CircularProgress,
  DialogContent,
  DialogActions,
  IconButton,
  Button,
  Backdrop,
  Typography,
  Stack,
} from '@mui/material';
import Step1 from './Step1';
import Step2 from './Step2';
import Step4 from './Step4';
import CloseIcon from '@mui/icons-material/Close';
import { PackageDto, SimDto } from '@/generated';
import { getSubTitle } from '../../../utils';

interface SubscriberDialogProps {
  open: boolean;
  onClose: () => void;
  submitButtonState: boolean;
  pkgList: PackageDto[];
  loading: boolean;
  sims: SimDto[];
  pSimCount: number | undefined;
  eSimCount: number | undefined;
  handleRoamingInstallation: Function;
  onSuccess: boolean;
  qrCode: string;
}

const AddSubscriberDialog: React.FC<SubscriberDialogProps> = ({
  open,
  onClose,
  submitButtonState,
  qrCode,
  pkgList,
  loading = false,
  handleRoamingInstallation,
  sims,
  pSimCount,
  eSimCount,
  onSuccess = false,
}) => {
  const [activeStep, setActiveStep] = useState(1);
  const [selectedSimType, setSelectedSimType] = useState<string>('eSim');
  const [name, setName] = useState<string>('');
  const [formData, setFormData] = useState<any>(null);

  const handleSimInstallation = async (values: any) => {
    setActiveStep((prevStep) => prevStep + 1);
    setSelectedSimType(values.selectedSimType);
    setName(values.name);
    setFormData(values);
  };

  const handleDialogClose = () => {
    setActiveStep(1);
    setName('');
    onClose();
    onSuccess = false;
  };
  const getSubscriberForm = (step: number) => {
    const commonProps = {
      onClose: () => {
        setActiveStep(1);
      },
      handleSimInstallation,
      pSimCount,
      eSimCount,
    };
    switch (step) {
      case 1:
        return <Step1 {...commonProps} />;
      case 2:
        return (
          <Step2
            goBack={() => setActiveStep((prevStep) => prevStep - 1)}
            {...commonProps}
            handlePlanInstallation={(plan: string, simIccid: string) => {
              handleRoamingInstallation({ ...formData, plan, simIccid });
            }}
            submitButtonState={submitButtonState}
            packages={pkgList}
            sims={sims}
            selectedSimType={selectedSimType}
          />
        );
      case 3:
        return <Step4 qrCode={qrCode} simType={selectedSimType} />;
      default:
        return <Step1 {...commonProps} />;
    }
  };

  useEffect(() => {
    return () => {
      setActiveStep(1);
    };
  }, []);
  const subTitle = getSubTitle(
    onSuccess ? 3 : activeStep,
    selectedSimType,
    name,
  );

  return (
    <Dialog open={open} onClose={handleDialogClose} fullWidth maxWidth="sm">
      <DialogTitle>
        {onSuccess ? `Successfully added ${name}` : `Add subscriber ${name}`}
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
      <DialogContent>{subTitle}</DialogContent>
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
      <Box sx={{ p: 2, width: '100%' }}>
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
