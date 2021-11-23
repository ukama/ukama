import { BillResponse, CurrentBillResponse } from "./types";

export interface IBillService {
    getCurrentBill(): Promise<BillResponse>;
}

export interface IBillMapper {
    dtoToDto(data: CurrentBillResponse): BillResponse;
}
