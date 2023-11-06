import { ColumnsWithOptions, MenuItemType } from '@/types';
import {
  Box,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Button,
  ListItem,
  ListItemText,
  Menu,
  MenuItem,
  Typography,
} from '@mui/material';
import { useState, useEffect } from 'react';
import Link from 'next/link';
import ArrowDropDown from '@mui/icons-material/ArrowDropDown';
import { SiteDto } from '@/generated';

import EmptyView from '../EmptyView';
import OptionsPopover from '../OptionsPopover';

interface DataTableWithOptionsInterface {
  icon?: any;
  dataset: any;
  isRowClickable?: boolean;
  emptyViewLabel?: string;
  onMenuItemClick: Function;
  menuOptions: MenuItemType[];
  columns: ColumnsWithOptions[];
  networkList?: SiteDto[];
  getSelectedNetwork?: Function;
}

type CellValueByTypeProps = {
  row: any;
  type: string;
  isRowClickable: boolean;
  menuOptions: MenuItemType[];
  onMenuItemClick: Function;
};

const CellValueByType = ({
  type,
  row,
  menuOptions,
  isRowClickable,
  onMenuItemClick,
}: CellValueByTypeProps) => {
  switch (type) {
    case 'name':
      return isRowClickable ? (
        // <Link href={`nodes/${row['id']}`} unselectable="on">
        <Link href={`nodes/uk-test36-hnode-a1-00ff`} unselectable="on">
          {row[type]}
        </Link>
      ) : (
        <Typography variant="body2">{row[type]}</Typography>
      );
    case 'actions':
      return (
        <OptionsPopover
          cid={'data-table-action-popover'}
          menuOptions={menuOptions}
          handleItemClick={onMenuItemClick}
        />
      );
    default:
      return <Typography variant="body2">{row[type]}</Typography>;
  }
};

const DataTableWithOptions = ({
  icon: Icon,
  columns,
  dataset,
  menuOptions,
  networkList,
  onMenuItemClick,
  getSelectedNetwork,
  emptyViewLabel = '',
  isRowClickable = true,
}: DataTableWithOptionsInterface) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const [selectedNetwork, setSelectedNetwork] = useState<any>();

  const handleOpenMenu = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };
  const handleNetworkSelect = (network: any) => {
    setSelectedNetwork(network);
    handleCloseMenu();
  };
  useEffect(() => {
    if (networkList && networkList.length > 0) {
      setSelectedNetwork(networkList[0].name);
    }
  }, []);

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
                    <b>
                      {column.label == 'network' ? (
                        <>
                          <Button
                            onClick={handleOpenMenu}
                            endIcon={<ArrowDropDown />}
                            aria-controls="network-menu"
                          >
                            <b>
                              {selectedNetwork
                                ? selectedNetwork
                                : 'networkName'}
                            </b>
                          </Button>
                          <Menu
                            id="network-menu"
                            anchorEl={anchorEl}
                            open={Boolean(anchorEl)}
                            onClose={handleCloseMenu}
                          >
                            {networkList &&
                              networkList.map(
                                ({ name, networkId, id }: SiteDto) => {
                                  if (getSelectedNetwork) {
                                    getSelectedNetwork(networkId); // Call your function only if it's defined
                                  }
                                  return (
                                    <MenuItem
                                      key={id}
                                      onClick={() => handleNetworkSelect(name)}
                                    >
                                      <ListItem>
                                        <ListItemText>
                                          <b>{name}</b>
                                        </ListItemText>
                                      </ListItem>
                                    </MenuItem>
                                  );
                                },
                              )}
                          </Menu>
                        </>
                      ) : (
                        column.label
                      )}
                    </b>
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
                        isRowClickable={isRowClickable}
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
        <EmptyView size="medium" title={emptyViewLabel} icon={Icon} />
      )}
    </Box>
  );
};

export default DataTableWithOptions;
