import colors from '@/styles/theme/colors';
import { EmptyView } from '@/ui/components';
import PageContainerHeader from '@/ui/components/PageContainerHeader';
import { getDataPlanUsage } from '@/utils';
import { PeopleAlt } from '@mui/icons-material';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { Grid, Paper, Stack, Typography } from '@mui/material';

interface IDataPlan {
  data: any;
  handleActionButon: Function;
}

const DataPlan = ({ data, handleActionButon }: IDataPlan) => {
  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        borderRadius: '5px',
        height: 'calc(100vh - 200px)',
      }}
    >
      <PageContainerHeader
        showSearch={false}
        title={'Data plans'}
        buttonTitle={'CREATE DATA PLAN'}
        handleButtonAction={handleActionButon}
      />
      <br />
      {data.length === 0 ? (
        <EmptyView icon={UpdateIcon} title="No data plan created yet!" />
      ) : (
        <Grid container rowSpacing={2} columnSpacing={2}>
          {data.map(
            ({
              uuid,
              name,
              duration,
              users,
              currency,
              dataVolume,
              dataUnit,
              amount,
            }: any) => (
              <Grid item xs={12} sm={6} md={4} key={uuid}>
                <Paper
                  variant="outlined"
                  sx={{
                    px: 3,
                    py: 2,
                    display: 'flex',
                    boxShadow: 'none',
                    borderRadius: '4px',
                    textAlign: 'center',
                    justifyContent: 'center',
                    borderTop: `4px solid ${colors.primaryMain}`,
                  }}
                >
                  <Stack spacing={1}>
                    <Typography variant="h5" sx={{ fontWeight: 400 }}>
                      {name}
                    </Typography>
                    <Typography variant="body2" fontWeight={400}>
                      {getDataPlanUsage(
                        duration,
                        currency,
                        amount,
                        dataVolume,
                        dataUnit,
                      )}
                    </Typography>
                    {false && (
                      <Stack
                        spacing={0.6}
                        direction={'row'}
                        alignItems={'flex-end'}
                        justifyContent={'center'}
                      >
                        <PeopleAlt htmlColor={colors.black54} />
                        <Typography variant="body2" fontWeight={400}>
                          {users}
                        </Typography>
                      </Stack>
                    )}
                  </Stack>
                </Paper>
              </Grid>
            ),
          )}
        </Grid>
      )}
    </Paper>
  );
};

export default DataPlan;
