/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { getDataPlanUsage, getDuration } from '@/utils';
import { Card, Grid, Stack, Typography } from '@mui/material';
import OptionsPopover from '../OptionsPopover';
interface IPlanCard {
  uuid: string;
  name: string;
  amount: string;
  duration: number;
  currency: string;
  dataUnit: string;
  dataVolume: string;
  isOptions?: boolean;
  handleOptionMenuItemAction?: (type: string) => void;
}

const PlanCard = ({
  name,
  amount,
  dataUnit,
  duration,
  currency,
  dataVolume,
  isOptions = true,
  handleOptionMenuItemAction,
}: IPlanCard) => {
  return (
    <Card
      variant="outlined"
      sx={{
        p: { xs: 1.5, md: 3 },
        display: 'flex',
        boxShadow: 'none',
        borderRadius: '4px',
        textAlign: 'center',
        justifyContent: 'center',
        borderTop: `4px solid ${colors.primaryMain}`,
      }}
    >
      <Stack spacing={1} width={'100%'}>
        <Grid xs={12} container direction={'row'} textAlign={'center'}>
          <Grid item xs={11} pl={3}>
            <Typography
              variant="h5"
              sx={{
                fontWeight: 400,
                whiteSpace: 'nowrap',
                overflow: 'hidden',
                textOverflow: 'ellipsis',
              }}
            >
              {name}
            </Typography>
          </Grid>
          <Grid item xs={1}>
            {isOptions && (
              <OptionsPopover
                cid={'data-table-action-popover'}
                menuOptions={[
                  { id: 0, title: 'Edit', route: 'edit', Icon: null },
                ]}
                handleItemClick={(type: string) =>
                  handleOptionMenuItemAction && handleOptionMenuItemAction(type)
                }
              />
            )}
          </Grid>
        </Grid>
        <Typography variant="body2" fontWeight={400}>
          {getDataPlanUsage(
            getDuration(duration),
            currency,
            amount,
            dataVolume,
            dataUnit,
          )}
        </Typography>
      </Stack>
    </Card>
  );
};

export default PlanCard;
