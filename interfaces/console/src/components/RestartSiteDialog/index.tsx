import React, { useState } from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  Button,
  TextField,
  Typography,
  Box,
  IconButton,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';

interface RestartSiteDialogProps {
  open: boolean;
  onClose: () => void;
  onConfirm: (siteName: string) => void;
  siteName: string;
}

const RestartSiteDialog: React.FC<RestartSiteDialogProps> = ({
  open,
  onClose,
  onConfirm,
  siteName,
}) => {
  const [inputValue, setInputValue] = useState('');

  const handleConfirm = () => {
    if (inputValue === siteName) {
      onConfirm(siteName);
      setInputValue('');
    }
  };

  return (
    <Dialog open={open} onClose={onClose}>
      <DialogTitle>
        Restart site
        <IconButton
          aria-label="close"
          onClick={onClose}
          style={{ position: 'absolute', right: '8px', top: '8px' }}
        >
          <CloseIcon />
        </IconButton>
      </DialogTitle>
      <DialogContent>
        <Typography>
          Restarting site will cause it to be down for 10 minutes. Please type
          in the site name to confirm restarting the site.
        </Typography>
        <Box mt={2}>
          <TextField
            fullWidth
            label="Site Name"
            value={inputValue}
            onChange={(e) => setInputValue(e.target.value)}
          />
        </Box>
      </DialogContent>
      <DialogActions>
        <Button onClick={onClose}>Cancel</Button>
        <Button
          onClick={handleConfirm}
          color="primary"
          variant="contained"
          disabled={inputValue !== siteName}
        >
          Restart
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default RestartSiteDialog;
