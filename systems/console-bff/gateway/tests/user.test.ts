import { ApolloServer } from "@apollo/server";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { THeaders } from "../../common/types";
import { parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../user/context";
import UserApi from "../../user/datasource/user_api";
import { GetUserResolver } from "../../user/resolver/getUser";
import { WhoamiResolver } from "../../user/resolver/whoami";
import { GET_USER, WHO_AM_I } from "./graphql";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};

const parsedHeaders = parseGatewayHeaders(headers);
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
  let contextValue: { dataSources: { dataSource: UserApi }; headers: THeaders };
  beforeAll(async () => {
    server = await startServer();
    contextValue = {
      dataSources: {
        dataSource: userApi,
      },
      headers: parsedHeaders,
    };
  });
  afterAll(async () => {
    await server.stop();
  });

  it("should test GetUser Resolver", async () => {
    const { userId } = parsedHeaders;
    const res = await server.executeOperation(
      {
        query: GET_USER,
        variables: { userId },
      },
      {
        contextValue: contextValue,
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
    const res = await server.executeOperation(
      {
        query: WHO_AM_I,
      },
      {
        contextValue: contextValue,
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
