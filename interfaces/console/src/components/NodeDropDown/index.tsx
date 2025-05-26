/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, NodeConnectivityEnum } from '@/client/graphql/generated';
import { colors } from '@/theme';
import { duration, hexToRGB } from '@/utils';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import InfoIcon from '@mui/icons-material/InfoOutlined';
import { MenuItem, SelectChangeEvent, Stack, Typography } from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';
import { PaperProps, SelectDisplayProps, SelectStyle } from './styles';

const getStatusIcon = (status: NodeConnectivityEnum) => {
  switch (status) {
    case NodeConnectivityEnum.Unknown:
      return <CheckCircleIcon fontSize={'small'} color="warning" />;
    case NodeConnectivityEnum.Online:
      return <CheckCircleIcon fontSize={'small'} color="success" />;
    case NodeConnectivityEnum.Offline:
      return <InfoIcon fontSize={'small'} color="error" />;
    default:
      return <CheckCircleIcon fontSize={'small'} color="disabled" />;
  }
};

interface INodeDropDown {
  uptime: number;
  loading: boolean;
  isNodeReady?: boolean;
  onAddNode?: () => void;
  nodes: Node[] | [];
  onNodeSelected: (node: Node) => void;
  selectedNode: Node | undefined;
}

const NodeDropDown = ({
  uptime,
  nodes = [],
  selectedNode,
  loading = true,
  onNodeSelected,
  isNodeReady = true,
}: INodeDropDown) => {
  const handleChange = (e: SelectChangeEvent<unknown>) => {
    const value = e.target.value as string;
    const node = nodes.find((item: Node) => item.name === value);
    if (node) {
      onNodeSelected(node);
    }
  };
  return (
    <Stack direction={'row'} spacing={1} alignItems="center">
      {selectedNode &&
        getStatusIcon(selectedNode.status.connectivity as NodeConnectivityEnum)}

      <LoadingWrapper radius="small" isLoading={loading} width={'fit-content'}>
        <SelectStyle
          disableUnderline
          variant="standard"
          onChange={handleChange}
          value={selectedNode?.name ?? ''}
          SelectDisplayProps={{
            style: { ...SelectDisplayProps.style, marginRight: '8px' },
          }}
          MenuProps={{
            disablePortal: true,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
            transformOrigin: {
              vertical: 'top',
              horizontal: 'left',
            },
            PaperProps: {
              sx: {
                width: '164px',
                ...PaperProps,
              },
            },
          }}
          renderValue={(selected: unknown) => selected as React.ReactNode}
        >
          {nodes.map(({ id, name }) => (
            <MenuItem
              key={id}
              value={name}
              sx={{
                m: 0,
                p: '6px 24px',
                backgroundColor: `${
                  id === selectedNode?.id
                    ? hexToRGB(colors.secondaryLight, 0.25)
                    : 'inherit'
                } !important`,
                ':hover': {
                  backgroundColor: `${hexToRGB(colors.secondaryLight, 0.25)} !important`,
                },
              }}
            >
              <Typography variant="body1" pr={2}>
                {name}
              </Typography>
            </MenuItem>
          ))}
        </SelectStyle>
      </LoadingWrapper>

      {selectedNode && (
        <Typography variant={'subtitle1'}>
          {isNodeReady && (
            <>
              Node is up for <b>{duration(uptime)}</b>
            </>
          )}
          {selectedNode.status.connectivity === NodeConnectivityEnum.Offline &&
            'Node is offline'}
        </Typography>
      )}
    </Stack>
  );
};

export default NodeDropDown;
