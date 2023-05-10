import { isDarkmode } from '@/app-recoil';
import { colors } from '@/styles/theme';
import { ColumnsWithOptions } from '@/types';
import {
  Button,
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
  maxHeight?: number;
  totalRows?: number;
  setSelectedRows?: any;
  selectedRows?: number[];
  rowSelection?: boolean;
  columns: ColumnsWithOptions[];
  isHistoryTab?: boolean;
  handleViewPdf?: any;
  totalAmount?: number;
}

const SimpleDataTable = ({
  columns,
  dataset,
  maxHeight,
  isHistoryTab = false,
  totalAmount,
  handleViewPdf,
}: SimpleDataTableInterface) => {
  const _isDarkMode = useRecoilValue(isDarkmode);
  return (
    <TableContainer
      sx={{
        mt: '24px',
        maxHeight: maxHeight ? maxHeight : '100%',
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
          {isHistoryTab && (
            <TableCell
              style={{
                padding: '0px 16px 12px 16px',
                fontStyle: '600',
                fontWeight: '600',
              }}
            >
              Invoice
            </TableCell>
          )}
        </TableRow>

        <TableBody>
          {dataset?.map((row: any) => (
            <TableRow
              key={row.id}
              sx={{
                '&:last-child th, &:last-child td': {
                  borderBottom: isHistoryTab ? 0 : null,
                },
                ':hover': {
                  backgroundColor: _isDarkMode
                    ? colors.nightGrey
                    : colors.hoverColor08,
                },
              }}
            >
              {columns?.map((column: ColumnsWithOptions, index: number) => (
                <TableCell
                  key={`${row.date}-${index}`}
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
              {isHistoryTab && (
                <TableCell
                  sx={{
                    '&:last-child th, &:last-child td': {
                      border: 0,
                    },
                  }}
                >
                  <a
                    href={'https://docdro.id/J2v6TJO'}
                    target="_blank"
                    rel="noreferrer"
                    style={{ textDecoration: 'none' }}
                  >
                    <Button
                      variant="text"
                      sx={{
                        color: colors.primaryMain,
                        textTransform: 'capitalize',
                      }}
                      onClick={handleViewPdf}
                    >
                      <Typography variant="body2">View as PDF</Typography>
                    </Button>
                  </a>
                </TableCell>
              )}
            </TableRow>
          ))}
          {isHistoryTab !== true && (
            <TableRow
              sx={{
                '&:last-child th, &:last-child td': {
                  borderBottom: 0,
                },
              }}
            >
              <TableCell colSpan={2} />
              <TableCell>
                <b>
                  {'$'}
                  {totalAmount}
                </b>
              </TableCell>
            </TableRow>
          )}
        </TableBody>
      </Table>
    </TableContainer>
  );
};

export default SimpleDataTable;
