/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { colors } from '@/theme';
import React, { useEffect, useState, useCallback } from 'react';

import TabsComponent from '@/components/TabsComponent';
import CloseIcon from '@mui/icons-material/Close';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import EditIcon from '@mui/icons-material/Edit';
import SaveIcon from '@mui/icons-material/Save';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogTitle,
  Divider,
  FormControl,
  Grid,
  IconButton,
  InputAdornment,
  InputLabel,
  Menu,
  MenuItem,
  OutlinedInput,
  Stack,
  Tab,
  Tabs,
  TextField,
  Typography,
} from '@mui/material';
import SimTable from './SimInfoTab';
import BillingCycle from './billingCycle';
import DataPlanComponent from './dataPlanInfo';
import UserInfo from './userInfo';

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
  loading,
  currentSite = '',
  simStatusLoading = false,
}) => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);
  const [firstName, setFirstName] = useState(subscriberInfo?.firstName || '');
  const [selectedsTab, setSelectedsTab] = useState(0);
  const [mobileNumber, setMobileNumber] = useState('');
  const [isEditingFirstName, setIsEditingFirstName] = useState(false);
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
    handleUpdateSubscriber(subscriberInfo.uuid, firstName, mobileNumber);
  }, [mobileNumber, firstName, handleUpdateSubscriber, subscriberInfo]);

  useEffect(() => {
    if (subscriberInfo) {
      setMobileNumber(subscriberInfo.phone);
      setFirstName(subscriberInfo.firstName);
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
          <Typography variant="h6">{subscriberInfo?.firstName}</Typography>
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
              {' '}
              {/* Add padding here */}
              <Stack spacing={2} direction="column">
                {/* First Name Field */}
                <Box sx={{ position: 'relative' }}>
                  <InputLabel
                    shrink
                    htmlFor="firstName"
                    sx={{
                      position: 'absolute',
                      top: isEditingFirstName ? '-8px' : '0px',
                      left: isEditingFirstName ? '14px' : '0px',
                      background: isEditingFirstName ? 'white' : 'transparent',
                      padding: isEditingFirstName ? '0 4px' : '0',
                      transition: 'all 0.2s',
                      zIndex: 1,
                    }}
                  >
                    FIRST NAME *
                  </InputLabel>
                  {isEditingFirstName ? (
                    <TextField
                      id="firstName"
                      required
                      value={firstName}
                      onChange={(e) => setFirstName(e.target.value)}
                      variant="outlined"
                      fullWidth
                      InputProps={{
                        endAdornment: (
                          <InputAdornment position="end">
                            <Button
                              variant="text"
                              onClick={() => {
                                setIsEditingFirstName(false),
                                  handleSaveSubscriber;
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
                        {firstName}
                      </Typography>
                      <IconButton
                        onClick={() => setIsEditingFirstName(true)}
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
