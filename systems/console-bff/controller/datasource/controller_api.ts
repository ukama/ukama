/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { RESTDataSource } from "@apollo/datasource-rest";

import { VERSION } from "../../common/configs";
import { TBooleanResponse } from "../../common/types";
import {
  RestartNodeInputDto,
  RestartNodesInputDto,
  RestartSiteInputDto,
  ToggleInternetSwitchInputDto,
  ToggleRFStatusInputDto,
} from "../resolvers/types";

const CONTROLLER = "controller";

class ControllerApi extends RESTDataSource {
  restartNode = async (
    baseURL: string,
    req: RestartNodeInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `RestartNode [POST]: ${baseURL}/${VERSION}/nodes/${req.nodeId}/restart`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/${CONTROLLER}/nodes/${req.nodeId}/restart`)
      .then(() => {
        return { success: true };
      })
      .catch(() => {
        return { success: false };
      });
  };

  restartNodes = async (
    baseURL: string,
    req: RestartNodesInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `RestartNodes [POST]: ${baseURL}/${VERSION}/network/${req.networkId}/restart-nodes`
    );
    this.baseURL = baseURL;
    return this.post(`/${VERSION}/network/${req.networkId}/restart-nodes`, {
      body: req.nodeIds,
    })
      .then(() => {
        return { success: true };
      })
      .catch(() => {
        return { success: false };
      });
  };

  restartSite = async (
    baseURL: string,
    req: RestartSiteInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `RestartSite [POST]: ${baseURL}/${VERSION}/networks/${req.networkId}/sites/${req.siteId}/restart`
    );
    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/networks/${req.networkId}/sites/${req.siteId}/restart`
    )
      .then(() => {
        return { success: true };
      })
      .catch(() => {
        return { success: false };
      });
  };

  toggleInternetSwitch = async (
    baseURL: string,
    req: ToggleInternetSwitchInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `ToggleInternetSwitch [POST]: ${baseURL}/${VERSION}/${CONTROLLER}/sites/${req.siteId}/toggle-internet-switch`
    );

    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/${CONTROLLER}/sites/${req.siteId}/toggle-internet-port`,
      {
        body: {
          port: req.port,
          status: req.status,
        },
      }
    )
      .then(() => {
        return { success: true };
      })
      .catch(() => {
        return { success: false };
      });
  };
  toggleRFStatus = async (
    baseURL: string,
    req: ToggleRFStatusInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `ToggleRFStatus [POST]: ${baseURL}/${VERSION}/${CONTROLLER}/nodes/${req.nodeId}/toggle-rf`
    );

    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/${CONTROLLER}/nodes/${req.nodeId}/toggle-rf`,
      {
        body: {
          status: req.status,
        },
      }
    )
      .then(() => {
        return { success: true };
      })
      .catch(() => {
        return { success: false };
      });
  };
}

export default ControllerApi;
