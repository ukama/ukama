import { BillingTabs } from '@/constants';
import { colors } from '@/styles/theme';
import BillHistoryTab from '@/ui/molecules/BillHistoryTab';
import CurrentBillTab from '@/ui/molecules/CurrentBillTab';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import TabPanel from '@/ui/molecules/TabPanel';
import { Box, Stack, Tab, Tabs, Typography } from '@mui/material';
import { useState } from 'react';

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

export default function Billing() {
  const [tab, setTab] = useState(0);
  const handleTabChange = (event: React.SyntheticEvent, newValue: number) =>
    setTab(newValue);
  return (
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
              sx={{ px: 3 }}
              id={`billing-tab-${value}`}
            />
          ))}
        </Tabs>
        <Typography variant="caption">
          <b>May overview </b>(06/14/2022 - 07/14/2022)
        </Typography>
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
  );
}
