import { ApolloServer } from "@apollo/server";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { openStore } from "../../common/storage";
import { THeaders } from "../../common/types";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../org/context";
import OrgApi from "../../org/datasource/org_api";
import { GetOrgResolver } from "../../org/resolver/getOrg";
import { GET_ORG, GET_ORGS } from "./graphql";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName, orgId, userId } = parsedHeaders;
const orgApi = new OrgApi();

const createSchema = async () => {
  return await buildSchema({
    resolvers: [GetOrgResolver, GetOrgResolver],
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
  const baseURL = await getBaseURL(SUB_GRAPHS.org.name, orgName, store);

  return {
    dataSources: { dataSource: orgApi },
    baseURL: baseURL.message,
    headers: parsedHeaders,
  };
};

describe("Org API integration test", () => {
  let server: ApolloServer<Context>;
  let contextValue: {
    dataSources: { dataSource: OrgApi };
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

  it("should get an org", async () => {
    const res = await server.executeOperation(
      {
        query: GET_ORG,
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
    expect(data.getOrg.id).toEqual(orgId);
    expect(data.getOrg.name).toEqual(orgName);
  });

  it("should get all orgs", async () => {
    const res = await server.executeOperation(
      {
        query: GET_ORGS,
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getOrgs.user.id).toEqual(userId);
  });
});
