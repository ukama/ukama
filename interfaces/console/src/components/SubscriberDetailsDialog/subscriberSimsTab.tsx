/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React, { useState } from 'react';
import {
  Box,
  Table,
  TableHead,
  TableRow,
  TableCell,
  TableBody,
  IconButton,
  Menu,
  MenuItem,
  Tooltip,
} from '@mui/material';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import { colors } from '@/theme';
import { SubscriberSimsDto } from '@/client/graphql/generated';

interface SubscriberSimsTabProps {
  sims?: SubscriberSimsDto[];
  onSimAction?: (action: string, simId: string, additionalData?: any) => void;
}

const SubscriberSimsTab: React.FC<SubscriberSimsTabProps> = ({
  sims,
  onSimAction,
}) => {
  const [simMenuAnchor, setSimMenuAnchor] = useState<{
    el: HTMLElement;
    id: string;
  } | null>(null);

  const handleOpenSimMenu = (
    event: React.MouseEvent<HTMLElement>,
    simId: string,
  ) => {
    setSimMenuAnchor({
      el: event.currentTarget,
      id: simId,
    });
  };

  const handleCloseSimMenu = () => {
    setSimMenuAnchor(null);
  };

  const handleDeleteSimRequest = (simId: string) => {
    const sim = sims?.find((s) => s.id === simId);
    if (sim && onSimAction) {
      const isLastSim = sims?.length === 1;

      // Pass additional data along with the action and simId
      onSimAction('deleteSim', sim.id, {
        iccid: sim.iccid,
        isLastSim: isLastSim,
      });

      setSimMenuAnchor(null);
    }
  };

  const handleToggleSimStatus = (simId: string) => {
    const sim = sims?.find((s) => s.id === simId);
    if (sim && onSimAction) {
      const action =
        sim.status.toLowerCase() === 'active' ? 'deactivateSim' : 'activateSim';

      onSimAction(action, simId);
      setSimMenuAnchor(null);
    }
  };

  const selectedSim = simMenuAnchor?.id
    ? sims?.find((s) => s.id === simMenuAnchor.id)
    : null;

  const canDeleteSim = selectedSim?.status.toLowerCase() === 'inactive';

  return (
    <>
      <Box sx={{ overflowX: 'auto' }}>
        <Table size="small">
          <TableHead>
            <TableRow>
              <TableCell>SIM ICCID</TableCell>
              <TableCell>Type</TableCell>
              <TableCell>Status</TableCell>
              <TableCell align="right">Actions</TableCell>
            </TableRow>
          </TableHead>
          <TableBody>
            {sims && sims.length > 0 ? (
              sims.map((sim) => (
                <TableRow key={sim.id}>
                  <TableCell sx={{ color: colors.black70 }}>
                    {sim.iccid}
                  </TableCell>
                  <TableCell sx={{ color: colors.black70 }}>
                    {sim.isPhysical ? 'pSIM' : 'eSIM'}
                  </TableCell>
                  <TableCell sx={{ color: colors.black70 }}>
                    {sim.status}
                  </TableCell>
                  <TableCell align="right">
                    <IconButton
                      size="small"
                      onClick={(e) => handleOpenSimMenu(e, sim.id)}
                    >
                      <MoreHorizIcon fontSize="small" />
                    </IconButton>
                  </TableCell>
                </TableRow>
              ))
            ) : (
              <TableRow>
                <TableCell colSpan={4} align="center">
                  No SIMs available for this subscriber
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </Box>

      <Menu
        anchorEl={simMenuAnchor?.el}
        open={Boolean(simMenuAnchor)}
        onClose={handleCloseSimMenu}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <MenuItem
          onClick={() => {
            if (simMenuAnchor?.id) {
              handleToggleSimStatus(simMenuAnchor.id);
            }
          }}
        >
          {selectedSim?.status.toLowerCase() === 'active'
            ? 'Deactivate SIM'
            : 'Activate SIM'}
        </MenuItem>

        <Tooltip
          title={
            !canDeleteSim
              ? 'SIM must be deactivated before it can be deleted'
              : ''
          }
          placement="left"
        >
          <span>
            <MenuItem
              onClick={() => {
                if (simMenuAnchor?.id && canDeleteSim) {
                  handleDeleteSimRequest(simMenuAnchor.id);
                }
              }}
              disabled={!canDeleteSim}
              sx={{
                color: canDeleteSim ? colors.error : colors.black38,
              }}
            >
              Delete SIM
            </MenuItem>
          </span>
        </Tooltip>
      </Menu>
    </>
  );
};

export default SubscriberSimsTab;
