import {
  Component,
  DataPlan,
  Network,
  OrgTreeRes,
  Site,
} from '@/client/graphql/generated';

export const convertToNamesJson = (data: OrgTreeRes): any => {
  const org = data.org;

  const convertComponents = (components: Component[]) => {
    return components.map((component) => ({
      name: component.componentName,
      size: 1,
      elementType: component.elementType,
    }));
  };

  const convertSites = (sites: Site[]) => {
    return sites.map((site) => ({
      name: site.siteName,
      elementType: site.elementType,
      children:
        site.components.length > 0 ? convertComponents(site.components) : [],
    }));
  };

  const convertNetworks = (networks: Network[]) => {
    return networks.map((network) => ({
      name: network.networkName,
      elementType: network.elementType,
      children:
        network.sites.length > 0
          ? [
              ...convertSites(network.sites),
              {
                name: 'Subscribers',
                elementType: 'SUBSCRIBERS',
                children: [
                  {
                    totalSubscribers:
                      network.subscribers?.totalSubscribers ?? `?`,
                    activeSubscribers:
                      network.subscribers?.activeSubscribers ?? `?`,
                    inactiveSubscribers:
                      network.subscribers?.inactiveSubscribers ?? `?`,
                    elementType: 'SUBSCRIBER_STATS',
                  },
                ],
              },
            ]
          : [],
    }));
  };

  const convertDataPlans = (dataplans: DataPlan[]) => {
    return dataplans.map((plan) => ({
      name: plan.planName,
      size: 1,
      elementType: plan.elementType,
    }));
  };

  const namesJson = {
    name: org.orgName,
    elementType: org.elementType,
    children: [
      {
        name: 'Members',
        elementType: 'MEMBERS',
        children: [
          {
            totalMembers: org?.members?.totalMembers ?? `?`,
            activeMembers: org?.members?.activeMembers ?? `?`,
            inactiveMembers: org?.members?.inactiveMembers ?? `?`,
            elementType: 'MEMBER_STATS',
          },
        ],
      },
      {
        name: 'sims',
        elementType: 'SIMS',
        children: [
          {
            totalSims: org.sims?.totalSims ?? `?`,
            consumed: org.sims?.consumed ?? `?`,
            availableSims: org.sims?.availableSims ?? `?`,
            elementType: 'SIM_STATS',
          },
        ],
      },
      {
        name: 'Data Plans',
        elementType: 'DATAPLAN',
        children:
          org.dataplans.length > 0 ? convertDataPlans(org.dataplans) : [],
      },
      {
        name: 'Networks',
        elementType: 'NETWORKS',
        children: org.networks.length > 0 ? convertNetworks(org.networks) : [],
      },
    ],
  };

  return namesJson;
};
