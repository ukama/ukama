import React, { useState } from 'react';
import {
  Grid,
  Typography,
  Select,
  Menu,
  MenuItem,
  IconButton,
  Stack,
  SelectChangeEvent,
} from '@mui/material';
import AttachMoneyIcon from '@mui/icons-material/AttachMoney';
import { RoundedCard } from '@/styles/global';
import PeopleAltIcon from '@mui/icons-material/PeopleAlt';
import BarChartIcon from '@mui/icons-material/BarChart';
import colors from '@/styles/theme/colors';

interface SalesMetric {
  label: string;
  value: number;
  icon: React.ReactNode; // React Node for custom icons
}

interface SubscriberSummaryProps {
  monthlySales: number;
  total: number;
  returning: number;
  averageSale: number;
}

const SubscriberSummary: React.FC<SubscriberSummaryProps> = ({
  monthlySales,
  total,
  returning,
  averageSale,
}) => {
  const [selectedDate, setSelectedDate] = useState<string | null>(null);

  const handleDateChange = (event: SelectChangeEvent<string | null>) => {
    const newDate = event.target.value;
    setSelectedDate(newDate);
  };

  const salesMetrics: SalesMetric[] = [
    {
      label: 'Monthly Sales',
      value: monthlySales,
      icon: <AttachMoneyIcon sx={{ color: colors.primaryMain }} />,
    },
    {
      label: 'Total Sales',
      value: total,
      icon: <PeopleAltIcon sx={{ color: colors.darkPurple }} />,
    },
    {
      label: 'Returning Customers',
      value: returning,
      icon: <PeopleAltIcon sx={{ color: colors.primaryMain }} />,
    },
    { label: 'Average Sale', value: averageSale, icon: <BarChartIcon /> },
  ];
  const filteredMetrics = selectedDate
    ? salesMetrics.map((metric) => ({
        ...metric,
        value: Math.round(Math.random() * 1000), // Replace this with actual data fetching based on the selected date
      }))
    : salesMetrics;

  return (
    <>
      <Grid container spacing={2} sx={{ marginTop: 2 }}>
        <Grid item xs={12}>
          <Stack direction="row" spacing={2} alignItems={'center'}>
            <Typography variant="body1"> overview</Typography>
            <Select
              value={selectedDate}
              onChange={handleDateChange}
              displayEmpty
              sx={{ width: '10%' }}
            >
              <MenuItem value="" disabled>
                Select Date
              </MenuItem>
              {/* Add your date options here */}
              <MenuItem value="2024-01-23">2024-01-23</MenuItem>
              <MenuItem value="2024-01-24">2024-01-24</MenuItem>
              {/* Add more date options as needed */}
            </Select>
          </Stack>
        </Grid>
        {filteredMetrics.map((metric, index) => (
          <Grid item xs={3} key={index}>
            <RoundedCard>
              <Stack direction="row" spacing={2}>
                <IconButton
                  style={{
                    width: '70px',
                    backgroundColor: '#d9eae4',
                    borderRadius: '8px',
                  }}
                  color="primary"
                >
                  {metric.icon}
                </IconButton>
                <Stack direction="column" spacing={2}>
                  <Typography variant="subtitle2">{metric.label}</Typography>
                  <Typography variant="h5">${metric.value}</Typography>
                </Stack>
              </Stack>
            </RoundedCard>
          </Grid>
        ))}
      </Grid>
    </>
  );
};

export default SubscriberSummary;
