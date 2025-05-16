import React from 'react';
import { Typography, Box, TextField, IconButton, styled } from '@mui/material';
import EditIcon from '@mui/icons-material/Edit';
import { colors } from '@/theme';

interface Subscriber {
  id: string;
  firstName: string;
  email: string;
}

interface SubscriberInfoTabProps {
  subscriber: Subscriber;
  onUpdateSubscriber: (updates: { name?: string; email?: string }) => void;
}

const FieldLabel = styled(Typography)(({ theme }) => ({
  color: colors.black38,
  fontSize: theme.typography.caption.fontSize,
  lineHeight: theme.typography.caption.lineHeight,
  marginBottom: theme.spacing(0.5),
}));

const FieldValue = styled(Typography)(({ theme }) => ({
  fontSize: theme.typography.body1.fontSize,
  marginBottom: theme.spacing(2),
}));

const FieldContainer = styled(Box)(() => ({
  position: 'relative',
  paddingRight: '32px',
  marginBottom: 16,
}));

const EditButton = styled(IconButton)(() => ({
  position: 'absolute',
  right: 0,
  top: '50%',
  transform: 'translateY(-50%)',
}));

const SubscriberInfoTab: React.FC<SubscriberInfoTabProps> = ({
  subscriber,
  onUpdateSubscriber,
}) => {
  const [editingField, setEditingField] = React.useState<
    'firstName' | 'email' | null
  >(null);
  const [inputValues, setInputValues] = React.useState({
    firstName: subscriber.firstName,
    email: subscriber.email,
  });

  React.useEffect(() => {
    setInputValues({
      firstName: subscriber.firstName,
      email: subscriber.email,
    });
  }, [subscriber]);

  const startEditing = (field: 'firstName' | 'email') => {
    setEditingField(field);
  };

  const saveChanges = (field: 'firstName' | 'email') => {
    if (inputValues[field] !== subscriber[field]) {
      onUpdateSubscriber({
        [field === 'firstName' ? 'name' : 'email']: inputValues[field],
      });
    }
    setEditingField(null);
  };

  const handleKeyDown = (
    e: React.KeyboardEvent,
    field: 'firstName' | 'email',
  ) => {
    if (e.key === 'Enter') {
      saveChanges(field);
    } else if (e.key === 'Escape') {
      setInputValues((prev) => ({
        ...prev,
        [field]: subscriber[field],
      }));
      setEditingField(null);
    }
  };

  return (
    <>
      <FieldContainer>
        <FieldLabel>FIRST NAME</FieldLabel>
        {editingField === 'firstName' ? (
          <TextField
            autoFocus
            fullWidth
            value={inputValues.firstName}
            onChange={(e) =>
              setInputValues((prev) => ({
                ...prev,
                firstName: e.target.value,
              }))
            }
            onBlur={() => saveChanges('firstName')}
            onKeyDown={(e) => handleKeyDown(e, 'firstName')}
            sx={{ mb: 2 }}
          />
        ) : (
          <FieldValue onClick={() => startEditing('firstName')}>
            {subscriber.firstName}
          </FieldValue>
        )}
        <EditButton onClick={() => startEditing('firstName')}>
          <EditIcon fontSize="small" />
        </EditButton>
      </FieldContainer>

      <FieldContainer>
        <FieldLabel>EMAIL</FieldLabel>
        {editingField === 'email' ? (
          <TextField
            autoFocus
            fullWidth
            value={inputValues.email}
            onChange={(e) =>
              setInputValues((prev) => ({
                ...prev,
                email: e.target.value,
              }))
            }
            onBlur={() => saveChanges('email')}
            onKeyDown={(e) => handleKeyDown(e, 'email')}
            sx={{ mb: 2 }}
          />
        ) : (
          <FieldValue onClick={() => startEditing('email')}>
            {subscriber.email}
          </FieldValue>
        )}
        <EditButton onClick={() => startEditing('email')}>
          <EditIcon fontSize="small" />
        </EditButton>
      </FieldContainer>
    </>
  );
};

export default SubscriberInfoTab;
