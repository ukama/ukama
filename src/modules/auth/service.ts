import { Service } from "typedi";
import { LoginDto } from "./types";
import { IAuthService } from "./interface";
import { checkError } from "../../errors";
import AuthMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";

@Service()
export class AuthService implements IAuthService {
    public getActionUrl = async (): Promise<string | null> => {
        const flow = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_FLOW_ID,
        });
        if (checkError(flow)) return null;
        if (!flow.ui && !flow.ui.action) return null;

        return flow.ui.action;
    };
    login = async (req: LoginDto, url: string): Promise<string | null> => {
        const body = AuthMapper.dtoToBody(req);

        const loginRes = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.POST,
            path: url,
            body,
        });
        if (checkError(loginRes)) return null;
        if (!loginRes.session_token) return null;
        return loginRes.session_token;
    };
}
