/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";
import { GraphQLError } from "graphql";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
import { CBooleanResponse } from "../../common/types";
import {
  AddMemberInputDto,
  MemberDto,
  MembersResDto,
  UpdateMemberInputDto,
} from "../resolver/types";
import { dtoToMemberResDto, dtoToMembersResDto } from "./mapper";

class MemberApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW;
  getMembers = async (): Promise<MembersResDto> => {
    return this.get(`/${VERSION}/members`)
      .then(res => dtoToMembersResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getMember = async (id: string): Promise<MemberDto> => {
    return this.get(`/${VERSION}/members/${id}`)
      .then(res => dtoToMemberResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  removeMember = async (id: string): Promise<CBooleanResponse> => {
    return this.delete(`/${VERSION}/members/${id}`)
      .then(() => {
        return {
          success: true,
        };
      })
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  addMember = async (data: AddMemberInputDto): Promise<MemberDto> => {
    return this.post(`/${VERSION}/members`, {
      body: { user_uuid: data.userId, role: data.role },
    })
      .then(res => dtoToMemberResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  updateMember = async (
    memberId: string,
    req: UpdateMemberInputDto
  ): Promise<CBooleanResponse> => {
    return this.post(`/${VERSION}/members/${memberId}`, {
      body: {
        isDeactivated: req.isDeactivated,
        role: req.role,
      },
    })
      .then(() => {
        return {
          success: true,
        };
      })
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default MemberApi;
