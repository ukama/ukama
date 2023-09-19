import { NODE_IMAGES } from '@/constants';
import { Box } from '@mui/material';

interface IDeviceModalView {
  nodeType: string | undefined;
}

const DeviceModalView = ({ nodeType = 'hnode' }: IDeviceModalView) => {
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
        src={NODE_IMAGES[nodeType as 'hnode' | 'anode' | 'tnode']}
        alt="node-img"
        style={{ maxWidth: '100%', maxHeight: '500px' }}
      />
    </Box>
  );
};

export default DeviceModalView;
