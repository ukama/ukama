import { ParsedCookie, MetricServiceRes } from "../../common/types";
import {
    ActivateUserResponse,
    UserResDto,
    AddUserServiceRes,
    ConnectedUserDto,
    DeactivateResponse,
    GetUserDto,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserDto,
    OrgUserResponse,
    OrgUsersResponse,
    UserInputDto,
} from "./types";

export interface IUserService {
    getConnectedUsers(cookie: ParsedCookie): Promise<ConnectedUserDto>;
    updateUser(
        userId: string,
        req: UserInputDto,
        cookie: ParsedCookie
    ): Promise<UserResDto>;
    deactivateUser(
        id: string,
        cookie: ParsedCookie
    ): Promise<DeactivateResponse>;
    getUser(userId: string, cookie: ParsedCookie): Promise<GetUserDto>;
    getUsersByOrg(cookie: ParsedCookie): Promise<GetUsersDto[]>;
    addUser(
        req: UserInputDto,
        cookie: ParsedCookie
    ): Promise<UserResDto | null>;
    deleteUser(
        userId: string,
        cookie: ParsedCookie
    ): Promise<ActivateUserResponse>;
}

export interface IUserMapper {
    dtoToAddUserDto(res: AddUserServiceRes): UserResDto | null;
    connectedUsersDtoToDto(res: MetricServiceRes[]): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    dtoToUsersDto(req: OrgUsersResponse): GetUsersDto[];
    dtoToUserDto(req: OrgUserResponse): GetUserDto;
    dtoToUserResDto(req: OrgUserDto): UserResDto;
}
