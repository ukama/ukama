import {
    Table,
    TableRow,
    TableBody,
    TableCell,
    TableHead,
    TableContainer,
    Link,
} from "@mui/material";
import OptionsPopover from "../OptionsPopover";
import { ColumnsWithOptions, MenuItemType } from "../../types";

interface DataTableWithOptionsInterface {
    columns: ColumnsWithOptions[];
    dataset: Object[];
    menuOptions: MenuItemType[];
    onMenuItemClick: Function;
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
            return <>{row[type]}</>;
    }
};

const DataTableWithOptions = ({
    columns,
    dataset,
    menuOptions,
    onMenuItemClick,
}: DataTableWithOptionsInterface) => {
    return (
        <>
            <TableContainer sx={{ maxHeight: 200 }}>
                <Table stickyHeader>
                    <TableHead>
                        <TableRow>
                            {columns.map(column => (
                                <TableCell
                                    key={column.id}
                                    align={column.align}
                                    style={{
                                        minWidth: column.minWidth,
                                        padding: "12px",
                                        fontSize: "0.875rem",
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
                                role="checkbox"
                                tabIndex={-1}
                                key={row.name}
                            >
                                {columns.map(
                                    (
                                        column: ColumnsWithOptions,
                                        index: number
                                    ) => (
                                        <TableCell
                                            key={`${row.name}-${index}`}
                                            align={column.align}
                                            sx={{
                                                padding: "14px",
                                                fontSize: "0.875rem",
                                            }}
                                        >
                                            <CellValueByType
                                                row={row}
                                                type={column.id}
                                                menuOptions={menuOptions}
                                                onMenuItemClick={
                                                    onMenuItemClick
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
        </>
    );
};

export default DataTableWithOptions;
