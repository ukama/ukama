/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Invitation_Status, NodeConnectivityEnum, Role_Type } from '@/client/graphql/generated';
import ChipDropdown from '@/components/ui/ChipDropDown';
import colors from '@/theme/colors';
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
  TableSortLabel,
  Typography,
} from '@mui/material';
import React, { useMemo, useState } from 'react';

interface SimpleDataTableInterface {
  dataKey?: string;
  dataset: Record<string, string | number | boolean | null | undefined>[];
  height?: string;
  isIdHyperlink?: boolean;
  columns: ColumnsWithOptions[];
  networkList?: Array<{ id: string; name: string }>;
  handleCreateNetwork?: () => void;
  handleDeleteElement?: (id: string) => void;
  showActionButton?: boolean;
  hyperlinkPrefix?: string;
}

interface TableCellProps {
  row: Record<string, string | number | boolean | null | undefined>;
  isIdHyperlink?: boolean;
  handleCreateNetwork?: () => void;
  hyperlinkPrefix?: string;
  column: ColumnsWithOptions;
  networkList: Array<{ id: string; name: string }> | undefined;
  handleDeleteElement?: (id: string) => void;
}

interface TableHeaderProps {
  columns: ColumnsWithOptions[];
  order: 'asc' | 'desc';
  sortedColumn: string;
  onSort: (column: ColumnsWithOptions) => void;
}

const MemoizedTableHeader = React.memo(
  ({ columns, order, sortedColumn, onSort }: TableHeaderProps) => {
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
              {column?.options?.isSortable ? (
                <TableSortLabel
                  active={sortedColumn === column.id}
                  direction={order}
                  onClick={() => onSort(column)}
                >
                  {column.label}
                </TableSortLabel>
              ) : (
                column.label
              )}
            </TableCell>
          ))}
        </TableRow>
      </TableHead>
    );
  },
);
MemoizedTableHeader.displayName = 'MemoizedTableHeader';

const renderCellContent = ({
  column,
  row,
  hyperlinkPrefix,
  isIdHyperlink,
  handleCreateNetwork,
  handleDeleteElement,
  networkList,
}: TableCellProps) => {
  const val = row[column.id];
  const strVal = String(val ?? '');

  const handleDeleteRow = () => {
    handleDeleteElement?.(String(row.id ?? ''));
  };
  switch (column.id) {
    case 'id':
      return isIdHyperlink ? (
        <Link href={`${hyperlinkPrefix}${strVal}`} unselectable="on">
          {strVal}
        </Link>
      ) : (
        <Typography variant="body2">{strVal}</Typography>
      );
    case 'iccid':
      return isIdHyperlink && row['isAllocated'] ? (
        <Link
          href={`${hyperlinkPrefix}iccid=${strVal}`}
          unselectable="on"
        >
          {strVal}
        </Link>
      ) : (
        <Typography variant="body2">{strVal}</Typography>
      );
    case 'role':
      return (
        <div>
          <Chip
            color="info"
            sx={{ color: 'white' }}
            label={roleEnumToString(strVal as Role_Type)}
          />
        </div>
      );
    case 'pdf':
      return (
        <Link target="_blank" underline="hover" href={strVal}>
          View as PDF
        </Link>
      );
    case 'network':
      return (
        <ChipDropdown
          onCreateNetwork={handleCreateNetwork}
          menu={
            (networkList && networkList.map((network) => network.name)) ??
            []
          }
        />
      );
    case 'edit':
      return (
        <IconButton onClick={() => {}}>
          <EditIcon />
        </IconButton>
      );
    case 'delete':
      return (
        <IconButton onClick={handleDeleteRow}>
          <DeleteIcon />
        </IconButton>
      );
    case 'status':
      return (
        <Chip
          sx={{
            p: 1,
            color: 'black',
            backgroundColor: provideStatusColor(strVal as Invitation_Status),
          }}
          label={inviteStatusEnumToString(strVal as Invitation_Status)}
        />
      );
    case 'simType':
      return (
        <Chip
          label={getSimValuefromSimType(strVal)}
          sx={{ color: 'white' }}
          color="info"
        />
      );
    case 'isPhysical':
      return (
        <Typography variant="body2" sx={{ padding: '8px' }}>
          {val === 'true' ? 'pSIM' : 'eSIM'}
        </Typography>
      );
    case 'connectivity':
      return (
        <Chip
          sx={{
            p: 1,
            color: (theme) => theme.palette.text.primary,
            backgroundColor:
              strVal === NodeConnectivityEnum.Online
                ? colors.primaryLight
                : colors.dullRed,
          }}
          label={strVal}
        />
      );
    case 'state':
      return (
        <Chip
          sx={{
            p: 1,
            color: (theme) => theme.palette.text.primary,
            backgroundColor: colors.secondaryLight,
          }}
          label={strVal}
        />
      );
    case 'isAllocated':
      return (
        <Typography variant="body2" sx={{ padding: '8px' }}>
          {val === true ? 'Assigned' : 'Unassigned'}
        </Typography>
      );
    default:
      return (
        <Typography variant="body2" sx={{ padding: '8px' }}>
          {strVal}
        </Typography>
      );
  }
};

