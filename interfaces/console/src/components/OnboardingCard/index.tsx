import React, { useState } from 'react';
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
}

const OnboardingCard: React.FC<OnboardingCardProps> = ({ open, onClose }) => {
  const [completedSteps, setCompletedSteps] = useState<number[]>([]);

  const handleStepClick = (step: number) => {
    if (!completedSteps.includes(step)) {
      setCompletedSteps([...completedSteps, step]);
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
          {steps.map(({ step, title, description }) => (
            <ListItem
              key={step}
              onClick={() => handleStepClick(step)}
              sx={{
                cursor: 'pointer',
                '&:hover': {
                  backgroundColor: `${colors.gray}`,
                },
              }}
            >
              <Box sx={{ display: 'flex', alignItems: 'center', mr: 2 }}>
                {completedSteps.includes(step) ? (
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
