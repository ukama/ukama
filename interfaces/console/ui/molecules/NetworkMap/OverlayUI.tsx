/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NodeStatusEnum } from '@/generated';
import { ExpandLess, ExpandMore } from '@mui/icons-material';
import TuneIcon from '@mui/icons-material/Tune';
import {
  Box,
  Checkbox,
  Collapse,
  List,
  ListItem,
  ListItemButton,
  ListItemIcon,
  ListItemText,
  Stack,
  Typography,
} from '@mui/material';
import { SimpleTreeView, TreeItem } from '@mui/x-tree-view';
import { useRouter } from 'next/router';
import { useState } from 'react';

export const LabelOverlayUI = ({ name }: { name: string }) => {
  return (
    <Box
      sx={{
        top: 24,
        left: 24,
        zIndex: 400,
        display: 'flex',
        padding: '4px 12px',
        borderRadius: '4px',
        width: 'fit-content',
        position: 'absolute',
        boxShadow: '2px 2px 6px 0px rgba(0, 0, 0, 0.05)',
        background: (theme) => theme.palette.background.paper,
      }}
    >
      <Typography variant="h6" fontWeight={500}>
        {name}
      </Typography>
    </Box>
  );
};

interface ISitesTree {
  sites: any;
}

export const SitesTree = ({ sites }: ISitesTree) => {
  const router = useRouter();
  return (
    <Box
      sx={{
        top: 24,
        right: 24,
        zIndex: 400,
        display: 'flex',
        padding: '16px 24px',
        borderRadius: '4px',
        width: '220px',
        position: 'absolute',
        boxShadow: '2px 2px 6px 0px rgba(0, 0, 0, 0.05)',
        background: (theme) => theme.palette.background.paper,
      }}
    >
      <Stack spacing={0.5}>
        <Typography variant="body1" fontWeight={600}>
          Network ({sites.length})
        </Typography>
        <SimpleTreeView
          aria-label="sites-tree"
          sx={{
            flexGrow: 1,
            overflowY: 'auto',
            height: 'fit-content',
            maxHeight: '400px',
          }}
        >
          {sites?.map((site: any) => {
            return (
              <TreeItem key={site.id} itemId={site.id} label={site.name}>
                <TreeItem
                  itemId={site.nodeId}
                  label={site.nodeName}
                  onClick={() => router.push(`/nodes/${site.nodeId}`)}
                />
              </TreeItem>
            );
          })}
        </SimpleTreeView>
      </Stack>
    </Box>
  );
};

interface ISitesSelection {
  filterState: NodeStatusEnum;
  handleFilterState: (value: NodeStatusEnum) => void;
}

export const SitesSelection = ({
  filterState,
  handleFilterState,
}: ISitesSelection) => {
  const [open, setIsOpen] = useState(false);

  const handleToggle = (value: NodeStatusEnum) => () => {
    handleFilterState(value);
  };
  return (
    <Box
      sx={{
        bottom: 12,
        left: 24,
        zIndex: 400,
        display: 'flex',
        minWidth: '200px',
        padding: '12px 18px',
        borderRadius: '4px',
        width: 'fit-content',
        position: 'absolute',
        boxShadow: '2px 2px 6px 0px rgba(0, 0, 0, 0.05)',
        background: (theme) => theme.palette.background.paper,
      }}
    >
      <List
        component="nav"
        sx={{ width: '100%', p: 0, bgcolor: 'background.paper' }}
        aria-labelledby="nested-list-subheader"
      >
        <ListItemButton
          sx={{ p: 0, m: 0, pb: 1 }}
          onClick={() => setIsOpen(!open)}
        >
          <ListItemIcon>
            <TuneIcon />
          </ListItemIcon>
          <ListItemText primary="Filters" />
          {open ? <ExpandLess /> : <ExpandMore />}
        </ListItemButton>
        <Collapse in={open} timeout="auto" unmountOnExit>
          <List component="div" disablePadding>
            {[
              { id: 0, label: 'All', value: NodeStatusEnum.Undefined },
              { id: 1, label: 'Configured', value: NodeStatusEnum.Configured },
              {
                id: 2,
                label: 'Maintenance',
                value: NodeStatusEnum.Maintenance,
              },
              { id: 3, label: 'Faulty', value: NodeStatusEnum.Faulty },
              { id: 4, label: 'Onboarded', value: NodeStatusEnum.Onboarded },
              { id: 5, label: 'Active', value: NodeStatusEnum.Active },
            ].map(({ id, label, value }) => {
              const labelId = `checkbox-list-label-${value}`;

              return (
                <ListItem key={id} disablePadding>
                  <ListItemButton
                    role={undefined}
                    sx={{ p: 0, m: 0 }}
                    onClick={handleToggle(value)}
                    dense
                  >
                    <ListItemIcon>
                      <Checkbox
                        edge="start"
                        checked={filterState === value}
                        tabIndex={-1}
                        disableRipple
                        inputProps={{ 'aria-labelledby': labelId }}
                      />
                    </ListItemIcon>
                    <ListItemText id={labelId} primary={label} />
                  </ListItemButton>
                </ListItem>
              );
            })}
          </List>
        </Collapse>
      </List>
    </Box>
  );
};
