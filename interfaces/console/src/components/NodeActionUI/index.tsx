/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NodeConnectivityEnum, NodeTypeEnum } from '@/client/graphql/generated';
import { NODE_IMAGES } from '@/constants';
import { colors } from '@/theme';
import { CircularProgress, Stack, Typography } from '@mui/material';
import DeviceModalView from '../DeviceModalView';

interface NodeActionUI {
  value: number;
  action: string;
  description: string;
  nodeType: NodeTypeEnum | undefined;
  connectivity: NodeConnectivityEnum | undefined;
}

export const NodeActionUI = ({
  value,
  action,
  description,
  nodeType = NodeTypeEnum.Tnode,
  connectivity = NodeConnectivityEnum.Online,
}: NodeActionUI) => {
  return (
    <Stack
      spacing={2}
      width={'100%'}
      height={'100%'}
      direction={'column'}
      alignItems={'center'}
      justifyContent={'center'}
    >
      {NODE_IMAGES[nodeType as 'hnode' | 'anode' | 'tnode'] && (
        <Stack direction={'column'} spacing={2} alignItems={'center'}>
          <Stack
            position={'relative'}
            alignItems={'center'}
            justifyContent={'center'}
            sx={{
              width: 300,
              height: 300,
              borderRadius: '50%',
              bgcolor:
                connectivity === NodeConnectivityEnum.Online
                  ? 'background.paper'
                  : colors.dullGrey,
            }}
          >
            <CircularProgress
              size={300}
              thickness={1}
              variant="determinate"
              sx={{
                position: 'absolute',
                color: 'primary.main',
              }}
              value={value}
            />
            <DeviceModalView
              size={200}
              nodeType={nodeType}
              connectivity={connectivity}
              image={NODE_IMAGES[nodeType as 'hnode' | 'anode' | 'tnode']}
            />
          </Stack>
          <Typography
            variant="subtitle1"
            color="text.secondary"
            fontWeight={500}
          >
            {description}
          </Typography>
        </Stack>
      )}
    </Stack>
  );
};
