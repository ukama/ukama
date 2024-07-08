import { ApolloServer } from "@apollo/server";
import { createClient } from "redis";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { logger } from "../../common/logger";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../package/context";
import PackageApi from "../../package/datasource/package_api";
import { AddPackageResolver } from "../../package/resolver/addPackage";
import { DeletePackageResolver } from "../../package/resolver/deletePackage";
import { GetPackageResolver } from "../../package/resolver/getPackage";
import { GetPackagesResolver } from "../../package/resolver/getPackages";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName } = parsedHeaders;

const packageApi = new PackageApi();
const redisClient = createClient().on("error", error => {
  logger.error(
    `Error creating redis for ${SUB_GRAPHS.package.name} service, Error: ${error}`
  );
});

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      AddPackageResolver,
      DeletePackageResolver,
      GetPackageResolver,
      GetPackagesResolver,
    ],
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

let packageId = "";

describe("Package API integration test", () => {
  let server: ApolloServer<Context>;

  beforeAll(async () => {
    server = await startServer();
  });
  afterAll(async () => {
    await server.stop();
  });

  const testPackage = {
    amount: 40.12,
    dataUnit: "asfaf",
    dataVolume: 10,
    duration: 30,
    name: "Test-Package",
  };

  it("should add a package", async () => {
    const ADD_PACKAGE = `mutation AddPackage($data: AddPackageInputDto!) {
  addPackage(data: $data) {
    uuid
    name
    active
    duration
    simType
    createdAt
    deletedAt
    updatedAt
    smsVolume
    dataVolume
    voiceVolume
    ulbr
    dlbr
    type
    dataUnit
    voiceUnit
    messageUnit
    flatrate
    currency
    from
    to
    country
    provider
    apn
    ownerId
    amount
    rate {
      sms_mo
      sms_mt
      data
      amount
    }
    markup {
      baserate
      markup
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: ADD_PACKAGE,
        variables: { data: testPackage },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.package.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: packageApi,
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
    expect(data.addPackage).toBeDefined();
    expect(data.addPackage.uuid).toBeDefined();
    packageId = data.addPackage.uuid;
    expect(data.addPackage.name).toEqual(testPackage.name);
  });

  it("should get a package using packageId", async () => {
    const GET_PACKAGE = `query GetPackage($packageId: String!) {
  getPackage(packageId: $packageId) {
    uuid
    name
    active
    duration
    simType
    createdAt
    deletedAt
    updatedAt
    smsVolume
    dataVolume
    voiceVolume
    ulbr
    dlbr
    type
    dataUnit
    voiceUnit
    messageUnit
    flatrate
    currency
    from
    to
    country
    provider
    apn
    ownerId
    amount
    rate {
      sms_mo
      sms_mt
      data
      amount
    }
    markup {
      baserate
      markup
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_PACKAGE,
        variables: { packageId: packageId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.package.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: packageApi,
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
    expect(data.getPackage).toBeDefined();
    expect(data.getPackage.uuid).toBeDefined();
    packageId = data.getPackage.uuid;
    expect(data.getPackage.name).toEqual(testPackage.name);
  });

  it("should get all packages", async () => {
    const GET_PACKAGES = `query GetPackages {
  getPackages {
    packages {
      uuid
      name
      active
      duration
      simType
      createdAt
      deletedAt
      updatedAt
      smsVolume
      dataVolume
      voiceVolume
      ulbr
      dlbr
      type
      dataUnit
      voiceUnit
      messageUnit
      flatrate
      currency
      from
      to
      country
      provider
      apn
      ownerId
      amount
      rate {
        sms_mo
        sms_mt
        data
        amount
      }
      markup {
        baserate
        markup
      }
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_PACKAGES,
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.package.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: packageApi,
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
    expect(data.getPackages.packages.length).toBeGreaterThanOrEqual(1);
  });

  it("should delete a package", async () => {
    const DELETE_PACKAGE = `mutation DeletePackage($packageId: String!) {
  deletePackage(packageId: $packageId) {
    uuid
  }
}`;
    const res = await server.executeOperation(
      {
        query: DELETE_PACKAGE,
        variables: { packageId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.package.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: packageApi,
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
    expect(data.deletePackage.uuid).toEqual(packageId);
  });
});
