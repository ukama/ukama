import { MetricServiceValueRes } from "../../common/types";
import {
    ConnectedUserDto,
    GetUserDto,
    GetUserResponseDto,
    UserAPIResDto,
    UserResDto,
    WhoamiAPIDto,
    WhoamiDto,
} from "./types";
// eslint-disable-next-line @typescript-eslint/no-empty-interface
export interface IUserService {}

export interface IUserMapper {
    connectedUsersDtoToDto(res: MetricServiceValueRes[]): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    dtoToUserResDto(res: UserAPIResDto): UserResDto;
    dtoToWhoamiResDto(res: WhoamiAPIDto): WhoamiDto;
}
