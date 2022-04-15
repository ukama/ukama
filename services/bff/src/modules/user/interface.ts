import {
    ParsedCookie,
    PaginationDto,
    MetricServiceRes,
} from "../../common/types";
import {
    ActivateUserDto,
    ActivateUserResponse,
    AddUserDto,
    AddUserResponse,
    AddUserServiceRes,
    ConnectedUserDto,
    DeactivateResponse,
    GetUserDto,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserResponse,
    OrgUsersResponse,
    ResidentResponse,
    ResidentsResponse,
    UpdateUserDto,
    UserResponse,
} from "./types";

export interface IUserService {
    getConnectedUsers(cookie: ParsedCookie): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(req: UpdateUserDto): Promise<UserResponse>;
    deactivateUser(id: string): Promise<DeactivateResponse>;
    getUser(userId: string, cookie: ParsedCookie): Promise<GetUserDto>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
    getUsersByOrg(cookie: ParsedCookie): Promise<GetUsersDto[]>;
    addUser(
        req: AddUserDto,
        cookie: ParsedCookie
    ): Promise<AddUserResponse | null>;
    deleteUser(
        userId: string,
        cookie: ParsedCookie
    ): Promise<ActivateUserResponse>;
}

export interface IUserMapper {
    dtoToAddUserDto(res: AddUserServiceRes): AddUserResponse | null;
    connectedUsersDtoToDto(res: MetricServiceRes[]): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    residentDtoToDto(res: GetUserResponseDto): ResidentResponse;
    dtoToUsersDto(req: OrgUsersResponse): GetUsersDto[];
    dtoToUserDto(req: OrgUserResponse): GetUserDto;
}
