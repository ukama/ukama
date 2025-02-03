import {
  Graphs_Type,
  MetricsRes,
} from '@/client/graphql/generated/subscriptions';
import { getComponentHealth, getMetricValue, isMetricValue } from '@/utils';
import {
  Box,
  Grid,
  Skeleton,
  Typography,
  useTheme,
  useMediaQuery,
  List,
  ListItem,
  ListItemText,
  Stack,
} from '@mui/material';
import LineChart from '../LineChart';
import {
  BatteryChartsConfig,
  SolarChartsConfig,
  ControllerChartsConfig,
  BackhaulChartsConfig,
  SwitchChartConfig,
} from '@/constants';
import { SiteHealth } from '@/../public/svg';
import { useState, useMemo, useCallback, useRef, useEffect } from 'react';
import { Nodes } from '@/client/graphql/generated';

interface ISiteOverallHealth {
  siteId: string;
  metrics: MetricsRes;
  loading?: boolean;
  metricFrom: number;
  onGraphTypeChange: (type: Graphs_Type) => void;
  nodes: Nodes | undefined;
  backhaulComponent: string;
}

const SiteOverallHealth: React.FC<ISiteOverallHealth> = ({
  siteId,
  metrics,
  loading,
  metricFrom,
  onGraphTypeChange,
  nodes,
  backhaulComponent,
}) => {
  const theme = useTheme();
  const isMobile = useMediaQuery(theme.breakpoints.down('md'));
  const [selectedComponent, setSelectedComponent] = useState<string>('battery');
  const [isTransitioning, setIsTransitioning] = useState(false);
  const transitionTimeoutRef = useRef<NodeJS.Timeout>();
  const lastClickTimeRef = useRef<number>(0);
  const DEBOUNCE_TIME = 600;

  const componentConfigMap = useMemo(
    () => ({
      battery: { config: BatteryChartsConfig, type: Graphs_Type.Battery },
      solar: { config: SolarChartsConfig, type: Graphs_Type.Solar },
      controller: {
        config: ControllerChartsConfig,
        type: Graphs_Type.Controller,
      },
      backhaul: { config: BackhaulChartsConfig, type: Graphs_Type.Backhaul },
      switch: { config: SwitchChartConfig, type: Graphs_Type.Switch },
    }),
    [],
  );

  const getComponentConfig = useCallback(
    (component: string) => {
      return (
        componentConfigMap[component as keyof typeof componentConfigMap] ||
        componentConfigMap.battery
      );
    },
    [componentConfigMap],
  );

  const currentConfig = useMemo(
    () => getComponentConfig(selectedComponent).config,
    [selectedComponent, getComponentConfig],
  );

  useEffect(() => {
    return () => {
      if (transitionTimeoutRef.current) {
        clearTimeout(transitionTimeoutRef.current);
      }
    };
  }, []);

  const handleComponentClick = useCallback(
    (component: string) => {
      const now = Date.now();
      if (now - lastClickTimeRef.current < DEBOUNCE_TIME) {
        return;
      }
      lastClickTimeRef.current = now;

      setIsTransitioning(true);
      const { type } = getComponentConfig(component);

      setSelectedComponent(component);
      onGraphTypeChange(type);

      if (transitionTimeoutRef.current) {
        clearTimeout(transitionTimeoutRef.current);
      }

      transitionTimeoutRef.current = setTimeout(() => {
        setIsTransitioning(false);
      }, 500);
    },
    [getComponentConfig, onGraphTypeChange],
  );

  const getComponentDisplayName = (component: string) => {
    const names = {
      battery: 'Battery',
      solar: 'Solar',
      controller: 'Controller',
      backhaul: 'Backhaul',
      switch: 'Switch',
      node: 'Nodes',
    };
    return `${names[component as keyof typeof names] || 'Component'} ${component === 'node' ? 'List' : 'KPIs'}`;
  };

  const getBackhaulStatus = useCallback(() => {
    const status = getMetricValue('backhaul_status', metrics);
    const latency = getMetricValue('backhaul_latency', metrics);
    const speed = getMetricValue('backhaul_speed', metrics);

    return {
      status: status ? 'connected' : 'disconnected',
      latency: latency || 0,
      speed: speed || 0,
    };
  }, [metrics]);

  const renderContent = () => {
    if (isTransitioning || loading) {
      return <Skeleton variant="rectangular" height={350} />;
    }

    if (selectedComponent === 'node') {
      return (
        <Box sx={{ p: 2, border: 1, borderColor: 'divider', borderRadius: 2 }}>
          <List>
            {nodes?.nodes.map((node) => (
              <ListItem key={node.id}>
                <ListItemText primary={`#${node.id}`} />
              </ListItem>
            ))}
          </List>
        </Box>
      );
    }

    if (selectedComponent === 'backhaul') {
      const backhaulStatus = getBackhaulStatus();
      return (
        <Box sx={{ p: 2, border: 1, borderColor: 'divider', borderRadius: 2 }}>
          <Stack direction="row" spacing={2} sx={{ mb: 2 }}>
            <Typography variant="body2" gutterBottom>
              {backhaulComponent}
            </Typography>
            <Typography variant="body2" gutterBottom>
              {backhaulStatus.status === 'connected' ? 'Online' : 'Offline'}
            </Typography>
          </Stack>
          <Typography variant="body2">
            Speed: {backhaulStatus.speed} Mbps
          </Typography>
        </Box>
      );
    }

    return (
      <Grid container spacing={2}>
        {currentConfig.map((chartConfig: any) => (
          <Grid item xs={12} key={chartConfig.id}>
            <Box
              sx={{
                minHeight: '350px',
                height: 'auto',
                border: 1,
                borderColor: 'divider',
                borderRadius: 2,
                p: 2,
                mb: 2,
                display: 'flex',
                flexDirection: 'column',
              }}
            >
              <Box sx={{ flexGrow: 1, minHeight: 0 }}>
                <LineChart
                  siteId={siteId}
                  loading={loading}
                  metricFrom={metricFrom}
                  topic={chartConfig.id}
                  title={chartConfig.name}
                  tabSection={getComponentConfig(selectedComponent).type}
                  initData={getMetricValue(chartConfig.id, metrics)}
                  hasData={isMetricValue(chartConfig.id, metrics)}
                />
              </Box>
            </Box>
          </Grid>
        ))}
      </Grid>
    );
  };
  const healthStatuses = useMemo(
    () => ({
      batteryHealth: getComponentHealth('battery', metrics),
      solarHealth: getComponentHealth('solar', metrics),
      switchHealth: getComponentHealth('switch', metrics),
      controllerHealth: getComponentHealth('controller', metrics),
      nodeHealth: getComponentHealth('node', metrics),
      backhaulHealth: getComponentHealth('backhaul', metrics),
    }),
    [metrics],
  );
  return (
    <Grid container spacing={2}>
      <Grid container item xs={12} spacing={2}>
        <Grid item xs={12} md={7}>
          <Box sx={{ height: '100%' }}>
            <Typography variant="h6" sx={{ mb: 2 }}>
              Site Overview
            </Typography>
            <Box>
              <SiteHealth
                solarHealth={healthStatuses.solarHealth}
                nodeHealth={healthStatuses.nodeHealth}
                switchHealth={healthStatuses.switchHealth}
                controllerHealth={healthStatuses.controllerHealth}
                batteryHealth={healthStatuses.batteryHealth}
                backhaulHealth={healthStatuses.backhaulHealth}
                onNodeClick={() => handleComponentClick('node')}
                onSolarClick={() => handleComponentClick('solar')}
                onSwitchClick={() => handleComponentClick('switch')}
                onControllerClick={() => handleComponentClick('controller')}
                onBatteryClick={() => handleComponentClick('battery')}
                onBackhaulClick={() => handleComponentClick('backhaul')}
                selectedKPI={selectedComponent}
              />
            </Box>
          </Box>
        </Grid>
        <Grid item xs={12} md={5}>
          <Box
            sx={{
              height: isMobile ? 'auto' : '40vh',
              maxHeight: '40vh',
              overflowY: 'auto',
              p: 1,
              '&::-webkit-scrollbar': {
                width: '8px',
              },
              '&::-webkit-scrollbar-track': {
                background: theme.palette.background.default,
              },
              '&::-webkit-scrollbar-thumb': {
                background: theme.palette.divider,
                borderRadius: '4px',
              },
              scrollbarWidth: 'thin',
              scrollBehavior: 'smooth',
            }}
          >
            <Typography
              variant="h6"
              sx={{
                mb: 2,
                pl: 1,
                position: 'sticky',
                top: 0,
                zIndex: 1,
              }}
            >
              {getComponentDisplayName(selectedComponent)}
            </Typography>
            {renderContent()}
          </Box>
        </Grid>
      </Grid>
    </Grid>
  );
};

export default SiteOverallHealth;
