import { ColumnsWithOptions, MenuItemType } from '@/types';
import UserIcon from '@mui/icons-material/Person';
import {
  Box,
  Link,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Typography,
} from '@mui/material';
import EmptyView from '../EmptyView';
import OptionsPopover from '../OptionsPopover';

interface DataTableWithOptionsInterface {
  dataset: any;
  emptyViewLabel?: string;
  onMenuItemClick: Function;
  menuOptions: MenuItemType[];
  columns: ColumnsWithOptions[];
}

type CellValueByTypeProps = {
  type: string;
  row: any;
  menuOptions: MenuItemType[];
  onMenuItemClick: Function;
};

const CellValueByType = ({
  type,
  row,
  menuOptions,
  onMenuItemClick,
}: CellValueByTypeProps) => {
  switch (type) {
    case 'name':
      return (
        <Link href="#" underline="hover">
          {row[type]}
        </Link>
      );
    case 'actions':
      return (
        <OptionsPopover
          cid={'data-table-action-popover'}
          menuOptions={menuOptions}
          handleItemClick={onMenuItemClick}
        />
      );
    // case 'dataUsage':
    //   return (
    //     <LoadingWrapper
    //       width="60px"
    //       height="23px"
    //       radius="small"
    //       variant="text"
    //       isLoading={!row.dataPlan}
    //     >
    //       {formatBytesToMB(parseInt(row[type] || '0'))} MB
    //     </LoadingWrapper>
    //   );
    default:
      return <Typography variant="caption">{row[type]}</Typography>;
  }
};

const DataTableWithOptions = ({
  columns,
  dataset,
  menuOptions,
  onMenuItemClick,
  emptyViewLabel = '',
}: DataTableWithOptionsInterface) => {
  return (
    <Box
      component="div"
      sx={{
        width: '100%',
        height: '100%',
        display: 'flex',
      }}
    >
      {dataset?.length > 0 ? (
        <TableContainer>
          <Table stickyHeader>
            <TableHead>
              <TableRow>
                {columns?.map((column) => (
                  <TableCell
                    key={`header-cell-${column.id}`}
                    align={column.align}
                    style={{
                      fontSize: '0.875rem',
                      minWidth: column.minWidth,
                      padding: '6px 12px 12px 0px',
                    }}
                  >
                    <b>{column.label}</b>
                  </TableCell>
                ))}
              </TableRow>
            </TableHead>
            <TableBody>
              {dataset?.map((row: any, id: number) => (
                <TableRow role="row" tabIndex={-1} key={`row-${id}`}>
                  {columns.map((column: ColumnsWithOptions, index: number) => (
                    <TableCell
                      key={`cell-${index}`}
                      align={column.align}
                      sx={{
                        padding: '13px 12px 13px 0px',
                        fontSize: '0.875rem',
                      }}
                    >
                      <CellValueByType
                        row={row}
                        type={column.id}
                        menuOptions={menuOptions}
                        onMenuItemClick={(type: string) =>
                          onMenuItemClick(row.id, type)
                        }
                      />
                    </TableCell>
                  ))}
                </TableRow>
              ))}
            </TableBody>
          </Table>
        </TableContainer>
      ) : (
        <EmptyView size="large" title={emptyViewLabel} icon={UserIcon} />
      )}
    </Box>
  );
};

export default DataTableWithOptions;