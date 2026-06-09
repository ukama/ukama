/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { Ctx, Query, Resolver } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { SIM_STATUS, SIM_TYPES } from "../../common/enums";
import { openStore } from "../../common/storage";
import { getBaseURL } from "../../common/utils";
import { ComponentDto } from "../../component/resolvers/types";
import type { AppContext } from "../../server/context";
import {
  Component,
  DataPlan,
  Network,
  Org,
  OrgTreeRes,
  Site,
  Subscribers,
} from "./types";

const DEFAULT_ORG_TREE = {
  orgId: "ORGID",
  orgName: "ORGNAME",
  country: "COUNTRY",
  currency: "CURRENCY",
  ownerName: "OWNERNAME",
  ownerEmail: "OWNEREMAIL",
  ownerId: "OWNERID",
  elementType: "ORG",
  networks: [
    {
      networkId: "NETWORKID",
      networkName: "NETWORKNAME",
      elementType: "NETWORK",
      sites: [
        {
          siteId: "SITEID",
          siteName: "SITENAME",
          elementType: "SITE",
          components: [
            {
              componentId: "COMPONENTID",
              componentName: "COMPONENTNAME",
              elementType: "ACCESS",
            },
            {
              componentId: "COMPONENTID",
              componentName: "COMPONENTNAME",
              elementType: "SWITCH",
            },
            {
              componentId: "COMPONENTID",
              componentName: "COMPONENTNAME",
              elementType: "POWER",
            },
            {
              componentId: "COMPONENTID",
              componentName: "COMPONENTNAME",
              elementType: "BACKHAUL",
            },
          ],
        },
      ],
      subscribers: {
        totalSubscribers: "TOTALSUBSCRIBERS",
        activeSubscribers: "ACTIVESUBSCRIBERS",
        inactiveSubscribers: "INACTIVESUBSCRIBERS",
      },
    },
  ],
  dataplans: [
    {
      planId: "PLANID",
      planName: "PLANNAME",
      elementType: "DATAPLAN",
    },
  ],
  sims: {
    availableSims: "AVAILABLESIMS",
    consumed: "CONSUMEDSIMS",
    totalSims: "TOTALSIMS",
  },
  members: {
    totalMembers: "TOTALMEMBERS",
    activeMembers: "ACTIVEMEMBERS",
    inactiveMembers: "INACTIVEMEMBERS",
  },
};

const toComponent = (
  elementType: string,
  component?: ComponentDto
): Component => ({
  componentId: component?.partNumber ? component.partNumber : undefined,
  componentName: component?.partNumber ? component.description : undefined,
  elementType,
});

