/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NodeConnectivityEnum } from '@/client/graphql/generated';
import { Box } from '@mui/material';
import Image from 'next/image';

interface IDeviceModalView {
  size?: number;
  image: string;
  nodeType?: string;
  connectivity?: NodeConnectivityEnum;
}

const DeviceModalView = ({
  image,
  nodeType = 'hnode',
  size = 300,
  connectivity = NodeConnectivityEnum.Online,
}: IDeviceModalView) => {
  return (
    <Box
      component={'div'}
      sx={{
        py: { xs: 2, md: 4, lg: 6 },
        height: '100%',
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <Image
        sizes="100vw"
        priority={true}
        placeholder="empty"
        style={{
          objectFit: 'contain',
          filter:
            connectivity === NodeConnectivityEnum.Offline
              ? 'grayscale(100%) opacity(70%)'
              : 'none',
          transition: 'filter 0.3s ease',
        }}
        src={image}
        width={size}
        height={size}
        alt={nodeType}
      />
    </Box>
  );
};

export default DeviceModalView;
