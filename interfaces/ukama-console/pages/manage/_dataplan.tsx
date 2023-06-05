import colors from '@/styles/theme/colors';
import { EmptyView, OptionsPopover } from '@/ui/components';
import PageContainerHeader from '@/ui/components/PageContainerHeader';
import { getDataPlanUsage } from '@/utils';
import { PeopleAlt } from '@mui/icons-material';
import UpdateIcon from '@mui/icons-material/SystemUpdateAltRounded';
import { Card, Grid, Paper, Stack, Typography } from '@mui/material';

interface IDataPlan {
  data: any;
  handleActionButon: Function;
  handleOptionMenuItemAction: Function;
}

const DataPlan = ({
  data,
  handleActionButon,
  handleOptionMenuItemAction,
}: IDataPlan) => {
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
                <Card
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
                  <Stack spacing={1} width={'100%'}>
                    <Grid
                      xs={12}
                      container
                      direction={'row'}
                      textAlign={'center'}
                    >
                      <Grid item xs={11} pl={3}>
                        <Typography variant="h5" sx={{ fontWeight: 400 }}>
                          {name}
                        </Typography>
                      </Grid>
                      <Grid item xs={1}>
                        <OptionsPopover
                          cid={'data-table-action-popover'}
                          menuOptions={[
                            { id: 0, title: 'Edit', route: 'edit' },
                            { id: 1, title: 'Delete', route: 'delete' },
                          ]}
                          handleItemClick={(type: string) =>
                            handleOptionMenuItemAction(uuid, type)
                          }
                        />
                      </Grid>
                    </Grid>
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
                </Card>
              </Grid>
            ),
          )}
        </Grid>
      )}
    </Paper>
  );
};

export default DataPlan;
