import { Service } from "typedi";
import { BillResponse, CurrentBillResponse } from "./types";
import { IBillService } from "./interface";
import BillMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";

@Service()
export class BillService implements IBillService {
    public getCurrentBill = async (): Promise<BillResponse> => {
        const res = await catchAsyncIOMethod<CurrentBillResponse>({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_CURRENT_BILL,
        });
        const bill = BillMapper.dtoToDto(res);
        return bill;
    };
}
