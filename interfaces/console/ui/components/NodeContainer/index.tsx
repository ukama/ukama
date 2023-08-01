import { NodeDto } from '@/generated';
import RouterIcon from '@mui/icons-material/Router';
import { Box } from '@mui/material';
import React from 'react';
import { EmptyView, NodeSlider } from '..';
type NodeContainerProps = {
  items: NodeDto[];
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
