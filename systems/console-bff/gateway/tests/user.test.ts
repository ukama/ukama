import { ApolloServer } from "@apollo/server";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { AuthType } from "../../common/types";
import { Context } from "../../user/context";
import UserApi from "../../user/datasource/user_api";
import { GetUserResolver } from "../../user/resolver/getUser";
import { WhoamiResolver } from "../../user/resolver/whoami";

const userId = process.env.USER_ID;
const orgId = process.env.ORG_ID;
const orgName = process.env.ORG_NAME;

if (!userId || !orgId || !orgName) {
  throw new Error(
    "Environment variables USER_ID, ORG_ID, and ORG_NAME must be set"
  );
}

const userApi = new UserApi();

const createSchema = async () => {
  return await buildSchema({
    resolvers: [GetUserResolver, WhoamiResolver],
    validate: true,
  });
};

async function startServer() {
  const schema = await createSchema();

  const server = new ApolloServer<Context>({
    schema,
  });
  await server.start();
  return server;
}

describe("USER API integration test", () => {
  let server: ApolloServer<Context>;
  beforeAll(async () => {
    server = await startServer();
  });
  afterAll(async () => {
    await server.stop();
  });

  it("should test GetUser Resolver", async () => {
    const GET_USER = `query GetUser($userId: String!) {
  getUser(userId: $userId) {
    name
    email
    uuid
    phone
    isDeactivated
    authId
    registeredSince
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_USER,
        variables: { userId },
      },
      {
        contextValue: {
          dataSources: {
            dataSource: userApi,
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

    const { singleResult } = JSON.parse(body);
    expect(singleResult.data).toHaveProperty("getUser");
    expect(singleResult.data.getUser).toHaveProperty("email");
    expect(singleResult.data.getUser).toHaveProperty("uuid");
    expect(singleResult.data.getUser).toHaveProperty("phone");
    expect(singleResult.data.getUser).toHaveProperty("isDeactivated");
    expect(singleResult.data.getUser).toHaveProperty("authId");
    expect(singleResult.data.getUser).toHaveProperty("registeredSince");
    expect(singleResult.errors).toBeUndefined();
  });
  it("should test WhoAmI Resolver", async () => {
    const WHO_AM_I = `query Whoami {
  whoami {
    user {
      name
      email
      uuid
      phone
      isDeactivated
      authId
      registeredSince
    }
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
        query: WHO_AM_I,
      },
      {
        contextValue: {
          dataSources: {
            dataSource: userApi,
          },
          headers: {
            auth: new AuthType(),
            token: userId,
            orgId: orgId,
            orgName: orgName,
            userId: userId,
          },
        },
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.data).toBeDefined();
    expect(singleResult.errors).toBeUndefined();
    expect(singleResult.data).toHaveProperty("whoami");
    expect(singleResult.data.whoami).toHaveProperty("ownerOf");
    expect(singleResult.data.whoami).toHaveProperty("memberOf");
  });
});
