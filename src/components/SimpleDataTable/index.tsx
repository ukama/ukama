import {
    Table,
    TableRow,
    TableBody,
    TableCell,
    TableHead,
    TableContainer,
    Typography,
} from "@mui/material";
import { colors } from "../../theme";
import { ColumnsWithOptions } from "../../types";

interface SimpleDataTableInterface {
    dataset: Object[];
    columns: ColumnsWithOptions[];
}

const SimpleDataTable = ({ columns, dataset }: SimpleDataTableInterface) => {
    return (
        <TableContainer sx={{ overflowX: "hidden", mt: "24px" }}>
            <Table stickyHeader>
                <TableHead>
                    <TableRow>
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
                            role="row"
                            key={row.name}
                            sx={{
                                ":hover": {
                                    backgroundColor: colors.solitude,
                                },
                            }}
                        >
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
