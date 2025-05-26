/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { SubscriberSimDto } from '@/client/graphql/generated';
import colors from '@/theme/colors';
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
} from '@mui/material';
import { styled } from '@mui/material/styles';
import React, { useState } from 'react';
import LoadingWrapper from '../LoadingWrapper';

interface SimTableProps {
  simData: SubscriberSimDto[];
  onSimAction: (action: string, id: string) => void;
  simLoading: boolean;
}

const StyledTableContainer = styled(TableContainer)({
  boxShadow: 'none',
  border: 'none',
});

const StyledTableCellHeader = styled(TableCell)({
  color: `${colors.black}`,
});

const StyledTableCellBody = styled(TableCell)({
  color: `${colors.black70}`,
});

const SimTable: React.FC<SimTableProps> = ({
  simData,
  onSimAction,
  simLoading,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [selectedSim, setSelectedSim] = useState<string | null>(null);

  const handleMenuClick = (
    event: React.MouseEvent<HTMLElement>,
    id: string,
  ) => {
    setAnchorEl(event.currentTarget);
    setSelectedSim(id);
  };

  const handleClose = () => {
    setAnchorEl(null);
    setSelectedSim(null);
  };

  const handleSimMenu = (action: string) => {
    if (selectedSim) {
      onSimAction(action, selectedSim);
    }

    handleClose();
  };

  return (
    <LoadingWrapper height={120} isLoading={simLoading}>
      <StyledTableContainer>
        <Table>
          <TableHead>
            <TableRow>
              <StyledTableCellHeader>SIM ICCID</StyledTableCellHeader>
              <StyledTableCellHeader>Type</StyledTableCellHeader>
              <StyledTableCellHeader>Status</StyledTableCellHeader>
              <StyledTableCellHeader>
                <></>
              </StyledTableCellHeader>
            </TableRow>
          </TableHead>

          <TableBody>
            {simData &&
              simData.map((sim) => (
                <TableRow key={sim.iccid}>
                  <StyledTableCellBody>{sim.iccid}</StyledTableCellBody>
                  <StyledTableCellBody>
                    {sim.isPhysical ? 'pSim' : 'eSim'}
                  </StyledTableCellBody>
                  <StyledTableCellBody>
                    {sim.status.charAt(0).toUpperCase() + sim.status.slice(1)}
                  </StyledTableCellBody>
                  <TableCell>
                    <IconButton
                      aria-controls="menu"
                      aria-haspopup="true"
                      onClick={(event) => handleMenuClick(event, sim.id)}
                    >
                      <MoreVertIcon sx={{ transform: 'rotate(90deg)' }} />
                    </IconButton>
                    <Menu
                      id="menu"
                      anchorEl={anchorEl}
                      open={Boolean(anchorEl) && selectedSim === sim.id}
                      onClose={handleClose}
                    >
                      <MenuItem
                        onClick={() => {
                          handleSimMenu(
                            sim.status === 'inactive'
                              ? 'activateSim'
                              : 'deactivateSim',
                          );
                        }}
                      >
                        {sim.status === 'active' ? 'Deactivate' : 'Activate'}
                      </MenuItem>
                    </Menu>
                  </TableCell>
                </TableRow>
              ))}
          </TableBody>
        </Table>
      </StyledTableContainer>
    </LoadingWrapper>
  );
};

export default SimTable;
