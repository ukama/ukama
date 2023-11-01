/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

// import { AlertDto } from '@/generated';
import { getColorByType } from '@/utils';
import CloudOffIcon from '@mui/icons-material/CloudOff';
import { Grid, List, ListItem, Typography } from '@mui/material';
import { format } from 'date-fns';

// type AlertsProps = {
//   alertOptions: AlertDto[] | undefined;
// };

const PROPS = {
  display: 'flex',
  alignItems: 'center',
};

const Alerts = ({ alertOptions = [] }: any) => {
  return (
    <List
      sx={{
        pr: 3,
        width: 372,
        maxHeight: 305,
        overflowY: 'scroll',
        overflowX: 'hidden',
        position: 'relative',
      }}
    >
      {alertOptions.map(({ id, alertDate, description, title, type }: any) => (
        <ListItem key={id} sx={{ p: 0, mb: 2 }}>
          <Grid container>
            <Grid item container>
              <Grid item xs={1.4} {...PROPS}>
                <CloudOffIcon fontSize="small" color={getColorByType(type)} />
              </Grid>
              <Grid item xs={6.6} {...PROPS}>
                <Typography variant="body1" sx={{ fontWeight: 500 }}>
                  {title}
                </Typography>
              </Grid>
              <Grid item xs={4} {...PROPS} justifyContent={'flex-end'}>
                <Typography variant="caption" color={'textSecondary'}>
                  {format(new Date(alertDate), 'MMM dd ha')}
                </Typography>
              </Grid>
            </Grid>
            <Grid item container>
              <Grid item xs={1.4} />
              <Grid item xs={10.6}>
                <Typography variant="body2" color={'textSecondary'} {...PROPS}>
                  {description}
                </Typography>
              </Grid>
            </Grid>
          </Grid>
        </ListItem>
      ))}
    </List>
  );
};
export default Alerts;
