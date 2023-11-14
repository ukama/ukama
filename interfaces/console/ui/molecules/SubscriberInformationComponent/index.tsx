// SubscriberInformationComponent.tsx
import React from 'react';
import { Stack, Typography, TextField, IconButton } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';

interface SubscriberInformationProps {
  loading: boolean;
  subscriberLoading: boolean;
  firstName: string;
  email: string;
  onEditName: boolean;
  onEditEmail: boolean;
  handleEditName: (event: React.ChangeEvent<HTMLInputElement>) => void;
  handleSimEdit: (event: React.ChangeEvent<HTMLInputElement>) => void;
  setOnEditName: (value: boolean) => void;
  setOnEditEmail: (value: boolean) => void;
}

const SubscriberInformationComponent: React.FC<SubscriberInformationProps> = ({
  loading,
  subscriberLoading,
  firstName,
  email,
  onEditName,
  onEditEmail,
  handleEditName,
  handleSimEdit,
  setOnEditName,
  setOnEditEmail,
}) => {
  return (
    <Stack direction="column" spacing={2}>
      <Typography variant="body1">Name</Typography>
      {/* ... Other subscriber information */}
    </Stack>
  );
};

export default SubscriberInformationComponent;
