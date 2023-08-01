import {
  BillHistoryDto,
  BillHistoryResponse,
  BillResponse,
  CurrentBillResponse,
} from "../resolvers/types";

export const dtoToDto = (res: CurrentBillResponse): BillResponse => {
  const bill = res.data;
  let total = 0;
  for (const sub of bill) {
    const subTotal = sub.subtotal;
    total = total + subTotal;
  }
  return {
    bill,
    total,
    dueDate: "10-10-2021",
    billMonth: "11-10-2021",
  };
};

export const billHistoryDtoToDto = (
  res: BillHistoryResponse
): BillHistoryDto[] => {
  return res.data;
};
