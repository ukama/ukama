import {
  UserAPIResDto,
  UserResDto,
  WhoamiAPIDto,
  WhoamiDto,
} from "../resolver/types";

export const dtoToUserResDto = (res: UserAPIResDto): UserResDto => {
  return {
    uuid: res.user.id,
    email: res.user.email,
    isDeactivated: res.user.is_deactivated,
    name: res.user.name,
    authId: res.user.auth_id,
    phone: res.user.phone,
    registeredSince: res.user.registered_since,
  };
};
export const dtoToWhoamiResDto = (res: WhoamiAPIDto): WhoamiDto => {
  return {
    id: res.id,
    email: res.email,
    name: res.name,
    role: "admin",
    isFirstVisit: res.first_visit,
  };
};
