/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { MetricsRes } from '@/client/graphql/generated/subscriptions';
import { SiteKpiConfig } from '@/constants';
import { getMetricValue, isMetricValue } from '@/utils';
import LineChart from '@/components/ui/LineChart';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Box,
  Stack,
  Switch,
  Typography,
} from '@mui/material';

export interface PortGroup {
  id: string;
  portNumber: number;
  description: string;
  metrics: SiteKpiConfig[];
}

interface SwitchPortItemProps {
  portGroup: PortGroup;
  isExpanded: boolean;
  isOn: boolean;
  isDisabled: boolean;
  metrics: MetricsRes;
  metricFrom: number;
  metricsLoading: boolean;
  onToggleExpand: (portId: string) => void;
  onToggleSwitch: (portNumber: number, currentIsOn: boolean) => void;
}

/**
 * Renders a single collapsible switch-port row with its status toggle and
 * metric line charts. Extracted from SiteComponents to keep that file focused
 * on orchestration rather than per-port rendering.
 */
const SwitchPortItem: React.FC<SwitchPortItemProps> = ({
  portGroup,
  isExpanded,
  isOn,
  isDisabled,
  metrics,
  metricFrom,
  metricsLoading,
  onToggleExpand,
  onToggleSwitch,
}) => {
  const portTitle = portGroup.description
    ? `Port ${portGroup.portNumber} (${portGroup.description})`
    : `Port ${portGroup.portNumber}`;

  const statusMetric = portGroup.metrics.find((m) =>
    m.id.includes('switch_port_status'),
  );

  return (
    <Box
      sx={{ borderBottom: '1px solid rgba(0, 0, 0, 0.12)', py: 2 }}
      data-testid={`port-${portGroup.id}-container`}
    >
      <Accordion
        expanded={isExpanded}
        onChange={() => onToggleExpand(portGroup.id)}
        sx={{ boxShadow: 'none' }}
        data-testid={`accordion-${portGroup.id}`}
      >
        <AccordionSummary
          expandIcon={<ExpandMoreIcon />}
          sx={{ display: 'flex', alignItems: 'center', p: 0 }}
          data-testid={`accordion-summary-${portGroup.id}`}
        >
          <Typography
            variant="subtitle1"
            fontWeight="medium"
            sx={{ flexGrow: 1 }}
            data-testid={`accordion-title-${portGroup.id}`}
          >
            {portTitle}
          </Typography>
        </AccordionSummary>

        <AccordionDetails
          sx={{ mt: 2, ml: 2, p: 0 }}
          data-testid={`accordion-details-${portGroup.id}`}
        >
          {statusMetric && (
            <Box
              sx={{
                display: 'flex',
                justifyContent: 'space-between',
                alignItems: 'center',
                mb: 3,
                pb: 2,
                borderBottom: '1px solid rgba(0, 0, 0, 0.06)',
              }}
              data-testid={`status-metric-${portGroup.id}`}
            >
              <Typography variant="body1">{statusMetric.name}</Typography>
              <Box display="flex" alignItems="center">
                <Typography variant="body2" sx={{ mr: 1 }}>
                  {isOn ? 'On' : 'Off'}
                </Typography>
                <Switch
                  checked={isOn}
                  onChange={() => onToggleSwitch(portGroup.portNumber, isOn)}
                  disabled={isDisabled}
                  color="primary"
                  data-testid={`toggle-switch-${portGroup.id}`}
                />
              </Box>
            </Box>
          )}

          <Stack spacing={3}>
            {portGroup.metrics
              .filter((m) => !m.id.includes('switch_port_status'))
              .map((metric) => (
                <Box key={metric.id} data-testid={`metric-item-${metric.id}`}>
                  <Typography variant="body1" sx={{ mb: 1 }}>
                    {metric.name}
                  </Typography>
                  <LineChart
                    from={metricFrom}
                    topic={metric.id}
                    yunit={metric.unit}
                    loading={metricsLoading}
                    tickInterval={metric.tickInterval}
                    tickPositions={metric.tickPositions}
                    hasData={isMetricValue(metric.id, metrics)}
                    initData={getMetricValue(metric.id, metrics)}
                    format={metric.format}
                    data-testid={`line-chart-${metric.id}`}
                  />
                </Box>
              ))}
          </Stack>
        </AccordionDetails>
      </Accordion>
    </Box>
  );
};

export default SwitchPortItem;
