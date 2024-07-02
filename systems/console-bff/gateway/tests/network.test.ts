import { ApolloServer } from "@apollo/server";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { AuthType } from "../../common/types";
import { Context } from "../../network/context";
import NetworkApi from "../../network/datasource/network_api";
import { AddNetworkResolver } from "../../network/resolvers/addNetwork";
import { AddSiteResolver } from "../../network/resolvers/addSite";
import { GetNetworkResolver } from "../../network/resolvers/getNetwork";
import { GetNetworksResolver } from "../../network/resolvers/getNetworks";
import { GetSiteResolver } from "../../network/resolvers/getSite";
import { GetSitesResolver } from "../../network/resolvers/getSites";

const userId = process.env.USER_ID;
const orgId = process.env.ORG_ID;
const orgName = process.env.ORG_NAME;

if (!userId || !orgId || !orgName) {
  throw new Error(
    "Environment variables USER_ID, ORG_ID, and ORG_NAME must be set"
  );
}

const networkApi = new NetworkApi();

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      AddNetworkResolver,
      AddSiteResolver,
      GetNetworkResolver,
      GetNetworksResolver,
      GetSiteResolver,
      GetSitesResolver,
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

describe("Network API integration test", () => {
  let server: ApolloServer<Context>;

  beforeAll(async () => {
    server = await startServer();
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
    const testNetwork = {
      budget: Math.floor(Math.random() * 10),
      countries: ["Country"],
      name: "notmynetwork1",
      networks: ["A3"],
      org: "ukama-test-org",
    };

    const res = await server.executeOperation(
      {
        query: ADD_NETWORK,
        variables: {
          data: testNetwork,
        },
      },
      {
        contextValue: {
          dataSources: {
            dataSource: networkApi,
          },
          headers: {
            auth: new AuthType(),
            token: "",
            orgId: orgId,
            orgName: orgName,
            userId: userId,
          },
        },
      }
    );
    const body = JSON.stringify(res.body);
    const { addNetwork } = JSON.parse(body);
    expect(addNetwork).toBeDefined();
    expect(addNetwork).toHaveProperty("id");
    expect(addNetwork.name).toEqual(testNetwork.name);
    expect(addNetwork.countries).toEqual(testNetwork.countries);
    expect(addNetwork.budget).toEqual(testNetwork.budget);
    expect(addNetwork.networks).toEqual(testNetwork.networks);
  });
});
