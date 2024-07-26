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
    this.logger.info(`[POST]: ${baseURL}/${VERSION}/${INVITATIONS}`);
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${INVITATIONS}`, {
      body: { ...req },
    }).then(res => inviteResToInvitationDto(res));
  };

  getInvitation = async (
    baseURL: string,
    id: string
  ): Promise<InvitationDto> => {
    this.logger.info(`[GET]: ${baseURL}/${VERSION}/${INVITATIONS}/${id}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}/${id}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };

  updateInvitation = async (
    baseURL: string,
    req: UpateInvitationInputDto
  ): Promise<UpdateInvitationResDto> => {
    this.logger.info(`[PATCH]: ${baseURL}/${VERSION}/${INVITATIONS}/${req.id}`);
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/${INVITATIONS}/${req.id}`, {
      body: { status: req.status },
    }).then(res => res);
  };

  deleteInvitation = async (
    baseURL: string,
    id: string
  ): Promise<DeleteInvitationResDto> => {
    this.logger.info(`[DELETE]: ${baseURL}/${VERSION}/${INVITATIONS}/${id}`);
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/${INVITATIONS}/${id}`).then(res => res);
  };

  getInvitationsByOrg = async (baseURL: string): Promise<InvitationsResDto> => {
    this.logger.info(`[GET]: ${baseURL}/${VERSION}/${INVITATIONS}`);
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}`).then(res =>
      dtoToInvitationsResDto(res)
    );
  };

  getInvitationsByEmail = async (
    baseURL: string,
    email: string
  ): Promise<InvitationDto> => {
    this.logger.info(
      `[GET]: ${baseURL}/${VERSION}/${INVITATIONS}/user/${email}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}/user/${email}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };
}

export default InvitationApi;
