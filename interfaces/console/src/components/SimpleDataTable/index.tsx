import {
    Table,
    TableRow,
    Checkbox,
    TableBody,
    TableCell,
    TableHead,
    Typography,
    Button,
    TableContainer,
} from "@mui/material";
import { format } from "date-fns";
import { useRecoilValue } from "recoil";
import { isDarkmode } from "../../recoil";
import { colors } from "../../theme";
import { ColumnsWithOptions } from "../../types";

interface SimpleDataTableInterface {
    dataset: Object[];
    maxHeight?: number;
    totalRows?: number;
    setSelectedRows?: any;
    selectedRows?: number[];
    rowSelection?: boolean;
    columns: ColumnsWithOptions[];
    isHistoryTab?: boolean;
    handleViewPdf?: any;
}

const SimpleDataTable = ({
    columns,
    dataset,
    maxHeight,
    totalRows = 0,
    isHistoryTab = false,
    setSelectedRows,
    handleViewPdf,
    selectedRows = [],
    rowSelection = false,
}: SimpleDataTableInterface) => {
    const _isDarkMode = useRecoilValue(isDarkmode);
    const onRowSelection = (id: number) => {
        setSelectedRows && setSelectedRows([...selectedRows, id]);
    };

    const onRowsSelection = () => {
        if (selectedRows.length === totalRows) setSelectedRows([]);
        else setSelectedRows(dataset.map((i: any) => i?.id));
    };

    return (
        <TableContainer
            sx={{
                mt: "24px",
                maxHeight: maxHeight ? maxHeight : "100%",
            }}
        >
            <Table stickyHeader>
                <TableHead>
                    <TableRow>
                        {columns?.map(column => (
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
                        {isHistoryTab && (
                            <TableCell
                                style={{
                                    padding: "0px 16px 12px 16px",
                                    fontStyle: "600",
                                    fontWeight: "600",
                                }}
                            >
                                Invoice
                            </TableCell>
                        )}
                    </TableRow>
                </TableHead>
                <TableBody>
                    {dataset?.map((row: any) => (
                        <TableRow
                            key={row.id}
                            sx={{
                                ":hover": {
                                    backgroundColor: _isDarkMode
                                        ? colors.nightGrey
                                        : colors.hoverColor08,
                                },
                            }}
                            selected={selectedRows.includes(row.id)}
                            role={rowSelection ? "checkbox" : "row"}
                            onClick={() => onRowSelection(row.id)}
                        >
                            {columns?.map(
                                (column: ColumnsWithOptions, index: number) => (
                                    <TableCell
                                        key={`${row.date}-${index}`}
                                        sx={{
                                            padding: 1,
                                            fontSize: "0.875rem",
                                        }}
                                    >
                                        <Typography
                                            variant={"body2"}
                                            sx={{ padding: "8px" }}
                                        >
                                            {row[column.id]}
                                        </Typography>
                                    </TableCell>
                                )
                            )}
                            {isHistoryTab && (
                                <TableCell>
                                    <Button
                                        variant="text"
                                        sx={{
                                            color: colors.primaryMain,
                                            textTransform: "capitalize",
                                        }}
                                        onClick={handleViewPdf}
                                    >
                                        View as PDF
                                    </Button>
                                </TableCell>
                            )}
                        </TableRow>
                    ))}
                </TableBody>
            </Table>
        </TableContainer>
    );
};

export default SimpleDataTable;
