import {
    BillHistoryDto,
    BillHistoryResponse,
    BillResponse,
    CurrentBillResponse,
} from "./types";

export interface IBillService {
    getCurrentBill(): Promise<BillResponse>;
    getBillHistory(): Promise<BillHistoryDto[]>;
}

export interface IBillMapper {
    dtoToDto(data: CurrentBillResponse): BillResponse;
    billHistoryDtoToDto(res: BillHistoryResponse): BillHistoryDto[];
}
