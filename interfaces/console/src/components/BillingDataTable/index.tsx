import Grid from "@mui/material/Grid";
import { RoundedCard } from "../../styles";
import { useState } from "react";
import CheckBoxOutlineBlankIcon from "@mui/icons-material/CheckBoxOutlineBlank";
import { BillingTableHeaderOptionsType, ExportOptionsType } from "../../types";
import { TABLE_EXPORT_OPTIONS } from "../../constants";
import MoreHorizIcon from "@mui/icons-material/MoreHoriz";

import {
    Typography,
    IconButton,
    Select,
    MenuItem,
    Table,
    TableCell,
    TableContainer,
    TableHead,
    TableRow,
    Paper,
} from "@mui/material";
import { colors } from "../../theme";
type BillingTableProps = {
    tableTitle: string;
    isBillingHistory: boolean;
    headerOptions: any[];
    children?: any;
};
const BillingDataTable = ({
    tableTitle,
    isBillingHistory,
    headerOptions,
    children,
}: BillingTableProps) => {
    const [exportOptions, setExportOptions] = useState("EXPORT");
    const handleExportOptions = (event: any) => {
        setExportOptions(event.target.value);
    };

    return (
        <>
            <RoundedCard>
                <Grid container spacing={1}>
                    <Grid item xs={4} container>
                        <Typography variant="h6">{tableTitle}</Typography>
                    </Grid>
                    <Grid item xs={8} container justifyContent="flex-end">
                        <Select
                            value={exportOptions}
                            onChange={handleExportOptions}
                            displayEmpty
                            sx={{
                                color: colors.primaryMain,
                                width: "20%",
                                textAlign: "center",
                            }}
                        >
                            {TABLE_EXPORT_OPTIONS.map(
                                ({ value, label }: ExportOptionsType) => (
                                    <MenuItem key={value} value={value}>
                                        <Typography variant="body1">
                                            {label}
                                        </Typography>
                                    </MenuItem>
                                )
                            )}
                        </Select>
                    </Grid>
                    <Grid container spacing={1}>
                        <Grid item xs={12}>
                            <TableContainer component={Paper} elevation={0}>
                                <Table
                                    sx={{ minWidth: 400 }}
                                    aria-label="spanning table"
                                >
                                    <TableHead>
                                        <TableRow>
                                            {isBillingHistory && (
                                                <TableCell align="left">
                                                    <IconButton
                                                        aria-label="CheckBoxOutlineBlankIcon"
                                                        size="small"
                                                    >
                                                        <CheckBoxOutlineBlankIcon />
                                                    </IconButton>
                                                </TableCell>
                                            )}

                                            {headerOptions.map(
                                                ({
                                                    id,
                                                    label,
                                                }: BillingTableHeaderOptionsType) => (
                                                    <TableCell
                                                        align="left"
                                                        key={id}
                                                    >
                                                        <Typography variant="h3"></Typography>
                                                        {label}
                                                    </TableCell>
                                                )
                                            )}
                                            {isBillingHistory && (
                                                <TableCell align="right">
                                                    <IconButton
                                                        aria-label="CheckBoxOutlineBlankIcon"
                                                        size="small"
                                                    >
                                                        <MoreHorizIcon />
                                                    </IconButton>
                                                </TableCell>
                                            )}
                                        </TableRow>
                                    </TableHead>
                                    {children}
                                </Table>
                            </TableContainer>
                        </Grid>
                    </Grid>
                </Grid>
            </RoundedCard>
        </>
    );
};
export default BillingDataTable;
