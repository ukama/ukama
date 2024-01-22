import React from 'react';
import { Card, CardContent, Stack, Typography, Box } from '@mui/material';
import GroupIcon from '@mui/icons-material/Group';
import { colors } from '@/styles/theme';
import { PackageDto } from '@/generated';

interface BasePlanProps {
  subscriberCount: number;
  dataPlans: PackageDto[];
}

const BasePlan: React.FC<BasePlanProps> = ({ subscriberCount, dataPlans }) => {
  if (dataPlans.length === 0) {
    return (
      <div style={{ minHeight: '200px', paddingTop: 60 }}>
        <Typography variant="body1" align="center">
          No data plan yet.
        </Typography>
      </div>
    );
  }

  return (
    <div>
      {dataPlans.map((dataPlan, index) => (
        <Card
          key={index}
          variant="outlined"
          sx={{
            position: 'relative',
            maxWidth: 300,
            borderRadius: 2,
            overflow: 'hidden',
          }}
        >
          <div
            style={{
              position: 'absolute',
              top: 0,
              left: 0,
              width: '100%',
              height: '8px',
              backgroundColor: colors.primaryMain,
            }}
          />

          <CardContent
            sx={{
              display: 'flex',
              flexDirection: 'column',
              alignItems: 'center',
              justifyContent: 'center',
              height: '100%',
            }}
          >
            <Stack direction="column" spacing={2}>
              <Typography variant="h5">{dataPlan.name}</Typography>
              <Typography variant="body1">
                {dataPlan.amount}$/{dataPlan.dataVolume} MB/{dataPlan.duration}
                Month
              </Typography>
              <Stack
                direction="row"
                spacing={1}
                alignItems={'center'}
                justifyContent={'center'}
              >
                <GroupIcon />
                <Typography variant="body1">{subscriberCount}</Typography>
              </Stack>
            </Stack>
          </CardContent>
        </Card>
      ))}
    </div>
  );
};

export default BasePlan;
