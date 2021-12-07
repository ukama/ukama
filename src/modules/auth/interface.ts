import { LoginDto, ServerLoginDto } from "./types";

export interface IAuthService {
    getActionUrl(): Promise<string | null>;
    login(req: LoginDto, url: string): Promise<string | null>;
}

export interface IAuthMapper {
    dtoToBody(req: LoginDto): ServerLoginDto;
}
