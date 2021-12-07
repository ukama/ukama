import { PaginationDto } from "../../common/types";
import { TIME_FILTER } from "../../constants";
import {
    ActivateUserDto,
    ActivateUserResponse,
    ConnectedUserDto,
    ConnectedUserResponse,
    DeactivateResponse,
    GetUserDto,
    GetUserPaginationDto,
    GetUserResponse,
    GetUserResponseDto,
    OrganisationDto,
    ResidentResponse,
    ResidentsResponse,
    UpdateUserDto,
    UserResponse,
} from "./types";

export interface IUserService {
    getConnectedUsers(filter: TIME_FILTER): Promise<ConnectedUserDto>;
    activateUser(req: ActivateUserDto): Promise<ActivateUserResponse>;
    updateUser(req: UpdateUserDto): Promise<UserResponse>;
    deactivateUser(id: string): Promise<DeactivateResponse>;
    getUser(id: string): Promise<GetUserDto>;
    getUsers(req: GetUserPaginationDto): Promise<GetUserResponse>;
    getResidents(req: PaginationDto): Promise<ResidentsResponse>;
    whoAmI(token: string): Promise<OrganisationDto | null>;
}

export interface IUserMapper {
    connectedUsersDtoToDto(res: ConnectedUserResponse): ConnectedUserDto;
    dtoToDto(res: GetUserResponseDto): GetUserDto[];
    residentDtoToDto(res: GetUserResponseDto): ResidentResponse;
}
