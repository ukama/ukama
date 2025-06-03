/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { MenuItemType } from '@/types';
import MenuDots from '@mui/icons-material/MoreHoriz';
import {
  IconButton,
  ListItemIcon,
  ListItemText,
  MenuItem,
  Popover,
} from '@mui/material';
import { useState } from 'react';

type ItemProps = {
  Icon: any;
  type: string;
  title: string;
  titleColor?: string;
  isShowUpdate: boolean;
  handleItemClick: (type: string) => void;
};

const OptionItem = ({
  type,
  Icon,
  title,
  isShowUpdate,
  handleItemClick,
  titleColor = 'inherit',
}: ItemProps) => (
  <MenuItem
    onClick={() => handleItemClick(type)}
    sx={{ display: isShowUpdate ? 'flex' : 'none' }}
  >
    {Icon && (
      <ListItemIcon>
        <Icon fontSize="small" />
      </ListItemIcon>
    )}
    <ListItemText
      sx={{
        color: titleColor,
        mr: 1,
        '.MuiTypography-root': {
          color: title.toLowerCase() === 'delete' ? colors.redMatt : 'inherit',
        },
      }}
    >
      {title}
    </ListItemText>
    {type === 'update' && (
      <div
        style={{
          width: 6,
          height: 6,
          borderRadius: '100%',
          backgroundColor: colors.primaryMain,
        }}
      />
    )}
  </MenuItem>
);

type OptionsPopoverProps = {
  cid: string;
  isShowUpdate?: boolean;
  menuOptions: MenuItemType[];
  handleItemClick: (type: string) => void;
  style?: any;
};

const OptionsPopover = ({
  cid,
  menuOptions,
  isShowUpdate = false,
  handleItemClick,
  style,
}: OptionsPopoverProps) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const handlePopoverClose = () => setAnchorEl(null);
  const handlePopoverOpen = (event: any) => setAnchorEl(event.currentTarget);

  const open = Boolean(anchorEl);
  const id = open ? cid : undefined;
  return (
    <>
      <IconButton
        id={cid}
        onClick={handlePopoverOpen}
        aria-describedby={id}
        style={style}
        sx={{ p: 0, position: 'relative' }}
      >
        <MenuDots fontSize="small" />
        {isShowUpdate && (
          <div
            style={{
              top: '2px',
              right: 0,
              width: 6,
              height: 6,
              position: 'absolute',
              borderRadius: '100%',
              backgroundColor: colors.primaryMain,
            }}
          />
        )}
      </IconButton>
      <Popover
        id={id}
        open={open}
        anchorEl={anchorEl}
        onClose={handlePopoverClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'right',
        }}
      >
        {menuOptions.map(({ id: optId, Icon, title, color, route }: any) => (
          <OptionItem
            Icon={Icon}
            type={route}
            title={title}
            titleColor={color}
            isShowUpdate={true}
            key={`${title}-${optId}`}
            handleItemClick={(type: string) => {
              handleItemClick(type);
              handlePopoverClose();
            }}
          />
        ))}
      </Popover>
    </>
  );
};

export default OptionsPopover;
