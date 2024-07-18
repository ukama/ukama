import { ApolloServer } from "@apollo/server";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { THeaders } from "../../common/types";
import { parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../init/context";
import InitAPI from "../../init/datasource/init_api";
import { GetCountriesResolver } from "../../init/resolver/getCountries";
import { GetTimezonesResolver } from "../../init/resolver/getTimezones";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};

const parsedHeaders = parseGatewayHeaders(headers);
const initApi = new InitAPI();

const createSchema = async () => {
  return await buildSchema({
    resolvers: [GetCountriesResolver, GetTimezonesResolver],
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

describe("Init API integration tests", () => {
  let server: ApolloServer<Context>;
  let contextValue: { dataSources: { dataSource: InitAPI }; headers: THeaders };

  beforeAll(async () => {
    server = await startServer();
    contextValue = {
      dataSources: { dataSource: initApi },
      headers: parsedHeaders,
    };
  });
  afterAll(async () => {
    await server.stop();
  });

  it("should get all countries", async () => {
    const GET_COUNTRIES = `query GetCountries {
  getCountries {
    countries {
      name
      code
    }
  }
}`;

    const res = await server.executeOperation(
      {
        query: GET_COUNTRIES,
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.data).toBeDefined();
    const { data } = singleResult;
    expect(data.getCountries).toBeDefined();
    expect(data.getCountries.countries.length).toBeGreaterThanOrEqual(1);
  });
  it("should get all time zones", async () => {
    const GET_TIMEZONES = `query GetTimezones {
  getTimezones {
    timezones {
      value
      abbr
      offset
      isdst
      text
      utc
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_TIMEZONES,
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.data).toBeDefined();
    const { data } = singleResult;
    expect(data.getTimezones).toBeDefined();
    expect(data.getTimezones.timezones.length).toBeGreaterThanOrEqual(1);
  });
});
