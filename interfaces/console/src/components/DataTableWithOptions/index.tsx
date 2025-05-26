/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { ColumnsWithOptions, MenuItemType } from '@/types';
import ArrowDropDown from '@mui/icons-material/ArrowDropDown';
import {
  Box,
  Button,
  Chip,
  ListItem,
  ListItemText,
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
import Link from 'next/link';
import { useEffect, useState } from 'react';

import {
  Invitation_Status,
  NetworkDto,
  NodeConnectivityEnum,
} from '@/client/graphql/generated';
import colors from '@/theme/colors';
import { getInvitationStatusColor, roleEnumToString } from '@/utils';
import EmptyView from '../EmptyView';
import OptionsPopover from '../OptionsPopover';

interface DataTableWithOptionsInterface {
  icon?: any;
  dataset: any;
  withStatusColumn?: boolean;
  isRowClickable?: boolean;
  emptyViewLabel?: string;
  onMenuItemClick: (id: string, type: string) => void;
  menuOptions: MenuItemType[];
  columns: ColumnsWithOptions[];
  networkList?: NetworkDto[];
  getSelectedNetwork?: (networkId: string) => void;
  emptyViewDescription?: string;
}

type CellValueByTypeProps = {
  row: any;
  type: string;
  isRowClickable: boolean;
  withStatusColumn: boolean;
  menuOptions: MenuItemType[];
  onMenuItemClick: (type: string) => void;
};

const CellValueByType = ({
  type,
  row,
  menuOptions,
  isRowClickable,
  onMenuItemClick,
  withStatusColumn,
}: CellValueByTypeProps) => {
  switch (type) {
    case 'id':
      return isRowClickable ? (
        <Link href={`nodes/${row[type]}`} unselectable="on">
          {row[type]}
        </Link>
      ) : (
        <Typography variant="body2">{row[type]}</Typography>
      );
    case 'role':
      return (
        <Chip
          label={roleEnumToString(row[type])}
          sx={{ color: 'white' }}
          color={'info'}
        />
      );
    case 'status':
      return getInvitationStatusColor(
        row[type],
        new Date(row['expireAt']) < new Date(),
      );
    case 'connectivity':
      return (
        <Chip
          sx={{
            p: 1,
            backgroundColor:
              row[type] === NodeConnectivityEnum.Online
                ? colors.primaryLight
                : colors.dullRed,
            color: (theme) => theme.palette.text.primary,
          }}
          label={row[type]}
        />
      );
    case 'actions':
      if (
        (withStatusColumn &&
          row['status'] === Invitation_Status.InviteAccepted) ||
        row['status'] === Invitation_Status.InviteDeclined ||
        new Date(row['expireAt']) < new Date()
      ) {
        return <div>-</div>;
      } else
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
  emptyViewDescription,
  isRowClickable = true,
  withStatusColumn = false,
}: DataTableWithOptionsInterface) => {
  const [anchorEl, setAnchorEl] = useState(null);
  const [selectedNetwork, setSelectedNetwork] = useState<any>();

  const handleOpenMenu = (event: any) => {
    setAnchorEl(event.currentTarget);
  };

  const handleCloseMenu = () => {
    setAnchorEl(null);
  };
  const handleNetworkSelect = (network: string, networkId: string) => {
    setSelectedNetwork(network);
    handleCloseMenu();
    if (getSelectedNetwork) {
      getSelectedNetwork(networkId);
    }
  };
  useEffect(() => {
    if (networkList && networkList.length > 0) {
      setSelectedNetwork(networkList[0].name);
    }
  }, [networkList]);

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
                            sx={{ p: 0, typography: 'body2', fontWeight: 700 }}
                            onClick={handleOpenMenu}
                            endIcon={<ArrowDropDown />}
                            aria-controls="network-menu"
                          >
                            {selectedNetwork || 'networkName'}
                          </Button>
                          <Menu
                            id="network-menu"
                            anchorEl={anchorEl}
                            open={Boolean(anchorEl)}
                            onClose={handleCloseMenu}
                          >
                            {networkList?.map(({ name, id }: NetworkDto) => {
                              return (
                                <MenuItem
                                  key={id}
                                  onClick={() => handleNetworkSelect(name, id)}
                                >
                                  <ListItem>
                                    <ListItemText sx={{ typography: 'body1' }}>
                                      {name}
                                    </ListItemText>
                                  </ListItem>
                                </MenuItem>
                              );
                            })}
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
                <TableRow role="row" tabIndex={-1} key={`tr-${id}`}>
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
                        withStatusColumn={withStatusColumn}
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
        <EmptyView
          icon={Icon}
          size="medium"
          title={emptyViewLabel}
          description={emptyViewDescription}
        />
      )}
    </Box>
  );
};

export default DataTableWithOptions;
