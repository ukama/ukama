import { RESTDataSource } from "@apollo/datasource-rest";

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
    }).then(res => res);
  };

  getInvitation = async (id: string): Promise<InvitationDto> => {
    return this.get(`/${VERSION}/invitation/${id}`).then(res =>
      dtoToInvitationResDto(res)
    );
  };

  updateInvitation = async (
    id: string,
    req: UpateInvitationInputDto
  ): Promise<UpdateInvitationResDto> => {
    return this.put(`/${VERSION}/invitation/${id}`, {
      body: { status: req.status },
    }).then(res => res);
  };

  deleteInvitation = async (id: string): Promise<DeleteInvitationResDto> => {
    return this.delete(`/${VERSION}/invitation/${id}`).then(res => res);
  };

  getInvitationsByOrg = async (
    orgName: string
  ): Promise<GetInvitationByOrgResDto> => {
    return this.get(`/${VERSION}/invitation/${orgName}`).then(res => res);
  };
}

export default InvitationApi;
