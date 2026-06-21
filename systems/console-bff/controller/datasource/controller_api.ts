/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { VERSION } from "../../common/configs";
import { BaseRESTDataSource } from "../../common/datasource";
import { TBooleanResponse } from "../../common/types";
import {
  RestartNodeInputDto,
  RestartNodesInputDto,
  RestartSiteInputDto,
  ToggleInternetSwitchInputDto,
  SetSiteInputDto,
  ToggleSiteStatusInputDto,
} from "../resolvers/types";

const CONTROLLER = "controller";

class ControllerApi extends BaseRESTDataSource {
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
      .catch(error => {
        this.logger.error(`Request failed: ${error}`);
        return { success: false };
      });
  };

  restartNodes = async (
    baseURL: string,
    req: RestartNodesInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `RestartNodes [POST]: ${baseURL}/${VERSION}/${CONTROLLER}/networks/${req.networkId}/restart-nodes`
    );
    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/${CONTROLLER}/networks/${req.networkId}/restart-nodes`,
      {
        body: req.nodeIds,
      }
    )
      .then(() => {
        return { success: true };
      })
      .catch(error => {
        this.logger.error(`Request failed: ${error}`);
        return { success: false };
      });
  };

  restartSite = async (
    baseURL: string,
    req: RestartSiteInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `RestartSite [POST]: ${baseURL}/${VERSION}/${CONTROLLER}/networks/${req.networkId}/sites/${req.siteId}/restart`
    );
    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/${CONTROLLER}/networks/${req.networkId}/sites/${req.siteId}/restart`
    )
      .then(() => {
        return { success: true };
      })
      .catch(error => {
        this.logger.error(`Request failed: ${error}`);
        return { success: false };
      });
  };

  toggleInternetSwitch = async (
    baseURL: string,
    req: ToggleInternetSwitchInputDto
  ): Promise<TBooleanResponse> => {
    this.logger.info(
      `ToggleInternetSwitch [POST]: ${baseURL}/${VERSION}/sites/${req.siteId}/internet-port`
    );

    this.baseURL = baseURL;
    return this.post(
      `/${VERSION}/sites/${req.siteId}/internet-port`,
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
      .catch(error => {
        this.logger.error(`Request failed: ${error}`);
        return { success: false };
      });
  };
  toggleRFStatus = async (
    baseURL: string,
    req: ToggleSiteStatusInputDto
  ): Promise<TBooleanResponse> => {
    const state = req.status ? "on" : "off";
    this.logger.info(
      `ToggleRFStatus [POST]: ${baseURL}/${VERSION}/sites/${req.siteId}/radio/${state}`
    );

    this.baseURL = baseURL;
    return this.post(`/${VERSION}/sites/${req.siteId}/radio/${state}`)
      .then(() => {
        return { success: true };
      })
      .catch(error => {
        this.logger.error(`Request failed: ${error}`);
        return { success: false };
      });
  };
  toggleService = async (
    baseURL: string,
    req: ToggleSiteStatusInputDto
  ): Promise<TBooleanResponse> => {
    const state = req.status ? "on" : "off";
    this.logger.info(
      `ToggleServiceStatus [POST]: ${baseURL}/${VERSION}/sites/${req.siteId}/service/${state}`
    );

    this.baseURL = baseURL;
    return this.post(`/${VERSION}/sites/${req.siteId}/service/${state}`)
      .then(() => {
        return { success: true };
      })
      .catch(error => {
        this.logger.error(`Request failed: ${error}`);
        return { success: false };
      });
  };
  setSite = async (
    baseURL: string,
    req: SetSiteInputDto
  ): Promise<TBooleanResponse> => {
    const state = req.status ? "on" : "off";
    this.logger.info(
      `SetSite [POST]: ${baseURL}/${VERSION}/sites/${req.siteId}/${state}`
    );

    this.baseURL = baseURL;
    return this.post(`/${VERSION}/sites/${req.siteId}/${state}`)
      .then(() => {
        return { success: true };
      })
      .catch(error => {
        this.logger.error(`Request failed: ${error}`);
        return { success: false };
      });
  };
}

export default ControllerApi;
