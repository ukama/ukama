/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { INVITATION_STATUS } from "../../common/enums";
import { logger } from "../../common/logger";
import { addInStore, getFromStore, openStore } from "../../common/storage";
import { getBaseURL } from "../../common/utils";
import InitAPI from "../../init/datasource/init_api";
import {
  CreateInvitationInputDto,
  DeleteInvitationResDto,
  InvitationDto,
  InvitationsDto,
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
    this.logger.info(
      `SendInvitation [POST]: ${baseURL}/${VERSION}/${INVITATIONS}`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${INVITATIONS}`, {
      body: { ...req },
    }).then(res => inviteResToInvitationDto(res));
  };

  getInvitation = async (
    baseURL: string,
    id: string
  ): Promise<InvitationDto> => {
    this.logger.info(
      `GetInvitation [GET]: ${baseURL}/${VERSION}/${INVITATIONS}/${id}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}/${id}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };

  updateInvitation = async (
    req: UpateInvitationInputDto
  ): Promise<UpdateInvitationResDto> => {
    const store = openStore();
    const baseURL = await getFromStore(store, `${req.email}/${req.id}`);
    this.logger.info(
      `UpdateInvitation [PATCH]: ${baseURL}/${VERSION}/${INVITATIONS}/${req.id}`
    );
    this.baseURL = baseURL;
    return this.patch(`/${VERSION}/${INVITATIONS}/${req.id}`, {
      body: { status: req.status },
    }).then(res => res);
  };

  deleteInvitation = async (
    baseURL: string,
    id: string
  ): Promise<DeleteInvitationResDto> => {
    this.logger.info(
      `DeleteInvitation [DELETE]: ${baseURL}/${VERSION}/${INVITATIONS}/${id}`
    );
    this.baseURL = baseURL;
    return this.delete(`/${VERSION}/${INVITATIONS}/${id}`).then(res => res);
  };

  getInvitations = async (baseURL: string): Promise<InvitationsResDto> => {
    this.logger.info(
      `GetInvitationByOrg [GET]: ${baseURL}/${VERSION}/${INVITATIONS}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}`).then(res =>
      dtoToInvitationsResDto(res)
    );
  };

  getAllInvitationsByEmail = async (email: string): Promise<InvitationsDto> => {
    const invitations: InvitationsDto = { invitations: [] };
    const init = new InitAPI();
    const orgsName = await init.getOrgs();
    if (orgsName?.orgs?.length > 0) {
      const store = openStore();
      for (const element of orgsName.orgs) {
        const baseURL = await getBaseURL("invitation", element.name, store);
        logger.info(`BaseURL: ${JSON.stringify(baseURL)}`);
        if (baseURL.status === 200) {
          const res = await this.getInvitationsByEmail(baseURL.message, email);
          logger.info(`Invitations res: ${JSON.stringify(res)}`);
          if (res && res.status !== INVITATION_STATUS.INVITE_ACCEPTED) {
            await addInStore(store, `${email}/${res.id}`, baseURL);
            invitations.invitations.push({
              id: res.id,
              name: res.name,
              link: res.link,
              role: res.role,
              email: res.email,
              status: res.status,
              userId: res.userId,
              expireAt: res.expireAt,
            });
          }
        }
      }
    }

    this.logger.info(`Invitations: ${JSON.stringify(invitations)}`);

    return invitations;
  };

  getInvitationsByEmail = async (
    baseURL: string,
    email: string
  ): Promise<InvitationDto> => {
    this.logger.info(
      `GetInvitationByEmail [GET]: ${baseURL}/${VERSION}/${INVITATIONS}/user/${email}`
    );
    this.baseURL = baseURL;
    return this.get(`/${VERSION}/${INVITATIONS}/user/${email}`).then(res =>
      inviteResToInvitationDto(res)
    );
  };
}

export default InvitationApi;
