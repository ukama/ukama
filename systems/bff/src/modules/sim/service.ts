import { Service } from "typedi";
import { catchAsyncIOMethod } from "../../common";
import { ParsedCookie } from "../../common/types";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";
import { ISimService } from "./interface";
import SimMapper from "./mapper";
import { AllocateSimInputDto, SimResDto } from "./types";

@Service()
export class SimService implements ISimService {
    allocateSim = async (
        req: AllocateSimInputDto,
        cookie: ParsedCookie,
    ): Promise<SimResDto> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.PUT,
            path: `${SERVER.SUBSCRIBER_REGISTRY_API_URL}`,
            body: { ...req },
            headers: cookie.header,
        });
        if (checkError(res)) throw new Error(res.message);
        return SimMapper.dtoToSimResDto(res);
    };
}
