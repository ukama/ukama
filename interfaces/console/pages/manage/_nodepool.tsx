/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { MANAGE_NODE_POOL_COLUMN } from '@/constants';
import EmptyView from '@/ui/molecules/EmptyView';
import PageContainerHeader from '@/ui/molecules/PageContainerHeader';
import SimpleDataTable from '@/ui/molecules/SimpleDataTable';
import RouterIcon from '@mui/icons-material/Router';
import { Box } from '@mui/material';

interface INodePool {
  data: any;
  search: string;
  setSearch: (value: string) => void;
  networkList: string[];
  handleCreateNetwork: () => void;
}

const NodePool = ({
  data,
  search,
  setSearch,
  networkList,
  handleCreateNetwork,
}: INodePool) => {
  return (
    <Box sx={{ width: '100%', height: '100%' }}>
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
        <SimpleDataTable
          dataset={data}
          columns={MANAGE_NODE_POOL_COLUMN}
          networkList={networkList}
          handleCreateNetwork={handleCreateNetwork}
        />
      )}
    </Box>
  );
};

export default NodePool;
