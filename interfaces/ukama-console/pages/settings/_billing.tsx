import { BillingTabs } from '@/constants';
import { colors } from '@/styles/theme';
import CurrentBillTab from '@/ui/molecules/CurrentBillTab';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import TabPanel from '@/ui/molecules/TabPanel';
import { Box, Stack, Tab, Tabs, Typography } from '@mui/material';
import { useState } from 'react';

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
            <CurrentBillTab />
          </TabPanel>
          <TabPanel id={'sites-power-tab'} value={tab} index={1}>
            <Typography>TWO</Typography>
          </TabPanel>
        </Box>
      </Stack>
    </LoadingWrapper>
  );
}
