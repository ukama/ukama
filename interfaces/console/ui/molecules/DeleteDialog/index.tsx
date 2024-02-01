import React, { useState } from 'react';
import { colors } from '@/styles/theme';

import {
  Dialog,
  Button,
  Typography,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  IconButton,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';



interface DeleteConfirmationProps {
  onDelete: () => void;
  onCancel: () => void;
  open: boolean;
  itemName: string;
  loading: boolean;
}

const DeleteConfirmation: React.FC<DeleteConfirmationProps> = ({
  onDelete,
  onCancel,
  open,
  itemName,
  loading = false,
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

  return (
  <Dialog
      open={open}
      onClose={handleClose}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
    >
      <DialogTitle id="alert-dialog-title">Delete Confirmation</DialogTitle>
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
        <Typography variant="body1" sx={{color:colors.black}}>
        Are you certain you wish to delete the following subscriber -
          <span style={{ color: "black" }}>{itemName}</span> ? 
        This action will also remove all SIMs associated with them from your network.
        </Typography>

          
      </DialogContentText>
    </DialogContent><DialogActions>
        <Button onClick={handleClose} color="primary" autoFocus size="medium">
          Cancel
        </Button>
        <Button
          variant="contained"
          onClick={handleDelete}
          sx={{ background: 'red' }}
          disabled={loading}
        >
          Delete
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default DeleteConfirmation;
