import { globalUseStyles } from '@/styles/global';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  TextField,
} from '@mui/material';
import React, { useState } from 'react';

type EditNodeProps = {
  title: string;
  isOpen: boolean;
  nodeName: string;
  isClosable?: boolean;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction?: any;
  labelNegativeBtn?: string;
};

const EditNode = ({
  title,
  isOpen,
  nodeName,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  isClosable = true,
  handleSuccessAction,
}: EditNodeProps) => {
  const gclasses = globalUseStyles();
  const [value, setValue] = useState(nodeName);
  return (
    <Dialog
      fullWidth
      open={isOpen}
      maxWidth="sm"
      onClose={handleCloseAction}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
      onBackdropClick={() => isClosable && handleCloseAction()}
    >
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <DialogTitle>{title}</DialogTitle>
        <IconButton
          onClick={handleCloseAction}
          sx={{ position: 'relative', right: 8 }}
        >
          <CloseIcon />
        </IconButton>
      </Stack>

      <DialogContent>
        <TextField
          required
          fullWidth
          value={value}
          label={'NODE NAME'}
          InputLabelProps={{ shrink: true }}
          InputProps={{
            classes: {
              input: gclasses.inputFieldStyle,
            },
          }}
          onChange={(e: any) => setValue(e.target.value)}
        />
      </DialogContent>

      <DialogActions>
        <Stack direction={'row'} alignItems="center" spacing={2}>
          {labelNegativeBtn && (
            <Button
              variant="text"
              color={'primary'}
              onClick={handleCloseAction}
            >
              {labelNegativeBtn}
            </Button>
          )}
          {labelSuccessBtn && (
            <Button
              variant="contained"
              disabled={!value}
              onClick={() => handleSuccessAction(value)}
            >
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default React.memo(EditNode);
