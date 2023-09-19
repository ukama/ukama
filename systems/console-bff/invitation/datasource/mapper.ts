import { InvitationAPIResDto, InvitationDto } from "../resolver/types";

export const dtoToInvitationResDto = (
  res: InvitationAPIResDto
): InvitationDto => {
  return {
    email: res.email,
    expiresAt: res.expires_at,
    id: res.id,
    link: res.link,
    userId: res.user_id,
    name: res.name,
    org: res.org,
    role: res.role,
    status: res.status,
  };
};
