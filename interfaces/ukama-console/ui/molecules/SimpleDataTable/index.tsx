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

interface SimpleDataTableInterface {
  dataKey?: string;
  dataset: any;
  height?: string;
  columns: ColumnsWithOptions[];
  networkList?: string[] | [] | undefined;
  handleCreateNetwork?: any;
}

const SimpleDataTable = ({
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
                    <div>
                      {row[column.id].map((role: string) => (
                        <Chip
                          key={role}
                          label={role}
                          sx={{
                            backgroundColor: colors.blueGray,
                            color: colors.black,
                            marginRight: '4px',
                            marginBottom: '4px',
                            borderRadius: '4px',
                          }}
                        />
                      ))}
                    </div>
                  ) : column.id === 'pdf' ? (
                    <Link
                      target="_blank"
                      underline="hover"
                      href={row[column.id]}
                    >
                      View as PDF
                    </Link>
                  ) : column.id === 'network' ? (
                    <ChipDropdown
                      onCreateNetwork={handleCreateNetwork}
                      menu={networkList}
                    />
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
