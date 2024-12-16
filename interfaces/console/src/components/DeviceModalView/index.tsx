/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { NODE_IMAGES } from '@/constants';
import { Box } from '@mui/material';

interface IDeviceModalView {
  nodeType: string | undefined;
}

const DeviceModalView = ({ nodeType = 'hnode' }: IDeviceModalView) => {
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
      <img
        src={NODE_IMAGES[nodeType as 'hnode' | 'anode' | 'tnode']}
        alt="node-img"
        style={{ maxWidth: '100%', maxHeight: '500px' }}
      />
    </Box>
  );
};

export default DeviceModalView;
