import { isDarkmode } from '@/app-recoil';
import { ColumnsWithOptions } from '@/types';
import ChipDropdown from '@/ui/components/ChipDropDown';
import {
  Chip,
  Link,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import { useRecoilValue } from 'recoil';
import { colors } from '@/styles/theme';
import React from 'react';
import { NetworkDto } from '@/generated';

interface SimpleDataTableInterface {
  dataKey?: string;
  dataset: any;
  height?: string;
  columns: ColumnsWithOptions[];
  networkList?: any;
  handleCreateNetwork?: any;
}

// eslint-disable-next-line react/display-name
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

// eslint-disable-next-line react/display-name
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

// eslint-disable-next-line react/display-name
const MemoizedTableCell = React.memo(
  ({
    column,
    row,
    handleCreateNetwork,
    networkList,
  }: {
    column: ColumnsWithOptions;
    row: any;
    handleCreateNetwork: any;
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
            {row[column.id].map((role: string) => (
              <MemoizedChip key={role} label={role} />
            ))}
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

const SimpleDataTable = React.memo(
  ({
    dataKey = 'id',
    columns,
    dataset,
    height,
    networkList,
    handleCreateNetwork,
  }: SimpleDataTableInterface) => {
    const _isDarkMode = useRecoilValue(isDarkmode);
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
                    networkList={networkList}
                  />
                ))}
              </TableRow>
            ))}
          </TableBody>
        </Table>
      </TableContainer>
    );
  },
);
SimpleDataTable.displayName = 'SimpleDataTable';

export default SimpleDataTable;
