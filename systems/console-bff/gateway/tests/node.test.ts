import { ApolloServer } from "@apollo/server";
import { faker } from "@faker-js/faker";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { NODE_STATE, NODE_TYPE } from "../../common/enums";
import { openStore } from "../../common/storage";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import NetworkApi from "../../network/datasource/network_api";
import { Context } from "../../node/context";
import NodeAPI from "../../node/dataSource/node-api";
import { AddNodeResolver } from "../../node/resolvers/addNode";
import { AddNodeToSiteResolver } from "../../node/resolvers/addNodeToSite";
import { AttachNodeResolver } from "../../node/resolvers/attachNode";
import { DeleteNodeFromOrgResolver } from "../../node/resolvers/deleteNodeFromOrg";
import { DetachNodeResolver } from "../../node/resolvers/detachNode";
import { GetAppsChangeLogResolver } from "../../node/resolvers/getAppsChangeLog";
import { GetNodeResolver } from "../../node/resolvers/getNode";
import { GetNodeAppsResolver } from "../../node/resolvers/getNodeApps";
import { GetNodesResolver } from "../../node/resolvers/getNodes";
import { GetNodesByNetworkResolver } from "../../node/resolvers/getNodesByNetwork";
import { GetNodesLocationResolver } from "../../node/resolvers/getNodesLocation";
import { ReleaseNodeFromSiteResolver } from "../../node/resolvers/releaseNodeFromSite";
import { UpdateNodeResolver } from "../../node/resolvers/updateNode";
import { UpdateNodeStateResolver } from "../../node/resolvers/updateNodeState";
import SiteApi from "../../site/datasource/site_api";
import {
  ADD_NODE,
  ADD_NODE_TO_SITE,
  ATTACH_NODE,
  DELETE_NODE,
  DETACH_NODE,
  GET_APPS_CHANGE,
  GET_NETWORK_NODES,
  GET_NODE,
  GET_NODES,
  GET_NODES_LOCATION,
  GET_NODE_APPS,
  GET_NODE_LOCATION,
  UPDATE_NODE,
  UPDATE_NODE_STATE,
} from "./graphql";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName, orgId } = parsedHeaders;

const nodeApi = new NodeAPI();

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      AddNodeResolver,
      AddNodeToSiteResolver,
      AttachNodeResolver,
      DeleteNodeFromOrgResolver,
      DetachNodeResolver,
      GetAppsChangeLogResolver,
      GetNodeResolver,
      GetNodesResolver,
      GetNodeAppsResolver,
      GetNodesByNetworkResolver,
      GetNodesLocationResolver,
      ReleaseNodeFromSiteResolver,
      UpdateNodeResolver,
      UpdateNodeStateResolver,
    ],
    validate: true,
  });
};

const startServer = async () => {
  const schema = await createSchema();
  const server = new ApolloServer<Context>({
    schema,
  });
  await server.start();
  return server;
};

const createContextValue = async () => {
  const store = openStore();
  const baseURL = await getBaseURL(SUB_GRAPHS.node.name, orgName, store);
  return {
    dataSources: { dataSource: nodeApi },
    baseURL: baseURL.message,
    headers: parsedHeaders,
  };
};

