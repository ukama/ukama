import colors from '@/styles/theme/colors';
import { Chip, Grid, Link, Paper, Stack, Typography } from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';

interface INodeDetailsCard {
  loading: boolean;
  nodeTitle: string;
  nodeType?: any;
  isUpdateAvailable: boolean;
  handleUpdateNode: Function;
  getNodeUpdateInfos: Function;
}

const NodeDetailsCard = ({
  loading,
  nodeTitle,
  isUpdateAvailable,
  getNodeUpdateInfos,
  nodeType = 'HOME',
}: INodeDetailsCard) => {
  return (
    <LoadingWrapper
      width="100%"
      height="100%"
      radius={'small'}
      isLoading={loading}
    >
      <Paper sx={{ p: 3, gap: 1 }}>
        <Stack spacing={3}>
          <Grid container>
            <Grid item xs={5}>
              <Typography variant="h6">{nodeTitle}</Typography>
            </Grid>
            {isUpdateAvailable && (
              <Grid item container xs={7} justifyContent="flex-end">
                <Chip
                  variant="outlined"
                  sx={{
                    color: colors.primaryMain,
                    border: `1px solid ${colors.primaryMain}`,
                  }}
                  label={
                    <Stack spacing={'4px'} direction="row" alignItems="center">
                      <Typography variant="body2">
                        Software update available — view
                      </Typography>
                      <Link
                        onClick={() => getNodeUpdateInfos()}
                        sx={{
                          cursor: 'pointer',
                          typography: 'body2',
                          color: colors.primaryDark,
                        }}
                      >
                        notes
                      </Link>
                    </Stack>
                  }
                />
              </Grid>
            )}
          </Grid>

          {/* <DeviceModalView nodeType={nodeType} /> */}
        </Stack>
      </Paper>
    </LoadingWrapper>
  );
};

export default NodeDetailsCard;
