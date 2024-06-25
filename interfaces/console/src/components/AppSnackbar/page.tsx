import { useAppContext } from '@/context';
import { Alert, AlertColor, Snackbar } from '@mui/material';

const SNACKBAR_TIMEOUT = 5000;

const AppSnackbar = () => {
  const { snackbarMessage, setSnackbarMessage } = useAppContext();

  const handleSnackbarClose = () =>
    setSnackbarMessage({ ...snackbarMessage, show: false });

  return (
    <Snackbar
      open={snackbarMessage.show}
      autoHideDuration={SNACKBAR_TIMEOUT}
      onClose={handleSnackbarClose}
    >
      <Alert
        id={snackbarMessage.id}
        severity={snackbarMessage.type as AlertColor}
        onClose={handleSnackbarClose}
      >
        {snackbarMessage.message}
      </Alert>
    </Snackbar>
  );
};

export default AppSnackbar;
