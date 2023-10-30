/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import EmptyView from '@/ui/molecules/EmptyView';
import RouterIcon from '@mui/icons-material/Router';
import { Box } from '@mui/material';
import React from 'react';
import NodeSlider from '../NodesSlider';
type NodeContainerProps = {
  items: any;
  handleItemAction: Function;
};

const NodeContainer = ({ items, handleItemAction }: NodeContainerProps) => {
  return (
    <Box
      component="div"
      sx={{
        display: 'flex',
        minHeight: '246px',
        alignItems: 'center',
      }}
    >
      {items.length > 0 ? (
        <NodeSlider items={items} handleItemAction={handleItemAction} />
      ) : (
        <EmptyView size="large" title="No nodes yet!" icon={RouterIcon} />
      )}
    </Box>
  );
};

export default React.memo(NodeContainer);
