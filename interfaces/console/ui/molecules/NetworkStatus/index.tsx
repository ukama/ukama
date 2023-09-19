import { colors } from '@/styles/theme';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import SignalCellularOffIcon from '@mui/icons-material/SignalCellularOff';
import { Stack, Tooltip, Typography } from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';

const getIconByStatus = (status: string, info: string) => {
  switch (status) {
    case 'DOWN':
      return (
        <Tooltip title={info}>
          <CheckCircleIcon sx={{ color: colors.red }} />
        </Tooltip>
      );

    case 'ONLINE':
      return (
        <Tooltip title={info}>
          <CheckCircleIcon sx={{ color: colors.green }} />
        </Tooltip>
      );

    default:
      return (
        <Tooltip title={info}>
          <SignalCellularOffIcon />
        </Tooltip>
      );
  }
};

type NetworkStatusProps = {
  loading?: boolean;
  tooltipInfo: string;
  availableNodes: number | undefined;
  statusType: string | undefined;
};

const NetworkStatus = ({
  loading,
  tooltipInfo,
  availableNodes = undefined,
  statusType = 'onboarding',
}: NetworkStatusProps) => {
  return (
    <LoadingWrapper
      isLoading={loading}
      height={loading ? '30px' : 'fit-content'}
      width={loading ? '40%' : 'fit-content'}
    >
      <Stack direction={'row'} alignItems="center" spacing={1}>
        {getIconByStatus(statusType, tooltipInfo)}
        <Typography variant={'h6'} sx={{ fontWeight: { xs: 400, md: 500 } }}>
          Joe’s network requires
        </Typography>
        <Typography variant={'h6'} sx={{ fontWeight: { xs: 400, md: 500 } }}>
          <u>node installation</u>.
        </Typography>
        {availableNodes && (
          <Typography
            variant={'h6'}
            color="secondary"
            sx={{
              fontWeight: {
                xs: 400,
                md: 500,
              },
              whiteSpace: 'break-spaces',
              ml: { xs: '28px', md: '8px' },
            }}
          >
            {`${availableNodes} available nodes`}
          </Typography>
        )}
      </Stack>
    </LoadingWrapper>
  );
};

export default NetworkStatus;
