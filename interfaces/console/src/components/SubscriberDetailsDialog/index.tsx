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
  styled,
  Menu,
  MenuItem,
} from '@mui/material';
import MoreHorizIcon from '@mui/icons-material/MoreHoriz';
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
  onDeleteSim?: (simId: string) => void;
  packageHistories?: any[];
  packagesData?: PackagesResDto;
  loadingPackageHistories?: boolean;
  dataUsage: string;
}

function TabPanel({ children, value, index }: any) {
  return (
    <div role="tabpanel" hidden={value !== index}>
      {value === index && <Box sx={{ pt: 2 }}>{children}</Box>}
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
  onDeleteSim,
  packageHistories,
  packagesData,
  loadingPackageHistories,
  dataUsage,
}) => {
  const [tabIndex, setTabIndex] = React.useState(0);
  const [menuAnchor, setMenuAnchor] = React.useState<null | HTMLElement>(null);

  const handleTabChange = (_: React.SyntheticEvent, newIndex: number) => {
    setTabIndex(newIndex);
  };

  const openMenu = (e: React.MouseEvent<HTMLElement>) =>
    setMenuAnchor(e.currentTarget);
  const closeMenu = () => setMenuAnchor(null);

  const handleDelete = () => {
    onDeleteSubscriber(subscriber.id);
    closeMenu();
  };

  return (
    <Dialog open={open} onClose={onClose} fullWidth maxWidth="sm">
      <DialogTitle sx={{ px: 3, pt: 2, pb: 1, position: 'relative' }}>
        <Box sx={{ display: 'flex', alignItems: 'center' }}>
          <Typography variant="h6" sx={{ fontWeight: 600 }}>
            {subscriber.firstName}
          </Typography>
          <IconButton size="small" onClick={openMenu} sx={{ ml: 0.5 }}>
            <MoreHorizIcon fontSize="small" />
          </IconButton>
        </Box>
        <IconButton
          aria-label="close"
          onClick={onClose}
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
          onClose={closeMenu}
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
          borderBottom: `1px solid ${colors.black10}`,
        }}
      >
        {/* //TODO billing cycle sitll under discussion */}
        {['INFORMATION', 'DATA PLANS', 'SIMS', 'HISTORY'].map((label, idx) => (
          <Tab key={label} label={label} {...a11yProps(idx)} />
        ))}
      </Tabs>

      <DialogContent>
        <TabPanel value={tabIndex} index={0}>
          <SubscriberInfoTab
            subscriber={subscriber}
            onUpdateSubscriber={onUpdateSubscriber}
          />
        </TabPanel>
        <TabPanel value={tabIndex} index={1}>
          <SubscriberDataPlansTab
            packageHistories={packageHistories}
            packagesData={packagesData}
            dataUsage={dataUsage}
          />
        </TabPanel>
        <TabPanel value={tabIndex} index={2}>
          <SubscriberSimsTab
            sims={sims}
            onSimAction={onSimAction}
            onDeleteSim={onDeleteSim}
          />
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
        <Button onClick={onClose} variant="contained">
          CLOSE
        </Button>
      </DialogActions>
    </Dialog>
  );
};

export default SubscriberDetailsDialog;
