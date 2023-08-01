import { colors } from '@/styles/theme';
import { Alert, AlertColor, Button, Stack, Typography } from '@mui/material';

interface IBillingAlerts {
  title: string;
  btnText: string;
  type: AlertColor;
  onActionClick: Function;
}

const BillingAlerts = ({
  type,
  title,
  btnText,
  onActionClick,
}: IBillingAlerts) => {
  return (
    <Alert
      icon={false}
      severity={type}
      sx={{
        background: type == 'info' ? colors.white : colors.lightRed,
        color: colors.black,
      }}
    >
      <Stack direction={'row'} p={0} px={1} spacing={1}>
        <Typography variant="body1">{title}</Typography>
        <Button
          variant="text"
          sx={{
            textTransform: 'none',
            color: colors.primaryMain,
            ':hover': {
              textDecoration: 'underline',
            },
          }}
          onClick={() => onActionClick()}
        >
          <Typography variant="body1"> {btnText}</Typography>
        </Button>
      </Stack>
    </Alert>
  );
};

export default BillingAlerts;
