import { colors } from '@/styles/theme';
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
  isShowUpdate: boolean;
  handleItemClick: Function;
};

const OptionItem = ({
  type,
  Icon,
  title,
  isShowUpdate,
  handleItemClick,
}: ItemProps) => (
  <MenuItem
    onClick={() => handleItemClick(type)}
    sx={{ display: isShowUpdate ? 'flex' : 'none' }}
  >
    <ListItemIcon>
      <Icon fontSize="small" />
    </ListItemIcon>
    <ListItemText sx={{ mr: 1 }}>{title}</ListItemText>
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
  handleItemClick: Function;
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
        {menuOptions.map(({ id: optId, Icon, title, route }: any) => (
          <OptionItem
            key={`${cid}-${optId}`}
            type={route}
            Icon={Icon}
            title={title}
            isShowUpdate={optId === 3 ? isShowUpdate : true}
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
