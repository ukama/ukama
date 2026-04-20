/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Node, SiteDto as Site } from '@/client/graphql/generated';
import { colors } from '@/theme';
import { duration, hexToRGB } from '@/utils';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import InfoIcon from '@mui/icons-material/InfoOutlined';
import { MenuItem, SelectChangeEvent, Stack, Typography } from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';
import { PaperProps, SelectDisplayProps, SelectStyle } from './styles';

const getStatusIcon = (status: boolean) => {
  switch (status) {
    case true:
      return <CheckCircleIcon fontSize={'small'} color="success" />;
    case false:
      return <InfoIcon fontSize={'small'} color="error" />;
  }
};

interface INodeDropDown {
  uptime: number;
  loading: boolean;
  isReady: boolean;
  objs: Node[] | Site[] | [];
  onSelected: (obj: Node | Site) => void;
  selected: Node | Site | undefined;
}

const NodeDropDown = ({
  uptime,
  objs = [],
  selected,
  loading = true,
  onSelected,
  isReady = true,
}: INodeDropDown) => {
  const handleChange = (e: SelectChangeEvent<unknown>) => {
    const value = e.target.value as string;
    const obj = objs.find((item: Node | Site) => item.name === value);
    if (obj) {
      onSelected(obj);
    }
  };
  return (
    <Stack direction={'row'} spacing={1} alignItems="center">
      {getStatusIcon(isReady)}

      <LoadingWrapper radius="small" isLoading={loading} width={'fit-content'}>
        <SelectStyle
          disableUnderline
          variant="standard"
          onChange={handleChange}
          value={selected?.name ?? ''}
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
          {objs.map(({ id, name }) => (
            <MenuItem
              key={id}
              value={name}
              sx={{
                m: 0,
                p: '6px 24px',
                backgroundColor: `${
                  id === selected?.id
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

      {selected && (
        <Typography variant={'subtitle1'}>
          {isReady ? (
            <>
              {selected.name} is up for <b>{duration(uptime)}</b>
            </>
          ) : (
            <>{selected.name} is offline</>
          )}
        </Typography>
      )}
    </Stack>
  );
};

export default NodeDropDown;
