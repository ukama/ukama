import { RoundedCard } from '@/styles/global';
import { SelectItemType } from '@/types';
import {
  Box,
  Grid,
  MenuItem,
  Select,
  Theme,
  Typography,
  useMediaQuery,
} from '@mui/material';
import { makeStyles } from '@mui/styles';
import LoadingWrapper from '../../molecules/LoadingWrapper';

const useStyles = makeStyles<Theme>((theme) => ({
  selectStyle: {
    width: '108px',
    textAlign: 'end',
    '& p': {
      color: theme?.palette?.text?.secondary,
      fontWeight: 500,
      fontSize: '14px',
      lineHeight: '157%',
    },
    '& .MuiSelect-iconStandard': {
      paddingBottom: '4px',
    },
    '& .MuiSelect-iconOpen': {
      paddingBottom: '0px',
    },
  },
}));

type StatusCardProps = {
  Icon: any;
  title: string;
  option: string;
  loading: boolean;
  subtitle1: string;
  subtitle2: string;
  handleSelect: Function;
  options: SelectItemType[];
};

const StatusCard = ({
  Icon,
  title,
  option,
  options,
  loading,
  subtitle1 = '0',
  subtitle2 = '',
  handleSelect,
}: StatusCardProps) => {
  const classes = useStyles();
  const isSmall = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));

  return (
    <LoadingWrapper height={isSmall ? 64 : 100} isLoading={loading}>
      {isSmall ? (
        <Box
          component="div"
          sx={{
            py: 1,
            px: '4px',
            borderRadius: '4px',
            bgcolor:
              title === 'Connected Users'
                ? 'rgba(33, 144, 246, 0.1)'
                : title === 'Data Usage'
                ? 'rgba(105, 116, 248, 0.1)'
                : 'rgba(3, 116, 75, 0.1)',
          }}
        >
          <Grid container alignItems="center">
            <Grid
              item
              xs={6}
              container
              sx={{ svg: { rect: { fill: 'transparent' } } }}
            >
              <Icon />
            </Grid>
            <Grid item xs={6}>
              <Typography variant="body2" paddingRight="6px">
                {`${subtitle1}${title === 'Data Usage' ? ' MBs' : ''}`}
              </Typography>
            </Grid>
          </Grid>
        </Box>
      ) : (
        <RoundedCard>
          <Grid spacing={2} container direction="row" justifyContent="center">
            <Grid item xs={2} display="flex" alignItems="center">
              <Icon />
            </Grid>
            <Grid xs={10} item sm container direction="column">
              <Grid
                sm
                item
                container
                spacing={2}
                display="flex"
                direction="row"
                alignItems="center"
              >
                <Grid item xs={12} mb={{ xs: 0.6, sm: 0 }}>
                  <Typography variant="subtitle2">{title}</Typography>
                </Grid>
                <Grid item xs={5} display="none" justifyContent="flex-end">
                  <Select
                    value={option}
                    disableUnderline
                    variant="standard"
                    className={classes.selectStyle}
                    MenuProps={{
                      sx: {
                        maxHeight: '194px',
                      },
                    }}
                    onChange={(e) => handleSelect(e.target.value)}
                  >
                    {options.map(({ id, label, value }: SelectItemType) => (
                      <MenuItem key={id} value={value}>
                        <Typography variant="body1">{label}</Typography>
                      </MenuItem>
                    ))}
                  </Select>
                </Grid>
              </Grid>
              <Grid item container alignItems="baseline">
                <Grid item>
                  <Typography variant="h5" paddingRight="6px">
                    {subtitle1}
                  </Typography>
                </Grid>
                {title === 'Data Usage' && (
                  <Grid item>
                    <Typography variant="body1" paddingRight="4px">
                      MBs
                    </Typography>
                  </Grid>
                )}
                <Grid item>
                  <Typography variant="body1" color="textSecondary">
                    {subtitle2}
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </RoundedCard>
      )}
    </LoadingWrapper>
  );
};
export default StatusCard;