describe("Node API integration tests", () => {
  let server: ApolloServer<Context>;
  let contextValue: any;

  beforeAll(async () => {
    server = await startServer();
    contextValue = await createContextValue();
  });
  afterAll(async () => {
    await server.stop();
  });
  const networkApi = new NetworkApi();
  const siteApi = new SiteApi();
  let siteId: string;
  let networkId: string;
  let nodeId: string;

  it("should add a node", async () => {
    const res = await server.executeOperation(
      {
        query: ADD_NODE,
        variables: {
          data: {
            id: faker.string.toString(),
            name: faker.number.toString(),
            orgId: orgId,
          },
        },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.addNode.id).toBeDefined();
    nodeId = data.addNode.id;
    expect(data.addNode.name).toBeDefined();
    expect(data.addNode.orgId).toBeDefined();
    expect(data.addNode.type).toBeDefined();
    expect(data.addNode.attached).toBeDefined();
    expect(data.addNode.site).toBeDefined();
    expect(data.addNode.status).toBeDefined();
  });

  it("should add node to a site", async () => {
    const store = openStore();
    const baseURL = await getBaseURL(SUB_GRAPHS.network.name, orgName, store);

    const testNetwork = {
      budget: faker.number.float(),
      countries: ["Country"],
      name: faker.person.fullName.toString(),
      networks: ["A3"],
    };

    const network = await networkApi.addNetwork(baseURL.message, testNetwork);
    networkId = network.id;

    const site = await siteApi.addSite(baseURL.message, {
      access_id: "",
      backhaul_id: "",
      install_date: "",
      latitude: "",
      location: "",
      longitude: "",
      name: "",
      network_id: "",
      power_id: "",
      spectrum_id: "",
      switch_id: "",
    });
    siteId = site.id;

    const res = await server.executeOperation(
      {
        query: ADD_NODE_TO_SITE,
        variables: {
          data: { networkId: networkId, nodeId: nodeId, siteId: siteId },
        },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.addNodeToSite.success).toBeTruthy();
  });

  it("should attach a node", async () => {
    const res = await server.executeOperation(
      {
        query: ATTACH_NODE,
        variables: {
          data: { anoder: "anoder", parentNode: nodeId, anodel: "anodel" },
        },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.attachNode.success).toBeTruthy();
  });

  it("should get apps change log", async () => {
    const res = await server.executeOperation(
      {
        query: GET_APPS_CHANGE,
        variables: { data: { type: NODE_TYPE.anode } },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getAppsChangeLog.logs).toBeDefined();
    expect(data.getAppsChangeLog.type).toEqual(NODE_TYPE.anode);
  });

  it("should get node using node id", async () => {
    const res = await server.executeOperation(
      {
        query: GET_NODE,
        variables: { data: { id: nodeId } },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getNode.id).toEqual(nodeId);
  });

  it("should get all free nodes", async () => {
    const res = await server.executeOperation(
      {
        query: GET_NODES,
        variables: { data: { isFree: true } },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getNodes.nodes.length).toBeGreaterThanOrEqual(1);
  });

  it("should get node apps", async () => {
    const res = await server.executeOperation(
      {
        query: GET_NODE_APPS,
        variables: { data: { type: NODE_TYPE.anode } },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getNodeApps.apps.length).toBeGreaterThanOrEqual(1);
    expect(data.getNodeApps.type).toEqual(NODE_TYPE.anode);
  });

  it("should get node location", async () => {
    const res = await server.executeOperation(
      {
        query: GET_NODE_LOCATION,
        variables: { data: { type: NODE_TYPE.anode } },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getNodeLocation.id).toEqual(nodeId);
    expect(data.getNodeLocation.lat).toBeDefined();
    expect(data.getNodeLocation.lng).toBeDefined();
  });

  it("should get nodes by network", async () => {
    const res = await server.executeOperation(
      {
        query: GET_NETWORK_NODES,
        variables: { networkId: networkId },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getNodesByNetwork.nodes.length).toBeGreaterThanOrEqual(1);
  });

  it("should get nodes location", async () => {
    const res = await server.executeOperation(
      {
        query: GET_NODES_LOCATION,
        variables: {
          data: {
            networkId: networkId,
            nodeFilterState: NODE_STATE.Configured,
          },
        },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getNodesLocation.networkId).toBeGreaterThanOrEqual(1);
  });

  it("should update node", async () => {
    const res = await server.executeOperation(
      {
        query: UPDATE_NODE,
        variables: { data: { id: nodeId, name: "updated node" } },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.updateNode.id).toEqual(nodeId);
    expect(data.updateNode.name).toEqual("updated node");
  });

  it("should update node state", async () => {
    const res = await server.executeOperation(
      {
        query: UPDATE_NODE_STATE,
        variables: { data: { id: nodeId, state: NODE_STATE.Faulty } },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.updateNodeState.id).toEqual(nodeId);
    expect(data.updateNodeState.status.state).toEqual(NODE_STATE.Faulty);
  });

  it("should detach a node", async () => {
    const res = await server.executeOperation(
      {
        query: DETACH_NODE,
        variables: { data: { id: nodeId } },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.detachNode.success).toBeTruthy();
  });

  it("should delete a node from the org", async () => {
    const res = await server.executeOperation(
      {
        query: DELETE_NODE,
        variables: { data: { id: nodeId } },
      },
      {
        contextValue: contextValue,
      }
    );

    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.deleteNodeFromOrg.id).toEqual(nodeId);
  });
});
