import { user } from '@/app-recoil';
import { BillingTabs } from '@/constants';
import { colors } from '@/styles/theme';
import BillHistoryTab from '@/ui/molecules/BillHistoryTab';
import CurrentBillTab from '@/ui/molecules/CurrentBillTab';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import TabPanel from '@/ui/molecules/TabPanel';
import {
  Box,
  FormControl,
  InputLabel,
  MenuItem,
  Select,
  Stack,
  Tab,
  Tabs,
  Typography,
} from '@mui/material';
import { useState } from 'react';
import { useRecoilValue } from 'recoil';

const HistoryData = [
  {
    id: '1',
    date: 'July 11 2022',
    usage: '3 GB',
    total: '$ 100',
    pdf: 'www.google.com',
  },
];

const CurrentBillData = [
  {
    id: '1',
    name: 'Tryphena Nelson',
    rate: '3 GB',
    subtotal: '200',
  },
];

const ORGS = ['Org 1', 'Org 2', 'Ukama'];
const NETWORKS = ['Network 1', 'Network 2', 'Star'];

export default function Billing() {
  const [tab, setTab] = useState(0);
  const [value, setValue] = useState('');
  const userInfo = useRecoilValue(user);
  const handleTabChange = (event: React.SyntheticEvent, newValue: number) =>
    setTab(newValue);

  return (
    <Stack spacing={2}>
      <FormControl variant="standard" sx={{ width: 94 }}>
        <InputLabel id="bill-for-dropdown-label">
          {userInfo?.role === 'admin' ? 'Organization' : 'Networks'}
        </InputLabel>
        <Select
          value={'Ukama'}
          disableUnderline
          onChange={(e) => setValue(e.target.value as string)}
          id="bill-for-dropdown"
          labelId="bill-for-dropdown-label"
        >
          {(userInfo?.role === 'admin' ? ORGS : NETWORKS).map((label, i) => (
            <MenuItem key={`${label}-${i}`} value={label} sx={{ fontSize: 14 }}>
              <Typography variant="body2">{label}</Typography>
            </MenuItem>
          ))}
        </Select>
      </FormControl>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          overflow: 'auto',
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <Stack spacing={2}>
          <Tabs value={tab} onChange={handleTabChange}>
            {BillingTabs.map(({ id, label, value }) => (
              <Tab
                key={id}
                label={label}
                sx={{ px: 3, fontSize: '15px' }}
                id={`billing-tab-${value}`}
              />
            ))}
          </Tabs>
          <Box
            sx={{
              width: '100%',
              height: 'calc(100vh - 200px)',
            }}
          >
            <TabPanel id={'sites-summary-tab'} value={tab} index={0}>
              <CurrentBillTab
                loading={false}
                totalAmount="$20.00"
                currentBill="$ 20.00"
                data={CurrentBillData}
                planName="Console cloud plan - [community/empowerment]"
              />
            </TabPanel>
            <TabPanel id={'sites-power-tab'} value={tab} index={1}>
              <BillHistoryTab loading={false} data={HistoryData} />
            </TabPanel>
          </Box>
        </Stack>
      </LoadingWrapper>
    </Stack>
  );
}
