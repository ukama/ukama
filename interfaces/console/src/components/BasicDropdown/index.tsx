/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import colors from '@/theme/colors';
import { SelectItemType } from '@/types';
import Add from '@mui/icons-material/Add';
import { Button, Divider, FormControl, MenuItem, Select } from '@mui/material';

interface IBasicDropdown {
  value: string;
  placeholder: string;
  isDisableAddOption?: boolean;
  isShowAddOption: boolean;
  handleOnChange: Function;
  handleAddNetwork: Function;
  list: SelectItemType[];
}
const BasicDropdown = ({
  value,
  list,
  placeholder,
  handleOnChange,
  isShowAddOption,
  handleAddNetwork,
  isDisableAddOption = false,
}: IBasicDropdown) => (
  <FormControl sx={{ width: '100%' }} size="small">
    <Select
      value={value}
      disableUnderline
      variant="standard"
      displayEmpty={true}
      renderValue={(value) =>
        value?.length
          ? Array.isArray(value)
            ? value.join(', ')
            : list?.find((item) => item.value === value)?.label
          : placeholder
      }
      onChange={(e) => handleOnChange(e.target.value)}
      sx={{
        color: value?.length
          ? Array.isArray(value)
            ? value.join(', ')
            : colors.primaryMain
          : colors.black38,
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
      {list?.length === 0 && <MenuItem disabled>No network found!</MenuItem>}
      {isShowAddOption && (
        <>
          <Divider sx={{ width: '100%', height: '1px' }} />
          <Button
            disabled={isDisableAddOption}
            startIcon={
              <Add
                sx={{
                  color: isDisableAddOption ? colors.silver : colors.black70,
                }}
              />
            }
            sx={{
              px: 2,
              py: 1,
              color: 'textPrimary',
              typography: 'body1',
              fontWeight: 400,
              display: 'flex',
              cursor: 'pointer',
              textTransform: 'none',
              alignItems: 'center',
              justifyContent: 'center',
              ':hover': {
                backgroundColor: colors.primaryMain02,
              },
              ':hover .MuiSvgIcon-root': {
                fill: colors.primaryMain,
              },
            }}
            onClick={(e) => {
              handleAddNetwork();
              e.stopPropagation();
            }}
          >
            Add new network
          </Button>
        </>
      )}
    </Select>
  </FormControl>
);

export default BasicDropdown;
