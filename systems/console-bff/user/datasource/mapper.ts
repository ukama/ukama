import { MetricServiceValueRes } from "../../common/types";
import {
  ConnectedUserDto,
  GetUserDto,
  GetUserResponseDto,
  UserAPIResDto,
  UserResDto,
  WhoamiAPIDto,
  WhoamiDto,
} from "../resolver/types";

export const connectedUsersDtoToDto = (
  res: MetricServiceValueRes[]
): ConnectedUserDto => {
  if (res.length > 0) {
    const value: any = res[0].value[1];
    return { totalUser: value };
  }
  return { totalUser: "0" };
};
export const dtoToDto = (res: GetUserResponseDto): GetUserDto[] => {
  return res.data;
};
export const dtoToUserResDto = (res: UserAPIResDto): UserResDto => {
  return {
    uuid: res.user.uuid,
    email: res.user.email,
    isDeactivated: res.user.is_deactivated,
    name: res.user.name,
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
