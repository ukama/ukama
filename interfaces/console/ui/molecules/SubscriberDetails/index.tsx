import React, { useState, useEffect } from 'react';
import { colors } from '@/styles/theme';

import {
  Dialog,
  Button,
  DialogActions,
  Skeleton,
  Box,
  Menu,
  MenuItem,
  Tabs,
  Tab,
  Stack,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  TextField,
  Typography,
  DialogTitle,
  IconButton,
} from '@mui/material';
import CloseIcon from '@mui/icons-material/Close';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import EditIcon from '@mui/icons-material/Edit';

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
  handleUpdateSubscriber: (subscriberId: string, email: string) => void;
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
  simStatusLoading = false,
}) => {
  const [selectedTab, setSelectedTab] = useState(0);
  const [simAction, setSimAction] = useState<any>();
  const [SubscriberLoading, setSubscriberLoading] = useState(true);
  const [email, setEmail] = useState<string>('');
  const [onEditEmail, setOnEditEmail] = useState<boolean>(false);
  const [onEditName, setOnEditName] = useState<boolean>(false);
  const [fistName, setFistName] = useState<string>('');
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
      setFistName(subscriberInfo.firstName);
    }
  }, [subscriberInfo]);

  const handleSimEdit = (event: React.ChangeEvent<HTMLInputElement>) => {
    setEmail(event.target.value);
  };
  const handleEditName = (event: React.ChangeEvent<HTMLInputElement>) => {
    setFistName(event.target.value);
  };
  const handleSaveSubscriber = () => {
    setOnEditEmail(false);
    setOnEditName(false);
    handleUpdateSubscriber(subscriberInfo.uuid, email);
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
          <Tabs
            value={selectedTab}
            onChange={handleTabChange}
            aria-label="tabs menu"
          >
            <Tab label="Information" />
            <Tab label="Data Usage" />
            <Tab label="SIMs" />
            <Tab label="History" />
          </Tabs>
        </Box>

        <Box sx={{ pl: 2 }}>
          <Typography
            component="div"
            role="tabpanel"
            hidden={selectedTab !== 0}
          >
            <Stack direction="column" spacing={2}>
              <Typography variant="body1" sx={{ color: colors.black }}>
                Name
              </Typography>
              <Stack
                direction="row"
                spacing={2}
                alignItems={'center'}
                justifyContent={'space-between'}
                sx={{ pr: 2 }}
              >
                {SubscriberLoading ? (
                  <Skeleton
                    variant="rectangular"
                    width={120}
                    height={24}
                    sx={{ backgroundColor: colors.black10 }}
                  />
                ) : (
                  <TextField
                    id="outlined-basic"
                    value={fistName}
                    variant="standard"
                    disabled={!onEditName}
                    size="small"
                    onChange={handleEditName}
                    sx={{ width: '100%' }}
                  />
                )}
                <IconButton
                  size="small"
                  color="primary"
                  onClick={() => setOnEditName(!onEditName)}
                >
                  <EditIcon fontSize="small" />
                </IconButton>
              </Stack>

              <Typography variant="body1" sx={{ color: colors.black }}>
                Email
              </Typography>
              <Stack
                direction="row"
                spacing={2}
                alignItems={'center'}
                justifyContent={'space-between'}
                sx={{ pr: 2 }}
              >
                {SubscriberLoading ? (
                  <Skeleton
                    variant="rectangular"
                    width={120}
                    height={24}
                    sx={{ backgroundColor: colors.black10 }}
                  />
                ) : (
                  <TextField
                    id="outlined-basic"
                    value={email}
                    variant="standard"
                    disabled={!onEditEmail}
                    size="small"
                    onChange={handleSimEdit}
                    sx={{ width: '100%' }}
                  />
                )}
                <IconButton
                  size="small"
                  color="primary"
                  onClick={() => {
                    setOnEditEmail(!onEditEmail);
                  }}
                >
                  <EditIcon fontSize="small" />
                </IconButton>
              </Stack>
            </Stack>
          </Typography>
          <Typography
            component="div"
            role="tabpanel"
            hidden={selectedTab !== 1}
          >
            <Stack direction="column" spacing={2}>
              <Stack direction="row" spacing={2}>
                <Typography variant="body1" sx={{ color: colors.black }}>
                  Data plan
                </Typography>
                <Typography variant="subtitle1" sx={{ color: colors.black }}>
                  {packageName && packageName.length ? (
                    packageName
                  ) : (
                    <Skeleton
                      variant="rectangular"
                      width={120}
                      height={24}
                      sx={{ backgroundColor: colors.black10 }}
                    />
                  )}
                </Typography>
              </Stack>
              <Stack direction="row" spacing={2}>
                <Typography variant="body1" sx={{ color: colors.black }}>
                  Current site
                </Typography>
                <Typography variant="subtitle1" sx={{ color: colors.black }}>
                  Pamoja site #
                </Typography>
              </Stack>
              <Stack direction="row" spacing={2}>
                <Typography variant="body1" sx={{ color: colors.black }}>
                  Month usage
                </Typography>
                <Typography variant="subtitle1" sx={{ color: colors.black }}>
                  {bundle && bundle.length ? (
                    bundle
                  ) : (
                    <Skeleton
                      variant="rectangular"
                      width={120}
                      height={24}
                      sx={{ backgroundColor: colors.black10 }}
                    />
                  )}
                </Typography>
              </Stack>
            </Stack>
          </Typography>
          <Typography
            component="div"
            role="tabpanel"
            hidden={selectedTab !== 2}
          >
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
                  {subscriberInfo &&
                    subscriberInfo.sim.map((sim: any) => (
                      <TableRow key={sim.iccid}>
                        <TableCell>{sim.msisdn}</TableCell>
                        <TableCell>
                          {sim.isPhysical ? 'pSim' : 'eSim'}
                        </TableCell>
                        <TableCell>
                          {simStatusLoading ? (
                            <Skeleton
                              variant="rectangular"
                              width={120}
                              height={24}
                              sx={{ backgroundColor: colors.black10 }}
                            />
                          ) : (
                            sim.status
                          )}
                        </TableCell>
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
                              onClick={() =>
                                handleSimMenu('deactivateSim', sim.id)
                              }
                            >
                              Deactivate SIM
                            </MenuItem>
                            <MenuItem
                              onClick={() => handleSimMenu('deleteSim', sim.id)}
                              sx={{ color: colors.red }}
                            >
                              Delete SIM
                            </MenuItem>
                          </Menu>
                        </TableCell>
                      </TableRow>
                    ))}
                </TableBody>
              </Table>
            </TableContainer>
          </Typography>
          <Typography
            component="div"
            role="tabpanel"
            hidden={selectedTab !== 3}
          >
            <TableContainer>
              <Table>
                <TableHead>
                  <TableRow>
                    <TableCell>
                      <strong style={{ fontWeight: 'bold' }}>
                        {' '}
                        Billing cycle
                      </strong>
                    </TableCell>
                    <TableCell>
                      <strong style={{ fontWeight: 'bold' }}>
                        {' '}
                        Data usage
                      </strong>
                    </TableCell>
                    <TableCell>
                      <strong style={{ fontWeight: 'bold' }}> Data plan</strong>
                    </TableCell>
                  </TableRow>
                </TableHead>
              </Table>
            </TableContainer>
          </Typography>
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
