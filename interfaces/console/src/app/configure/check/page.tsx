/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';

import {
  useGetNodeQuery,
  useGetSitesLazyQuery,
} from '@/client/graphql/generated';
import InstallSiteLoading from '@/components/InstallSiteLoading';
import {
  CHECK_SITE_FLOW,
  INSTALLATION_FLOW,
  NETWORK_FLOW,
  ONBOARDING_FLOW,
} from '@/constants';
import { useAppContext } from '@/context';
import { usePathname, useRouter, useSearchParams } from 'next/navigation';
import { useEffect, useState } from 'react';

const Check = () => {
  const router = useRouter();
  const pathname = usePathname();
  const searchParams = useSearchParams();
  const flow = searchParams.get('flow') ?? INSTALLATION_FLOW;
  const nodeId = searchParams.get('nid') ?? '';
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
  const { network } = useAppContext();

  const setQueryParam = (key: string, value: string) => {
    const p = new URLSearchParams(searchParams.toString());
    p.set(key, value);
    window.history.replaceState({}, '', `${pathname}?${p.toString()}`);
    return p;
  };

  const [getSites] = useGetSitesLazyQuery({
    onCompleted: (data) => {
      if (data.getSites.sites.length > 0 && !nodeId) {
        // TODO: CHECK IF ANY SITE IS AVAILABLE FOR CONFIGURE & REDIRECT TO SITE CONFIGURE STEP 3
        router.push(`/console/home`);
        // router.push(
        //     `/configure/node/uk-sa9001-tnode-a1-1234?step=1&flow=${INSTALLATION_FLOW}`,
        //   );
      }
    },
  });

  const { data: nodeData, loading: nodeLoading } = useGetNodeQuery({
    skip: !nodeId,
    variables: {
      data: {
        id: nodeId,
      },
    },
    onCompleted: async (data) => {
      if (data.getNode.latitude && data.getNode.longitude && nodeId) {
        setTimeout(() => {}, 2000);
        let p = setQueryParam('lat', data.getNode.latitude.toString());
        p.set('lng', data.getNode.longitude.toString());
        p.delete('nid');
        router.push(`/configure/node/${data.getNode.id}?${p.toString()}`);
      }
    },
  });

  useEffect(() => {
    getSites({ variables: { networkId: network.id } });
  }, []);

  const onInstallProgressComplete = () => {
    if (flow !== NETWORK_FLOW) {
      setSubtitle('❗ Site not detected');
      setDescription(
        'It is taking longer than usual to load up your site. Please check on your site to make sure that all parts are installed correctly.',
      );
    } else {
      // TODO: CHECK IF ANY SITE IS AVAILABLE FOR CONFIGURE & REDIRECT TO SITE CONFIGURE STEP 3
      router.push(`/configure?step=2&flow=${ONBOARDING_FLOW}`);
    }
  };

  const handleBack = () => {
    router.push(
      `/configure/sims?step=${flow === ONBOARDING_FLOW ? 5 : 4}&flow=${flow === NETWORK_FLOW ? ONBOARDING_FLOW : flow === CHECK_SITE_FLOW ? INSTALLATION_FLOW : flow}`,
    );
  };

  return (
    <InstallSiteLoading
      duration={10}
      title={title}
      subtitle={subtitle}
      handleBack={handleBack}
      description={description}
      onCompleted={onInstallProgressComplete}
    />
  );
};

export default Check;
