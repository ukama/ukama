import { ColumnsWithOptions } from '@/types';
import {
  Skeleton,
  Table,
  TableBody,
  TableCell,
  TableHead,
  TableRow,
} from '@mui/material';

interface DataTableSkeltonProps {
  header: ColumnsWithOptions[];
}

const DataTableSkelton = ({ header }: DataTableSkeltonProps) => {
  return (
    <Table>
      <TableHead>
        <TableRow>
          {header.map((column) => (
            <TableCell
              key={column.id}
              sx={{
                minWidth: column.minWidth,
                padding: '6px 12px 12px 0px',
              }}
            >
              <b> {column.label}</b>
            </TableCell>
          ))}
        </TableRow>
      </TableHead>
      <TableBody>
        <TableRow>
          <TableCell>
            <Skeleton />
          </TableCell>
          <TableCell>
            <Skeleton />
          </TableCell>
          <TableCell>
            <Skeleton />
          </TableCell>
          <TableCell>
            <Skeleton />
          </TableCell>
          <TableCell>
            <Skeleton />
          </TableCell>
        </TableRow>
      </TableBody>
    </Table>
  );
};

export default DataTableSkelton;
