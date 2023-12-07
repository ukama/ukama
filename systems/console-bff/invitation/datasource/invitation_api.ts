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
import {
  DeleteInvitationResDto,
  GetInvitationByOrgResDto,
  InvitationDto,
  SendInvitationInputDto,
  SendInvitationResDto,
  UpateInvitationInputDto,
  UpdateInvitationResDto,
} from "../resolver/types";
import { dtoToInvitationResDto } from "./mapper";

const version = "/v1/invitation";

class InvitationApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW + version;

  sendInvitation = async (
    req: SendInvitationInputDto
  ): Promise<SendInvitationResDto> => {
    return this.post(``, {
      body: { ...req },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  getInvitation = async (id: string): Promise<InvitationDto> => {
    return this.get(`/${VERSION}/invitation/${id}`)
      .then(res => dtoToInvitationResDto(res))
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  updateInvitation = async (
    id: string,
    req: UpateInvitationInputDto
  ): Promise<UpdateInvitationResDto> => {
    return this.put(`/${VERSION}/invitation/${id}`, {
      body: { status: req.status },
    })
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };

  deleteInvitation = async (id: string): Promise<DeleteInvitationResDto> => {
    return this.delete(`/${VERSION}/invitation/${id}`).then(res => res);
  };

  getInvitationsByOrg = async (
    orgName: string
  ): Promise<GetInvitationByOrgResDto> => {
    return this.get(`/${VERSION}/invitation/${orgName}`)
      .then(res => res)
      .catch(err => {
        throw new GraphQLError(err);
      });
  };
}

export default InvitationApi;
