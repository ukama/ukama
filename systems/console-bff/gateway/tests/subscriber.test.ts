import { ApolloServer } from "@apollo/server";
import { createClient } from "redis";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { logger } from "../../common/logger";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import NetworkApi from "../../network/datasource/network_api";
import { Context } from "../../subscriber/context";
import SubscriberApi from "../../subscriber/datasource/subscriber_api";
import { AddSubscriberResolver } from "../../subscriber/resolver/addSubscriber";
import { DeleteSubscriberResolver } from "../../subscriber/resolver/deleteSubscriber";
import { GetSubscriberResolver } from "../../subscriber/resolver/getSubscriber";
import { GetSubscriberMetricsByNetworkResolver } from "../../subscriber/resolver/getSubscriberMetricsByNetwork";
import { GetSubscribersByNetworkResolver } from "../../subscriber/resolver/getSubscribersByNetwork";
import UserApi from "../../user/datasource/user_api";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};
const parsedHeaders = parseGatewayHeaders(headers);

const { orgName, userId } = parsedHeaders;

const subscriberApi = new SubscriberApi();
const redisClient = createClient().on("error", error => {
  logger.error(
    `Error creating redis for ${SUB_GRAPHS.subscriber.name} service, Error: ${error}`
  );
});

const generateRandomName = (length = 10) => {
  const characters = "abcdefghijklmnopqrstuvwxyz-";
  return Array.from(
    { length },
    () => characters[Math.floor(Math.random() * characters.length)]
  ).join("");
};

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      AddSubscriberResolver,
      DeleteSubscriberResolver,
      GetSubscriberResolver,
      GetSubscriberMetricsByNetworkResolver,
      GetSubscribersByNetworkResolver,
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

describe("subscriber API integration tests", () => {
  let server: ApolloServer<Context>;

  beforeAll(async () => {
    server = await startServer();
  });
  afterAll(async () => {
    await server.stop();
  });
  const userApi = new UserApi();
  const networkApi = new NetworkApi();

  let subscriberId = "";
  let networkId = "";

  it("should add a subscriber", async () => {
    const ADD_SUBSCRIBER = `mutation AddSubscriber($data: SubscriberInputDto!) {
  addSubscriber(data: $data) {
    uuid
    address
    dob
    email
    firstName
    lastName
    gender
    idSerial
    networkId
    phone
    proofOfIdentification
    sim {
      id
      subscriberId
      networkId
      iccid
      msisdn
      imsi
      type
      status
      firstActivatedOn
      lastActivatedOn
      activationsCount
      deactivationsCount
      allocatedAt
      isPhysical
      package
    }
  }
}
`;
    const networkURL = await getBaseURL(
      SUB_GRAPHS.network.name,
      orgName,
      redisClient.isOpen ? redisClient : null
    );
    const network = await networkApi.addNetwork(networkURL.message, {
      budget: Math.floor(Math.random() * 10),
      countries: ["Country"],
      name: generateRandomName(),
      networks: ["A3"],
      org: orgName,
    });

    networkId = network.id;
    const user = await userApi.whoami(userId);

    const { email, phone } = user.user;

    const res = await server.executeOperation(
      {
        query: ADD_SUBSCRIBER,
        variables: {
          data: {
            email: email,
            first_name: "First",
            last_name: "Second",
            network_id: networkId,
            phone: phone,
          },
        },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.subscriber.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: subscriberApi,
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
    expect(data.addSubscriber).toHaveProperty("uuid");
    subscriberId = data.addSubscriber.uuid;
    expect(data.addSubscriber.networkId).toEqual(networkId);
  });

  it("should get a subscriber", async () => {
    const GET_SUBSCRIBER = `query GetSubscriber($subscriberId: String!) {
    getSubscriber(subscriberId: $subscriberId) {
      uuid
      address
      dob
      email
      firstName
      lastName
      gender
      idSerial
      networkId
      phone
      proofOfIdentification
      sim {
        id
        subscriberId
        networkId
        iccid
        msisdn
        imsi
        type
        status
        firstActivatedOn
        lastActivatedOn
        activationsCount
        deactivationsCount
        allocatedAt
        isPhysical
        package
      }
    }
  }`;

    const res = await server.executeOperation(
      {
        query: GET_SUBSCRIBER,
        variables: { subscriberId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.subscriber.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: subscriberApi,
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
    expect(data.getSubscriber.uuid).toEqual(subscriberId);
  });

  it("should get subscriber metrics by network", async () => {
    const GET_SUBSCRIBER = `query GetSubscriberMetricsByNetwork($networkId: String!) {
  getSubscriberMetricsByNetwork(networkId: $networkId) {
    total
    active
    inactive
    terminated
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_SUBSCRIBER,
        variables: { networkId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.subscriber.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: subscriberApi,
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
    expect(data.getSubscriberMetricsByNetwork).toHaveProperty("total");
    expect(data.getSubscriberMetricsByNetwork).toHaveProperty("active");
    expect(data.getSubscriberMetricsByNetwork).toHaveProperty("inactive");
    expect(data.getSubscriberMetricsByNetwork).toHaveProperty("terminated");
  });

  it("should get all subscribers by network", async () => {
    const GET_SUBSCRIBERS = `query GetSubscribersByNetwork($networkId: String!) {
  getSubscribersByNetwork(networkId: $networkId) {
    subscribers {
      uuid
      address
      dob
      email
      firstName
      lastName
      gender
      idSerial
      networkId
      phone
      proofOfIdentification
      sim {
        id
        subscriberId
        networkId
        iccid
        msisdn
        imsi
        type
        status
        firstActivatedOn
        lastActivatedOn
        activationsCount
        deactivationsCount
        allocatedAt
        isPhysical
        package
      }
    }
  }
}`;

    const res = await server.executeOperation(
      {
        query: GET_SUBSCRIBERS,
        variables: { networkId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.subscriber.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: subscriberApi,
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
    expect(
      data.getSubscribersByNetwork.subscribers.length
    ).toBeGreaterThanOrEqual(1);
  });

  it("should be able to update a subscriber", async () => {
    const UPDATE_SUBSCRIBER = `mutation UpdateSubscriber($data: UpdateSubscriberInputDto!, $subscriberId: String!) {
  updateSubscriber(data: $data, subscriberId: $subscriberId) {
    success
  }
}`;
    const res = await server.executeOperation(
      {
        query: UPDATE_SUBSCRIBER,
        variables: {
          data: {
            first_name: "New First Name",
          },
          subscriberId: subscriberId,
        },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.subscriber.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: subscriberApi,
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
    expect(data.updateSubscriber.success).toBeTruthy();
  });

  it("should delete a subscriber", async () => {
    const DELETE_SUBSCRIBER = `mutation DeleteSubscriber($subscriberId: String!) {
  deleteSubscriber(subscriberId: $subscriberId) {
    success
  }
}`;
    const res = await server.executeOperation(
      {
        query: DELETE_SUBSCRIBER,
        variables: { subscriberId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.subscriber.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: subscriberApi,
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
    expect(data.deleteSubscriber.success).toBeTruthy();
  });
});
