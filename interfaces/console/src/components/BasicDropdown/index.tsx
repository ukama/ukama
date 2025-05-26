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
        {isShowAddOption && (
          <>
            <Divider sx={{ width: '100%', height: '1px' }} />
            <Button
              id={`${id}-add`}
              data-testid={`${id}-add-button`}
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
};

export default BasicDropdown;