const SimpleTableCell = (tprops: TableCellProps) => (
  <TableCell
    sx={{
      padding: 1,
      fontSize: '0.875rem',
    }}
  >
    {renderCellContent(tprops)}
  </TableCell>
);

const SimpleDataTable = React.memo(
  ({
    dataKey = 'id',
    columns,
    dataset,
    height,
    hyperlinkPrefix = '/',
    isIdHyperlink = false,
    showActionButton = false,
    networkList,
    handleCreateNetwork,
    handleDeleteElement,
  }: SimpleDataTableInterface) => {
    const [order, setOrder] = useState<'asc' | 'desc'>('asc');
    const [sortedColumn, setSortedColumn] = useState<string>('');

    const handleSort = (column: ColumnsWithOptions) => {
      const isSameColumn = sortedColumn === column.id;
      const nextOrder = isSameColumn && order === 'asc' ? 'desc' : 'asc';
      setOrder(nextOrder);
      setSortedColumn(column.id);
    };

    const sortedData = useMemo(() => {
      if (!sortedColumn) return dataset;
      return [...dataset].sort((a, b) => {
        const aVal = a[sortedColumn];
        const bVal = b[sortedColumn];
        if (aVal == null || bVal == null) return 0;
        if (typeof aVal === 'number' && typeof bVal === 'number') {
          return order === 'asc' ? aVal - bVal : bVal - aVal;
        }
        const aString = String(aVal);
        const bString = String(bVal);
        return order === 'asc'
          ? aString.localeCompare(bString)
          : bString.localeCompare(aString);
      });
    }, [dataset, sortedColumn, order]);

    const [menuAnchorEl, setMenuAnchorEl] = useState<null | HTMLElement>(null);
    const handleMenuOpen = (event: React.MouseEvent<HTMLButtonElement>) => {
      setMenuAnchorEl(event.currentTarget);
    };

    return (
      <TableContainer
        sx={{
          mt: '16px',
          maxHeight: height ? height : '100%',
          overflow: 'auto',
          display: 'block',
        }}
      >
        <Table stickyHeader sx={{ width: '100%' }}>
          <MemoizedTableHeader
            columns={columns}
            order={order}
            sortedColumn={sortedColumn}
            onSort={handleSort}
          />
          <TableBody>
            {sortedData?.map((row) => (
              <TableRow key={String(row[dataKey] ?? '')}>
                {columns?.map((column: ColumnsWithOptions, index: number) => (
                  <SimpleTableCell
                    row={row}
                    column={column}
                    networkList={networkList}
                    isIdHyperlink={isIdHyperlink}
                    hyperlinkPrefix={hyperlinkPrefix}
                    key={`$cell-${index}-${column.id}`}
                    handleCreateNetwork={handleCreateNetwork}
                    handleDeleteElement={handleDeleteElement}
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
