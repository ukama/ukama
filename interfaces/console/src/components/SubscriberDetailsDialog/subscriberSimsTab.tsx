import React from 'react';
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
} from '@mui/material';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
import { colors } from '@/theme';
import { SubscriberSimsDto } from '@/client/graphql/generated';

interface SubscriberSimsTabProps {
  sims?: SubscriberSimsDto[];
  onSimAction?: (action: string, simId: string) => void;
  onDeleteSim?: (simId: string) => void;
}

const SubscriberSimsTab: React.FC<SubscriberSimsTabProps> = ({
  sims,
  onSimAction,
  onDeleteSim,
}) => {
  const [simMenuAnchor, setSimMenuAnchor] = React.useState<{
    el: HTMLElement;
    id: string;
  } | null>(null);

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
                      onClick={(e) => {
                        setSimMenuAnchor({
                          el: e.currentTarget,
                          id: sim.id,
                        });
                      }}
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

      {/* Menu for SIM actions */}
      <Menu
        anchorEl={simMenuAnchor?.el}
        open={Boolean(simMenuAnchor)}
        onClose={() => setSimMenuAnchor(null)}
        anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
        transformOrigin={{ vertical: 'top', horizontal: 'right' }}
      >
        <MenuItem
          onClick={() => {
            const sim = sims?.find((s) => s.id === simMenuAnchor?.id);
            if (sim && simMenuAnchor) {
              const action =
                sim.status.toLowerCase() === 'active'
                  ? 'deactivateSim'
                  : 'activateSim';
              onSimAction?.(action, simMenuAnchor.id);
              setSimMenuAnchor(null);
            }
          }}
        >
          {sims
            ?.find((s) => s.id === simMenuAnchor?.id)
            ?.status.toLowerCase() === 'active'
            ? 'Deactivate SIM'
            : 'Activate SIM'}
        </MenuItem>
        <MenuItem
          onClick={() => {
            if (simMenuAnchor?.id && onDeleteSim) {
              onDeleteSim(simMenuAnchor.id);
              setSimMenuAnchor(null);
            }
          }}
          sx={{ color: colors.error }}
        >
          Delete SIM
        </MenuItem>
      </Menu>
    </>
  );
};

export default SubscriberSimsTab;
