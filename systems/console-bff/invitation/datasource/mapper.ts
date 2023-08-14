import {
  InvitationAPIResDto,
  InvitationDto,
} from "../resolver/types";

export const dtoToInvitationResDto = (res: InvitationAPIResDto): InvitationDto => {
  return {
    email: res.email,
    expiresAt: res.expires_at,
    id: res.id,
    link: res.link,
    org: res.org,
    status: res.status,
  };
};

