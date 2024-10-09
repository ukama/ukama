/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { colors } from '@/theme';
import React, { useEffect, useState, useCallback } from 'react';
import CloseIcon from '@mui/icons-material/Close';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import EditIcon from '@mui/icons-material/Edit';
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
import SimTable from './SimInfoTab';
import BillingCycle from './billingCycle';
import DataPlanComponent from './dataPlanInfo';

interface SubscriberProps {
  ishowSubscriberDetails: boolean;
  subscriberId: string;
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
    mobileNumber: string,
    firstName: string,
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
  const [name, setName] = useState(subscriberInfo?.name || '');
  const [selectedsTab, setSelectedsTab] = useState(0);
  const [mobileNumber, setMobileNumber] = useState('');
  const [isEditingName, setIsEditingName] = useState(false);
  const [isEditingMobile, setIsEditingMobile] = useState(false);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCloseItem = useCallback(() => setAnchorEl(null), []);

  const handleMenuItemClick = useCallback(
    (action: string) => {
      handleCloseItem();
      handleDeleteSubscriber(action, subscriberInfo.uuid);
    },
    [subscriberInfo, handleCloseItem, handleDeleteSubscriber],
  );

  const handleSaveSubscriber = useCallback(() => {
    handleUpdateSubscriber(subscriberInfo.uuid, name, mobileNumber);
    handleClose();
  }, [mobileNumber, name, handleUpdateSubscriber, subscriberInfo]);

  useEffect(() => {
    if (subscriberInfo) {
      setMobileNumber(subscriberInfo.phone);
      setName(subscriberInfo.name);
    }
  }, [subscriberInfo]);

  const handleTabsChange = (_: React.SyntheticEvent, newValue: number) => {
    setSelectedsTab(newValue);
  };
  const handleSimAction = (action: string, iccid: string) => {
    if (action === 'deactivateSim') {
      handleSimActionOption(action, iccid, subscriberInfo.uuid);
    }
    if (action === 'deleteSim') {
      handleSimActionOption(action, iccid, subscriberInfo.uuid);
    }
  };
  return (
    <Dialog
      open={ishowSubscriberDetails}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle id="alert-dialog-title">
        <Stack direction="row" spacing={1} alignItems="center">
          <Typography variant="h6">{subscriberInfo?.name}</Typography>
          <IconButton
            aria-controls="menu"
            aria-haspopup="true"
            onClick={handleClick}
          >
            <MoreVertIcon />
          </IconButton>
          <Menu
            id="menu"
            anchorEl={anchorEl}
            open={Boolean(anchorEl)}
            onClose={handleCloseItem}
          >
            <MenuItem onClick={() => handleMenuItemClick('pauseService')}>
              Pause service
            </MenuItem>
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
                    htmlFor="firstName"
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
                    FIRST NAME *
                  </InputLabel>
                  {isEditingName ? (
                    <TextField
                      id="name"
                      required
                      value={name}
                      onChange={(e) => setName(e.target.value)}
                      variant="outlined"
                      fullWidth
                      InputProps={{
                        endAdornment: (
                          <InputAdornment position="end">
                            <Button
                              variant="text"
                              onClick={() => {
                                setIsEditingName(false), handleSaveSubscriber;
                              }}
                            >
                              Save
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
                      top: isEditingMobile ? '-8px' : '0px',
                      left: isEditingMobile ? '14px' : '0px',
                      background: isEditingMobile ? 'white' : 'transparent',
                      padding: isEditingMobile ? '0 4px' : '0',
                      transition: 'all 0.2s',
                      zIndex: 1,
                    }}
                  >
                    MOBILE NUMBER
                  </InputLabel>
                  {isEditingMobile ? (
                    <TextField
                      id="mobileNumber"
                      value={mobileNumber}
                      onChange={(e) => setMobileNumber(e.target.value)}
                      variant="outlined"
                      fullWidth
                      InputProps={{
                        endAdornment: (
                          <InputAdornment position="end">
                            <Button
                              variant="text"
                              onClick={() => setIsEditingMobile(false)}
                            >
                              Save
                            </Button>
                          </InputAdornment>
                        ),
                      }}
                    />
                  ) : (
                    <Box sx={{ display: 'flex', alignItems: 'center', mt: 2 }}>
                      <Typography variant="body1" sx={{ flexGrow: 1 }}>
                        {mobileNumber}
                      </Typography>
                      <IconButton
                        onClick={() => setIsEditingMobile(true)}
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
        <Button variant="text" onClick={handleClose}>
          CANCEL
        </Button>
        <Button variant="contained" onClick={handleSaveSubscriber}>
          DONE
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default SubscriberDetails;
