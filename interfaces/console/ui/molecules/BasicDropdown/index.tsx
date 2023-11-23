/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import colors from '@/styles/theme/colors';
import { SelectItemType } from '@/types';
import { FormControl, InputLabel, MenuItem, Select } from '@mui/material';

interface IBasicDropdown {
  value: string;
  placeholder: string;
  isLoading?: boolean;
  handleOnChange: Function;
  list: SelectItemType[];
}
const BasicDropdown = ({
  value,
  list,
  placeholder,
  handleOnChange,
}: IBasicDropdown) => (
  <FormControl sx={{ width: '100%' }} size="small">
    {!value && (
      <InputLabel sx={{ fontSize: '16px !important', pt: '10px' }}>
        {placeholder}
      </InputLabel>
    )}
    <Select
      value={value}
      disableUnderline
      variant="standard"
      onChange={(e) => handleOnChange(e.target.value)}
      sx={{
        p: 0,
        color: colors.primaryMain,
      }}
      SelectDisplayProps={{
        style: {
          fontWeight: 400,
          fontSize: '16px',
        },
      }}
    >
      {list?.map((item: SelectItemType) => (
        <MenuItem key={item.id} value={item.value}>
          {item.label}
        </MenuItem>
      ))}
    </Select>
  </FormControl>
);

export default BasicDropdown;
