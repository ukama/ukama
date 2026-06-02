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
import { Divider, FormControl, MenuItem, Select } from '@mui/material';

interface IBasicDropdown {
  id: string;
  value: string;
  placeholder: string;
  isDisableAddOption?: boolean;
  isShowAddOption: boolean;
  handleOnChange: (value: string) => void;
  handleAddNetwork: () => void;
  list: SelectItemType[];
}

const BasicDropdown = ({
  id,
  value,
  list,
  placeholder,
  handleOnChange,
  isShowAddOption,
  handleAddNetwork,
  isDisableAddOption = false,
}: IBasicDropdown) => {
  const getDisplayValue = (value: string) => {
    if (!value?.length) return placeholder;
    return Array.isArray(value)
      ? value.join(', ')
      : list?.find((item) => item.value === value)?.label;
  };

  return (
    <FormControl id={id} sx={{ width: '100%' }} size="small">
      <Select
        id={id}
        name={id}
        displayEmpty
        value={value}
        disableUnderline
        variant="standard"
        data-testid={`${id}-select`}
        renderValue={getDisplayValue}
        onChange={(e) => handleOnChange(e.target.value)}
        sx={{
          color: value?.length ? colors.primaryMain : colors.black38,
        }}
        SelectDisplayProps={{
          style: {
            fontWeight: 400,
            fontSize: '16px',
          },
        }}
      >
        {list?.map((item: SelectItemType, index: number) => (
          <MenuItem
            key={`network-option-${index}`}
            value={item.value}
            data-testid={`${id}-option-${item.value}`}
          >
            {item.label}
          </MenuItem>
        ))}
        {list?.length === 0 && (
          <MenuItem disabled data-testid={`${id}-no-options`}>
            No network found!
          </MenuItem>
        )}
        {isShowAddOption && <Divider component="li" />}
        {isShowAddOption && (
          <MenuItem
            id={`${id}-add`}
            data-testid={`${id}-add-button`}
            disabled={isDisableAddOption}
            onClick={(e) => {
              handleAddNetwork();
              e.stopPropagation();
            }}
            sx={{
              gap: 1,
              color: 'textPrimary',
              typography: 'body1',
              fontWeight: 400,
              ':hover .MuiSvgIcon-root': { fill: colors.primaryMain },
            }}
          >
            <Add
              fontSize="small"
              sx={{ color: isDisableAddOption ? colors.silver : colors.black70 }}
            />
            Add new network
          </MenuItem>
        )}
      </Select>
    </FormControl>
  );
};

export default BasicDropdown;
