import { colors } from '@/styles/theme';
import FormControlCheckboxes from '@/ui/components/FormControlCheckboxes';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import { Divider, Grid, Paper, Typography } from '@mui/material';
import { useCallback, useState } from 'react';

export default function Alerts() {
  const [alertList, setAlertList] = useState<Object>({});

  const handleAlertChange = useCallback((key: string, value: boolean) => {
    setAlertList((prevState: any) => ({
      ...prevState,
      [key]: value,
    }));
  }, []);

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
      <Paper
        sx={{
          py: 3,
          px: 4,
          width: '100%',
          borderRadius: '5px',
          height: 'calc(100vh - 200px)',
        }}
      >
        <Grid container spacing={2}>
          <Grid item container spacing={2}>
            <Grid item xs={12} md={3}>
              <Typography variant="h6">Common Events</Typography>
            </Grid>
            <Grid item xs={12} md={8}>
              {[1, 2].map((i) => (
                <Grid key={`${i}-`} item xs={12} sm={10} md={9}>
                  <FormControlCheckboxes
                    values={alertList}
                    handleChange={handleAlertChange}
                    checkboxList={[
                      {
                        id: 1,
                        label: `Event Log ${i}`,
                        value: `event${i}`,
                      },
                      {
                        id: 2,
                        label: `Alerts ${i}`,
                        value: `alert${i}`,
                      },
                      {
                        id: 3,
                        label: `Email ${i}`,
                        value: `email${i}`,
                      },
                    ]}
                  />
                </Grid>
              ))}
            </Grid>
          </Grid>
          <Divider sx={{ width: '100%' }} />
          <Grid item container spacing={2}>
            <Grid item xs={12} md={3}>
              <Typography variant="h6">Cloud Events</Typography>
            </Grid>
            <Grid item container xs={12} md={9}>
              {[3, 4].map((i) => (
                <Grid key={`${i}-`} item xs={12} sm={10} md={8}>
                  <FormControlCheckboxes
                    values={alertList}
                    handleChange={handleAlertChange}
                    checkboxList={[
                      {
                        id: 1,
                        label: `Event Log ${i}`,
                        value: `event${i}`,
                      },
                      {
                        id: 2,
                        label: `Alerts ${i}`,
                        value: `alert${i}`,
                      },
                      {
                        id: 3,
                        label: `Email ${i}`,
                        value: `email${i}`,
                      },
                    ]}
                  />
                </Grid>
              ))}
            </Grid>
          </Grid>
          <Divider sx={{ width: '100%' }} />
          <Grid item container spacing={2}>
            <Grid item xs={12} md={3}>
              <Typography variant="h6">AP Events</Typography>
            </Grid>
            <Grid item container xs={12} md={9}>
              {[3, 4].map((i) => (
                <Grid key={`${i}-`} item xs={12} sm={10} md={8}>
                  <FormControlCheckboxes
                    values={alertList}
                    handleChange={handleAlertChange}
                    checkboxList={[
                      {
                        id: 1,
                        label: `Event Log ${i}`,
                        value: `event${i}`,
                      },
                      {
                        id: 2,
                        label: `Alerts ${i}`,
                        value: `alert${i}`,
                      },
                      {
                        id: 3,
                        label: `Email ${i}`,
                        value: `email${i}`,
                      },
                    ]}
                  />
                </Grid>
              ))}
            </Grid>
          </Grid>
        </Grid>
      </Paper>
    </LoadingWrapper>
  );
}
