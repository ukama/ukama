import { ApolloServer } from "@apollo/server";
import { faker } from "@faker-js/faker";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { openStore } from "../../common/storage";
import { THeaders } from "../../common/types";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../network/context";
import NetworkApi from "../../network/datasource/network_api";
import { AddNetworkResolver } from "../../network/resolvers/addNetwork";
import { GetNetworkResolver } from "../../network/resolvers/getNetwork";
import { GetNetworksResolver } from "../../network/resolvers/getNetworks";
import { SetDefaultNetworkResolver } from "../../network/resolvers/setDefaultNetwork";
import { ADD_NETWORK, GET_NETWORK, GET_NETWORKS, SET_DEFAULT } from "./graphql";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName } = parsedHeaders;
const testNetwork = {
  budget: faker.number.int,
  countries: ["Country"],
  name: faker.person.fullName,
  networks: ["A3"],
};

let networkId = "";
const networkApi = new NetworkApi();

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
  const store = openStore();
  const baseURL = await getBaseURL(SUB_GRAPHS.network.name, orgName, store);
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
