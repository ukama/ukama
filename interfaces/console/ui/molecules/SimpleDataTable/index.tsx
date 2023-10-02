import { isDarkmode } from '@/app-recoil';
import { ColumnsWithOptions } from '@/types';
import ChipDropdown from '@/ui/molecules/ChipDropDown';
import {
  Chip,
  Link,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  Menu,
  TableHead,
  MenuItem,
  TableRow,
  Typography,
  IconButton,
} from '@mui/material';
import { useRecoilValue } from 'recoil';
import { colors } from '@/styles/theme';
import React, { useState } from 'react';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import EditIcon from '@mui/icons-material/Edit';
import DeleteIcon from '@mui/icons-material/Delete';

interface SimpleDataTableInterface {
  dataKey?: string;
  dataset: any;
  height?: string;
  columns: ColumnsWithOptions[];
  networkList?: any;
  handleCreateNetwork?: any;
  handleDeleteElement?: any;
  showActionButton?: boolean;
}

const MemoizedTableHeader = React.memo(
  ({ columns }: { columns: ColumnsWithOptions[] }) => {
    return (
      <TableHead>
        <TableRow>
          {columns?.map((column) => (
            <TableCell
              key={`row-${column.id}`}
              align={column.align}
              style={{
                fontWeight: 600,
                padding: '0px 16px 12px 16px',
                fontSize: '0.875rem',
                minWidth: column.minWidth,
              }}
            >
              {column.label}
            </TableCell>
          ))}
        </TableRow>
      </TableHead>
    );
  },
);
MemoizedTableHeader.displayName = 'MemoizedTableHeader';

const MemoizedChip = React.memo(({ label }: { label: string }) => {
  return (
    <Chip
      label={label}
      sx={{
        backgroundColor: colors.blueGray,
        color: colors.black,
        marginRight: '4px',
        marginBottom: '4px',
        borderRadius: '4px',
      }}
    />
  );
});
MemoizedChip.displayName = 'MemoizedChip';

const MemoizedTableCell = React.memo(
  ({
    column,
    row,
    handleCreateNetwork,
    handleDeleteElement,
    networkList,
  }: {
    column: ColumnsWithOptions;
    row: any;
    handleCreateNetwork: any;
    handleDeleteElement: any;
    networkList: string[] | [] | undefined;
  }) => {
    return (
      <TableCell
        sx={{
          padding: 1,
          fontSize: '0.875rem',
        }}
      >
        {column.id === 'role' ? (
          <div>
            <MemoizedChip label={row[column.id]} />
          </div>
        ) : column.id === 'pdf' ? (
          <Link target="_blank" underline="hover" href={row[column.id]}>
            View as PDF
          </Link>
        ) : column.id === 'network' ? (
          <ChipDropdown
            onCreateNetwork={handleCreateNetwork}
            menu={
              (networkList &&
                networkList.map((network: any) => network.name)) ||
              []
            }
          />
        ) : column.id === 'edit' ? (
          <IconButton onClick={() => console.log(row)}>
            <EditIcon />
          </IconButton>
        ) : column.id === 'delete' ? (
          <IconButton onClick={() => handleDeleteElement(row.id)}>
            <DeleteIcon />
          </IconButton>
        ) : (
          <Typography
            variant={'body2'}
            sx={{ padding: '8px' }}
            color={row[column.id] === 'false' ? 'primary' : ''}
          >
            {row[column.id] === 'true'
              ? 'Assigned'
              : row[column.id] === 'false'
              ? 'Unassigned'
              : row[column.id]}
          </Typography>
        )}
      </TableCell>
    );
  },
);
MemoizedTableCell.displayName = 'MemoizedTableCell';

const SimpleDataTable = React.memo(
  ({
    dataKey = 'id',
    columns,
    dataset,
    height,
    showActionButton = false,
    networkList,
    handleCreateNetwork,
    handleDeleteElement,
  }: SimpleDataTableInterface) => {
    const _isDarkMode = useRecoilValue(isDarkmode);
    const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);

    const handleMenuOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
      setMenuAnchorEl(event.currentTarget);
    };

    const handleMenuClose = () => {
      setMenuAnchorEl(null);
    };

    return (
      <TableContainer
        sx={{
          mt: '24px',
          height: height ? height : '100%',
        }}
      >
        <Table stickyHeader>
          <MemoizedTableHeader columns={columns} />

          <TableBody>
            {dataset?.map((row: any) => (
              <TableRow key={row[dataKey]} sx={{}}>
                {columns?.map((column: ColumnsWithOptions, index: number) => (
                  <MemoizedTableCell
                    key={`$cell-${index}-${column.id}`}
                    column={column}
                    row={row}
                    handleCreateNetwork={handleCreateNetwork}
                    handleDeleteElement={handleDeleteElement}
                    networkList={networkList}
                  />
                ))}
              </TableRow>
            ))}

            {showActionButton && (
              <TableRow>
                <TableCell colSpan={columns.length}>
                  <IconButton
                    aria-label="menu"
                    aria-controls="simple-menu"
                    aria-haspopup="true"
                    onClick={handleMenuOpen}
                  >
                    <MoreVertIcon />
                  </IconButton>
                  <Menu
                    id="simple-menu"
                    anchorEl={menuAnchorEl}
                    keepMounted
                    open={Boolean(menuAnchorEl)}
                    onClose={handleMenuClose}
                  >
                    <MenuItem onClick={() => alert('Resend email')}>
                      Resend Email
                    </MenuItem>
                    <MenuItem onClick={() => alert('Delete')}>Delete</MenuItem>
                  </Menu>
                </TableCell>
              </TableRow>
            )}
          </TableBody>
        </Table>
      </TableContainer>
    );
  },
);
SimpleDataTable.displayName = 'SimpleDataTable';

export default SimpleDataTable;