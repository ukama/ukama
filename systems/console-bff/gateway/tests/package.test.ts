import { ApolloServer } from "@apollo/server";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { openStore } from "../../common/storage";
import { THeaders } from "../../common/types";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../package/context";
import PackageApi from "../../package/datasource/package_api";
import { AddPackageResolver } from "../../package/resolver/addPackage";
import { DeletePackageResolver } from "../../package/resolver/deletePackage";
import { GetPackageResolver } from "../../package/resolver/getPackage";
import { GetPackagesResolver } from "../../package/resolver/getPackages";
import {
  ADD_PACKAGE,
  DELETE_PACKAGE,
  GET_PACKAGE,
  GET_PACKAGES,
} from "./graphql";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName } = parsedHeaders;

const packageApi = new PackageApi();

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

const createContextValue = async () => {
  const store = openStore();
  const baseURL = await getBaseURL(SUB_GRAPHS.package.name, orgName, store);
  return {
    dataSources: { dataSource: packageApi },
    baseURL: baseURL.message,
    headers: parsedHeaders,
  };
};

let packageId = "";

describe("Package API integration test", () => {
  let server: ApolloServer<Context>;
  let contextValue: {
    dataSources: { dataSource: PackageApi };
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

  const testPackage = {
    amount: 40.12,
    dataUnit: "asfaf",
    dataVolume: 10,
    duration: 30,
    name: "Test-Package",
  };

  it("should add a package", async () => {
    const res = await server.executeOperation(
      {
        query: ADD_PACKAGE,
        variables: { data: testPackage },
      },
      {
        contextValue: contextValue,
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
    const res = await server.executeOperation(
      {
        query: GET_PACKAGE,
        variables: { packageId: packageId },
      },
      {
        contextValue: contextValue,
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
    const res = await server.executeOperation(
      {
        query: GET_PACKAGES,
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getPackages.packages.length).toBeGreaterThanOrEqual(1);
  });

  it("should delete a package", async () => {
    const res = await server.executeOperation(
      {
        query: DELETE_PACKAGE,
        variables: { packageId },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.deletePackage.uuid).toEqual(packageId);
  });
});
