import React, { useState, useEffect } from 'react';
import {
  Dialog,
  DialogContent,
  DialogTitle,
  List,
  ListItem,
  ListItemText,
  Box,
  Typography,
} from '@mui/material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import colors from '@/theme/colors';

interface OnboardingCardProps {
  open: boolean;
  onClose: () => void;
  onStepClick: (step: number) => void;
  status: boolean[];
}

const OnboardingCard: React.FC<OnboardingCardProps> = ({
  open,
  onClose,
  onStepClick,
  status,
}) => {
  const [activeStep, setActiveStep] = useState(0);

  useEffect(() => {
    const nextIncompleteStep = status.findIndex((stepStatus) => !stepStatus);
    setActiveStep(
      nextIncompleteStep === -1 ? status.length : nextIncompleteStep,
    );
  }, [status]);

  const handleStepClick = (step: number) => {
    if (step === activeStep) {
      onStepClick(step);
    }
  };

  const steps = [
    {
      step: 1,
      title: 'Install and configure your site(s)',
      description:
        'Install sites in either test or real locations, then complete digital configuration.',
    },
    {
      step: 2,
      title: 'Create data plan(s)',
      description:
        'Create custom data plans to define how your subscribers use the network.',
    },
    {
      step: 3,
      title: 'Add subscriber(s)',
      description:
        'Digitally assign a SIM and data plan to your user, so that they can start using the network.',
    },
  ];

  return (
    <Dialog
      open={open}
      onClose={onClose}
      sx={{
        '& .MuiDialog-paper': {
          position: 'fixed',
          bottom: 70,
          left: 30,
          margin: 0,
          borderRadius: '10px',
          maxWidth: 400,
          boxShadow: `0px 4px 12px ${colors.nightGrey12}`,
          border: `1px solid  ${colors.gray}`,
          overflow: 'hidden',
        },
      }}
    >
      <div
        style={{
          height: '15px',
          width: '100%',
          background: 'linear-gradient(to right, #00f, #00c6ff)',
        }}
      />
      <DialogTitle sx={{ ml: 2 }}>Setup and use your network</DialogTitle>
      <DialogContent>
        <List>
          {steps.map(({ step, title, description }, index) => (
            <ListItem
              key={step}
              onClick={() => handleStepClick(index)}
              sx={{
                cursor: index === activeStep ? 'pointer' : 'default',
                opacity: index === activeStep ? 1 : 0.5,
                '&:hover': {
                  backgroundColor:
                    index === activeStep ? `${colors.gray}` : 'transparent',
                },
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
                {status[index] ? (
                  <CheckCircleIcon color="success" />
                ) : (
                  <Box
                    sx={{
                      display: 'flex',
                      alignItems: 'center',
                      justifyContent: 'center',
                      width: 24,
                      height: 24,
                      border: '1px solid #ccc',
                      borderRadius: '50%',
                    }}
                  >
                    <Typography variant="body2">{step}</Typography>
                  </Box>
                )}
              </Box>
              <ListItemText
                primary={
                  <Typography variant="subtitle1" fontWeight="bold">
                    {title}
                  </Typography>
                }
                secondary={description}
              />
            </ListItem>
          ))}
        </List>
      </DialogContent>
    </Dialog>
  );
};

export default OnboardingCard;
