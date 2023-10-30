/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Divider, Grid, Typography } from '@mui/material';
import { useCallback, useState } from 'react';
import FormControlCheckboxes from '../FormControlCheckboxes';

const LineDivider = () => (
  <Grid item xs={12}>
    <Divider sx={{ width: '100%' }} />
  </Grid>
);

const AlertSettings = () => {
  const [alertList, setAlertList] = useState<Object>({});

  const handleAlertChange = useCallback((key: string, value: boolean) => {
    setAlertList((prevState) => ({
      ...prevState,
      [key]: value,
    }));
  }, []);

  return (
    <Grid container spacing={2}>
      <Grid item container spacing={2}>
        <Grid item xs={12} md={3}>
          <Typography variant="h6">Common Events</Typography>
        </Grid>
        <Grid item xs={12} md={8}>
          {[1, 2].map((i) => (
            <Grid key={`${i}-`} item xs={12} sm={10} md={9}>
              <FormControlCheckboxes
                values={alertList}
                handleChange={handleAlertChange}
                checkboxList={[
                  {
                    id: 1,
                    label: `Event Log ${i}`,
                    value: `event${i}`,
                  },
                  {
                    id: 2,
                    label: `Alerts ${i}`,
                    value: `alert${i}`,
                  },
                  {
                    id: 3,
                    label: `Email ${i}`,
                    value: `email${i}`,
                  },
                ]}
              />
            </Grid>
          ))}
        </Grid>
      </Grid>
      <LineDivider />
      <Grid item container spacing={2}>
        <Grid item xs={12} md={3}>
          <Typography variant="h6">Cloud Events</Typography>
        </Grid>
        <Grid item container xs={12} md={9}>
          {[3, 4].map((i) => (
            <Grid key={`${i}-`} item xs={12} sm={10} md={8}>
              <FormControlCheckboxes
                values={alertList}
                handleChange={handleAlertChange}
                checkboxList={[
                  {
                    id: 1,
                    label: `Event Log ${i}`,
                    value: `event${i}`,
                  },
                  {
                    id: 2,
                    label: `Alerts ${i}`,
                    value: `alert${i}`,
                  },
                  {
                    id: 3,
                    label: `Email ${i}`,
                    value: `email${i}`,
                  },
                ]}
              />
            </Grid>
          ))}
        </Grid>
      </Grid>
    </Grid>
  );
};

export default AlertSettings;
