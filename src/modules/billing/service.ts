import { Service } from "typedi";
import { BillHistoryDto, BillResponse } from "./types";
import { IBillService } from "./interface";
import BillMapper from "./mapper";
import { catchAsyncIOMethod } from "../../common";
import { API_METHOD_TYPE } from "../../constants";
import { SERVER } from "../../constants/endpoints";
import { checkError } from "../../errors";

@Service()
export class BillService implements IBillService {
    public getCurrentBill = async (): Promise<BillResponse> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_CURRENT_BILL,
        });
        if (checkError(res)) throw new Error(res.message);
        const bill = BillMapper.dtoToDto(res);
        return bill;
    };

    public getBillHistory = async (): Promise<BillHistoryDto[]> => {
        const res = await catchAsyncIOMethod({
            type: API_METHOD_TYPE.GET,
            path: SERVER.GET_BILL_HISTORY,
        });
        if (checkError(res)) throw new Error(res.message);
        const bill = BillMapper.billHistoryDtoToDto(res);
        return bill;
    };
}
