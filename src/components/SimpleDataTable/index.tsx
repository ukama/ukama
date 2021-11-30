import {
    Table,
    TableRow,
    Checkbox,
    TableBody,
    TableCell,
    TableHead,
    Typography,
    TableContainer,
} from "@mui/material";
import { colors } from "../../theme";
import { ColumnsWithOptions } from "../../types";

interface SimpleDataTableInterface {
    dataset: Object[];
    totalRows?: number;
    setSelectedRows?: any;
    selectedRows?: number[];
    rowSelection?: boolean;
    columns: ColumnsWithOptions[];
}

const SimpleDataTable = ({
    columns,
    dataset,
    totalRows = 0,
    setSelectedRows,
    selectedRows = [],
    rowSelection = false,
}: SimpleDataTableInterface) => {
    const onRowSelection = (id: number) => {
        setSelectedRows([...selectedRows, id]);
    };

    const onRowsSelection = () => {
        if (selectedRows.length === totalRows) setSelectedRows([]);
        else setSelectedRows(dataset.map((i: any) => i?.id));
    };

    return (
        <TableContainer sx={{ overflowX: "hidden", mt: "24px" }}>
            <Table stickyHeader>
                <TableHead>
                    <TableRow>
                        {rowSelection && (
                            <TableCell padding="checkbox">
                                <Checkbox
                                    color="primary"
                                    indeterminate={
                                        selectedRows.length > 0 &&
                                        selectedRows.length < totalRows
                                    }
                                    checked={
                                        totalRows > 0 &&
                                        selectedRows.length === totalRows
                                    }
                                    onChange={onRowsSelection}
                                    inputProps={{
                                        "aria-label": "select all desserts",
                                    }}
                                />
                            </TableCell>
                        )}
                        {columns.map(column => (
                            <TableCell
                                key={column.id}
                                align={column.align}
                                style={{
                                    padding: "0px 16px 12px 16px",
                                    fontSize: "0.875rem",
                                    minWidth: column.minWidth,
                                }}
                            >
                                <b>{column.label}</b>
                            </TableCell>
                        ))}
                    </TableRow>
                </TableHead>
                <TableBody>
                    {dataset.map((row: any) => (
                        <TableRow
                            key={row.id}
                            sx={{
                                ":hover": {
                                    backgroundColor: colors.solitude,
                                },
                            }}
                            selected={selectedRows.includes(row.id)}
                            role={rowSelection ? "checkbox" : "row"}
                            onClick={() => onRowSelection(row.id)}
                        >
                            {rowSelection && (
                                <TableCell padding="checkbox">
                                    <Checkbox
                                        color="primary"
                                        inputProps={{
                                            "aria-labelledby": row.id,
                                        }}
                                        checked={selectedRows.includes(row.id)}
                                    />
                                </TableCell>
                            )}
                            {columns.map(
                                (column: ColumnsWithOptions, index: number) => (
                                    <TableCell
                                        key={`${row.name}-${index}`}
                                        sx={{
                                            padding: 1,
                                            fontSize: "0.875rem",
                                        }}
                                    >
                                        <Typography
                                            variant={"body2"}
                                            sx={{ padding: "8px" }}
                                        >
                                            {column.id === "name" ? (
                                                <u>{row[column.id]}</u>
                                            ) : (
                                                row[column.id]
                                            )}
                                        </Typography>
                                    </TableCell>
                                )
                            )}
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
};

export default SimpleDataTable;
