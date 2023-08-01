import { Paper } from '@mui/material';
import dynamic from 'next/dynamic';
import { LabelOverlayUI, SitesSelection, SitesTree } from './OverlayUI';
const DynamicMap = dynamic(() => import('./DynamicMap'), {
  ssr: false,
});
const NetworkMap = () => {
  return (
    <Paper
      sx={{
        borderRadius: '5px',
        height: 'calc(100vh - 310px)',
      }}
    >
      <DynamicMap id="network-map" zoom={6} className="network-map">
        {() => (
          <>
            <LabelOverlayUI />
            <SitesTree />
            <SitesSelection />
          </>
        )}
      </DynamicMap>
    </Paper>
  );
};

export default NetworkMap;
