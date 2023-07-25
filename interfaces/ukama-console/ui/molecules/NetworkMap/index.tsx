import { Paper, Typography } from '@mui/material';
import dynamic from 'next/dynamic';
const DynamicMap = dynamic(() => import('./DynamicMap'), {
  ssr: false,
});
const NetworkMap = () => {
  return (
    <Paper
      sx={{
        borderRadius: '5px',
        height: 'calc(100vh - 210px)',
      }}
    >
      <DynamicMap id="network-map" zoom={6} className="network-map">
        {() => (
          <>
            <Typography variant="h6">Network Map</Typography>
          </>
        )}
      </DynamicMap>
    </Paper>
  );
};

export default NetworkMap;
