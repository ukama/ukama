/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { REGISTRY_API_GW, VERSION } from "../../common/configs";
import {
  CreateInvitationInputDto,
  DeleteInvitationResDto,
  InvitationDto,
  InvitationsResDto,
  UpateInvitationInputDto,
  UpdateInvitationResDto,
} from "../resolver/types";
import { dtoToInvitationsResDto, inviteResToInvitationDto } from "./mapper";

const version = "/v1/invitation";

class InvitationApi extends RESTDataSource {
  baseURL = REGISTRY_API_GW + version;

  sendInvitation = async (
    req: CreateInvitationInputDto
  ): Promise<InvitationDto> => {
    return this.post(`/${VERSION}/invitations`, {
      body: { ...req },
    }).then(res => inviteResToInvitationDto(res));
  };

  getInvitation = async (id: string): Promise<InvitationDto> => {
    return this.get(`/${VERSION}/invitations/${id}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };

  updateInvitation = async (
    id: string,
    req: UpateInvitationInputDto
  ): Promise<UpdateInvitationResDto> => {
    return this.put(`/${VERSION}/invitations/${id}`, {
      body: { status: req.status },
    }).then(res => res);
  };

  deleteInvitation = async (id: string): Promise<DeleteInvitationResDto> => {
    return this.delete(`/${VERSION}/invitations/${id}`).then(res => res);
  };

  getInvitationsByOrg = async (): Promise<InvitationsResDto> => {
    return this.get(`/${VERSION}/invitations`).then(res =>
      dtoToInvitationsResDto(res)
    );
  };

  getInvitationsByEmail = async (email: string): Promise<InvitationDto> => {
    return this.get(`/${VERSION}/invitations/${email}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };
}

export default InvitationApi;
