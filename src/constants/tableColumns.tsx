import { ColumnsWithOptions } from "../types";

const DataTableWithOptionColumns: ColumnsWithOptions[] = [
    { id: "name", label: "Name", minWidth: 150 },
    { id: "dataUsage", label: "Usage", minWidth: 100 },
    {
        label: "",
        minWidth: 50,
        id: "actions",
        align: "right",
    },
];

const CurrentBillColumns: ColumnsWithOptions[] = [
    {
        id: "name",
        label: "Name",
    },
    {
        id: "dataUsage",
        label: "Data Used",
    },
    {
        id: "rate",
        label: "Rate",
    },
    {
        id: "subTotal",
        label: "SubTotal",
    },
];

export { DataTableWithOptionColumns, CurrentBillColumns };
