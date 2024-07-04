import { ApolloServer } from "@apollo/server";
import { createClient } from "redis";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { INVITATION_STATUS } from "../../common/enums";
import { logger } from "../../common/logger";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../invitation/context";
import InvitationApi from "../../invitation/datasource/invitation_api";
import { CreateInvitationResolver } from "../../invitation/resolver/createInvitation";
import { DeleteInvitationResolver } from "../../invitation/resolver/deleteInvitation";
import { GetInvitationResolver } from "../../invitation/resolver/getInvitation";
import { GetInVitationsByOrgResolver } from "../../invitation/resolver/getInvitationByOrg";
import { GetInvitationsResolver } from "../../invitation/resolver/getInvitations";
import { UpdateInvitationResolver } from "../../invitation/resolver/updateInvitation";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};

const parsedHeaders = parseGatewayHeaders(headers);
const { orgName } = parsedHeaders;

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      CreateInvitationResolver,
      DeleteInvitationResolver,
      GetInvitationResolver,
      GetInvitationsResolver,
      GetInVitationsByOrgResolver,
      UpdateInvitationResolver,
    ],
    validate: true,
  });
};

const redisClient = createClient().on("error", error => {
  logger.error(
    `Error creating redis for ${SUB_GRAPHS.network.name} service, Error: ${error}`
  );
});

const invitationAPi = new InvitationApi();
let invitationId = "";
const invitationData = {
  email: "test-email@mail.com",
  name: "testUser2",
  role: "ROLE_VENDOR",
};

const startServer = async () => {
  const schema = await createSchema();
  const server = new ApolloServer<Context>({
    schema,
  });
  await server.start();
  return server;
};

describe("Invitation API integration test", () => {
  let server: ApolloServer<Context>;

  beforeAll(async () => {
    server = await startServer();
  });
  afterAll(async () => {
    await server.stop();
  });
  it("should create an invitation", async () => {
    const CREATE_INVITATION = `mutation CreateInvitation($data: CreateInvitationInputDto!) {
    createInvitation(data: $data) {
      email
      expireAt
      id
      name
      role
      link
      userId
      status
    }
  }`;
    const res = await server.executeOperation(
      {
        query: CREATE_INVITATION,
        variables: {
          data: invitationData,
        },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.network.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: invitationAPi,
            },
            baseURL: baseURL.message,
            headers: parsedHeaders,
          };
        })(),
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    const { data } = singleResult;
    invitationId = data.createInvitation.id;
    expect(data.createInvitation.email).toEqual(invitationData.email);
    expect(data.createInvitation.name).toEqual(invitationData.name);
    expect(data.createInvitation).toHaveProperty("link");
    expect(data.createInvitation).toHaveProperty("expireAt");
  });

  it("should get invitation", async () => {
    const GET_INVITATION = `query GetInvitation($getInvitationId: String!) {
  getInvitation(id: $getInvitationId) {
    email
    expireAt
    id
    name
    role
    link
    userId
    status
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_INVITATION,
        variables: { getInvitationId: invitationId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.network.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: invitationAPi,
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
    expect(data).toHaveProperty("getInvitation");
    expect(data.getInvitation.email).toEqual(invitationData.email);
    expect(data.getInvitation.name).toEqual(invitationData.name);
    expect(data.getInvitation.role).toEqual(invitationData.role);
    expect(data.getInvitation).toHaveProperty("link");
    expect(data.getInvitation).toHaveProperty("expireAt");
  });

  it("should get invitation by org", async () => {
    const GET_ORG_INVITATION = `query GetInvitationsByOrg {
  getInvitationsByOrg {
    invitations {
      email
      expireAt
      id
      name
      role
      link
      userId
      status
    }
  }
}`;

    const res = await server.executeOperation(
      {
        query: GET_ORG_INVITATION,
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.network.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: invitationAPi,
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
    expect(data.getInvitationsByOrg).toHaveProperty("invitations");
    expect(data.getInvitationsByOrg.invitations.length).toBeGreaterThanOrEqual(
      1
    );
  });

  it("should get invitation by email", async () => {
    const GET_EMAIL_INVITATIONS = `query GetInvitations($email: String!) {
  getInvitations(email: $email) {
    email
    expireAt
    id
    name
    role
    link
    userId
    status
  }
}`;
    const res = await server.executeOperation(
      {
        query: GET_EMAIL_INVITATIONS,
        variables: { email: invitationData.email },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.network.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: invitationAPi,
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
    expect(data.getInvitations.email).toEqual(invitationData.email);
  });

  it("should update invitation", async () => {
    const UPDATE_INVITATION = `mutation UpdateInvitation($data: UpateInvitationInputDto!) {
  updateInvitation(data: $data) {
    id
  }
}`;

    const res = await server.executeOperation(
      {
        query: UPDATE_INVITATION,
        variables: {
          data: {
            id: invitationId,
            status: INVITATION_STATUS.INVITE_ACCEPTED,
          },
        },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.network.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: invitationAPi,
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
    expect(data.updateInvitation.id).toEqual(invitationId);
  });

  it("should delete invation using id", async () => {
    const DELETE_INVITATION = `mutation DeleteInvitation($deleteInvitationId: String!) {
    deleteInvitation(id: $deleteInvitationId) {
      id
    }
  }`;
    const res = await server.executeOperation(
      {
        query: DELETE_INVITATION,
        variables: { deleteInvitationId: invitationId },
      },
      {
        contextValue: await (async () => {
          const baseURL = await getBaseURL(
            SUB_GRAPHS.network.name,
            orgName,
            redisClient.isOpen ? redisClient : null
          );
          return {
            dataSources: {
              dataSource: invitationAPi,
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
    expect(data.deleteInvitation.id).toEqual(invitationId);
  });
});
