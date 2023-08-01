import { MANAGE_SIM_POOL_COLUMN } from '@/constants';
import EmptyView from '@/ui/molecules/EmptyView';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import SimCardIcon from '@mui/icons-material/SimCard';
import { Paper } from '@mui/material';

interface ISimPool {
  data: any;
  handleActionButon: Function;
}

const SimPool = ({ data, handleActionButon }: ISimPool) => {
  return (
    <Paper
      sx={{
        py: 3,
        px: 4,
        width: '100%',
        overflow: 'hidden',
        borderRadius: '5px',
        height: 'calc(100vh - 200px)',
      }}
    >
      <PageContainerHeader
        subtitle={data.length || '0'}
        showSearch={false}
        title={'My SIM pool'}
        buttonTitle={'IMPORT SIMS'}
        handleButtonAction={handleActionButon}
      />
      <br />
      {data.length === 0 ? (
        <EmptyView icon={SimCardIcon} title="No sims in sim pool!" />
      ) : (
        <SimpleDataTable dataset={data} columns={MANAGE_SIM_POOL_COLUMN} />
      )}
    </Paper>
  );
};

export default SimPool;
