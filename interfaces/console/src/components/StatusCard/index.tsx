/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { RoundedCard } from '@/styles/global';
import { colors } from '@/theme';
import { SelectItemType } from '@/types';
import {
  Box,
  Grid,
  Skeleton,
  Stack,
  Theme,
  Typography,
  useMediaQuery,
} from '@mui/material';
import { makeStyles } from '@mui/styles';

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
  iconColor?: string;
  handleSelect: Function;
  options: SelectItemType[];
};

const StatusCard = ({
  Icon,
  title,
  option,
  options,
  loading,
  iconColor,
  subtitle1 = '0',
  subtitle2 = '',
  handleSelect,
}: StatusCardProps) => {
  const classes = useStyles();
  const isSmall = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));

  return (
    <>
      {isSmall ? (
        <Box
          component="div"
          sx={{
            py: 2,
            px: 1.5,
            borderRadius: '4px',
            bgcolor: colors.white,
            // bgcolor:
            //   title === 'Connected Users'
            //     ? 'rgba(33, 144, 246, 0.1)'
            //     : title === 'Data Usage'
            //       ? 'rgba(105, 116, 248, 0.1)'
            //       : 'rgba(3, 116, 75, 0.1)',
          }}
        >
          <Grid container alignItems="center">
            <Grid container xs={12} alignItems="center">
              <Stack direction="row" alignItems="center" spacing={1.5}>
                <Box
                  sx={{
                    svg: {
                      fill: iconColor,
                    },
                  }}
                >
                  <Icon />
                </Box>
                <Typography variant="body1">
                  {`${subtitle1}${title === 'Data Usage' ? ' MBs' : ''}`}
                </Typography>
              </Stack>
              {/* <Grid item xs={6}>
                <Typography variant="body2" paddingRight="6px">
                  {`${subtitle1}${title === 'Data Usage' ? ' MBs' : ''}`}
                </Typography>
              </Grid> */}
            </Grid>
            <Grid item xs={12}>
              <Typography variant="caption">{title}</Typography>
            </Grid>
          </Grid>
        </Box>
      ) : (
        <RoundedCard>
          <Grid spacing={2} container direction="row" justifyContent="center">
            <Grid item xs={2} display="flex" alignItems="center">
              <Box
                sx={{
                  svg: {
                    fill: iconColor,
                  },
                }}
              >
                <Icon />
              </Box>
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
              </Grid>
              <Grid item container alignItems="baseline">
                <Grid item>
                  <Typography variant="h5" paddingRight="6px">
                    {loading ? (
                      <Skeleton variant="text" width={64} />
                    ) : (
                      subtitle1
                    )}
                  </Typography>
                </Grid>
                {title === 'Data Usage' && (
                  <Grid item>
                    <Typography variant="body1" paddingRight="4px">
                      {loading ? <Skeleton variant="text" width={64} /> : 'MBs'}
                    </Typography>
                  </Grid>
                )}
                <Grid item>
                  <Typography variant="body1" color="textSecondary">
                    {loading ? (
                      <Skeleton variant="text" width={64} />
                    ) : (
                      subtitle2
                    )}
                  </Typography>
                </Grid>
              </Grid>
            </Grid>
          </Grid>
        </RoundedCard>
      )}
    </>
  );
};
export default StatusCard;
