import colors from '@/styles/theme/colors';
import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogContentText,
  DialogTitle,
  FormControlLabel,
  IconButton,
  Typography,
} from '@mui/material';
import React from 'react';

export interface DialogTitleProps {
  id: string;
  children?: React.ReactNode;
  onClose: () => void;
}
type softwareUpdateModalProps = {
  isOpen: boolean;
  handleClose: any;
  title: string;
  content: string;
  submit: any;
  btnLabel?: string;
};
const BootstrapDialogTitle = (props: DialogTitleProps) => {
  const { children, onClose, ...other } = props;

  return (
    <DialogTitle sx={{ m: 0, p: 2 }} {...other}>
      {children}
      {onClose ? (
        <IconButton
          aria-label="close"
          onClick={onClose}
          sx={{
            position: 'absolute',
            right: 8,
            top: 8,
            color: (theme) => theme.palette.grey[500],
          }}
        >
          <CloseIcon />
        </IconButton>
      ) : null}
    </DialogTitle>
  );
};

const SoftwareUpdateModal = ({
  isOpen,
  handleClose,
  submit,
  content,
  title,
  btnLabel = 'UPDATE ALL',
}: softwareUpdateModalProps) => {
  const [checked, setChecked] = React.useState(false);

  const handleChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    setChecked(event.target.checked);
  };

  return (
    <div>
      <Dialog
        open={isOpen}
        onClose={handleClose}
        aria-labelledby="alert-dialog-title"
        aria-describedby="alert-dialog-description"
        sx={{ with: '100%' }}
      >
        <BootstrapDialogTitle
          id="customized-dialog-title"
          onClose={handleClose}
        >
          {title}
        </BootstrapDialogTitle>
        <DialogContent>
          <DialogContentText id="alert-dialog-description">
            <Typography sx={{ color: colors.black }}>{content}</Typography>
          </DialogContentText>
        </DialogContent>
        <DialogActions>
          <FormControlLabel
            control={<Checkbox checked={checked} onChange={handleChange} />}
            label="Don't show again"
          />
          <div
            style={{
              flex: '1 0 0',
            }}
          />
          <Button
            onClick={handleClose}
            sx={{
              marginRight: 3,
            }}
          >
            CANCEL
          </Button>
          <Button
            variant="contained"
            onClick={() => submit(checked)}
            sx={{ position: 'relative', right: 10 }}
          >
            {btnLabel}
          </Button>
        </DialogActions>
      </Dialog>
    </div>
  );
};
export default SoftwareUpdateModal;
