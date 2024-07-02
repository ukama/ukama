/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import {
  CreateInvitationInputDto,
  DeleteInvitationResDto,
  InvitationDto,
  InvitationsResDto,
  UpateInvitationInputDto,
  UpdateInvitationResDto,
} from "../resolver/types";
import { dtoToInvitationsResDto, inviteResToInvitationDto } from "./mapper";

const INVITATIONS = "invitations";

class InvitationApi extends RESTDataSource {
  sendInvitation = async (
    baseURL: string,
    req: CreateInvitationInputDto
  ): Promise<InvitationDto> => {
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${INVITATIONS}`, {
      body: { ...req },
    }).then(res => inviteResToInvitationDto(res));
  };

  getInvitation = async (
    baseURL: string,
    id: string
  ): Promise<InvitationDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}/${id}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };

  updateInvitation = async (
    baseURL: string,
    req: UpateInvitationInputDto
  ): Promise<UpdateInvitationResDto> => {
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/${INVITATIONS}/${req.id}`, {
      body: { status: req.status },
    }).then(res => res);
  };

  deleteInvitation = async (
    baseURL: string,
    id: string
  ): Promise<DeleteInvitationResDto> => {
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/${INVITATIONS}/${id}`).then(res => res);
  };

  getInvitationsByOrg = async (baseURL: string): Promise<InvitationsResDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}`).then(res =>
      dtoToInvitationsResDto(res)
    );
  };

  getInvitationsByEmail = async (
    baseURL: string,
    email: string
  ): Promise<InvitationDto> => {
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}/user/${email}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };
}

export default InvitationApi;
