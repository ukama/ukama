/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { colors } from '@/theme';
import React, { useState } from 'react';

import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
  Typography,
} from '@mui/material';

interface DeleteConfirmationProps {
  onDelete: () => void;
  onCancel: () => void;
  open: boolean;
  itemName: string;
  itemType?: 'subscriber' | 'sim';
  loading?: boolean;
  title?: string;
  description?: string;
}

const DeleteConfirmation: React.FC<DeleteConfirmationProps> = ({
  onDelete,
  onCancel,
  open,
  itemName,
  itemType = 'subscriber',
  loading = false,
  title,
  description,
}) => {
  const [isDeleting, setIsDeleting] = useState(false);

  const handleDelete = () => {
    setIsDeleting(true);
    onDelete();
  };

  const handleClose = () => {
    if (!isDeleting) {
      onCancel();
    }
  };

  const dialogTitle =
    title ||
    `Delete ${itemType === 'subscriber' ? 'Subscriber' : 'SIM'} Confirmation`;

  let dialogContent;

  if (itemType === 'subscriber') {
    dialogContent = (
      <Typography variant="body1" sx={{ color: colors.black70 }}>
        Are you certain you wish to delete the following subscriber -{' '}
        <span style={{ fontWeight: 'bold' }}>{itemName}</span>? This action will
        also remove all SIMs associated with them from your network.
      </Typography>
    );
  } else {
    dialogContent = (
      <>
        <Typography variant="body1" sx={{ color: colors.black70, mb: 2 }}>
          Are you sure you want to delete the SIM <strong>{itemName}</strong>?
          This will permanently remove all associated packages and usage data.
        </Typography>
      </>
    );
  }

  if (description) {
    dialogContent = (
      <Typography variant="body1" sx={{ color: colors.black }}>
        {description}
      </Typography>
    );
  }

  return (
    <Dialog
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">{dialogTitle}</DialogTitle>
      <IconButton
        aria-label="close"
        onClick={handleClose}
        sx={{
          position: 'absolute',
          right: 8,
          top: 8,
        }}
      >
        <CloseIcon />
      </IconButton>
      <DialogContent>
        <DialogContentText id="alert-dialog-description">
          {dialogContent}
        </DialogContentText>
      </DialogContent>
      <DialogActions>
        <Button onClick={handleClose} color="primary" autoFocus size="medium">
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleDelete}
          sx={{ background: colors.error }}
          disabled={loading}
        >
          Delete
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DeleteConfirmation;
