import { SIM_TYPES } from '@/constants';
import colors from '@/styles/theme/colors';
import { fileToBase64 } from '@/utils';
import CloseIcon from '@mui/icons-material/Close';
import DeleteOutlineOutlinedIcon from '@mui/icons-material/DeleteOutlineOutlined';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  FormControl,
  IconButton,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  Typography,
} from '@mui/material';
import { useEffect, useState } from 'react';
import { useDropzone } from 'react-dropzone';

type FileDropBoxDialogProps = {
  title: string;
  isOpen: boolean;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction: any;
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
  const [file, setFile] = useState<any>();
  const [simType, setSimType] = useState('');
  const { acceptedFiles, getRootProps, getInputProps } = useDropzone({
    accept: {
      'text/html': ['.csv'],
    },
    maxFiles: 1,
  });

  useEffect(() => {
    if (acceptedFiles.length > 0) {
      setFile(acceptedFiles[0]);
    }
  }, [acceptedFiles]);

  const handleUploadAction = () => {
    if (acceptedFiles && acceptedFiles.length > 0) {
      const file: any = acceptedFiles[0];
      fileToBase64(file)
        .then((base64String) => {
          handleSuccessAction('success', base64String, simType);
          handleCloseAction();
        })
        .catch((error) => {
          handleSuccessAction('error', error, simType);
          handleCloseAction();
        });
    }
  };

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
        <Stack
          direction="column"
          alignItems="flex-start"
          justifyContent="center"
          spacing={2}
        >
          <FormControl fullWidth>
            <InputLabel id={'simpool-sim-type-label'} shrink>
              Sim Type
            </InputLabel>
            <Select
              notched
              required
              label="Sim Type"
              value={simType}
              id={'simpool-sim-type-select'}
              labelId="simpool-sim-type-label"
              onChange={(e) => setSimType(e.target.value as string)}
            >
              {SIM_TYPES.map(({ id, label, value }) => (
                <MenuItem key={id} value={value}>
                  {label}
                </MenuItem>
              ))}
            </Select>
          </FormControl>

          {file ? (
            <Stack direction={'row'} spacing={2} alignItems={'center'}>
              <Typography variant="body1">{acceptedFiles[0].name}</Typography>
              <IconButton
                onClick={() => {
                  setFile(null);
                  acceptedFiles.pop();
                }}
                size="small"
              >
                <DeleteOutlineOutlinedIcon fontSize="small" />
              </IconButton>
            </Stack>
          ) : (
            <Box
              sx={{
                width: '100%',
                display: 'flex',
                padding: '4rem',
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
          )}
        </Stack>
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
            <Button variant="contained" onClick={handleUploadAction}>
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default FileDropBoxDialog;
