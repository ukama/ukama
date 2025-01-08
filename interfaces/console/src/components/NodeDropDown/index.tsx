/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, NodeStateEnum } from '@/client/graphql/generated';
import { colors } from '@/theme';
import { hexToRGB } from '@/utils';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CircleIcon from '@mui/icons-material/Circle';
import InfoIcon from '@mui/icons-material/InfoOutlined';
import {
  MenuItem,
  Select,
  SelectChangeEvent,
  Stack,
  Typography,
} from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';
import { PaperProps, SelectDisplayProps, useStyles } from './styles';

const getStatus = (status: NodeStateEnum, time: number) => {
  let str = '';
  switch (status) {
    case NodeStateEnum.Unknown:
      str = 'Active';
    case NodeStateEnum.Configured:
      str = 'Configured';
    case NodeStateEnum.Operational:
      str = 'Onboarded';
    case NodeStateEnum.Faulty:
      str = 'Faulty';
    default:
      str = 'Unknown';
  }
  return (
    <Typography variant={'h6'} mr={'6px'}>
      {str}
    </Typography>
  );
};

const getStatusIcon = (status: NodeStateEnum) => {
  switch (status.toUpperCase()) {
    case NodeStateEnum.Unknown:
      return (
        <CheckCircleIcon sx={{ color: colors.black54 }} fontSize={'small'} />
      );
    case NodeStateEnum.Configured:
      return <InfoIcon sx={{ color: colors.orange }} fontSize={'small'} />;
    case NodeStateEnum.Operational:
      return <InfoIcon htmlColor={colors.darkGreen05} fontSize={'small'} />;
    case NodeStateEnum.Faulty:
      return <InfoIcon htmlColor={colors.red} fontSize={'small'} />;
    default:
      return <CircleIcon htmlColor={colors.black54} fontSize={'small'} />;
  }
};

interface INodeDropDown {
  loading: boolean;
  onAddNode: Function;
  nodes: Node[] | [];
  onNodeSelected: Function;
  selectedNode: Node | undefined;
}

const NodeDropDown = ({
  nodes = [],
  onAddNode,
  selectedNode,
  loading = true,
  onNodeSelected,
}: INodeDropDown) => {
  const classes = useStyles();
  const handleChange = (e: SelectChangeEvent<string>) => {
    const { target } = e;
    target.value &&
      onNodeSelected(nodes.find((item: Node) => item.name === target.value));
  };
  return (
    <Stack direction={'row'} spacing={1} alignItems="center">
      {selectedNode &&
        getStatusIcon(selectedNode.status.state as NodeStateEnum)}

      <LoadingWrapper radius="small" isLoading={loading} width={'244px'}>
        <Select
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
          className={classes.selectStyle}
          renderValue={(selected) => selected}
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
        </Select>
      </LoadingWrapper>
    </Stack>
  );
};

export default NodeDropDown;
