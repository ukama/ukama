import { ApolloServer } from "@apollo/server";
import { readFile } from "fs";
import path from "path";
import { createClient } from "redis";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { SIM_TYPES } from "../../common/enums";
import { logger } from "../../common/logger";
import {
  generateNetworkName,
  getBaseURL,
  parseGatewayHeaders,
} from "../../common/utils";
import NetworkApi from "../../network/datasource/network_api";
import PackageApi from "../../package/datasource/package_api";
import { Context } from "../../sim/context";
import SimApi from "../../sim/datasource/sim_api";
import { AddPackageToSimResolver } from "../../sim/resolver/addPackagetoSim";
import { AllocateSimResolver } from "../../sim/resolver/allocateSim";
import { DeleteSimResolver } from "../../sim/resolver/delete";
import { GetSimByNetworkResolver } from "../../sim/resolver/getByNetwork";
import { GetSimBySubscriberResolver } from "../../sim/resolver/getBySubscriber";
import { GetDataUsageResolver } from "../../sim/resolver/getDataUsage";
import { GetPackagesForSimResolver } from "../../sim/resolver/getPackagesForSim";
import { GetSimResolver } from "../../sim/resolver/getSim";
import { GetSimPoolStatsResolver } from "../../sim/resolver/getSimPoolStats";
import { GetSimsResolver } from "../../sim/resolver/getSims";
import { GetSimsBySubscriberResolver } from "../../sim/resolver/getSimsBySubscriber";
import SubscriberApi from "../../subscriber/datasource/subscriber_api";
import UserApi from "../../user/datasource/user_api";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};

const parsedHeaders = parseGatewayHeaders(headers);
const { orgName, userId } = parsedHeaders;

const simApi = new SimApi();

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      AddPackageToSimResolver,
      AllocateSimResolver,
      DeleteSimResolver,
      GetSimByNetworkResolver,
      GetSimBySubscriberResolver,
      GetDataUsageResolver,
      GetPackagesForSimResolver,
      GetSimResolver,
      GetSimPoolStatsResolver,
      GetSimsResolver,
      GetSimsBySubscriberResolver,
    ],
    validate: true,
  });
};

const csvToBase64 = (filePath: string) => {
  readFile(filePath, (err, data) => {
    if (err) {
      console.error("Error reading file: ", err);
      return;
    }
    return data.toString("base64");
  });
};

const redisClient = createClient().on("error", error => {
  logger.error(
    `Error creating redis for ${SUB_GRAPHS.sim.name} service, Error: ${error}`
  );
});

const startServer = async () => {
  const schema = await createSchema();
  const server = new ApolloServer<Context>({
    schema,
  });
  await server.start();
  return server;
};

