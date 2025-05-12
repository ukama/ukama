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
import ComponentApi from "../../component/datasource/component_api";
import MemberApi from "../../member/datasource/member_api";
import NetworkApi from "../../network/datasource/network_api";
import PackageApi from "../../package/datasource/package_api";
import SimApi from "../../sim/datasource/sim_api";
import SiteApi from "../../site/datasource/site_api";
import SubscriberApi from "../../subscriber/datasource/subscriber_api";
import UserApi from "../../user/datasource/user_api";
import { Context } from "../context";
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

@Resolver()
export class GetOrgTreeResolver {
  @Query(() => OrgTreeRes)
  async getOrgTree(@Ctx() ctx: Context): Promise<OrgTreeRes> {
    const { dataSources, headers } = ctx;
    const store = openStore();
    const registryUrl = await getBaseURL(
      SUB_GRAPHS.network.name,
      headers.orgName,
      store
    );
    const dataPlanUrl = await getBaseURL(
      SUB_GRAPHS.package.name,
      headers.orgName,
      store
    );
    const subscriberUrl = await getBaseURL(
      SUB_GRAPHS.sim.name,
      headers.orgName,
      store
    );
    const userAPI = new UserApi();
    const siteAPI = new SiteApi();
    const memberAPI = new MemberApi();
    const simsAPI = new SimApi();
    const networkAPI = new NetworkApi();
    const packageAPI = new PackageApi();
    const subscriberAPI = new SubscriberApi();
    const componentAPI = new ComponentApi();

    const res: Org = DEFAULT_ORG_TREE;
    const nr: Network = DEFAULT_ORG_TREE.networks[0];

    const orgRes = await dataSources.dataSource.getOrg(headers.orgName);
    const userRes = await userAPI.getUser(orgRes.owner);
    res.orgName = headers.orgName;
    res.orgId = headers.orgId;
    res.ownerId = orgRes.owner;
    res.ownerEmail = userRes.email;
    res.ownerName = userRes.name;
    res.country = orgRes.country;
    res.currency = orgRes.currency;
    res.elementType = "ORG";

    const networksRes = await networkAPI.getNetworks(registryUrl.message);
    if (networksRes.networks.length > 0) {
      nr.networkId = networksRes.networks[0].id;
      nr.networkName = networksRes.networks[0].name;
      nr.elementType = "NETWORK";
      res.networks = [nr];

      const siteRes = await siteAPI.getSites(registryUrl.message, {
        networkId: networksRes.networks[0].id,
      });
      const sr: Site[] = [];
      if (siteRes.sites.length > 0) {
        const compRes = await componentAPI.getComponentsByUserId(
          headers,
          "ALL"
        );
        for (const site of siteRes.sites) {
          const cr: Component[] = [];
          const accessRes = compRes.components.find(
            comp => comp.id === site.accessId
          );
          if (accessRes?.partNumber) {
            cr.push({
              componentId: accessRes.partNumber,
              componentName: accessRes.description,
              elementType: "ACCESS",
            });
          } else {
            cr.push({
              componentId: undefined,
              componentName: undefined,
              elementType: "ACCESS",
            });
          }

          const powerRes = compRes.components.find(
            comp => comp.id === site.powerId
          );
          if (powerRes?.partNumber) {
            cr.push({
              componentId: powerRes.partNumber,
              componentName: powerRes.description,
              elementType: "POWER",
            });
          } else {
            cr.push({
              componentId: undefined,
              componentName: undefined,
              elementType: "POWER",
            });
          }

          const backhaulRes = compRes.components.find(
            comp => comp.id === site.backhaulId
          );
          if (backhaulRes?.partNumber) {
            cr.push({
              componentId: backhaulRes.partNumber,
              componentName: backhaulRes.description,
              elementType: "BACKHAUL",
            });
          } else {
            cr.push({
              componentId: undefined,
              componentName: undefined,
              elementType: "BACKHAUL",
            });
          }

          const switchRes = compRes.components.find(
            comp => comp.id === site.switchId
          );
          if (switchRes?.partNumber) {
            cr.push({
              componentId: switchRes.partNumber,
              componentName: switchRes.description,
              elementType: "SWITCH",
            });
          } else {
            cr.push({
              componentId: undefined,
              componentName: undefined,
              elementType: "SWITCH",
            });
          }

          sr.push({
            siteId: site.id,
            siteName: site.name,
            elementType: "SITE",
            components: cr,
          });
        }
        nr.sites = sr;
      } else {
        nr.sites = [];
      }

      const subscriberRes = await subscriberAPI.getSubscribersByNetwork(
        subscriberUrl.message,
        networksRes.networks[0].id
      );

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

    const plansRes = await packageAPI.getPackages(dataPlanUrl.message);
    if (plansRes.packages.length > 0) {
      const dr: DataPlan[] = [];
      plansRes.packages.forEach(plan => {
        dr.push({
          planId: plan.uuid,
          planName: plan.name,
          elementType: "DATAPLAN",
        });
      });
      res.dataplans = dr;
    } else {
      res.dataplans = [];
    }

    const simsRes = await simsAPI.getSimsFromPool(subscriberUrl.message, {
      type: SIM_TYPES.ukama_data,
      status: SIM_STATUS.ALL,
    });
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

    const membersRes = await memberAPI.getMembers(registryUrl.message);
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
