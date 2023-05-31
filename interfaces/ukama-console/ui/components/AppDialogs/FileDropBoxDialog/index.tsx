import colors from '@/styles/theme/colors';
import CloseIcon from '@mui/icons-material/Close';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  Typography,
} from '@mui/material';
import { useDropzone } from 'react-dropzone';

type FileDropBoxDialogProps = {
  title: string;
  isOpen: boolean;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction?: any;
  labelNegativeBtn?: string;
};

const FileDropBoxDialog = ({
  title,
  isOpen,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
}: FileDropBoxDialogProps) => {
  const { acceptedFiles, getRootProps, getInputProps } = useDropzone();
  const files = acceptedFiles.map((file: any) => (
    <li key={file.path}>{file.path}</li>
  ));
  return (
    <Dialog
      fullWidth
      open={isOpen}
      maxWidth="sm"
      onClose={handleCloseAction}
      aria-labelledby="alert-dialog-title"
      aria-describedby="alert-dialog-description"
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
        <Box
          sx={{
            width: '100%',
            display: 'flex',
            padding: '2rem',
            cursor: 'pointer',
            justifyContent: 'center',
            border: '1px dashed grey',
            backgroundColor: colors.white38,
            ':hover': {
              border: '1px dashed black',
            },
          }}
        >
          <div {...getRootProps({ className: 'dropzone' })}>
            <input {...getInputProps()} />
            <Typography variant="body2" sx={{ cursor: 'inherit' }}>
              Drag & Drop file here Or click to select file.
            </Typography>
          </div>
        </Box>
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
              onClick={() =>
                handleSuccessAction
                  ? handleSuccessAction()
                  : handleCloseAction()
              }
            >
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default FileDropBoxDialog;
