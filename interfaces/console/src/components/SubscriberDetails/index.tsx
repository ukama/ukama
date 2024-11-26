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
import BillingCycle from './billingCycle';
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
  packageName?: string;
  bundle?: string;
  currentSite?: string;
  handleUpdateSubscriber: (
    subscriberId: string,
    updates: { name?: string; phone?: string },
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
  packageName = '',
  bundle = '',
  currentSite = '',
  simStatusLoading = false,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const open = Boolean(anchorEl);
  const [name, setName] = useState(subscriberInfo?.name || '');
  const [selectedsTab, setSelectedsTab] = useState(0);
  const [email, setEmail] = useState(subscriberInfo?.email || '');
  const [isEditingName, setIsEditingName] = useState(false);
  const [isEditingEmail, setIsEditingEmail] = useState(false);
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
      const updates: { name?: string; email?: string } = {};
      if (name !== subscriberInfo.name) updates.name = name;
      if (email !== subscriberInfo.phone) updates.email = email;

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

  useEffect(() => {
    if (subscriberInfo) {
      setEmail(subscriberInfo.phone);
      setName(subscriberInfo.name);
      setHasChanges(false);
    }
  }, [subscriberInfo]);

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

  const handleMobileChange = (e: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(e.target.value);
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
        <Stack direction="row" sx={{ ml: 1 }} justifyItems={'center'}>
          <Typography variant="h6">{localSubscriberInfo?.name}</Typography>
          <IconButton onClick={handleClick}>
            <MoreVertIcon />
          </IconButton>
          <Menu
            anchorEl={anchorEl}
            open={open}
            onClose={() => setAnchorEl(null)}
          >
            <MenuItem
              onClick={() => handleMenuItemClick('deleteSubscriber')}
              sx={{ color: colors.red }}
              disabled={true}
            >
              Delete subscriber
            </MenuItem>
          </Menu>
        </Stack>
      </DialogTitle>
      <IconButton
        aria-label="close"
        onClick={handleClose}
        sx={{
          position: 'absolute',
          right: 8,
          top: 8,
        }}
      >
        <CloseIcon />
      </IconButton>
      <Box sx={{ px: 4 }}>
        <Box sx={{ borderBottom: 1, borderColor: 'divider' }}>
          <Tabs
            value={selectedsTab}
            onChange={handleTabsChange}
            variant="scrollable"
            scrollButtons="auto"
            indicatorColor="primary"
            textColor="primary"
            sx={{ alignItems: 'flex-start' }}
          >
            <Tab label="INFORMATION" />
            <Tab label="DATA USAGE" />
            <Tab label="SIMS" />
            <Tab label="HISTORY" />
          </Tabs>
        </Box>

        <Box sx={{ p: 3, width: '100%' }}>
          {selectedsTab === 0 && (
            <Box sx={{ position: 'relative', right: 22 }}>
              <Stack spacing={2} direction="column">
                <Box sx={{ position: 'relative' }}>
                  <InputLabel
                    shrink
                    htmlFor="name"
                    sx={{
                      position: 'absolute',
                      top: isEditingName ? '-8px' : '0px',
                      left: isEditingName ? '14px' : '0px',
                      background: isEditingName ? 'white' : 'transparent',
                      padding: isEditingName ? '0 4px' : '0',
                      transition: 'all 0.2s',
                      zIndex: 1,
                    }}
                  >
                    NAME *
                  </InputLabel>
                  {isEditingName ? (
                    <TextField
                      id="name"
                      required
                      value={name}
                      onChange={handleNameChange}
                      variant="outlined"
                      fullWidth
                      InputProps={{
                        endAdornment: (
                          <InputAdornment position="end">
                            <Button
                              variant="text"
                              onClick={() => setIsEditingName(false)}
                            >
                              SAVE
                            </Button>
                          </InputAdornment>
                        ),
                      }}
                    />
                  ) : (
                    <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
                      <Typography variant="body1" sx={{ flexGrow: 1 }}>
                        {name}
                      </Typography>
                      <IconButton
                        onClick={() => setIsEditingName(true)}
                        size="small"
                      >
                        <EditIcon />
                      </IconButton>
                    </Box>
                  )}
                </Box>

                <Box sx={{ position: 'relative' }}>
                  <InputLabel
                    shrink
                    htmlFor="mobileNumber"
                    sx={{
                      position: 'absolute',
                      top: isEditingEmail ? '-8px' : '0px',
                      left: isEditingEmail ? '14px' : '0px',
                      background: isEditingEmail ? 'white' : 'transparent',
                      padding: isEditingEmail ? '0 4px' : '0',
                      transition: 'all 0.2s',
                      zIndex: 1,
                    }}
                  >
                    EMAIL
                  </InputLabel>
                  {isEditingEmail ? (
                    <TextField
                      id="email"
                      value={email}
                      onChange={handleMobileChange}
                      variant="outlined"
                      fullWidth
                      placeholder="@@@-@@@-@@@@"
                      InputProps={{
                        endAdornment: (
                          <InputAdornment position="end">
                            <Button
                              variant="text"
                              onClick={() => setIsEditingEmail(false)}
                            >
                              Done
                            </Button>
                          </InputAdornment>
                        ),
                      }}
                    />
                  ) : (
                    <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
                      <Typography
                        variant="body1"
                        sx={{
                          flexGrow: 1,
                          color: email ? 'inherit' : 'text.secondary',
                        }}
                      >
                        {email || '@@@-@@@-@@@@'}
                      </Typography>
                      <IconButton
                        onClick={() => setIsEditingEmail(true)}
                        size="small"
                      >
                        <EditIcon />
                      </IconButton>
                    </Box>
                  )}
                </Box>
              </Stack>
            </Box>
          )}
          {selectedsTab === 1 && (
            <Box sx={{ position: 'relative', right: 22 }}>
              <DataPlanComponent
                packageName={packageName ?? ''}
                currentSite={currentSite ?? ''}
                bundle={bundle ?? ''}
              />
            </Box>
          )}
          {selectedsTab === 2 && (
            <Box sx={{ position: 'relative', right: 22, top: -8 }}>
              <SimTable
                simData={subscriberInfo?.sim}
                onSimAction={handleSimAction}
                simLoading={simStatusLoading}
              />
            </Box>
          )}
          {selectedsTab === 3 && (
            <Box sx={{ position: 'relative', right: 22, top: -30 }}>
              <BillingCycle />
            </Box>
          )}
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
