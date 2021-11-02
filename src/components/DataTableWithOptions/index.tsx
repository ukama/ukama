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
                                {columns.map((column: ColumnsWithOptions) => (
                                    <TableCell
                                        key={column.id}
                                        align={column.align}
                                        sx={{
                                            padding: "14px",
                                            fontSize: "0.875rem",
                                        }}
                                    >
                                        {column.id === "actions" ? (
                                            <OptionsPopover
                                                cid={
                                                    "data-table-action-popover"
                                                }
                                                menuOptions={menuOptions}
                                                handleItemClick={
                                                    onMenuItemClick
                                                }
                                            />
                                        ) : column.id === "name" ? (
                                            <Link href="#" underline="hover">
                                                {row[column.id]}
                                            </Link>
                                        ) : (
                                            row[column.id]
                                        )}
                                    </TableCell>
                                ))}
                            </TableRow>
                        ))}
                    </TableBody>
                </Table>
            </TableContainer>
        </>
    );
};

export default DataTableWithOptions;
