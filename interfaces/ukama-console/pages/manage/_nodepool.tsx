import { MANAGE_NODE_POOL_COLUMN } from '@/constants';
import EmptyView from '@/ui/molecules/EmptyView';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import RouterIcon from '@mui/icons-material/Router';
import { Paper } from '@mui/material';

interface INodePool {
  data: any;
  search: string;
  setSearch: (value: string) => void;
}

const NodePool = ({ data, search, setSearch }: INodePool) => {
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
        subtitle={data.length || '0'}
        search={search}
        title={'My node pool'}
        buttonTitle={'CLAIM NODE'}
        onSearchChange={(e: string) => setSearch(e)}
        handleButtonAction={() => {}}
      />
      <br />
      {data.length === 0 ? (
        <EmptyView icon={RouterIcon} title="No node in nodes pool!" />
      ) : (
        <SimpleDataTable dataset={data} columns={MANAGE_NODE_POOL_COLUMN} />
      )}
    </Paper>
  );
};

export default NodePool;
