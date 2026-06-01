/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { SiteDto } from '@/client/graphql/generated';
import { MetricsStateRes } from '@/client/graphql/generated/subscriptions';
import LoadingWrapper from '@/components/ui/LoadingWrapper';
import SiteCard from '@/app/console/sites/_components/SiteCard';
import UnassignedNodeCard from '@/app/console/nodes/_components/UnassingedNodecard';
import { Box, Divider, Grid, Typography } from '@mui/material';
import React from 'react';

interface NodeStatus {
  connectivity: string;
  state: string;
}

interface NodeDto {
  id: string;
  name: string;
  type: string;
  status: NodeStatus;
  site: { nodeId?: string | null; siteId?: string | null } | null;
  attached?: NodeDto[];
  latitude?: string | number;
  longitude?: string | number;
}

interface SitesWrapperProps {
  sites: SiteDto[];
  unassignedNodes: NodeDto[];
  loading: boolean;
  nodesLoading?: boolean;
  handleAddSite?: () => void;
  handleSiteNameUpdate: (siteId: string, siteName: string) => void;
  handleConfigureNode: (nodeId: string) => void;
  siteMetricsStatData: MetricsStateRes;
}

const SitesWrapper: React.FC<SitesWrapperProps> = ({
  sites,
  unassignedNodes,
  loading,
  nodesLoading = false,
  siteMetricsStatData,
  handleSiteNameUpdate,
  handleConfigureNode,
}) => {
  const showEmptyState =
    sites?.length === 0 && unassignedNodes?.length === 0 && !loading;
  if (showEmptyState) {
    return (
      <Box
        sx={{
          height: '100%',
          display: 'flex',
          alignItems: 'center',
          justifyContent: 'center',
          flexDirection: 'column',
          gap: 2,
          padding: '0 20px',
          textAlign: 'center',
        }}
      >
        <Typography variant="h6" color="textSecondary">
          No sites yet!
        </Typography>
        <Typography
          variant="body1"
          color="textSecondary"
          sx={{ maxWidth: '450px' }}
        >
          A site is a complete connection point to the network, made up of your
          Ukama node, and the power and backhaul components. Install a site to
          get started.
        </Typography>
      </Box>
    );
  }

  return (
    <Box
      sx={{
        height: '100%',
        overflowY: 'auto',
      }}
    >
      {sites && sites.length > 0 && (
        <LoadingWrapper isLoading={loading} height="auto">
          <Box sx={{ padding: '10px' }}>
            <Grid container spacing={2}>
              {sites.map((site) => (
                <Grid item xs={12} md={4} lg={4} key={site.id}>
                  <SiteCard
                    siteId={site.id}
                    name={site.name}
                    address={site.location}
                    loading={loading}
                    handleSiteNameUpdate={handleSiteNameUpdate}
                    metricsData={siteMetricsStatData}
                  />
                </Grid>
              ))}
            </Grid>
          </Box>
        </LoadingWrapper>
      )}

      {unassignedNodes && unassignedNodes.length > 0 && (
        <>
          {sites && sites.length > 0 && <Divider sx={{ my: 3 }} />}
          <LoadingWrapper isLoading={nodesLoading} height="auto">
            <Box sx={{ padding: '10px' }}>
              <Typography
                variant="subtitle1"
                color="initial"
                sx={{ paddingLeft: '12px', mb: 2, fontWeight: 'bold' }}
              >
                Unassigned Nodes
              </Typography>
              <Grid container spacing={2}>
                {unassignedNodes.map((node) => (
                  <Grid item xs={12} md={4} lg={4} key={node.id}>
                    <UnassignedNodeCard
                      id={node.id}
                      name={node.name || `Node-${node.id.substring(0, 8)}`}
                      loading={nodesLoading}
                      handleConfigureNode={handleConfigureNode}
                    />
                  </Grid>
                ))}
              </Grid>
            </Box>
          </LoadingWrapper>
        </>
      )}
    </Box>
  );
};

export default SitesWrapper;
