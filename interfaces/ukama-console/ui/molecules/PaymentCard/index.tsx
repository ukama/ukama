import { isDarkmode } from '@/app-recoil';
import colors from '@/styles/theme/colors';
import {
  Box,
  Button,
  Divider,
  FormControl,
  InputLabel,
  MenuItem,
  OutlinedInput,
  Select,
  Stack,
  Typography,
} from '@mui/material';
import { makeStyles } from '@mui/styles';
import { useRecoilValue } from 'recoil';
const useStyles = makeStyles(() => ({
  '&.MuiFormHelperText-root.Mui-error': {
    color: 'red',
  },
  selectStyle: () => ({
    width: '100%',
    height: '48px',
  }),
  formControl: {
    width: '100%',
    height: '48px',
    paddingBottom: '55px',
  },
}));

interface IPaymentProps {
  title: string;
  onChangePM: any;
  selectedPM: string;
  paymentMethodData: any;
  onAddPaymentMethod: any;
}
const PaymentCard = ({
  title,
  onChangePM,
  selectedPM,
  paymentMethodData,
  onAddPaymentMethod,
}: IPaymentProps) => {
  const classes = useStyles();

  const _isDarkMod = useRecoilValue(isDarkmode);

  const isDiable = () =>
    paymentMethodData.length === 1 &&
    paymentMethodData[0].value === 'no_payment_method_Set'
      ? true
      : false;

  return (
    <Box>
      <Typography variant="h6" sx={{ pb: 3 }}>
        {title}
      </Typography>
      <FormControl variant="outlined" className={classes.formControl}>
        <InputLabel
          shrink
          variant="outlined"
          htmlFor="outlined-age-always-notched"
        >
          PAYMENT METHOD
        </InputLabel>
        <Select
          value={selectedPM}
          variant="outlined"
          onChange={onChangePM}
          IconComponent={() => null}
          sx={{
            '& legend': { width: '135px' },
            '& #add-payment-method': {
              color: `${colors.primaryMain} !important`,
              '-webkit-text-fill-color': `${colors.primaryMain} !important`,
              ':hover': {
                color: (theme) => `${theme.palette.text.primary} !important`,
                '-webkit-text-fill-color': (theme) =>
                  `${theme.palette.text.primary} !important`,
              },
            },
          }}
          input={
            <OutlinedInput
              notched
              label="NODE TYPE"
              name="node_type"
              id="outlined-age-always-notched"
            />
          }
          disabled={isDiable()}
          MenuProps={{
            disablePortal: false,
            PaperProps: {
              sx: {
                boxShadow:
                  '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
                borderRadius: '4px',
              },
            },
          }}
          className={classes.selectStyle}
        >
          {paymentMethodData.map(({ id, value, label }: any) => (
            <MenuItem
              key={id}
              value={value}
              sx={{
                m: 0,
                p: '6px 16px',
              }}
            >
              <Stack direction="row" spacing={1}>
                <Typography variant="body1">{label}</Typography>
              </Stack>
            </MenuItem>
          ))}
          {!isDiable() && (
            <Box>
              <Divider />
              <Button
                variant="text"
                sx={{
                  padding: '6px 16px',
                  typography: 'body1',
                  textTransform: 'none',
                }}
                onClick={(e) => {
                  onAddPaymentMethod();
                  e.stopPropagation();
                }}
              >
                Add new payment method
              </Button>
            </Box>
          )}
        </Select>
      </FormControl>
      <Typography
        variant="caption"
        sx={{ color: _isDarkMod ? colors.white : colors.black54 }}
      >
        *Automatically charged to card EOD on the last day of the billing cycle
      </Typography>
    </Box>
  );
};

export default PaymentCard;
