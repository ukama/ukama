import { HealtChartsConfigure, TooltipsText } from '@/constants';
import { Node, NodeTypeEnum } from '@/generated';
import { MetricsRes } from '@/generated/metrics';
import { Grid, Paper, Stack, Typography, capitalize } from '@mui/material';
import { useEffect, useState } from 'react';
import LineChart from '../LineChart';
import NodeDetailsCard from '../NodeDetailsCard';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';
import { getMetricValue } from '@/utils';

interface INodeOverviewTab {
  metricFrom: number;
  metrics: MetricsRes;
  loading: boolean;
  metricsLoading: boolean;
  onNodeSelected: Function;
  uptime: number | undefined;
  isUpdateAvailable: boolean;
  handleUpdateNode: Function;
  selectedNode: Node | undefined;
  connectedUsers: string | undefined;
  getNodeSoftwareUpdateInfos: Function;
}

const NodeOverviewTab = ({
  metrics,
  uptime,
  loading,
  metricFrom,
  selectedNode,
  metricsLoading,
  connectedUsers = '0',
  onNodeSelected,
  handleUpdateNode,
  isUpdateAvailable,
  getNodeSoftwareUpdateInfos,
}: INodeOverviewTab) => {
  const nodeType = selectedNode?.type || NodeTypeEnum.Hnode;
  const [selected, setSelected] = useState<number>(0);
  useEffect(() => {
    setSelected(0);
  }, [selectedNode]);

  const handleOnSelected = (value: number) => setSelected(value);

  return (
    <Grid container columnSpacing={3}>
      <Grid item xs={12} md={4}>
        <Stack spacing={2}>
          <NodeStatsContainer
            index={0}
            loading={loading}
            isClickable={true}
            selected={selected}
            title={'Node Information'}
            handleAction={handleOnSelected}
          >
            <NodeStatItem
              value={`${capitalize(
                selectedNode?.type.toLowerCase() || 'HOME',
              )} Node`}
              name={'Model type'}
            />

            <NodeStatItem
              value={selectedNode?.id.toLowerCase() || '-'}
              name={'Serial #'}
            />
            {/* {selectedNode?.type === 'TOWER' && (
                <Grid item xs={12}>
                  <NodeGroup
                    nodes={nodeGroupData?.attached || []}
                    loading={nodeGroupLoading}
                    handleNodeAction={onNodeSelected}
                  />
                </Grid>
              )} */}
          </NodeStatsContainer>
          <NodeStatsContainer
            index={1}
            loading={loading}
            isClickable={true}
            selected={selected}
            title={'Node Health'}
            handleAction={handleOnSelected}
          >
            {HealtChartsConfigure[nodeType][0].show && (
              <NodeStatItem
                value={'24 °C'}
                name={HealtChartsConfigure[nodeType][0].name}
                showAlertInfo={false}
                nameInfo={TooltipsText.TRX}
              />
            )}
            {HealtChartsConfigure[nodeType][1].show && (
              <NodeStatItem
                value={'22 °C'}
                name={HealtChartsConfigure[nodeType][1].name}
                nameInfo={TooltipsText.COM}
              />
            )}
            {HealtChartsConfigure[nodeType][2].show && (
              <NodeStatItem
                name={HealtChartsConfigure[nodeType][2].name}
                nameInfo={TooltipsText.COM}
                value={uptime ? `${Math.floor(uptime / 60 / 60)} hours` : 'NA'}
              />
            )}
          </NodeStatsContainer>
          {selectedNode?.type !== NodeTypeEnum.Anode && (
            <NodeStatsContainer
              index={2}
              loading={loading}
              isClickable={true}
              selected={selected}
              title={'Subscribers'}
              handleAction={handleOnSelected}
            >
              <NodeStatItem
                name={'Attached'}
                value={connectedUsers}
                nameInfo={TooltipsText.ATTACHED}
              />
              <NodeStatItem
                name={'Active'}
                value={`${
                  connectedUsers === '0'
                    ? parseInt(connectedUsers)
                    : parseInt(connectedUsers) - 1
                }`}
                nameInfo={TooltipsText.ACTIVE}
              />
            </NodeStatsContainer>
          )}
        </Stack>
      </Grid>
      <Grid item xs={12} md={8}>
        {selected === 0 && (
          <NodeDetailsCard
            nodeType={selectedNode?.type || undefined}
            getNodeUpdateInfos={getNodeSoftwareUpdateInfos}
            loading={loading}
            nodeTitle={selectedNode?.name || 'HOME'}
            handleUpdateNode={handleUpdateNode}
            isUpdateAvailable={isUpdateAvailable}
          />
        )}
        {selected === 1 && (
          <Paper sx={{ p: 3 }}>
            <Stack spacing={4}>
              <Typography variant="h6">Node Health</Typography>
              {HealtChartsConfigure[nodeType][0].show && (
                <LineChart
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][0].id,
                    metrics,
                  )}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  hasData={metrics.metrics.length > 0 || false}
                  topic={HealtChartsConfigure[nodeType][0].id}
                  title={HealtChartsConfigure[nodeType][0].name}
                />
              )}
              {HealtChartsConfigure[nodeType][1].show && (
                <LineChart
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][1].id,
                    metrics,
                  )}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  hasData={metrics.metrics.length > 0}
                  topic={HealtChartsConfigure[nodeType][1].id}
                  title={HealtChartsConfigure[nodeType][1].name}
                />
              )}
              {HealtChartsConfigure[nodeType][2].show && (
                <LineChart
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][2].id,
                    metrics,
                  )}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  hasData={metrics.metrics.length > 0}
                  topic={HealtChartsConfigure[nodeType][2].id}
                  title={HealtChartsConfigure[nodeType][2].name}
                />
              )}
            </Stack>
          </Paper>
        )}
        {selected === 2 && nodeType !== NodeTypeEnum.Anode && (
          <Paper sx={{ p: 3 }}>
            <Stack spacing={4}>
              <Typography variant="h6">Subscribers</Typography>
              {HealtChartsConfigure[
                (selectedNode?.type as string) || 'hnode'
              ][4].show && (
                <LineChart
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][4].id,
                    metrics,
                  )}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  hasData={metrics.metrics.length > 0}
                  topic={HealtChartsConfigure[nodeType][4].id}
                  title={HealtChartsConfigure[nodeType][4].name}
                />
              )}
              {HealtChartsConfigure[
                (selectedNode?.type as string) || 'hnode'
              ][5].show && (
                <LineChart
                  initData={getMetricValue(
                    HealtChartsConfigure[nodeType][5].id,
                    metrics,
                  )}
                  metricFrom={metricFrom}
                  loading={metricsLoading}
                  hasData={metrics.metrics.length > 0}
                  topic={HealtChartsConfigure[nodeType][5].id}
                  title={HealtChartsConfigure[nodeType][5].name}
                />
              )}
            </Stack>
          </Paper>
        )}
      </Grid>
    </Grid>
  );
};

export default NodeOverviewTab;
