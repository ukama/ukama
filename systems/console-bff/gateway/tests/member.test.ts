import { ApolloServer } from "@apollo/server";
import { faker } from "@faker-js/faker";
import "reflect-metadata";
import { buildSchema } from "type-graphql";

import { SUB_GRAPHS } from "../../common/configs";
import { ROLE_TYPE } from "../../common/enums";
import { openStore } from "../../common/storage";
import { THeaders } from "../../common/types";
import { getBaseURL, parseGatewayHeaders } from "../../common/utils";
import { Context } from "../../member/context";
import MemberApi from "../../member/datasource/member_api";
import { AddMemberResolver } from "../../member/resolver/addMember";
import { GetMemberResolver } from "../../member/resolver/getMember";
import { GetMemberByUserIdResolver } from "../../member/resolver/getMemberByUserId";
import { GetMembersResolver } from "../../member/resolver/getMembers";
import { RemoveMemberResolver } from "../../member/resolver/removeMember";
import { UpdateMemberResolver } from "../../member/resolver/updateMember";
import {
  ADD_MEMBER,
  GET_MEMBER,
  GET_MEMBERS,
  GET_MEMBER_BY_ID,
  REMOVE_MEMBER,
  UPDATE_MEMBER,
} from "./graphql";

const token = process.env.TOKEN;
const headers = {
  cookie: "ukama_session=random-session",
  token: token,
};

const parsedHeaders = parseGatewayHeaders(headers);
const { orgName, userId } = parsedHeaders;

const createSchema = async () => {
  return await buildSchema({
    resolvers: [
      AddMemberResolver,
      GetMemberResolver,
      GetMembersResolver,
      GetMemberByUserIdResolver,
      RemoveMemberResolver,
      UpdateMemberResolver,
    ],
    validate: true,
  });
};

const memberApi = new MemberApi();

const testMember = {
  role: ROLE_TYPE.ROLE_USER,
  userId: faker.string.uuid(),
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
  const baseURL = await getBaseURL(SUB_GRAPHS.member.name, orgName, store);
  return {
    dataSources: { dataSource: memberApi },
    baseURL: baseURL.message,
    headers: parsedHeaders,
  };
};

describe("Member API integration test", () => {
  let server: ApolloServer<Context>;
  let contextValue: {
    dataSources: { dataSource: MemberApi };
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
  let memberId = "";

  it("should get all members", async () => {
    const res = await server.executeOperation(
      {
        query: GET_MEMBERS,
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getMembers).toBeDefined();
  });

  it("should add a member", async () => {
    const res = await server.executeOperation(
      {
        query: ADD_MEMBER,
        variables: { data: testMember },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.addMember).toBeDefined();
    expect(data.addMember.memberId).toBeDefined();
    memberId = data.addMember.memberId;
    expect(data.addMember.userId).toEqual(testMember.userId);
    expect(data.addMember.role).toEqual(testMember.role);
  });

  it("should get member using user id", async () => {
    const res = await server.executeOperation(
      {
        query: GET_MEMBER,
        variables: { userId: userId },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getMemberByUserId.userId).toEqual(userId);
  });

  it("should get member using member id", async () => {
    const res = await server.executeOperation(
      {
        query: GET_MEMBER_BY_ID,
        variables: { getMemberId: memberId },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.getMember).toBeDefined();
    expect(data.getMember.memberId).toEqual(memberId);
  });

  it("should be able to update a member", async () => {
    const res = await server.executeOperation(
      {
        query: UPDATE_MEMBER,
        variables: {
          data: {
            role: ROLE_TYPE.ROLE_VENDOR,
            isDeactivated: true,
          },
          memberId,
        },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.updateMember.success).toBeTruthy();
  });

  it("should remove a member", async () => {
    const res = await server.executeOperation(
      {
        query: REMOVE_MEMBER,
        variables: { removeMemberId: memberId },
      },
      {
        contextValue: contextValue,
      }
    );
    const body = JSON.stringify(res.body);
    const { singleResult } = JSON.parse(body);
    expect(singleResult.errors).toBeUndefined();
    const { data } = singleResult;
    expect(data.removeMember.success).toBeTruthy();
  });
});
