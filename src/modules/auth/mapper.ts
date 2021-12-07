import { IAuthMapper } from "./interface";
import { LoginDto, ServerLoginDto } from "./types";

class AuthMapper implements IAuthMapper {
    dtoToBody = (req: LoginDto): ServerLoginDto => {
        return {
            password_identifier: req.email,
            password: req.password,
            method: "password",
        };
    };
}
export default <IAuthMapper>new AuthMapper();
