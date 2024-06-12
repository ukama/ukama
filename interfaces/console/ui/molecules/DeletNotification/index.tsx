import { Box, Typography } from '@mui/material';

const DeleteNotification = () => {
  return (
    <Box p={'6px 16px'} display="flex" alignItems="center" borderRadius="4px">
      <Typography
        variant="body1"
        fontFamily="Work Sans"
        fontSize="16px"
        fontWeight="400"
        sx={{ cursor: 'pointer' }}
      >
        Delete Notification
      </Typography>
    </Box>
  );
};
export default DeleteNotification;
