import { ColumnsWithOptions } from "../types";

const DataTableWithOptionColumns: ColumnsWithOptions[] = [
    { id: "name", label: "Name", minWidth: 150 },
    { id: "usage", label: "Usage", minWidth: 100 },
    {
        id: "actions",
        label: "...",
        minWidth: 50,
        align: "right",
    },
];

export { DataTableWithOptionColumns };
