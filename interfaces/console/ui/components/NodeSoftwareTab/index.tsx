import { NodeAppsColumns } from '@/constants/tableColumns';
import { colors } from '@/styles/theme';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import {
  Button,
  Card,
  CardActions,
  CardContent,
  Grid,
  Paper,
  Stack,
  Typography,
} from '@mui/material';
interface INodeRadioTab {
  loading: boolean;
  nodeApps: any;
  NodeLogs: any;
  getNodeAppDetails: Function;
}

const NodeSoftwareTab = ({
  getNodeAppDetails,
  loading,
  NodeLogs,
  nodeApps,
}: INodeRadioTab) => {
  return (
    <LoadingWrapper isLoading={loading} height={400}>
      <Paper
        sx={{
          p: 3,
          height: '100%',
          borderRadius: '4px',
          marginBottom: 2,
        }}
      >
        <Typography variant="h6" sx={{ marginBottom: 3 }}>
          Change Logs
        </Typography>
        <SimpleDataTable
          height={'400'}
          dataset={NodeLogs}
          columns={NodeAppsColumns}
        />
      </Paper>
      <Paper sx={{ height: '100%', p: 3, borderRadius: '4px' }}>
        <Typography variant="h6" sx={{ mb: 4 }}>
          Node Apps
        </Typography>
        <Grid container spacing={3}>
          {nodeApps?.map(({ id, title, cpu, memory, version }: any) => (
            <Grid item xs={12} md={6} lg={3} key={id}>
              <Card variant="outlined">
                <CardContent>
                  <Stack
                    spacing={1}
                    direction="row"
                    sx={{ alignItems: 'center' }}
                  >
                    <CheckCircleIcon
                      htmlColor={colors.green}
                      fontSize="medium"
                    />
                    <Typography variant="h5" textTransform={'capitalize'}>
                      {title}
                    </Typography>
                  </Stack>
                  <Typography
                    variant="body2"
                    color="text.secondary"
                    gutterBottom
                  >
                    Version: {version}
                  </Typography>
                  <Stack direction="row" spacing={1 / 2} mt={'12px'}>
                    <Typography variant="body2">CPU:</Typography>
                    <Typography variant="body2" sx={{ color: colors.darkBlue }}>
                      {parseFloat(cpu).toFixed(2)} %
                    </Typography>
                  </Stack>
                  <Stack direction="row" spacing={1 / 2}>
                    <Typography variant="body2">MEMORY:</Typography>
                    <Typography variant="body2" sx={{ color: colors.darkBlue }}>
                      {parseFloat(memory).toFixed(2)} KB
                    </Typography>
                  </Stack>
                </CardContent>
                <CardActions sx={{ ml: 1 }}>
                  <Button onClick={() => getNodeAppDetails(id)}>
                    VIEW MORE
                  </Button>
                </CardActions>
              </Card>
            </Grid>
          ))}
        </Grid>
      </Paper>
    </LoadingWrapper>
  );
};

export default NodeSoftwareTab;
