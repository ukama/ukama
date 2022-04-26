import {
    ParsedCookie,
    PaginationDto,
    MetricServiceValueRes,
} from "../../common/types";
import {
    ActivateUserDto,
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
    ResidentResponse,
    ResidentsResponse,
    UserInputDto,
} from "./types";

export interface IUserService {
    getConnectedUsers(cookie: ParsedCookie): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(
        userId: string,
        req: UserInputDto,
        cookie: ParsedCookie
    ): Promise<UserResDto>;
    deactivateUser(id: string): Promise<DeactivateResponse>;
    getUser(userId: string, cookie: ParsedCookie): Promise<GetUserDto>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
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
    connectedUsersDtoToDto(res: MetricServiceValueRes[]): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    residentDtoToDto(res: GetUserResponseDto): ResidentResponse;
    dtoToUsersDto(req: OrgUsersResponse): GetUsersDto[];
    dtoToUserDto(req: OrgUserResponse): GetUserDto;
    dtoToUserResDto(req: OrgUserDto): UserResDto;
}
