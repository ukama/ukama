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
        id: "date",
        label: "Date",
    },
    {
        id: "dataUsage",
        label: "Data Used",
    },
    {
        id: "subtotal",
        label: "subtotal",
    },
    {
        id: "total",
        label: "total",
    },
];
const NodeAppsColumns = [
    {
        id: "version",
        label: "Version",
        minWidth: 200,
    },
    {
        id: "date",
        label: "Date",
        minWidth: 200,
    },
    {
        id: "notes",
        label: "Notes",
        minWidth: 600,
    },
];

export { DataTableWithOptionColumns, NodeAppsColumns, CurrentBillColumns };
