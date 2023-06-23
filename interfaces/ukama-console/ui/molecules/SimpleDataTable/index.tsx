import { isDarkmode } from '@/app-recoil';
import { ColumnsWithOptions } from '@/types';
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

interface SimpleDataTableInterface {
  dataKey?: string;
  dataset: any;
  height?: string;
  columns: ColumnsWithOptions[];
}

const SimpleDataTable = ({
  dataKey = 'id',
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

        <TableBody>
          {dataset?.map((row: any) => (
            <TableRow key={row[dataKey]} sx={{}}>
              {columns?.map((column: ColumnsWithOptions, index: number) => (
                <TableCell
                  key={`$cell-${index}-${column.id}`}
                  sx={{
                    padding: 1,
                    fontSize: '0.875rem',
                  }}
                >
                  {column.id === 'role' ? (
                    <Chip label={row[column.id]} variant="outlined" />
                  ) : column.id === 'pdf' ? (
                    <Link
                      target="_blank"
                      underline="hover"
                      href={row[column.id]}
                    >
                      View as PDF
                    </Link>
                  ) : (
                    <Typography variant={'body2'} sx={{ padding: '8px' }}>
                      {row[column.id]}
                    </Typography>
                  )}
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
