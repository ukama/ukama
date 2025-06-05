/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useEffect, useState } from 'react';
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
  isEditing: boolean;
  setIsEditing: React.Dispatch<React.SetStateAction<boolean>>;
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
  isEditing,
  setIsEditing,
}) => {
  const [editingField, setEditingField] = useState<
    'firstName' | 'email' | null
  >(null);
  const [inputValues, setInputValues] = useState({
    firstName: subscriber.firstName,
    email: subscriber.email,
  });

  useEffect(() => {
    setInputValues({
      firstName: subscriber.firstName,
      email: subscriber.email,
    });
  }, [subscriber]);

  useEffect(() => {
    if (!isEditing) {
      setEditingField(null);
    }
  }, [isEditing]);

  const startEditing = (field: 'firstName' | 'email') => {
    setEditingField(field);
    setIsEditing(true);
  };

  const handleFieldChange = (field: 'firstName' | 'email', value: string) => {
    setInputValues((prev) => ({
      ...prev,
      [field]: value,
    }));

    const updateKey = field === 'firstName' ? 'name' : 'email';

    if (value !== subscriber[field]) {
      const updates = { [updateKey]: value };
      onUpdateSubscriber(updates);
    } else {
      onUpdateSubscriber({ [updateKey]: undefined });
    }
  };

  const finishEditingField = () => {
    setEditingField(null);
  };

  const handleKeyDown = (
    e: React.KeyboardEvent,
    field: 'firstName' | 'email',
  ) => {
    if (e.key === 'Enter') {
      finishEditingField();
    } else if (e.key === 'Escape') {
      setInputValues((prev) => ({
        ...prev,
        [field]: subscriber[field],
      }));

      const updateKey = field === 'firstName' ? 'name' : 'email';
      onUpdateSubscriber({ [updateKey]: undefined });

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
            onChange={(e) => handleFieldChange('firstName', e.target.value)}
            onBlur={finishEditingField}
            onKeyDown={(e) => handleKeyDown(e, 'firstName')}
            sx={{ mb: 2 }}
          />
        ) : (
          <FieldValue onClick={() => startEditing('firstName')}>
            {inputValues.firstName}
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
            onChange={(e) => handleFieldChange('email', e.target.value)}
            onBlur={finishEditingField}
            onKeyDown={(e) => handleKeyDown(e, 'email')}
            sx={{ mb: 2 }}
          />
        ) : (
          <FieldValue onClick={() => startEditing('email')}>
            {inputValues.email}
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