@Resolver()
export class GetOrgTreeResolver {
  @Query(() => OrgTreeRes)
  async getOrgTree(@Ctx() ctx: AppContext): Promise<OrgTreeRes> {
    // Uses ctx.dataSources (per-request instances) so RESTDataSource GET
    // deduplication applies — do not instantiate datasources inline.
    const { dataSources, headers } = ctx;
    const store = openStore();
    const [registryUrl, dataPlanUrl, subscriberUrl] = await Promise.all([
      getBaseURL(SUB_GRAPHS.network.name, headers.orgName, store),
      getBaseURL(SUB_GRAPHS.package.name, headers.orgName, store),
      getBaseURL(SUB_GRAPHS.sim.name, headers.orgName, store),
    ]);

    const res: Org = DEFAULT_ORG_TREE;
    const nr: Network = DEFAULT_ORG_TREE.networks[0];

    // Independent root fetches — run in parallel.
    const [orgRes, networksRes, plansRes, simsRes, membersRes] =
      await Promise.all([
        dataSources.org.getOrg(headers.orgName),
        dataSources.network.getNetworks(registryUrl.message),
        dataSources.package.getPackages(dataPlanUrl.message),
        dataSources.sim.getSimsFromPool(subscriberUrl.message, {
          type: SIM_TYPES.ukama_data,
          status: SIM_STATUS.ALL,
        }),
        dataSources.member.getMembers(registryUrl.message),
      ]);

    const userRes = await dataSources.user.getUser(orgRes.owner);
    res.orgName = headers.orgName;
    res.orgId = headers.orgId;
    res.ownerId = orgRes.owner;
    res.ownerEmail = userRes.email;
    res.ownerName = userRes.name;
    res.country = orgRes.country;
    res.currency = orgRes.currency;
    res.elementType = "ORG";

    if (networksRes.networks.length > 0) {
      const networkId = networksRes.networks[0].id;
      nr.networkId = networkId;
      nr.networkName = networksRes.networks[0].name;
      nr.elementType = "NETWORK";

      // Sites + subscribers for the network are independent of each other.
      const [siteRes, subscriberRes] = await Promise.all([
        dataSources.site.getSites(registryUrl.message, { networkId }),
        dataSources.subscriber.getSubscribersByNetwork(
          subscriberUrl.message,
          networkId
        ),
      ]);

      if (siteRes.sites.length > 0) {
        const compRes = await dataSources.component.getComponentsByUserId(
          headers,
          "ALL"
        );
        // Index components once (O(n)) instead of find() 4x per site.
        const componentsById = new Map<string, ComponentDto>(
          compRes.components.map(comp => [comp.id, comp])
        );
        nr.sites = siteRes.sites.map(
          (site): Site => ({
            siteId: site.id,
            siteName: site.name,
            elementType: "SITE",
            components: [
              toComponent("ACCESS", componentsById.get(site.accessId)),
              toComponent("POWER", componentsById.get(site.powerId)),
              toComponent("BACKHAUL", componentsById.get(site.backhaulId)),
              toComponent("SWITCH", componentsById.get(site.switchId)),
            ],
          })
        );
      } else {
        nr.sites = [];
      }

      if (subscriberRes.subscribers.length > 0) {
        const ssr: Subscribers = DEFAULT_ORG_TREE.networks[0].subscribers;
        ssr.totalSubscribers = subscriberRes.subscribers.length.toString();
        let activeSubCount = 0;
        let inactiveSubCount = 0;
        subscriberRes.subscribers.forEach(sub => {
          if ((sub.sim?.length ?? 0) > 0) {
            activeSubCount++;
          } else {
            inactiveSubCount++;
          }
        });
        ssr.activeSubscribers = activeSubCount.toString();
        ssr.inactiveSubscribers = inactiveSubCount.toString();
        nr.subscribers = ssr;
      } else {
        nr.subscribers = undefined;
      }
      res.networks = [nr];
    } else {
      res.networks = [];
    }

    if (plansRes.packages.length > 0) {
      res.dataplans = plansRes.packages.map(
        (plan): DataPlan => ({
          planId: plan.uuid,
          planName: plan.name,
          elementType: "DATAPLAN",
        })
      );
    } else {
      res.dataplans = [];
    }

    if (simsRes.sims.length > 0) {
      const simr = DEFAULT_ORG_TREE.sims;
      simr.totalSims = simsRes.sims.length.toString();
      const consumed = simsRes.sims.filter(sim => sim.isAllocated === true);
      simr.consumed = consumed.length.toString();
      simr.availableSims = (simsRes.sims.length - consumed.length).toString();
      res.sims = simr;
    } else {
      res.sims = undefined;
    }

    if (membersRes.members.length > 0) {
      const mr = DEFAULT_ORG_TREE.members;
      mr.totalMembers = membersRes.members.length.toString();
      let activeMembers = 0;
      let inactiveMembers = 0;
      membersRes.members.forEach(member => {
        if (member.isDeactivated === false) {
          activeMembers++;
        } else {
          inactiveMembers++;
        }
      });
      mr.activeMembers = activeMembers.toString();
      mr.inactiveMembers = inactiveMembers.toString();
      res.members = mr;
    } else {
      res.members = undefined;
    }

    return {
      org: res,
    };
  }
}
