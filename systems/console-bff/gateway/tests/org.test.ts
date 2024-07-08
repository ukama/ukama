import { ApolloServer } from "@apollo/server";
import { createClient } from "redis";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { logger } from "../../common/logger";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../org/context";
import OrgApi from "../../org/datasource/org_api";
import { GetOrgResolver } from "../../org/resolver/getOrg";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName, orgId, userId } = parsedHeaders;
const orgApi = new OrgApi();
const redisClient = createClient().on("error", error => {
  logger.error(
    `Error creating redis for ${SUB_GRAPHS.org.name} service, Error: ${error}`
  );
});

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

describe("Org API integration test", () => {
  let server: ApolloServer<Context>;

  beforeAll(async () => {
    server = await startServer();
  });
  afterAll(async () => {
    await server.stop();
  });

  it("should get an org", async () => {
    const GET_ORG = `query GetOrg {
  getOrg {
    id
    name
    owner
    certificate
    isDeactivated
    createdAt
  }
}`;

    const res = await server.executeOperation(
      {
        query: GET_ORG,
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.org.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: orgApi,
            },
            baseURL: baseURL.message,
            headers: parsedHeaders,
          };
        })(),
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
    const GET_ORGS = `query GetOrgs {
  getOrgs {
    user
    ownerOf {
      id
      name
      owner
      certificate
      isDeactivated
      createdAt
    }
    memberOf {
      id
      name
      owner
      certificate
      isDeactivated
      createdAt
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_ORGS,
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.org.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: orgApi,
            },
            baseURL: baseURL.message,
            headers: parsedHeaders,
          };
        })(),
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getOrgs.user.id).toEqual(userId);
  });
});
