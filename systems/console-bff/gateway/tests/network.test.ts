import { ApolloServer } from "@apollo/server";
import { createClient } from "redis";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { logger } from "../../common/logger";
import { THeaders } from "../../common/types";
import {
  generateNetworkName,
  getBaseURL,
  parseGatewayHeaders,
} from "../../common/utils";
import { Context } from "../../network/context";
import NetworkApi from "../../network/datasource/network_api";
import { AddNetworkResolver } from "../../network/resolvers/addNetwork";
import { GetNetworkResolver } from "../../network/resolvers/getNetwork";
import { GetNetworksResolver } from "../../network/resolvers/getNetworks";
import { SetDefaultNetworkResolver } from "../../network/resolvers/setDefaultNetwork";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName } = parsedHeaders;
const testNetwork = {
  budget: Math.floor(Math.random() * 10),
  countries: ["Country"],
  name: generateNetworkName(),
  networks: ["A3"],
};

let networkId = "";
const networkApi = new NetworkApi();
const redisClient = createClient().on("error", error => {
  logger.error(
    `Error creating redis for ${SUB_GRAPHS.network.name} service, Error: ${error}`
  );
});

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      AddNetworkResolver,
      GetNetworkResolver,
      GetNetworksResolver,
      SetDefaultNetworkResolver,
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
  const baseURL = await getBaseURL(
    SUB_GRAPHS.network.name,
    orgName,
    redisClient.isOpen ? redisClient : null
  );
  return {
    dataSources: {
      dataSource: networkApi,
    },
    baseURL: baseURL.message,
    headers: parsedHeaders,
  };
};

describe("Network API integration test", () => {
  let server: ApolloServer<Context>;
  let contextValue: {
    dataSources: { dataSource: NetworkApi };
    baseURL: string;
    headers: THeaders;
  };

  beforeAll(async () => {
    server = await startServer();
    contextValue = await createContextValue();
  });
  afterAll(async () => {
    await server.stop();
  });

  it("should add network", async () => {
    const ADD_NETWORK = `mutation AddNetwork($data: AddNetworkInputDto!) {
  addNetwork(data: $data) {
    id
    name
    isDefault
    budget
    overdraft
    trafficPolicy
    isDeactivated
    paymentLinks
    createdAt
    countries
    networks
  }
}`;

    const res = await server.executeOperation(
      {
        query: ADD_NETWORK,
        variables: {
          data: testNetwork,
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
    expect(data).toBeDefined();
    expect(data.addNetwork).toBeDefined();
    expect(data.addNetwork).toHaveProperty("id");
    networkId = data.addNetwork.id;
    expect(data.addNetwork.name).toEqual(testNetwork.name);
    expect(data.addNetwork.countries).toEqual(testNetwork.countries);
    expect(data.addNetwork.budget).toEqual(testNetwork.budget);
    expect(data.addNetwork.networks).toEqual(testNetwork.networks);
  });
  it("should get all networks", async () => {
    const GET_NETWORKS = `query GetNetworks {
  getNetworks {
    networks {
      id
      name
      isDefault
      budget
      overdraft
      trafficPolicy
      isDeactivated
      paymentLinks
      createdAt
      countries
      networks
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_NETWORKS,
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getNetworks.networks.length).toBeGreaterThanOrEqual(1);
  });

  it("should get a network using network-id", async () => {
    const GET_NETWORK = `query GetNetwork($networkId: String!) {
  getNetwork(networkId: $networkId) {
    id
    name
    isDefault
    budget
    overdraft
    trafficPolicy
    isDeactivated
    paymentLinks
    createdAt
    countries
    networks
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_NETWORK,
        variables: { networkId },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data).toHaveProperty("getNetwork");
    expect(data.getNetwork.id).toEqual(networkId);
    expect(data.getNetwork.name).toEqual(testNetwork.name);
  });
  it("should set default network", async () => {
    const SET_DEFAULT = `mutation SetDefaultNetwork($data: SetDefaultNetworkInputDto!) {
  setDefaultNetwork(data: $data) {
    success
  }
}`;
    const res = await server.executeOperation(
      {
        query: SET_DEFAULT,
        variables: { data: { id: networkId } },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.setDefaultNetwork).toHaveProperty("success");
    expect(data.setDefaultNetwork.success).toBeTruthy();
  });
});
