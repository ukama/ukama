import { GetUsersDto } from '@/generated';
import { colors } from '@/styles/theme';
import { formatBytesToMB } from '@/utils';
import EditIcon from '@mui/icons-material/Edit';
import {
  Grid,
  IconButton,
  LinearProgress,
  Stack,
  Typography,
} from '@mui/material';
import { useEffect, useState } from 'react';
import LoadingWrapper from '../LoadingWrapper';

type UserCardProps = {
  loading: boolean;
  user: GetUsersDto;
  // eslint-disable-next-line no-unused-vars
  handleMoreUserdetails: (user: GetUsersDto) => void;
};

const UserCard = ({ user, loading, handleMoreUserdetails }: UserCardProps) => {
  const [dataLoading, setDataLoading] = useState(user?.dataPlan === '');

  useEffect(() => {
    if (
      user.dataPlan !== '' &&
      (user.dataPlan !== '0' || user.dataPlan === '0')
    ) {
      setDataLoading(false);
    }
  }, [loading, user]);
  return (
    <Grid container spacing={{ xs: 1.5 }}>
      <Grid item xs={12} container>
        <Grid
          item
          xs={10}
          sx={{
            overflow: 'hidden',
            textOverflow: 'ellipsis',
            width: '11rem',
            whiteSpace: 'normal',
          }}
        >
          <Typography
            variant="body2"
            color="textSecondary"
            sx={{
              width: '200px',
              whiteSpace: 'nowrap',
              overflow: 'hidden',
              textOverflow: 'ellipsis',
            }}
          >
            {user.id}
          </Typography>
        </Grid>
        <Grid
          item
          xs={2}
          container
          justifyContent="flex-end"
          sx={{ position: 'relative', bottom: 10 }}
        >
          <IconButton edge="end" onClick={() => handleMoreUserdetails(user)}>
            <EditIcon sx={{ fontSize: 25 }} />
          </IconButton>
        </Grid>
      </Grid>
      <Stack direction="row"></Stack>
      <Grid item xs={12}>
        <Typography variant="h5">{user.name}</Typography>
      </Grid>
      <Grid item xs={4}>
        <LoadingWrapper
          width="100%"
          height="36px"
          radius="small"
          variant="text"
          isLoading={dataLoading}
        >
          <Stack direction="row" spacing={1} alignItems="baseline">
            <Typography variant="h5">
              {formatBytesToMB(parseInt(user?.dataUsage || '0'))}
            </Typography>
            <Typography variant="body2" textAlign={'end'}>
              MB
            </Typography>
          </Stack>
        </LoadingWrapper>
      </Grid>
      <Grid item xs={8} alignSelf="end" sx={{ position: 'relative', top: 8 }}>
        <LoadingWrapper
          width="100%"
          height="23px"
          radius="small"
          variant="text"
          isLoading={dataLoading}
        >
          <Stack direction="column">
            <Typography variant="body2" textAlign={'end'}>
              {`${formatBytesToMB(
                parseInt(user?.dataPlan || '0') -
                  parseInt(user?.dataUsage || '0'),
              )} MB free `}
            </Typography>

            <Typography variant="body2" textAlign={'end'}>
              {`data left`}
            </Typography>
          </Stack>
        </LoadingWrapper>
      </Grid>
      <Grid item xs={12} display="grid" sx={{ pb: 2 }}>
        <LoadingWrapper
          width="100%"
          height="8px"
          radius="small"
          variant="text"
          isLoading={dataLoading}
        >
          <LinearProgress
            variant="determinate"
            value={
              user.dataPlan && user.dataPlan !== '0'
                ? (parseInt(user?.dataUsage || '0') * 100) /
                  parseInt(user?.dataPlan || '0')
                : 0
            }
            sx={{
              height: '8px',
              borderRadius: '2px',
              backgroundColor: colors.silver,
            }}
          />
        </LoadingWrapper>
      </Grid>
    </Grid>
  );
};
export default UserCard;
