/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import MenuOpenIcon from '@mui/icons-material/MenuOpen';
import { IconButton, Paper, Stack, Typography } from '@mui/material';
import React from 'react';
import LoadingWrapper from '../LoadingWrapper';
interface INodeStatsContainer {
  index: number;
  title: string;
  loading: boolean;
  selected?: number;
  isAlert?: boolean; //Pass true to show red border
  isCollapse?: boolean;
  isClickable?: boolean;
  onCollapse?: () => void;
  isCollapsable?: boolean;
  handleAction?: (index: number) => void;
  children: React.ReactNode;
}

const NodeStatsContainer = ({
  index,
  title,
  loading,
  children,
  onCollapse,
  handleAction,
  selected = -1,
  isAlert = false,
  isCollapse = false,
  isClickable = false,
  isCollapsable = false,
}: INodeStatsContainer) => {
  return (
    <LoadingWrapper
      width="100%"
      height="fit-content"
      radius="small"
      isLoading={loading}
    >
      <Paper
        sx={{
          padding: '24px 24px 24px 0px',
          cursor: (isCollapsable ?? !isClickable) ? 'default' : 'pointer',
          paddingLeft: isAlert && selected !== index ? '16px' : '24px',
          borderLeft: {
            md:
              selected === index
                ? `8px solid ${colors.secondaryMain}`
                : isAlert
                  ? `1px solid ${colors.error}`
                  : `8px solid ${colors.silver}`,
          },
          border: isAlert ? `0.5px solid ${colors.error}` : 'none',
        }}
        onClick={() => isClickable && handleAction && handleAction(index)}
      >
        <Stack
          direction="row"
          justifyContent="space-between"
          alignItems="center"
          spacing={2}
          sx={{ mb: 1 }}
        >
          {!isCollapse && <Typography variant="h6">{title}</Typography>}
          {isCollapsable && (
            <IconButton
              sx={{
                p: 0,
                position: isCollapse ? 'relative' : null,
                right: isCollapse ? 10 : null,
                transform: isCollapse ? 'rotate(180deg)' : 'none',
              }}
              onClick={() => onCollapse && onCollapse()}
            >
              <MenuOpenIcon fontSize="medium" />
            </IconButton>
          )}
        </Stack>
        {!isCollapse && children}
      </Paper>
    </LoadingWrapper>
  );
};

export default NodeStatsContainer;
