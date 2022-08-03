import {
    Table,
    TableRow,
    TableBody,
    TableCell,
    TableHead,
    Typography,
    Button,
    Stack,
    TableContainer,
} from "@mui/material";
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

    return (
        <>
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
                                    "&:last-child th, &:last-child td": {
                                        borderBottom: isHistoryTab ? 0 : null,
                                    },
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
                                    (
                                        column: ColumnsWithOptions,
                                        index: number
                                    ) => (
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
                                        <a
                                            href={"https://docdro.id/J2v6TJO"}
                                            target="_blank"
                                            rel="noreferrer"
                                            style={{ textDecoration: "none" }}
                                        >
                                            <Button
                                                variant="text"
                                                sx={{
                                                    color: colors.primaryMain,
                                                }}
                                                onClick={handleViewPdf}
                                            >
                                                <Stack
                                                    direction="row"
                                                    spacing={1}
                                                >
                                                    <Typography
                                                        variant="body2"
                                                        sx={{
                                                            textTransform:
                                                                "capitalize",
                                                        }}
                                                    >
                                                        {` View `}
                                                    </Typography>
                                                    <Typography
                                                        variant="body2"
                                                        sx={{
                                                            textTransform:
                                                                "lowercase",
                                                        }}
                                                    >
                                                        {` as `}
                                                    </Typography>
                                                    <Typography variant="body2">
                                                        {` PDF`}
                                                    </Typography>
                                                </Stack>
                                            </Button>
                                        </a>
                                    </TableCell>
                                )}
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </>
    );
};

export default SimpleDataTable;
