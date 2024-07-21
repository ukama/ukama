'use client';

import colors from '@/theme/colors';
import LoadingWrapper from '@/components/LoadingWrapper';
import SiteCard from '@/components/SiteCard';
import { Grid, Paper, Typography, Button, AlertColor } from '@mui/material';
import SiteConfigurationStepperDialog from '@/components/SiteConfigurationStepperDialog';
import StepperDialog from '@/components/StepperTest';

import { useEffect, useState } from 'react';
import {
  useGetSitesLazyQuery,
  useGetNetworksQuery,
  useGetComponentsByUserIdLazyQuery,
} from '@/client/graphql/generated';
import { useAppContext } from '@/context';
const dummyComponents = [
  {
    id: 'switch-1',
    inventory_id: 'SW001',
    category: 'SWITCH',
    type: 'switch',
    user_id: 'user-123',
    description: '8 port switch',
    datasheet_url: 'http://example.com/switch1-datasheet',
    images_url: 'http://example.com/switch1-image',
    part_number: 'SW-8P-001',
    manufacturer: 'NetworkGear Inc.',
    managed: 'yes',
    warranty: 24,
    specification: '8 ports, 1Gbps per port',
  },
  {
    id: 'switch-2',
    inventory_id: 'SW002',
    category: 'SWITCH',
    type: 'switch',
    user_id: 'user-123',
    description: '16 port switch',
    datasheet_url: 'http://example.com/switch2-datasheet',
    images_url: 'http://example.com/switch2-image',
    part_number: 'SW-16P-001',
    manufacturer: 'NetworkGear Inc.',
    managed: 'yes',
    warranty: 24,
    specification: '16 ports, 1Gbps per port',
  },
  {
    id: 'power-1',
    inventory_id: 'PW001',
    category: 'POWER',
    type: 'power',
    user_id: 'user-123',
    description: 'Battery Pack',
    datasheet_url: 'http://example.com/battery-datasheet',
    images_url: 'http://example.com/battery-image',
    part_number: 'BAT-12V-001',
    manufacturer: 'PowerSolutions Ltd.',
    managed: 'yes',
    warranty: 12,
    specification: '12V, 100Ah',
  },
  {
    id: 'power-2',
    inventory_id: 'PW002',
    category: 'POWER',
    type: 'power',
    user_id: 'user-123',
    description: 'AC Power Supply',
    datasheet_url: 'http://example.com/ac-power-datasheet',
    images_url: 'http://example.com/ac-power-image',
    part_number: 'ACP-001',
    manufacturer: 'PowerSolutions Ltd.',
    managed: 'yes',
    warranty: 24,
    specification: '110-240V AC input, 12V DC output',
  },
  {
    id: 'backhaul-1',
    inventory_id: 'BH001',
    category: 'BACKHAUL',
    type: 'backhaul',
    user_id: 'user-123',
    description: 'ViaSAT Satellite Modem',
    datasheet_url: 'http://example.com/viasat-datasheet',
    images_url: 'http://example.com/viasat-image',
    part_number: 'VS-SAT-001',
    manufacturer: 'ViaSAT',
    managed: 'yes',
    warranty: 36,
    specification: 'Up to 100Mbps downlink, 20Mbps uplink',
  },
  {
    id: 'backhaul-2',
    inventory_id: 'BH002',
    category: 'BACKHAUL',
    type: 'backhaul',
    user_id: 'user-123',
    description: '4G LTE Modem',
    datasheet_url: 'http://example.com/lte-datasheet',
    images_url: 'http://example.com/lte-image',
    part_number: 'LTE-4G-001',
    manufacturer: 'TelecomTech Corp.',
    managed: 'yes',
    warranty: 24,
    specification: 'Up to 150Mbps downlink, 50Mbps uplink',
  },
  {
    id: 'access-1',
    inventory_id: 'AC001',
    category: 'ACCESS',
    type: 'access',
    user_id: 'user-123',
    description: 'Wi-Fi 6 Access Point',
    datasheet_url: 'http://example.com/wifi6-datasheet',
    images_url: 'http://example.com/wifi6-image',
    part_number: 'WIFI6-AP-001',
    manufacturer: 'NetworkGear Inc.',
    managed: 'yes',
    warranty: 24,
    specification: 'Wi-Fi 6, up to 1.2Gbps',
  },
  {
    id: 'access-2',
    inventory_id: 'AC002',
    category: 'ACCESS',
    type: 'access',
    user_id: 'user-123',
    description: '5G Small Cell',
    datasheet_url: 'http://example.com/5g-smallcell-datasheet',
    images_url: 'http://example.com/5g-smallcell-image',
    part_number: '5G-SC-001',
    manufacturer: 'TelecomTech Corp.',
    managed: 'yes',
    warranty: 36,
    specification: '5G NR, up to 1Gbps',
  },
];
const dummysComponents = [
  {
    id: '123e4567-e89b-12d3-a456-426614175000',
    inventory_id: '100',
    category: 'ACCESS',
    type: 'access',
    user_id: '3cda3e88-cde2-447c-8b29-a2c931c3fe2f',
    description: 'A 10k ohm resistor',
    datasheet_url: 'http://example.com/datasheet1',
    images_url: 'http://example.com/image1',
    part_number: 'R-10K',
    manufacturer: 'Resistors Inc.',
    managed: 'yes',
    warranty: 24,
    specification: '10k ohm, 5% tolerance',
  },
  {
    id: '123e4567-e89b-12d3-a456-426614174002',
    inventory_id: '200',
    category: 'BACKHAUL',
    type: 'backhaul',
    user_id: '3cda3e88-cde2-447c-8b29-a2c931c3fe2f',
    description: 'A 100uF capacitor',
    datasheet_url: 'http://example.com/datasheet2',
    images_url: 'http://example.com/image2',
    part_number: 'C-100uF',
    manufacturer: 'Capacitors Ltd.',
    managed: 'no',
    warranty: 12,
    specification: '100uF, 16V',
  },
  {
    id: '123e4567-e89b-12d3-a456-426614176004',
    inventory_id: '150',
    category: 'POWER',
    type: 'power',
    user_id: '3cda3e88-cde2-447c-8b29-a2c931c3fe2f',
    description: 'A 10mH inductor',
    datasheet_url: 'http://example.com/datasheet3',
    images_url: 'http://example.com/image3',
    part_number: 'L-10mH',
    manufacturer: 'Inductors Co.',
    managed: 'yes',
    warranty: 36,
    specification: '10mH, 10% tolerance',
  },
  {
    id: '123e4567-e89b-12d3-a456-426614174006',
    inventory_id: '250',
    category: 'SWITCH',
    type: 'switch',
    user_id: '3cda3e88-cde2-447c-8b29-a2c931c3fe2f',
    description: 'A 1N4007 diode',
    datasheet_url: 'http://example.com/datasheet4',
    images_url: 'http://example.com/image4',
    part_number: 'D-1N4007',
    manufacturer: 'Diodes Inc.',
    managed: 'no',
    warranty: 12,
    specification: '1N4007, 1A, 1000V',
  },
];

