import { isDarkmode } from '@/app-recoil';
import { ColumnsWithOptions } from '@/types';
import {
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableRow,
  Typography,
} from '@mui/material';
import { useRecoilValue } from 'recoil';

interface SimpleDataTableInterface {
  dataset: any;
  height?: string;
  columns: ColumnsWithOptions[];
}

const SimpleDataTable = ({
  columns,
  dataset,
  height,
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
        <TableRow>
          {columns?.map((column) => (
            <TableCell
              key={column.id}
              align={column.align}
              style={{
                padding: '0px 16px 12px 16px',
                fontSize: '0.875rem',
                minWidth: column.minWidth,
              }}
            >
              <b>{column.label}</b>
            </TableCell>
          ))}
        </TableRow>

        <TableBody>
          {dataset?.map((row: any) => (
            <TableRow key={row.id} sx={{}}>
              {columns?.map((column: ColumnsWithOptions, index: number) => (
                <TableCell
                  key={`$cell-${index}-${column.id}`}
                  sx={{
                    padding: 1,
                    fontSize: '0.875rem',
                  }}
                >
                  <Typography variant={'body2'} sx={{ padding: '8px' }}>
                    {row[column.id]}
                  </Typography>
                </TableCell>
              ))}
            </TableRow>
          ))}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default SimpleDataTable;
