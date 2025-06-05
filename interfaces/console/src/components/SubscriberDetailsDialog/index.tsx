/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import React from 'react';
import {
  Dialog,
  DialogTitle,
  DialogContent,
  DialogActions,
  IconButton,
  Typography,
  Tabs,
  Tab,
  Box,
  Button,
  Menu,
  MenuItem,
} from '@mui/material';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import CloseIcon from '@mui/icons-material/Close';
import { colors } from '@/theme';
import { PackagesResDto, SubscriberSimsDto } from '@/client/graphql/generated';
import SubscriberInfoTab from './subscriberInfoTab';
import SubscriberDataPlansTab from './subscriberDataPlansTab';
import SubscriberSimsTab from './subscriberSimsTab';
import SubscriberHistoryTab from './subscriberHistoryTab';

interface Subscriber {
  id: string;
  firstName: string;
  email: string;
}

interface SubscriberDialogProps {
  open: boolean;
  onClose: () => void;
  subscriber: Subscriber;
  onUpdateSubscriber: (updates: { name?: string; email?: string }) => void;
  onDeleteSubscriber: (id: string) => void;
  onTopUpPlan: (subscriberId: string) => void;
  sims?: SubscriberSimsDto[];
  onSimAction?: (action: string, simId: string) => void;
  packageHistories?: any[];
  packagesData?: PackagesResDto;
  loadingPackageHistories?: boolean;
  dataUsage: string;
  currencySymbol?: string;
}

function TabPanel({ children, value, index }: any) {
  return (
    <div role="tabpanel" hidden={value !== index}>
      {value === index && <Box>{children}</Box>}
    </div>
  );
}

function a11yProps(index: number) {
  return {
    id: `subscriber-tab-${index}`,
    'aria-controls': `subscriber-tabpanel-${index}`,
  };
}

const SubscriberDetailsDialog: React.FC<SubscriberDialogProps> = ({
  open,
  onClose,
  subscriber,
  onUpdateSubscriber,
  onDeleteSubscriber,
  onTopUpPlan,
  sims,
  onSimAction,
  packageHistories,
  packagesData,
  loadingPackageHistories,
  dataUsage,
  currencySymbol,
}) => {
  const [tabIndex, setTabIndex] = React.useState(0);
  const [menuAnchor, setMenuAnchor] = React.useState<null | HTMLElement>(null);
  const [isEditing, setIsEditing] = React.useState(false);
  const [pendingChanges, setPendingChanges] = React.useState<{
    name?: string;
    email?: string;
  }>({});

  const handleTabChange = (_: React.SyntheticEvent, newIndex: number) => {
    if (isEditing) return;
    setTabIndex(newIndex);
  };

  const openMenu = (e: React.MouseEvent<HTMLElement>) =>
    setMenuAnchor(e.currentTarget);
  const closeMenu = () => setMenuAnchor(null);

  const handleDelete = () => {
    onDeleteSubscriber(subscriber.id);
    closeMenu();
  };

  const handleDoneClick = () => {
    if (Object.keys(pendingChanges).length > 0) {
      onUpdateSubscriber(pendingChanges);
    }
    setPendingChanges({});
    setIsEditing(false);
    onClose();
  };

  const handleCloseClick = () => {
    setMenuAnchor(null);
    if (isEditing) {
      setPendingChanges({});
      setIsEditing(false);
    }
    setMenuAnchor(null);
    onClose();
  };

  const updatePendingChanges = (updates: { name?: string; email?: string }) => {
    setPendingChanges((prev) => ({ ...prev, ...updates }));
  };

  return (
    <Dialog open={open} onClose={handleCloseClick} fullWidth maxWidth="sm">
      <DialogTitle sx={{ px: 3, pt: 2, pb: 1, position: 'relative' }}>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Typography variant="h6">{subscriber.firstName}</Typography>
          <IconButton size="small" onClick={openMenu} sx={{ ml: 0.5 }}>
            <MoreVertIcon fontSize="medium" />
          </IconButton>
        </Box>
        <IconButton
          aria-label="close"
          onClick={handleCloseClick}
          sx={{
            position: 'absolute',
            right: 12,
            top: 12,
          }}
        >
          <CloseIcon />
        </IconButton>
        <Menu
          anchorEl={menuAnchor}
          open={Boolean(menuAnchor)}
          onClose={() => {
            closeMenu();
            setMenuAnchor(null);
          }}
          anchorOrigin={{ vertical: 'bottom', horizontal: 'right' }}
          transformOrigin={{ vertical: 'top', horizontal: 'right' }}
          MenuListProps={{ dense: true }}
        >
          <MenuItem onClick={() => onTopUpPlan(subscriber.id)}>
            Top Up Data
          </MenuItem>
          <MenuItem onClick={handleDelete} sx={{ color: colors.error }}>
            Delete Subscriber
          </MenuItem>
        </Menu>
      </DialogTitle>

      <Tabs
        value={tabIndex}
        onChange={handleTabChange}
        indicatorColor="primary"
        textColor="primary"
        variant="fullWidth"
        sx={{
          px: 3,
          position: 'relative',
          '&::after': {
            content: '""',
            position: 'absolute',
            bottom: 0,
            left: '24px',
            right: '24px',
            height: '1px',
            backgroundColor: colors.black10,
          },
        }}
      >
        {/* //TODO billing cycle sitll under discussion */}
        {['INFORMATION', 'DATA PLANS', 'SIMS', 'HISTORY'].map((label, idx) => (
          <Tab
            key={label}
            label={label}
            {...a11yProps(idx)}
            disabled={isEditing && tabIndex !== idx}
          />
        ))}
      </Tabs>

      <DialogContent>
        <TabPanel value={tabIndex} index={0}>
          <SubscriberInfoTab
            subscriber={subscriber}
            onUpdateSubscriber={updatePendingChanges}
            isEditing={isEditing}
            setIsEditing={setIsEditing}
          />
        </TabPanel>
        <TabPanel value={tabIndex} index={1}>
          <SubscriberDataPlansTab
            packageHistories={packageHistories}
            packagesData={packagesData}
            dataUsage={dataUsage}
            currencySymbol={currencySymbol}
          />
        </TabPanel>
        <TabPanel value={tabIndex} index={2}>
          <SubscriberSimsTab sims={sims} onSimAction={onSimAction} />
        </TabPanel>
        <TabPanel value={tabIndex} index={3}>
          <SubscriberHistoryTab
            packageHistories={packageHistories}
            packagesData={packagesData}
            loadingPackageHistories={loadingPackageHistories}
          />
        </TabPanel>
      </DialogContent>

      <DialogActions sx={{ px: 3, pb: 2 }}>
        <Button
          onClick={isEditing ? handleDoneClick : handleCloseClick}
          variant="contained"
        >
          {isEditing ? 'DONE' : 'CLOSE'}
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default SubscriberDetailsDialog;
