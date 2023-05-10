import { NODE_IMAGES } from '@/constants';
import { Node_Type } from '@/generated';
import { Box } from '@mui/material';

interface IDeviceModalView {
  nodeType: Node_Type | undefined;
}

const DeviceModalView = ({ nodeType = Node_Type.Home }: IDeviceModalView) => {
  return (
    <Box
      component={'div'}
      sx={{
        height: { xs: '80vh', md: '62vh' },
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        marginTop: 6,
      }}
    >
      <img
        src={NODE_IMAGES[nodeType]}
        alt="node-img"
        style={{ maxWidth: '100%', maxHeight: '500px' }}
      />
    </Box>
  );
};

export default DeviceModalView;
