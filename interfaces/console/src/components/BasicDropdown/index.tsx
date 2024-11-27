/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import colors from '@/theme/colors';
import { FormControl, Typography } from '@mui/material';

interface IBasicDropdown {
  network: string;
  isShowAddOption: boolean;
}
const BasicDropdown = ({ network, isShowAddOption }: IBasicDropdown) => (
  <FormControl sx={{ width: '100%' }} size="small">
    {isShowAddOption && (
      <>
        <Typography variant="body1" sx={{ color: colors.primaryMain }}>
          {network || ''}
        </Typography>
      </>
    )}
  </FormControl>
);

export default BasicDropdown;
