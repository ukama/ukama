/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  Node,
  NodeConnectivityEnum,
  NodeStateEnum,
  useGetNodeLazyQuery,
  useGetNodesByStateLazyQuery,
  useGetNodeStateLazyQuery,
} from '@/client/graphql/generated';
import InstallSiteLoading from '@/components/InstallSiteLoading';
import {
  CHECK_SITE_FLOW,
  INSTALLATION_FLOW,
  NETWORK_FLOW,
  ONBOARDING_FLOW,
} from '@/constants';
import { useAppContext } from '@/context';
import { HorizontalContainerJustify } from '@/styles/global';
import { Button, Stack } from '@mui/material';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';

const Check = () => {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const [node, setNode] = useState<Node | undefined>(undefined);
  const nodeId = searchParams.get('nid') ?? '';
  const [duration, setDuration] = useState(5);
  const flow = searchParams.get('flow') ?? INSTALLATION_FLOW;
  const [showReturn, setShowReturn] = useState(false);
  const [title] = useState(
    flow === NETWORK_FLOW
      ? 'Creating your network...'
      : flow === CHECK_SITE_FLOW
        ? 'Checking for site availability to configure'
        : 'Loading up your site...',
  );
  const [subtitle, setSubtitle] = useState(
    flow === NETWORK_FLOW ? 'Loading up your network...' : '',
  );
  const [description, setDescription] = useState('');
  const { setSnackbarMessage } = useAppContext();

  const setQueryParam = (key: string, value: string) => {
    const p = new URLSearchParams(searchParams.toString());
    p.set(key, value);
    window.history.replaceState({}, '', `${pathname}?${p.toString()}`);
    return p;
  };

  const [getNodesByState] = useGetNodesByStateLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      const filterNodes = data.getNodesByState.nodes.filter(
        (node) =>
          node.latitude !== 0 &&
          node.longitude !== 0 &&
          node.status.connectivity === NodeConnectivityEnum.Online &&
          node.status.state === NodeStateEnum.Unknown,
      );
      if (
        filterNodes.length > 0 &&
        filterNodes[0].latitude !== 0 &&
        filterNodes[0].longitude !== 0 &&
        filterNodes[0].status.connectivity === NodeConnectivityEnum.Online &&
        filterNodes[0].status.state === NodeStateEnum.Unknown
      ) {
        setNode(filterNodes[0]);
      }
    },
  });

  const [getNode] = useGetNodeLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: async (data) => {
      if (data.getNode.latitude && data.getNode.longitude && nodeId) {
        if (
          data.getNode.status.connectivity === NodeConnectivityEnum.Online &&
          data.getNode.status.state === NodeStateEnum.Unknown
        ) {
          setNode(data.getNode);
        } else {
          setSnackbarMessage({
            id: 'node-configured-warn',
            message: `Node ${data.getNode.id} is already configured.`,
            type: 'warning',
            show: true,
          });
          router.push(`/console/home`);
        }
      }
    },
  });

  const [getNodeState] = useGetNodeStateLazyQuery({
    fetchPolicy: 'network-only',
    onCompleted: (data) => {
      if (node && data.getNodeState.currentState === NodeStateEnum.Unknown) {
        let p = setQueryParam('lat', node.latitude.toString());
        p.set('lng', node.longitude.toString());
        p.set(
          'flow',
          flow === NETWORK_FLOW
            ? ONBOARDING_FLOW
            : flow === CHECK_SITE_FLOW
              ? INSTALLATION_FLOW
              : flow,
        );
        p.delete('nid');
        router.push(`/configure/node/${node.id}?${p.toString()}`);
      }
    },
  });

  useEffect(() => {
    if (nodeId) {
      getNode({
        variables: {
          data: {
            id: nodeId,
          },
        },
      });
      getNodesByState({
        variables: {
          data: {
            state: NodeStateEnum.Unknown,
            connectivity: NodeConnectivityEnum.Online,
          },
        },
      });
    }
  }, [nodeId]);

  useEffect(() => {
    if (node?.id) {
      getNodeState({
        variables: {
          getNodeStateId: node.id,
        },
      });
    }
  }, [node]);

  const onInstallProgressComplete = () => {
    if (flow !== NETWORK_FLOW) {
      setShowReturn(true);
      setSubtitle('❗ Site not detected');
      setDescription(
        'It is taking longer than usual to load up your site. Please check on your site to make sure that all parts are installed correctly.',
      );
    } else {
      const p = setQueryParam('flow', ONBOARDING_FLOW);
      router.push(`/configure?step=2&${p}`);
    }
  };

  const handleBack = () => {
    router.push(
      `/configure/sims?step=${flow === ONBOARDING_FLOW ? 5 : 4}&flow=${flow === NETWORK_FLOW ? ONBOARDING_FLOW : flow === CHECK_SITE_FLOW ? INSTALLATION_FLOW : flow}`,
    );
  };

  const handleRetry = () => {
    if (nodeId) {
      setSubtitle(flow === NETWORK_FLOW ? 'Loading up your network...' : '');
      setDescription('');
      setDuration((prev) => prev + 2);
      setShowReturn(false);
      getNode({
        variables: {
          data: {
            id: nodeId,
          },
        },
      });
      getNodesByState({
        variables: {
          data: {
            state: NodeStateEnum.Unknown,
            connectivity: NodeConnectivityEnum.Online,
          },
        },
      });
      getNodeState({
        variables: {
          getNodeStateId: nodeId,
        },
      });
    }
  };

  return (
    <Stack spacing={{ xs: 4, md: 6 }}>
      <InstallSiteLoading
        duration={duration}
        title={title}
        subtitle={subtitle}
        handleBack={handleBack}
        description={description}
        onCompleted={onInstallProgressComplete}
      />
      {showReturn && (
        <HorizontalContainerJustify>
          <Button variant="text" sx={{ p: 0 }} onClick={handleRetry}>
            Retry
          </Button>
          <Button
            variant="contained"
            sx={{ width: 'fit-content', alignSelf: 'flex-end' }}
            onClick={() => {
              flow === INSTALLATION_FLOW
                ? router.push('/console/home')
                : router.push(`/configure/sims?flow=${ONBOARDING_FLOW}`);
            }}
          >
            {flow === INSTALLATION_FLOW
              ? 'Return to home'
              : 'Skip site configuration'}
          </Button>
        </HorizontalContainerJustify>
      )}
    </Stack>
  );
};

export default Check;