describe("Sim API integration tests", () => {
  let server: ApolloServer<Context>;

  beforeAll(async () => {
    server = await startServer();
  });
  afterAll(async () => {
    await server.stop();
  });

  const packageApi = new PackageApi();
  const networkApi = new NetworkApi();
  const subscriberApi = new SubscriberApi();
  const userApi = new UserApi();

  let packageId: string;
  let networkId: string;
  let subscriberId: string;
  let iccid: string;
  let simId: string;

  it("should upload sims", async () => {
    const UPLOAD_SIMS = `mutation UploadSims($data: UploadSimsInputDto!) {
  uploadSims(data: $data) {
    iccid
  }
}`;
    const simData = csvToBase64(path.join(__dirname, "SimPool.csv"));

    const res = await server.executeOperation(
      {
        query: UPLOAD_SIMS,
        variables: { data: { data: simData, simType: SIM_TYPES.TEST } },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.uploadSims.iccid).toBeDefined();
    iccid = data.uploadSims.iccid;
  });

  it("should add a package to sim", async () => {
    const ADD_PACKAGE = `mutation AddPackageToSim($data: AddPackageToSimInputDto!) {
  addPackageToSim(data: $data) {
    packageId
  }
}`;
    const packageURL = await getBaseURL(
      SUB_GRAPHS.package.name,
      orgName,
      redisClient.isOpen ? redisClient : null
    );
    const testPackage = await packageApi.addPackage(
      packageURL.message,
      {
        amount: 40.12,
        dataUnit: "asfaf",
        dataVolume: 10,
        duration: 30,
        name: "Test-Package",
      },
      parsedHeaders
    );
    packageId = testPackage.uuid;

    const res = await server.executeOperation(
      {
        query: ADD_PACKAGE,
        variables: {
          data: {
            package_id: packageId,
            sim_id: "0000000-0000000-0000000-00000",
            start_date: "12-02-2021",
          },
        },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.addPackageToSim.packageId).toEqual(packageId);
  });

  it("should allocate sim", async () => {
    const ALLOCATE_SIM = `mutation AllocateSim($data: AllocateSimInputDto!) {
  allocateSim(data: $data) {
    id
    subscriber_id
    network_id
    package {
      id
      packageId
      startDate
      endDate
      isActive
    }
    iccid
    msisdn
    imsi
    type
    status
    is_physical
    traffic_policy
    firstActivatedOn
    lastActivatedOn
    activationsCount
    deactivationsCount
    allocated_at
    sync_status
  }
}`;
    const networkURL = await getBaseURL(
      SUB_GRAPHS.network.name,
      orgName,
      redisClient.isOpen ? redisClient : null
    );
    const testNetwork = await networkApi.addNetwork(networkURL.message, {
      budget: Math.floor(Math.random() * 10),
      countries: ["Country"],
      name: generateNetworkName(),
      networks: ["A3"],
      org: orgName,
    });
    networkId = testNetwork.id;

    const subscriberURL = await getBaseURL(
      SUB_GRAPHS.subscriber.name,
      orgName,
      redisClient.isOpen ? redisClient : null
    );
    const user = await userApi.whoami(userId);
    const { email, phone } = user.user;
    const testSubscriber = await subscriberApi.addSubscriber(
      subscriberURL.message,
      {
        email: email,
        network_id: networkId,
        first_name: "First Name",
        last_name: "Last Name",
        phone: phone,
      }
    );
    subscriberId = testSubscriber.uuid;

    const res = await server.executeOperation(
      {
        query: ALLOCATE_SIM,
        variables: {
          data: {
            iccid: iccid,
            network_id: networkId,
            sim_type: SIM_TYPES.TEST,
            package_id: packageId,
            subscriber_id: subscriberId,
            traffic_policy: 123,
          },
        },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.allocateSim.id).toBeDefined();
    simId = data.allocateSim.id;
  });

  it("should get sim using simId", async () => {
    const GET_SIM = `query GetSim($data: GetSimInputDto!) {
  getSim(data: $data) {
    activationCode
    createdAt
    iccid
    id
    isAllocated
    isPhysical
    msisdn
    qrCode
    simType
    smapAddress
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_SIM,
        variables: { data: { simId: simId } },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.getSim.id).toEqual(simId);
    expect(data.getSim.iccid).toEqual(iccid);
  });

  it("should get sim pool stats", async () => {
    const GET_SIM_STATS = `query GetSimPoolStats($type: String!) {
  getSimPoolStats(type: $type) {
    total
    available
    consumed
    failed
    esim
    physical
  }
}`;

    const res = await server.executeOperation(
      {
        query: GET_SIM_STATS,
        variables: { type: SIM_TYPES.TEST },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.getSimPoolStats).toBeDefined();
  });

  it("should get the data usage", async () => {
    const GET_DATA_USAGE = `query GetDataUsage($simId: String!) {
  getDataUsage(simId: $simId) {
    usage
  }
}`;

    const res = await server.executeOperation(
      {
        query: GET_DATA_USAGE,
        variables: { data: { simId } },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.getDataUsage.usage).toBeDefined();
  });

  it("should get packages for sim", async () => {
    const GET_PACKAGES = `query GetPackagesForSim($data: GetPackagesForSimInputDto!) {
  getPackagesForSim(data: $data) {
    sim_id
    packages {
      id
      package_id
      start_date
      end_date
      is_active
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_PACKAGES,
        variables: { data: { simId } },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.getPackagesForSim.sim_id).toEqual(simId);
    expect(data.getPackagesForSim.packages).toBeDefined();
  });

  it("should get sims by subscriber", async () => {
    const GET_SIM = `query GetSimsBySubscriber($data: GetSimBySubscriberInputDto!) {
  getSimsBySubscriber(data: $data) {
    subscriber_id
    sims {
      activationCode
      createdAt
      iccid
      id
      isAllocated
      isPhysical
      msisdn
      qrCode
      simType
      smapAddress
    }
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_SIM,
        variables: { data: { subscriberId } },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.getSimsBySubscriber.subscriber_id).toEqual(subscriberId);
    expect(data.getSimsBySubscriber.sims).toBeDefined();
  });

  it("should delete a sim using simId", async () => {
    const DELETE_SIM = `mutation DeleteSim($data: DeleteSimInputDto!) {
  deleteSim(data: $data) {
    simId
  }
}`;

    const res = await server.executeOperation(
      {
        query: DELETE_SIM,
        variables: { data: { simId: simId } },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.sim.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: simApi,
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
    expect(data.deleteSim.simId).toEqual(simId);
  });
});