const Sites = () => {
  const [open, setOpen] = useState(false);
  const [sitesList, setSitesList] = useState<any[]>([]);
  const [componentsList, setComponentsList] = useState<any[]>([]);
  const [isLoading, setIsLoading] = useState(true);
  const { setSnackbarMessage } = useAppContext();
  const [opens, setSopen] = useState(false);

  const handleOpenTest = () => setSopen(true);
  const handleCloseTest = () => setSopen(false);

  const handleOpen = () => {
    setOpen(true);
  };

  const handleClose = () => {
    setOpen(false);
  };

  const { data: networkList, loading: networkLoading } = useGetNetworksQuery({
    fetchPolicy: 'cache-and-network',
    onError: (error) => {
      setSnackbarMessage({
        id: 'networks-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });

  const [getSites] = useGetSitesLazyQuery({
    onError: (error) => {
      setSnackbarMessage({
        id: 'sites-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  const [getComponents] = useGetComponentsByUserIdLazyQuery({
    onCompleted: (res) => {
      if (res.getComponentsByUserId) {
        setComponentsList(res.getComponentsByUserId.components);
      }
    },
    onError: (error) => {
      setSnackbarMessage({
        id: 'components-msg',
        message: error.message,
        type: 'error' as AlertColor,
        show: true,
      });
    },
  });
  useEffect(() => {
    getComponents({ variables: { category: 'switch' } });
  }, []);
  const dummyNetworks = [
    { id: '1', name: 'Network A' },
    { id: '2', name: 'Network B' },
    { id: '3', name: 'Network C' },
  ];

  useEffect(() => {
    const fetchAllSites = async () => {
      if (networkList && networkList.getNetworks.networks) {
        setIsLoading(true);
        const allSitesPromises = networkList.getNetworks.networks.map(
          (network) => getSites({ variables: { networkId: network.id } }),
        );

        try {
          const results = await Promise.all(allSitesPromises);
          const allSites = results.flatMap(
            (result) => result.data?.getSites.sites || [],
          );
          setSitesList(allSites);
        } catch (error) {
          console.error('Error fetching sites:', error);
          setSnackbarMessage({
            id: 'sites-fetch-error',
            message: 'Error fetching sites',
            type: 'error' as AlertColor,
            show: true,
          });
        } finally {
          setIsLoading(false);
        }
      }
    };

    fetchAllSites();
  }, [networkList, getSites]);

  const handleMenuClick = (siteId: string) => {
    console.log(`Menu clicked for siteId: ${siteId}`);
  };

  const handleFormDataSubmit = (formData: any) => {
    console.log('SUBMIITED');
  };
  console.log('SITES LIST', sitesList);
  return (
    <Grid container spacing={0} sx={{ mt: 1 }}>
      <Grid item xs={12}>
        <Paper sx={{ p: 4 }}>
          <Grid container spacing={0} sx={{ mb: 2 }}>
            <Grid item xs={6}>
              <Typography variant="h6" color="initial">
                My sites
              </Typography>
            </Grid>
            <Grid
              item
              xs={6}
              container
              justifyItems={'center'}
              justifyContent={'flex-end'}
            >
              <Button
                variant="contained"
                color="primary"
                onClick={handleOpenTest}
              >
                ADD SITE
              </Button>
            </Grid>
          </Grid>

          <Grid container spacing={2}>
            {sitesList.length === 0 ? (
              <Grid item xs={12} style={{ textAlign: 'center' }}>
                <Typography variant="body1">No sites available.</Typography>
              </Grid>
            ) : (
              sitesList.map((site, index) => (
                <Grid item xs={12} md={4} lg={4} key={index}>
                  <SiteCard
                    siteId={site.id}
                    name={site.name}
                    address={site.location}
                    users={site.users || ''}
                    siteStatus={site.isDeactivated}
                    onClickMenu={handleMenuClick}
                    loading={isLoading || networkLoading}
                    status={{
                      online: false,
                      charging: false,
                      signal: '',
                    }}
                  />
                </Grid>
              ))
            )}
          </Grid>
        </Paper>
      </Grid>
      <StepperDialog
        open={opens}
        onClose={handleCloseTest}
        components={dummysComponents}
        networks={dummyNetworks}
      />

      <SiteConfigurationStepperDialog
        open={open}
        handleClose={handleClose}
        components={dummyComponents}
        handleFormDataSubmit={handleFormDataSubmit}
      />
    </Grid>
  );
};

export default Sites;
