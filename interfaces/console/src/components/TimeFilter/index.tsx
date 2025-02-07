/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { StatsPeriodItemType } from '@/types';
import {
  Box,
  ToggleButton,
  ToggleButtonGroup,
  Typography,
} from '@mui/material';

interface ITimeFilter {
  filter?: string;
  options?: StatsPeriodItemType[];
  handleFilterSelect: Function;
}

const TimeFilter = ({
  filter = 'LIVE',
  handleFilterSelect,
  options = [{ id: '1', label: 'LIVE' }],
}: ITimeFilter) => {
  return (
    <Box component="div">
      <ToggleButtonGroup
        exclusive
        size="small"
        color="primary"
        value={filter}
        onChange={(_, value: string) => handleFilterSelect(value)}
      >
        {options.map(({ id, label }: StatsPeriodItemType) => (
          <ToggleButton
            fullWidth
            key={id}
            value={label}
            style={{
              height: '32px',
              color: colors.hoverColor,
              border: `1px solid ${colors.hoverColor}`,
            }}
          >
            <Typography
              variant="body2"
              sx={{
                p: '0px 2px',
                fontWeight: 600,
              }}
            >
              {label}
            </Typography>
          </ToggleButton>
        ))}
      </ToggleButtonGroup>
    </Box>
  );
};
export default TimeFilter;
