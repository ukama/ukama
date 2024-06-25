import LoadingWrapper from '@/components/LoadingWrapper';
import { colors } from '@/theme';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import {
  IconButton,
  Menu,
  MenuItem,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import React from 'react';

interface SimInfoProps {
  selectedTab: number;
  subscriberInfo: any;
  simStatusLoading: boolean;
  handleSimAction: any;
  simAction: any;
  handleCloseSimAction: () => void;
  handleSimMenu: Function;
}

const SimInfoTab: React.FC<SimInfoProps> = ({
  selectedTab,
  subscriberInfo,
  simStatusLoading,
  handleSimAction,
  simAction,
  handleCloseSimAction,
  handleSimMenu,
}) => (
  <Typography component="div" role="tabpanel" hidden={selectedTab !== 2}>
    <TableContainer>
      <Table>
        <TableHead>
          <TableRow>
            <TableCell>
              <strong style={{ fontWeight: 'bold' }}> SIM ICCID</strong>
            </TableCell>
            <TableCell>
              <strong style={{ fontWeight: 'bold' }}> Type</strong>
            </TableCell>
            <TableCell>
              <strong style={{ fontWeight: 'bold' }}> Status</strong>
            </TableCell>
            <TableCell></TableCell>
          </TableRow>
        </TableHead>
        <TableBody>
          <LoadingWrapper
            radius="small"
            width={'100%'}
            isLoading={simStatusLoading}
            cstyle={{
              overflow: 'auto',
              backgroundColor: false ? colors.white : 'transparent',
            }}
          >
            {subscriberInfo &&
              subscriberInfo.sim.map((sim: any) => (
                <TableRow key={sim.iccid}>
                  <TableCell>{sim.msisdn}</TableCell>
                  <TableCell>{sim.isPhysical ? 'pSim' : 'eSim'}</TableCell>
                  <TableCell>{sim.status}</TableCell>
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
          </LoadingWrapper>
        </TableBody>
      </Table>
    </TableContainer>
  </Typography>
);

export default SimInfoTab;
