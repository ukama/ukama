import React, { useState } from 'react';
import { Chip, Menu, MenuItem } from '@mui/material';
import { colors } from '@/styles/theme';
import { ArrowDropDown } from '@mui/icons-material';
import AddIcon from '@mui/icons-material/Add';

interface ChipComponentProps {
  menu: string[] | [] | undefined;
  onCreateNetwork: () => void | undefined;
}

const ChipDropdown: React.FC<ChipComponentProps> = ({
  menu,
  onCreateNetwork,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedNetwork, setSelectedNetwork] = useState<string>('');

  const handleMenuClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleMenuClose = () => {
    setAnchorEl(null);
  };
  const handleNetworkSelect = (network: string) => {
    setSelectedNetwork(network);
    handleMenuClose();
  };

  return (
    <>
      <Chip
        label={selectedNetwork || menu?.length ? menu && menu[0] : 'Not added'}
        component="div"
        variant="outlined"
        sx={{ border: `1px solid ${colors.black70}`, color: colors.black70 }}
        onDelete={handleMenuClick}
        deleteIcon={menu?.length ? <ArrowDropDown /> : <AddIcon />}
      />
      <Menu
        anchorEl={anchorEl}
        open={Boolean(anchorEl)}
        onClose={handleMenuClose}
      >
        {menu &&
          menu.map((network) => (
            <MenuItem
              key={network}
              onClick={() => handleNetworkSelect(network)}
            >
              {network}
            </MenuItem>
          ))}
        <MenuItem onClick={onCreateNetwork}>Create Network</MenuItem>
      </Menu>
    </>
  );
};

export default ChipDropdown;
