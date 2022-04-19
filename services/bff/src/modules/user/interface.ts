import {
    ParsedCookie,
    PaginationDto,
    MetricServiceRes,
} from "../../common/types";
import {
    ActivateUserDto,
    ActivateUserResponse,
    AddUserDto,
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
    UpdateUserDto,
} from "./types";

export interface IUserService {
    getConnectedUsers(cookie: ParsedCookie): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(req: UpdateUserDto, cookie: ParsedCookie): Promise<UserResDto>;
    deactivateUser(id: string): Promise<DeactivateResponse>;
    getUser(userId: string, cookie: ParsedCookie): Promise<GetUserDto>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
    getUsersByOrg(cookie: ParsedCookie): Promise<GetUsersDto[]>;
    addUser(req: AddUserDto, cookie: ParsedCookie): Promise<UserResDto | null>;
    deleteUser(
        userId: string,
        cookie: ParsedCookie
    ): Promise<ActivateUserResponse>;
}

export interface IUserMapper {
    dtoToAddUserDto(res: AddUserServiceRes): UserResDto | null;
    connectedUsersDtoToDto(res: MetricServiceRes[]): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    residentDtoToDto(res: GetUserResponseDto): ResidentResponse;
    dtoToUsersDto(req: OrgUsersResponse): GetUsersDto[];
    dtoToUserDto(req: OrgUserResponse): GetUserDto;
    dtoToUserResDto(req: OrgUserDto): UserResDto;
}
