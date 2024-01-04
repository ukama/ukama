// ReusableDialog.tsx

import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
} from '@mui/material';

interface ReusableDialogProps {
  open: boolean;
  onClose: () => void;
  title: string;
  content: React.ReactNode;
  onConfirm: () => void;
  onCancel: () => void;
  buttonText: string;
}

const NodeActionDialog: React.FC<ReusableDialogProps> = ({
  open,
  onClose,
  title,
  content,
  onConfirm,
  buttonText,
  onCancel,
}) => {
  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>{title}</DialogTitle>
      <DialogContent>{content}</DialogContent>
      <DialogActions>
        <Button onClick={onCancel} color="primary">
          Cancel
        </Button>
        <Button onClick={onConfirm} color="primary" variant="contained">
          {buttonText}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default NodeActionDialog;
