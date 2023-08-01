import { NODE_TABLE_COLUMNS, NODE_TABLE_MENU } from '@/constants';
import { PageContainer } from '@/styles/global';
import { colors } from '@/styles/theme';
import AddNodeDialog from '@/ui/molecules/AddNode';
import DataTableWithOptions from '@/ui/molecules/DataTableWithOptions';
import LoadingWrapper from '@/ui/molecules/LoadingWrapper';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import RouterIcon from '@mui/icons-material/Router';
import { Stack } from '@mui/material';
import { useEffect, useState } from 'react';

const DATA = [
  {
    id: '1',
    name: 'Node 1',
    type: 'Home',
    state: 'active',
    network: 'Network 1',
  },
  {
    id: '2',
    name: 'Node 2',
    type: 'Tower',
    state: 'active',
    network: 'Network 1',
  },
];

const AVAILABLE_NODES = [
  { id: 'node-1', name: 'Node 1', isChecked: false },
  { id: 'node-2', name: 'Node 2', isChecked: false },
];

const structureData = (data: any) => {};

export default function Page() {
  const [search, setSearch] = useState<string>('');
  const [nodes, setNodes] = useState(DATA);
  const [availableNodes, setAvailableNodes] = useState(AVAILABLE_NODES);
  const [isShowAddNodeDialog, setIsShowAddNodeDialog] =
    useState<boolean>(false);

  useEffect(() => {
    if (search.length > 3) {
      const nodes = DATA.filter((node) => {
        const s = search.toLowerCase();
        if (
          node.name.toLowerCase().includes(s) ||
          node.name.toLowerCase().includes(s)
        )
          return node;
      });
      setNodes(() => nodes);
    } else if (search.length === 0) {
      setNodes(() => DATA);
    }
  }, [search]);

  const handleSearchChange = (str: string) => {
    setSearch(str);
  };

  const handleAddNode = () => {};

  const handleNodeCheck = (id: string, isChecked: boolean) => {
    setAvailableNodes((prev) => {
      const nodes = prev.map((node) => {
        if (node.id === id) {
          return { ...node, isChecked };
        }
        return node;
      });
      return nodes;
    });
  };

  const handleCloseAddNodeDialog = () => setIsShowAddNodeDialog(false);
  return (
    <>
      <LoadingWrapper
        radius="small"
        width={'100%'}
        isLoading={false}
        cstyle={{
          backgroundColor: false ? colors.white : 'transparent',
        }}
      >
        <PageContainer>
          <Stack
            spacing={2}
            height={'100%'}
            direction={'column'}
            alignItems={'center'}
            justifyContent={'flex-start'}
          >
            <PageContainerHeader
              subtitle="3"
              search={search}
              title={'My Nodes'}
              showSearch={true}
              buttonTitle="Add Nodes"
              onSearchChange={handleSearchChange}
              handleButtonAction={() => setIsShowAddNodeDialog(true)}
            />
            <DataTableWithOptions
              dataset={nodes}
              icon={RouterIcon}
              onMenuItemClick={() => {}}
              columns={NODE_TABLE_COLUMNS}
              menuOptions={NODE_TABLE_MENU}
              emptyViewLabel={'No node yet!'}
            />
          </Stack>
        </PageContainer>
      </LoadingWrapper>
      <AddNodeDialog
        data={availableNodes}
        labelNegativeBtn="Close"
        labelSuccessBtn="Add Nodes"
        isOpen={isShowAddNodeDialog}
        handleNodeCheck={handleNodeCheck}
        title="Add nodes to your network"
        handleSuccessAction={handleAddNode}
        handleCloseAction={handleCloseAddNodeDialog}
        description="Add nodes to your network to start managing them here. If you cannot find a desired node, check to make sure it is not already added to another network."
      />
    </>
  );
}
