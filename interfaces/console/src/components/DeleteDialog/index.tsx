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
  Box,
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
  isLastSim?: boolean;
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
  isLastSim = false, // Default to false
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

        <Typography variant="body1" sx={{ color: colors.black70 }}>
          {isLastSim &&
            'This is the last SIM for this subscriber. Deleting it will also remove the subscriber record.'}
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
