/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { logger } from "../../common/logger";
import { CBooleanResponse } from "../../common/types";
import {
  AddMemberInputDto,
  MemberDto,
  MembersResDto,
  UpdateMemberInputDto,
} from "../resolver/types";
import { dtoToMemberResDto, dtoToMembersResDto } from "./mapper";

class MemberApi extends RESTDataSource {
  getMembers = async (baseURL: string): Promise<MembersResDto> => {
    this.logger.info(`GetMembers [GET]: ${baseURL}/${VERSION}/members`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/members`).then(res => dtoToMembersResDto(res));
  };

  getMember = async (baseURL: string, id: string): Promise<MemberDto> => {
    this.logger.info(`GetMember [GET]: ${baseURL}/${VERSION}/members/${id}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/members/${id}`).then(res =>
      dtoToMemberResDto(res)
    );
  };

  getMemberByUserId = async (
    baseURL: string,
    userId: string
  ): Promise<MemberDto> => {
    this.logger.info(
      `GetMemberByUserId [GET]: ${baseURL}/${VERSION}/members/user/${userId}`
    );
    this.baseURL = baseURL;
    logger.info(`Request Url: ${baseURL}/${VERSION}/members/user/${userId}`);
    return this.get(`/${VERSION}/members/user/${userId}`)
      .then(res => {
        return dtoToMemberResDto(res);
      })
      .catch(err => {
        logger.error(`Error: ${err}`);
        return {
          id: "",
          userId: "",
          role: "",
          isDeactivated: false,
          orgId: "",
          createdAt: "",
          updatedAt: "",
          memberId: "",
          memberSince: "",
          email: "",
          name: "",
        };
      });
  };

  removeMember = async (
    baseURL: string,
    id: string
  ): Promise<CBooleanResponse> => {
    this.logger.info(
      `RemoveMember [DELETE]: ${baseURL}/${VERSION}/members/${id}`
    );
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/members/${id}`).then(() => {
      return {
        success: true,
      };
    });
  };

  addMember = async (
    baseURL: string,
    data: AddMemberInputDto
  ): Promise<MemberDto> => {
    this.logger.info(`AddMember [POST]: ${baseURL}/${VERSION}/members`);
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/members`, {
      body: { user_uuid: data.userId, role: data.role },
    }).then(res => dtoToMemberResDto(res));
  };

  updateMember = async (
    baseURL: string,
    memberId: string,
    req: UpdateMemberInputDto
  ): Promise<CBooleanResponse> => {
    this.logger.info(
      `UpdateMember [PATCH]: ${baseURL}/${VERSION}/members/${memberId}`
    );
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/members/${memberId}`, {
      body: {
        isDeactivated: req.isDeactivated,
        role: req.role,
      },
    }).then(() => {
      return {
        success: true,
      };
    });
  };
}

export default MemberApi;
