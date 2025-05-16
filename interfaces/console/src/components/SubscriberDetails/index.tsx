/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { colors } from '@/theme';
import CloseIcon from '@mui/icons-material/Close';
import EditIcon from '@mui/icons-material/Edit';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogTitle,
  IconButton,
  InputAdornment,
  InputLabel,
  Menu,
  MenuItem,
  Stack,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@mui/material';
import React, { useCallback, useEffect, useState } from 'react';
import SimTable from './SimInfoTab';
import DataPlanComponent from './dataPlanInfo';

interface SubscriberProps {
  ishowSubscriberDetails: boolean;
  handleClose: () => void;
  subscriberInfo: any;
  handleSimActionOption: (
    action: string,
    simId: string,
    subscriberId: string,
  ) => void;
  handleUpdateSubscriber: (
    subscriberId: string,
    updates: { name?: string; email?: string },
  ) => void;
  handleDeleteSubscriber: (action: string, subscriberId: string) => void;
  loading: boolean;
  simStatusLoading: boolean;
}

const SubscriberDetails: React.FC<SubscriberProps> = ({
  ishowSubscriberDetails,
  subscriberInfo,
  handleUpdateSubscriber,
  handleClose,
  handleDeleteSubscriber,
  handleSimActionOption,
  simStatusLoading = false,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const [name, setName] = useState(subscriberInfo?.name || '');
  const [selectedsTab, setSelectedsTab] = useState(0);
  const [email, setEmail] = useState(subscriberInfo?.email || '');
  const [isEditingName, setIsEditingName] = useState(false);
  const [hasChanges, setHasChanges] = useState(false);
  const [localSubscriberInfo, setLocalSubscriberInfo] =
    useState<any>(subscriberInfo);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCloseItem = () => {
    setAnchorEl(null);
  };

  const handleMenuItemClick = useCallback(
    (action: string) => {
      handleCloseItem();
      handleDeleteSubscriber(action, subscriberInfo.uuid);
    },
    [subscriberInfo, handleCloseItem, handleDeleteSubscriber],
  );

  useEffect(() => {
    if (subscriberInfo) {
      setLocalSubscriberInfo(subscriberInfo);
      setName(subscriberInfo.name || '');
      setEmail(subscriberInfo.email || '');
      setHasChanges(false);
    }
  }, [subscriberInfo]);

  const handleSaveSubscriber = useCallback(() => {
    if (hasChanges) {
      const updates: { name?: string } = {};
      if (name !== subscriberInfo.name) updates.name = name;
      handleUpdateSubscriber(subscriberInfo.uuid, updates);
    }
    handleClose();
  }, [
    name,
    email,
    hasChanges,
    handleUpdateSubscriber,
    subscriberInfo,
    handleClose,
  ]);

  const handleTabsChange = (_: React.SyntheticEvent, newValue: number) => {
    setSelectedsTab(newValue);
  };

  const handleSimAction = (action: string, iccid: string) => {
    if (action === 'deactivateSim' || action === 'activateSim') {
      handleSimActionOption(action, iccid, subscriberInfo.uuid);
      handleCloseItem();
    }
  };

  const handleNameChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setName(e.target.value);
    setHasChanges(true);
  };

  return (
    <Dialog
      open={ishowSubscriberDetails}
      onClose={() => {
        handleCloseItem();
        handleClose();
      }}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle id="alert-dialog-title">
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="center"
        >
          <Stack direction="row" spacing={1} alignItems={'center'}>
            <Typography variant="h6">{localSubscriberInfo?.name}</Typography>
            <IconButton onClick={handleClick}>
              <MoreVertIcon />
            </IconButton>
          </Stack>
          <Stack direction="row" alignItems="center">
            <IconButton
              aria-label="close"
              onClick={handleClose}
              sx={{
                position: 'relative',
                right: 0,
              }}
            >
              <CloseIcon />
            </IconButton>
          </Stack>
        </Stack>
      </DialogTitle>
      <Menu anchorEl={anchorEl} open={open} onClose={handleCloseItem}>
        <MenuItem
          onClick={() => handleMenuItemClick('deleteSubscriber')}
          sx={{ color: colors.red }}
        >
          Delete subscriber
        </MenuItem>
      </Menu>
      <Box sx={{ px: 4, py: 2 }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs
            value={selectedsTab}
            onChange={handleTabsChange}
            variant="scrollable"
            scrollButtons="auto"
            indicatorColor="primary"
            textColor="primary"
          >
            <Tab label="INFORMATION" />
            <Tab label="DATA USAGE" />
            <Tab label="SIMS" />
            {/* <Tab label="HISTORY" /> */}
          </Tabs>
        </Box>

        <Box sx={{ pt: 3 }}>
          {selectedsTab === 0 && (
            <Box>
              <Stack spacing={3} direction="column">
                <Box>
                  <TextField
                    id="name"
                    required
                    value={name}
                    label="NAME"
                    onChange={handleNameChange}
                    variant={isEditingName ? 'outlined' : 'standard'}
                    fullWidth
                    InputProps={{
                      disableUnderline: !isEditingName,
                      endAdornment: isEditingName ? (
                        <InputAdornment position="end">
                          <Button
                            variant="text"
                            onClick={() => {
                              setIsEditingName(false);
                              setHasChanges(true);
                            }}
                          >
                            SAVE
                          </Button>
                        </InputAdornment>
                      ) : (
                        <IconButton
                          onClick={() => setIsEditingName(true)}
                          size="small"
                        >
                          <EditIcon />
                        </IconButton>
                      ),
                      style: {
                        height: '53px',
                      },
                    }}
                  />
                </Box>

                <Box>
                  <InputLabel
                    shrink
                    htmlFor="email"
                    sx={{
                      transition: 'all 0.2s',
                      zIndex: 1,
                    }}
                  >
                    EMAIL
                  </InputLabel>

                  <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
                    <Typography
                      variant="body1"
                      sx={{
                        flexGrow: 1,
                        color: email ? 'inherit' : 'text.secondary',
                      }}
                    >
                      {email}
                    </Typography>
                  </Box>
                </Box>
              </Stack>
            </Box>
          )}
          {selectedsTab === 1 && (
            <Box>
              <DataPlanComponent
                packageName={subscriberInfo?.dataPlan ?? ''}
                currentSite={'-'}
                bundle={subscriberInfo?.dataUsage}
              />
            </Box>
          )}
          {selectedsTab === 2 && (
            <Box>
              <SimTable
                simData={subscriberInfo?.sim}
                onSimAction={handleSimAction}
                simLoading={simStatusLoading}
              />
            </Box>
          )}
          {/* TODO: Need more discussion
            {selectedsTab === 3 && (
            <Box>
              <BillingCycle />
            </Box>
          )} */}
        </Box>
      </Box>
      <DialogActions>
        {selectedsTab === 2 || selectedsTab === 3 || selectedsTab === 1 ? (
          <Button variant="contained" onClick={handleClose}>
            CLOSE
          </Button>
        ) : (
          <>
            <Button variant="text" onClick={handleClose}>
              CANCEL
            </Button>
            <Button variant="contained" onClick={handleSaveSubscriber}>
              DONE
            </Button>
          </>
        )}
      </DialogActions>
    </Dialog>
  );
};

export default SubscriberDetails;
