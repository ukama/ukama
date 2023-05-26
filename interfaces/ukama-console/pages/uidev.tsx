import { HorizontalContainerJustify } from '@/styles/global';
import { colors } from '@/styles/theme';
import { Search } from '@mui/icons-material';
import { Button, Grid, TextField, Typography } from '@mui/material';
import { useRouter } from 'next/router';

const Page = () => {
  const router = useRouter();
  return (
    <>
      <HorizontalContainerJustify>
        <Grid container justifyContent={'space-between'} spacing={1}>
          <Grid container item xs={12} md="auto" alignItems={'center'}>
            <Grid item xs={'auto'}>
              <Typography variant="h6" mr={1}>
                My subscribers
              </Typography>
            </Grid>
            <Grid item xs={'auto'}>
              <Typography variant="subtitle2" mr={1.4}>{`(0)`}</Typography>
            </Grid>
            <Grid item xs={12} md={'auto'}>
              <TextField
                id="subscriber-search"
                variant="outlined"
                size="small"
                placeholder="Search"
                value={''}
                onChange={(e) => {}}
                sx={{ width: { xs: '100%', lg: '250px' } }}
                InputLabelProps={{
                  shrink: false,
                }}
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
    </>
  );
};
export default Page;
