import {
    HeaderType,
    MetricServiceRes,
    PaginationDto,
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
    GetUserPaginationDto,
    GetUserResponse,
    GetUserResponseDto,
    GetUsersDto,
    OrgUserResponse,
    OrgUsersResponse,
    ResidentResponse,
    ResidentsResponse,
    UpdateUserDto,
    UserInput,
    UserResponse,
} from "./types";

export interface IUserService {
    getConnectedUsers(
        orgId: string,
        header: HeaderType
    ): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(req: UpdateUserDto): Promise<UserResponse>;
    deactivateUser(id: string): Promise<DeactivateResponse>;
    getUser(data: UserInput, header: HeaderType): Promise<GetUserDto>;
    getUsers(req: GetUserPaginationDto): Promise<GetUserResponse>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
    getUsersByOrg(orgId: string, header: HeaderType): Promise<GetUsersDto[]>;
    addUser(
        orgId: string,
        req: AddUserDto,
        header: HeaderType
    ): Promise<AddUserResponse | null>;
    deleteUser(
        orgId: string,
        userId: string,
        header: HeaderType
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
