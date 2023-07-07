import { RoundedCard } from '@/styles/global';
import { Stack, TextField, Typography } from '@mui/material';

const NotificationContainer = () => {
  return (
    <RoundedCard radius="4px">
      <Stack spacing={2}>
        <Typography variant="body2" fontWeight={600}>
          Notification settings
        </Typography>
        <Typography variant="caption">
          All billing invoices will be sent to the primary email address.
        </Typography>
        <TextField
          label="Primary Email"
          id="primary-notification-email"
          defaultValue="default@email.com"
          sx={{ maxWidth: 350 }}
          InputProps={{
            readOnly: true,
          }}
        />
      </Stack>
    </RoundedCard>
  );
};

export default NotificationContainer;
