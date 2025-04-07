/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import ChipDropdown from '@/components/ChipDropDown';
import { useAppContext } from '@/context';
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
  Typography,
} from '@mui/material';
import React, { useState } from 'react';

interface SimpleDataTableInterface {
  dataKey?: string;
  dataset: any;
  height?: string;
  isIdHyperlink?: boolean;
  columns: ColumnsWithOptions[];
  networkList?: any;
  handleCreateNetwork?: any;
  handleDeleteElement?: any;
  showActionButton?: boolean;
  hyperlinkPrefix?: string;
}

interface TableCellProps {
  row: any;
  isIdHyperlink?: boolean;
  handleCreateNetwork: any;
  hyperlinkPrefix?: string;
  column: ColumnsWithOptions;
  networkList: string[] | [] | undefined;
  handleDeleteElement: (id: string) => void;
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

const renderCellContent = ({
  column,
  row,
  hyperlinkPrefix,
  isIdHyperlink,
  handleCreateNetwork,
  handleDeleteElement,
  networkList,
}: TableCellProps) => {
  const handleDeleteRow = () => {
    handleDeleteElement(row.id);
  };
  switch (column.id) {
    case 'id':
      return isIdHyperlink ? (
        <Link href={`${hyperlinkPrefix}/${row[column.id]}`} unselectable="on">
          {row[column.id]}
        </Link>
      ) : (
        <Typography variant="body2">{row[column.id]}</Typography>
      );
    case 'role':
      return (
        <div>
          <Chip
            color="info"
            sx={{ color: 'white' }}
            label={roleEnumToString(row[column.id])}
          />
        </div>
      );
    case 'pdf':
      return (
        <Link target="_blank" underline="hover" href={row[column.id]}>
          View as PDF
        </Link>
      );
    case 'network':
      return (
        <ChipDropdown
          onCreateNetwork={handleCreateNetwork}
          menu={
            (networkList && networkList.map((network: any) => network.name)) ??
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
            backgroundColor: provideStatusColor(row[column.id]),
          }}
          label={inviteStatusEnumToString(row[column.id])}
        />
      );
    case 'simType':
      return (
        <Chip
          label={getSimValuefromSimType(row[column.id])}
          sx={{ color: 'white' }}
          color="info"
        />
      );
    case 'isPhysical':
      return (
        <Typography variant="body2" sx={{ padding: '8px' }}>
          {row[column.id] === 'true' ? 'pSIM' : 'eSIM'}
        </Typography>
      );
    case 'connectivity':
      return (
        <Chip
          sx={{
            p: 1,
            color: (theme) => theme.palette.text.primary,
            backgroundColor: colors.primaryLight,
          }}
          label={row[column.id]}
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
          label={row[column.id]}
        />
      );
    case 'isAllocated':
      return (
        <Typography variant="body2" sx={{ padding: '8px' }}>
          {row[column.id] === true ? 'Assigned' : 'Unassigned'}
        </Typography>
      );
    default:
      return (
        <Typography variant="body2" sx={{ padding: '8px' }}>
          {row[column.id]}
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
