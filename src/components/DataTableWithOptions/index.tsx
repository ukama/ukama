import {
    Box,
    Link,
    Table,
    TableRow,
    TableBody,
    TableCell,
    TableHead,
    TableContainer,
} from "@mui/material";
import { EmptyView } from "..";
import OptionsPopover from "../OptionsPopover";
import UserIcon from "@mui/icons-material/Person";
import { ColumnsWithOptions, MenuItemType } from "../../types";

interface DataTableWithOptionsInterface {
    dataset: any;
    onMenuItemClick: Function;
    menuOptions: MenuItemType[];
    columns: ColumnsWithOptions[];
}

type CellValueByTypeProps = {
    type: string;
    row: any;
    menuOptions: MenuItemType[];
    onMenuItemClick: Function;
};

const CellValueByType = ({
    type,
    row,
    menuOptions,
    onMenuItemClick,
}: CellValueByTypeProps) => {
    switch (type) {
        case "name":
            return (
                <Link href="#" underline="hover">
                    {row[type]}
                </Link>
            );
        case "actions":
            return (
                <OptionsPopover
                    cid={"data-table-action-popover"}
                    menuOptions={menuOptions}
                    handleItemClick={onMenuItemClick}
                />
            );
        default:
            return <>{`${row[type]} MB`}</>;
    }
};

const DataTableWithOptions = ({
    columns,
    dataset,
    menuOptions,
    onMenuItemClick,
}: DataTableWithOptionsInterface) => {
    return (
        <Box
            component="div"
            mt={2}
            sx={{
                height: "100%",
                minHeight: 230,
                display: "flex",
                alignItems: "center",
            }}
        >
            {dataset?.length > 0 ? (
                <TableContainer sx={{ maxHeight: 230 }}>
                    <Table stickyHeader>
                        <TableHead>
                            <TableRow>
                                {columns?.map(column => (
                                    <TableCell
                                        key={column.id}
                                        align={column.align}
                                        style={{
                                            fontSize: "0.875rem",
                                            minWidth: column.minWidth,
                                            padding: "6px 12px 12px 0px",
                                        }}
                                    >
                                        <b>{column.label}</b>
                                    </TableCell>
                                ))}
                            </TableRow>
                        </TableHead>
                        <TableBody>
                            {dataset?.map((row: any) => (
                                <TableRow role="row" tabIndex={-1} key={row.id}>
                                    {columns.map(
                                        (
                                            column: ColumnsWithOptions,
                                            index: number
                                        ) => (
                                            <TableCell
                                                key={`${row.name}-${index}`}
                                                align={column.align}
                                                sx={{
                                                    padding:
                                                        "13px 12px 13px 0px",
                                                    fontSize: "0.875rem",
                                                }}
                                            >
                                                <CellValueByType
                                                    row={row}
                                                    type={column.id}
                                                    menuOptions={menuOptions}
                                                    onMenuItemClick={(
                                                        type: string
                                                    ) =>
                                                        onMenuItemClick(
                                                            row.id,
                                                            type
                                                        )
                                                    }
                                                />
                                            </TableCell>
                                        )
                                    )}
                                </TableRow>
                            ))}
                        </TableBody>
                    </Table>
                </TableContainer>
            ) : (
                <EmptyView
                    size="large"
                    title="No residents yet!"
                    icon={UserIcon}
                />
            )}
        </Box>
    );
};

export default DataTableWithOptions;
