import { NodeResourcesTabConfigure, TooltipsText } from '@/constants';
import { Node, NodeTypeEnum } from '@/generated';
import { Grid, Paper } from '@mui/material';
import { useState } from 'react';
import NodeStatItem from '../NodeStatItem';
import NodeStatsContainer from '../NodeStatsContainer';

const PLACEHOLDER_VALUE = 'NA';
interface INodeResourcesTab {
  metrics: any;
  loading: boolean;
  selectedNode: Node | undefined;
}
const NodeResourcesTab = ({
  metrics,
  loading,
  selectedNode,
}: INodeResourcesTab) => {
  const nodeType = selectedNode?.type || NodeTypeEnum.Hnode;
  const [isCollapse, setIsCollapse] = useState<boolean>(false);
  const handleCollapse = () => setIsCollapse((prev) => !prev);
  return (
    <Grid container spacing={3}>
      <Grid item lg={!isCollapse ? 4 : 1} md xs>
        <NodeStatsContainer
          index={0}
          selected={0}
          loading={loading}
          title={'Resources'}
          isCollapsable={true}
          isCollapse={isCollapse}
          onCollapse={handleCollapse}
        >
          {NodeResourcesTabConfigure[nodeType][0].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][0].name}
              nameInfo={TooltipsText.MTRX}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][1].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][1].name}
              nameInfo={TooltipsText.MCOM}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][2].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              name={NodeResourcesTabConfigure[nodeType][2].name}
              variant={'large'}
              nameInfo={TooltipsText.CPUTRX}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][3].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              name={NodeResourcesTabConfigure[nodeType][3].name}
              variant={'large'}
              nameInfo={TooltipsText.CPUCOM}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][4].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][4].name}
              nameInfo={TooltipsText.DISKTRX}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][5].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              variant={'large'}
              name={NodeResourcesTabConfigure[nodeType][5].name}
              nameInfo={TooltipsText.DISKCOM}
            />
          )}
          {NodeResourcesTabConfigure[nodeType][6].show && (
            <NodeStatItem
              value={PLACEHOLDER_VALUE}
              name={NodeResourcesTabConfigure[nodeType][6].name}
              variant={'large'}
              nameInfo={TooltipsText.POWER}
            />
          )}
        </NodeStatsContainer>
      </Grid>
      <Grid item lg={isCollapse ? 11 : 8} md xs>
        <Paper sx={{ p: 3, width: '100%' }}>
          {/* <Stack spacing={4}>
            {NodeResourcesTabConfigure[nodeType][0]
              .show && (
              <ApexLineChart
                data={
                  metrics[
                    NodeResourcesTabConfigure[
                      nodeType
                    ][0].id
                  ]
                }
              />
            )}
            {NodeResourcesTabConfigure[nodeType][1]
              .show && (
              <ApexLineChart
                data={
                  metrics[
                    NodeResourcesTabConfigure[
                      nodeType
                    ][1].id
                  ]
                }
              />
            )}
            {NodeResourcesTabConfigure[nodeType][2]
              .show && (
              <ApexLineChart
                data={
                  metrics[
                    NodeResourcesTabConfigure[
                      nodeType
                    ][2].id
                  ]
                }
              />
            )}
            {NodeResourcesTabConfigure[nodeType][3]
              .show && (
              <ApexLineChart
                data={
                  metrics[
                    NodeResourcesTabConfigure[
                      nodeType
                    ][3].id
                  ]
                }
              />
            )}
            {NodeResourcesTabConfigure[nodeType][4]
              .show && (
              <ApexLineChart
                data={
                  metrics[
                    NodeResourcesTabConfigure[
                      nodeType
                    ][4].id
                  ]
                }
              />
            )}
            {NodeResourcesTabConfigure[nodeType][5]
              .show && (
              <ApexLineChart
                data={
                  metrics[
                    NodeResourcesTabConfigure[
                      nodeType
                    ][5].id
                  ]
                }
              />
            )}
            {NodeResourcesTabConfigure[nodeType][6]
              .show && (
              <ApexLineChart
                data={
                  metrics[
                    NodeResourcesTabConfigure[
                      nodeType
                    ][6].id
                  ]
                }
              />
            )}
          </Stack> */}
        </Paper>
      </Grid>
    </Grid>
  );
};

export default NodeResourcesTab;
