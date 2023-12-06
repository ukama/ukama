import { COUNTRIES } from '@/utils';
import CloseIcon from '@mui/icons-material/Close';
import {
  Autocomplete,
  Button,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  Stack,
  TextField,
  Typography,
} from '@mui/material';

type AddNetworkDialogProps = {
  title: string;
  isOpen: boolean;
  description: string;
  isClosable?: boolean;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction?: any;
  labelNegativeBtn?: string;
};

const AddNetworkDialog = ({
  title,
  isOpen,
  description,
  labelSuccessBtn,
  labelNegativeBtn,
  handleCloseAction,
  isClosable = true,
  handleSuccessAction,
}: AddNetworkDialogProps) => {
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
          <CloseIcon fontSize="small" />
        </IconButton>
      </Stack>

      <DialogContent sx={{ maxHeight: '400px', overflow: 'auto' }}>
        <Stack spacing={2} direction={'column'} alignItems="start">
          {description && (
            <Typography variant="body1">{description}</Typography>
          )}
          <Stack direction={'row'} width={'100%'} spacing={2} mt={1}>
            <TextField
              sx={{ width: '50%' }}
              name={'name'}
              size="medium"
              placeholder="Mesh"
              label={'Network name'}
              InputLabelProps={{
                shrink: true,
              }}
              onChange={() => {}}
            />
            <TextField
              sx={{ width: '50%' }}
              name={'budget'}
              size="medium"
              placeholder="100"
              label={'Network budget'}
              InputLabelProps={{
                shrink: true,
              }}
              onChange={() => {}}
            />
          </Stack>
          <Autocomplete
            multiple
            defaultValue={[]}
            options={COUNTRIES}
            id="country-select"
            getOptionLabel={(option) => option.name}
            renderInput={(params) => (
              <TextField
                {...params}
                label="Allowed Countries"
                placeholder="Country"
                InputLabelProps={{
                  shrink: true,
                }}
              />
            )}
            sx={{
              width: '100%',
              mt: 1,
              p: 0,
              '.MuiOutlinedInput-root': {
                p: '4px',
              },
            }}
          />
          <Autocomplete
            multiple
            options={[]}
            defaultValue={[]}
            id="select-network"
            getOptionLabel={(option: any) =>
              option ? option?.name : 'No network available'
            }
            renderInput={(params) => (
              <TextField
                {...params}
                label="Allowed Networks"
                placeholder="Network"
                InputLabelProps={{
                  shrink: true,
                }}
              />
            )}
            sx={{
              width: '100%',
              mt: 1,
              '.MuiOutlinedInput-root': {
                p: '4px',
              },
            }}
          />
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

export default AddNetworkDialog;
