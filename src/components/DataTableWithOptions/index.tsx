import {
    TableContainer,
    Table,
    TableHead,
    TableRow,
    TableCell,
    TableBody,
    TablePagination,
} from "@mui/material";
import React from "react";
import { ColumnsWithOptions } from "../../types";

const columns: readonly ColumnsWithOptions[] = [
    { id: "name", label: "Name", minWidth: 170 },
    { id: "usage", label: "Usage", minWidth: 100 },
    // {
    //     id: "actions",
    //     label: "...",
    //     minWidth: 74,
    //     align: "right",
    // },
];

interface Data {
    name: string;
    usage: string;
}

function createData(name: string, usage: string): Data {
    return { name, usage };
}

const rows = [
    createData("India", "IN"),
    createData("China", "CN"),
    createData("Italy", "IT"),
    createData("United States", "US"),
    createData("Canada", "CA"),
    createData("Australia", "AU"),
    createData("Germany", "DE"),
    createData("Ireland", "IE"),
    createData("Mexico", "MX"),
    createData("Japan", "JP"),
    createData("France", "FR"),
    createData("United Kingdom", "GB"),
    createData("Russia", "RU"),
    createData("Nigeria", "NG"),
    createData("Brazil", "BR"),
];

const DataTableWithOptions = () => {
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(10);

    const handleChangePage = (event: unknown, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (
        event: React.ChangeEvent<HTMLInputElement>
    ) => {
        setRowsPerPage(+event.target.value);
        setPage(0);
    };

    return (
        <>
            <TableContainer sx={{ maxHeight: 440 }}>
                <Table stickyHeader aria-label="sticky table">
                    <TableHead>
                        <TableRow>
                            {columns.map(column => (
                                <TableCell
                                    key={column.id}
                                    align={column.align}
                                    style={{ minWidth: column.minWidth }}
                                >
                                    {column.label}
                                </TableCell>
                            ))}
                        </TableRow>
                    </TableHead>
                    <TableBody>
                        {rows
                            .slice(
                                page * rowsPerPage,
                                page * rowsPerPage + rowsPerPage
                            )
                            .map((row: Data) => {
                                return (
                                    <TableRow
                                        hover
                                        role="checkbox"
                                        tabIndex={-1}
                                        key={row.name}
                                    >
                                        {columns.map(
                                            (column: ColumnsWithOptions) => {
                                                const value = row[column.id];
                                                return (
                                                    <TableCell
                                                        key={column.id}
                                                        align={column.align}
                                                    >
                                                        {column.format &&
                                                        typeof value ===
                                                            "number"
                                                            ? column.format(
                                                                  value
                                                              )
                                                            : value}
                                                    </TableCell>
                                                );
                                            }
                                        )}
                                    </TableRow>
                                );
                            })}
                    </TableBody>
                </Table>
            </TableContainer>
            <TablePagination
                rowsPerPageOptions={[10, 25, 100]}
                component="div"
                count={rows.length}
                rowsPerPage={rowsPerPage}
                page={page}
                onPageChange={handleChangePage}
                onRowsPerPageChange={handleChangeRowsPerPage}
            />
        </>
    );
};

export default DataTableWithOptions;
