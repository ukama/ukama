import { SUBSCRIBER_TABLE_COLUMNS, SUBSCRIBER_TABLE_MENU } from '@/constants';
import {
  ContainerMax,
  HorizontalContainerJustify,
  PageContainer,
  VerticalContainer,
} from '@/styles/global';
import { colors } from '@/styles/theme';
import { DataTableWithOptions } from '@/ui/components';
import { Search } from '@mui/icons-material';
import { Button, Grid, TextField, Typography } from '@mui/material';

const SUB_COUNT = 1;
export default function Page() {
  const onResidentsTableMenuItem = (id: string, type: string) => {
    console.log(id, type);
  };
  return (
    <PageContainer>
      <HorizontalContainerJustify>
        <Grid container justifyContent={'space-between'} spacing={1}>
          <Grid container item xs={12} md="auto" alignItems={'center'}>
            <Grid item xs={'auto'}>
              <Typography variant="h6" mr={1}>
                My subscribers
              </Typography>
            </Grid>
            <Grid item xs={'auto'}>
              <Typography variant="subtitle2" mr={1.4}>{`(${0})`}</Typography>
            </Grid>
            <Grid item xs={12} md={'auto'}>
              <TextField
                id="subscriber-search"
                label="Search"
                variant="outlined"
                size="small"
                sx={{ width: { xs: '100%', lg: '250px' } }}
                InputProps={{
                  endAdornment: <Search htmlColor={colors.black54} />,
                }}
              />
            </Grid>
          </Grid>
          <Grid item xs={12} md={'auto'}>
            <Button
              variant="contained"
              color="primary"
              size="medium"
              sx={{ width: { xs: '100%', md: '250px' } }}
            >
              Add Subscriber
            </Button>
          </Grid>
        </Grid>
      </HorizontalContainerJustify>
      <VerticalContainer>
        <ContainerMax mt={4.5}>
          <DataTableWithOptions
            columns={SUBSCRIBER_TABLE_COLUMNS}
            dataset={[
              {
                name: 'John Doe',
                network: 'Globe',
                dataUsage: '1.2 GB',
                dataPlan: '1.5 GB',
                actions: 'actions',
              },
              {
                name: 'John Do',
                network: 'Earth',
                dataUsage: '1.1 GB',
                dataPlan: '1.9 GB',
                actions: 'actions',
              },
            ]}
            menuOptions={SUBSCRIBER_TABLE_MENU}
            onMenuItemClick={onResidentsTableMenuItem}
            emptyViewLabel={'No subscribers yet! [100] SIMs left in pool.'}
          />
        </ContainerMax>
      </VerticalContainer>
    </PageContainer>
  );
}
