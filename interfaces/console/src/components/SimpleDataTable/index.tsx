/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import ChipDropdown from '@/components/ChipDropDown';
import { useAppContext } from '@/context';
import { colors } from '@/theme';
import { ColumnsWithOptions } from '@/types';
import {
  getSimValuefromSimType,
  inviteStatusEnumToString,
  provideStatusColor,
  roleEnumToString,
} from '@/utils';
import DeleteIcon from '@mui/icons-material/Delete';
import EditIcon from '@mui/icons-material/Edit';
import MoreVertIcon from '@mui/icons-material/MoreVert';
import {
  Chip,
  IconButton,
  Link,
  Menu,
  MenuItem,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import React, { useState } from 'react';

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

const SimpleTableCell = ({
  column,
  row,
  handleCreateNetwork,
  handleDeleteElement,
  networkList,
}: {
  column: ColumnsWithOptions;
  row: any;
  handleCreateNetwork: any;
  handleDeleteElement: (id: string) => void;
  networkList: string[] | [] | undefined;
}) => {
  const handleDeleteRow = () => {
    handleDeleteElement(row.id);
  };
  return (
    <TableCell
      sx={{
        padding: 1,
        fontSize: '0.875rem',
      }}
    >
      {column.id === 'role' ? (
        <div>
          <MemoizedChip label={roleEnumToString(row[column.id])} />
        </div>
      ) : column.id === 'pdf' ? (
        <Link target="_blank" underline="hover" href={row[column.id]}>
          View as PDF
        </Link>
      ) : column.id === 'network' ? (
        <ChipDropdown
          onCreateNetwork={handleCreateNetwork}
          menu={
            (networkList && networkList.map((network: any) => network.name)) ??
            []
          }
        />
      ) : column.id === 'edit' ? (
        <IconButton onClick={() => {}}>
          <EditIcon />
        </IconButton>
      ) : column.id === 'delete' ? (
        <IconButton onClick={handleDeleteRow}>
          <DeleteIcon />
        </IconButton>
      ) : column.id === 'status' ? (
        <Chip
          sx={{ color: 'white' }}
          label={inviteStatusEnumToString(row[column.id])}
          color={provideStatusColor(row[column.id])}
        />
      ) : column.id === 'simType' ? (
        <Chip
          label={getSimValuefromSimType(row[column.id])}
          sx={{ color: 'white' }}
          color={'info'}
        />
      ) : column.id === 'isPhysical' ? (
        <Typography variant={'body2'} sx={{ padding: '8px' }}>
          {row[column.id] === 'true' ? 'Yes' : 'No'}
        </Typography>
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
};

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
    const { isDarkMode } = useAppContext();
    const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);

    const handleMenuOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
      setMenuAnchorEl(event.currentTarget);
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
                  <SimpleTableCell
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
                    onClose={() => setMenuAnchorEl(null)}
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
