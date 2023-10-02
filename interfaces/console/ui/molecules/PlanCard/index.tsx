import colors from '@/styles/theme/colors';
import { convertToWeeksOrMonths, getDataPlanUsage } from '@/utils';
import { PeopleAlt } from '@mui/icons-material';
import { Card, Grid, Stack, Typography } from '@mui/material';
import OptionsPopover from '../OptionsPopover';

interface IPlanCard {
  uuid: string;
  name: string;
  users: string;
  amount: string;
  duration: number;
  currency: string;
  dataUnit: string;
  dataVolume: string;
  isOptions?: boolean;
  handleOptionMenuItemAction?: Function;
}

const PlanCard = ({
  uuid,
  name,
  users,
  amount,
  dataUnit,
  duration,
  currency,
  dataVolume,
  isOptions = true,
  handleOptionMenuItemAction,
}: IPlanCard) => {
  return (
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
        <Grid xs={12} container direction={'row'} textAlign={'center'}>
          <Grid item xs={11} pl={3}>
            <Typography variant="h5" sx={{ fontWeight: 400 }}>
              {name}
            </Typography>
          </Grid>
          <Grid item xs={1}>
            {isOptions && (
              <OptionsPopover
                cid={'data-table-action-popover'}
                menuOptions={[
                  { id: 0, title: 'Edit', route: 'edit', Icon: null },
                  {
                    id: 1,
                    title: 'Delete',
                    route: 'delete',
                    Icon: null,
                  },
                ]}
                handleItemClick={(type: string) =>
                  handleOptionMenuItemAction &&
                  handleOptionMenuItemAction(uuid, type)
                }
              />
            )}
          </Grid>
        </Grid>
        <Typography variant="body2" fontWeight={400}>
          {getDataPlanUsage(
            convertToWeeksOrMonths(duration),
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
  );
};

export default PlanCard;
