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
import { useEffect, useState } from 'react';

type StatusCardProps = {
  Icon: any;
  topic: string;
  title: string;
  option: string;
  loading: boolean;
  subtitle2: string;
  iconColor?: string;
  handleSelect: () => void;
  options: SelectItemType[];
};

const StatusCard = ({
  Icon,
  topic,
  title,
  loading,
  iconColor,
  subtitle2 = '',
}: StatusCardProps) => {
  const [subtitle1, setSubtitle1] = useState<string>('0');
  const isSmall = useMediaQuery((theme: Theme) => theme.breakpoints.down('md'));

  useEffect(() => {
    const c = PubSub.subscribe(topic, (_, data) => {
      setSubtitle1(data);
    });
    return () => {
      PubSub.unsubscribe(c);
    };
  }, [topic]);

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
          }}
        >
          <Grid container alignItems="center">
            <Grid container xs={12} alignItems="center">
              <Stack direction="row" alignItems="center" spacing={1.5}>
                <Box
                  sx={{
                    svg: {
                      fill: iconColor,
                      width: '28px',
                      height: '28px',
                    },
                  }}
                >
                  <Icon />
                </Box>
                <Typography variant="body1">{subtitle1}</Typography>
              </Stack>
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
                    width: '28px',
                    height: '28px',
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
