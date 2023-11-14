import React, { useState } from 'react';
import TableRow from '@mui/material/TableRow';
import TableCell from '@mui/material/TableCell';
import IconButton from '@mui/material/IconButton';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import Skeleton from '@mui/material/Skeleton';
import { colors } from '@/styles/theme';

interface SimInfo {
  iccid: string;
  msisdn: string;
  isPhysical: boolean;
  status: string;
  id: string;
}

interface Props {
  subscriberInfo: {
    sim: SimInfo[];
  } | null;
  simStatusLoading: boolean;
}

const SubscriberMenu: React.FC<Props> = ({
  subscriberInfo,
  simStatusLoading,
}) => {
  const [simAction, setSimAction] = useState<HTMLElement | null>(null);

  const handleSimAction = (event: React.MouseEvent<HTMLElement>) => {
    setSimAction(event.currentTarget);
  };

  const handleCloseSimAction = () => {
    setSimAction(null);
  };

  const handleSimMenu = (action: string, simId: string) => {
    // Handle the selected action
    console.log(`Performing ${action} on SIM with ID: ${simId}`);
    handleCloseSimAction();
    // Implement your logic for each action
  };

  return (
    <>
      {subscriberInfo &&
        subscriberInfo.sim.map((sim: SimInfo) => (
          <TableRow key={sim.iccid}>
            <TableCell>{sim.msisdn}</TableCell>
            <TableCell>{sim.isPhysical ? 'pSim' : 'eSim'}</TableCell>
            <TableCell>
              {simStatusLoading ? (
                <Skeleton
                  variant="rectangular"
                  width={120}
                  height={24}
                  sx={{ backgroundColor: colors.black10 }}
                />
              ) : (
                sim.status
              )}
            </TableCell>
            <TableCell>
              <IconButton
                aria-controls="menu"
                aria-haspopup="true"
                onClick={handleSimAction}
              >
                <MoreVertIcon sx={{ transform: 'rotate(90deg)' }} />
              </IconButton>
              <Menu
                id="menu"
                anchorEl={simAction}
                open={Boolean(simAction)}
                onClose={handleCloseSimAction}
              >
                <MenuItem
                  onClick={() => handleSimMenu('deactivateSim', sim.id)}
                >
                  Deactivate SIM
                </MenuItem>
                <MenuItem
                  onClick={() => handleSimMenu('deleteSim', sim.id)}
                  sx={{ color: colors.red }}
                >
                  Delete SIM
                </MenuItem>
                <MenuItem
                  onClick={() => handleSimMenu('topUp', sim.id)}
                  sx={{ color: colors.red }}
                >
                  Top up data
                </MenuItem>
              </Menu>
            </TableCell>
          </TableRow>
        ))}
    </>
  );
};

export default SubscriberMenu;
