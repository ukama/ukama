import { colors } from '@/theme';
import React, { useEffect, useState } from 'react';

import TabsComponent from '@/components/TabsComponent';
import CloseIcon from '@mui/icons-material/Close';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import {
  Box,
  Button,
  Dialog,
  DialogActions,
  DialogTitle,
  IconButton,
  Menu,
  MenuItem,
  Stack,
  Typography,
} from '@mui/material';
import SimInfoTab from './SimInfoTab';
import BillingCycle from './billingCycle';
import DataPlanComponent from './dataPlanInfo';
import UserInfo from './userInfo';
interface SubscriberProps {
  onCancel: () => void;
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
    email: string,
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
  packageName,
  bundle,
  loading,
  currentSite,
  simStatusLoading = false,
}) => {
  const [selectedTab, setSelectedTab] = useState(0);
  const [simAction, setSimAction] = useState<any>();
  const [subscriberLoading, setSubscriberLoading] = useState(true);
  const [email, setEmail] = useState<string>('');
  const [onEditEmail, setOnEditEmail] = useState<boolean>(false);
  const [onEditName, setOnEditName] = useState<boolean>(false);
  const [firstName, setFirstName] = useState<string>('');
  const handleTabChange = (event: React.SyntheticEvent, newValue: number) => {
    setSelectedTab(newValue);
  };
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleClick = (event: React.MouseEvent<HTMLButtonElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCloseItem = () => {
    setAnchorEl(null);
  };

  const handleMenuItemClick = (action: string, subscriberId: string) => {
    handleCloseItem();

    handleDeleteSubscriber(action, subscriberId);
  };
  const handleSimAction = (event: React.MouseEvent<HTMLButtonElement>) => {
    setSimAction(event.currentTarget);
  };
  const handleCloseSimAction = () => {
    setSimAction(null);
  };
  const handleSimMenu = (action: string, simId: string) => {
    handleCloseSimAction();

    handleSimActionOption(action, simId, subscriberInfo.uuid);
  };
  useEffect(() => {
    if (subscriberInfo) {
      setSubscriberLoading(false);
      setEmail(subscriberInfo.email);
      setFirstName(subscriberInfo.firstName);
    }
  }, [subscriberInfo]);

  const handleSimEdit = (event: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(event.target.value);
  };
  const handleEditName = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFirstName(event.target.value);
  };
  const handleSaveSubscriber = () => {
    setOnEditEmail(false);
    setOnEditName(false);
    handleUpdateSubscriber(subscriberInfo.uuid, email, firstName);
  };

  return (
    <Dialog
      open={ishowSubscriberDetails}
      onClose={handleClose}
      maxWidth="sm"
      fullWidth
    >
      <DialogTitle id="alert-dialog-title">
        <Stack direction="row" spacing={2} alignItems={'center'}>
          <Typography variant="h6">
            {subscriberInfo && subscriberInfo.firstName}
          </Typography>
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
            <MenuItem
              onClick={() =>
                handleMenuItemClick('pauseService', subscriberInfo.uuid)
              }
            >
              Pause service
            </MenuItem>
            <MenuItem
              onClick={() =>
                handleMenuItemClick('deleteSubscriber', subscriberInfo.uuid)
              }
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
      <Box sx={{ width: '100%' }}>
        <Box sx={{ p: 2 }}>
          <TabsComponent
            selectedTab={selectedTab}
            handleTabChange={handleTabChange}
          />
        </Box>

        <Box sx={{ pl: 2 }}>
          <Box component="div" role="tabpanel" hidden={selectedTab !== 0}>
            <UserInfo
              subscriberLoading={subscriberLoading}
              onEditName={onEditName}
              firstName={firstName}
              handleEditName={() => setOnEditName(!onEditName)}
              onEditEmail={onEditEmail}
              email={email}
              handleSimEdit={() => setOnEditEmail(!onEditEmail)}
              setOnEditName={setOnEditName}
              setOnEditEmail={setOnEditEmail}
            />
          </Box>

          <Box component="div" role="tabpanel" hidden={selectedTab !== 1}>
            <DataPlanComponent
              packageName={packageName || ''}
              currentSite={currentSite || ''}
              bundle={bundle || ''}
            />
          </Box>

          <Box component="div" role="tabpanel" hidden={selectedTab !== 2}>
            <SimInfoTab
              selectedTab={selectedTab}
              subscriberInfo={subscriberInfo}
              simStatusLoading={simStatusLoading}
              handleSimAction={handleSimAction}
              simAction={simAction}
              handleCloseSimAction={handleCloseSimAction}
              handleSimMenu={handleSimMenu}
            />
          </Box>

          <Box component="div" role="tabpanel" hidden={selectedTab !== 3}>
            <BillingCycle />
          </Box>
        </Box>
      </Box>
      <DialogActions>
        {onEditName || onEditEmail ? (
          <Button
            variant="contained"
            onClick={handleSaveSubscriber}
            disabled={loading}
          >
            SAVE
          </Button>
        ) : (
          <Button variant="contained" onClick={handleClose}>
            CLOSE
          </Button>
        )}
      </DialogActions>
    </Dialog>
  );
};

export default SubscriberDetails;
