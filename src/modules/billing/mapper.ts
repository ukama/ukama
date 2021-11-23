import * as defaultCasual from "casual";
import { IBillMapper } from "./interface";
import {
    CurrentBillResponse,
    BillResponse,
    BillHistoryDto,
    BillHistoryResponse,
} from "./types";

class BillMapper implements IBillMapper {
    dtoToDto = (res: CurrentBillResponse): BillResponse => {
        const bill = res.data;
        let total = 0;
        for (let i = 0; i < bill.length; i++) {
            const subTotal = bill[i].subtotal;
            total = total + subTotal;
        }
        return {
            bill,
            total,
            dueDate: defaultCasual.date("10-10-2021"),
            billMonth: defaultCasual.date("11-10-2021"),
        };
    };
    billHistoryDtoToDto = (res: BillHistoryResponse): BillHistoryDto[] => {
        return res.data;
    };
}
export default <IBillMapper>new BillMapper();
