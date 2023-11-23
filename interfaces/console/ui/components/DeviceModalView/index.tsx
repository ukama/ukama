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
        height: { xs: '80vh', md: '62vh' },
        display: 'flex',
        alignItems: 'center',
        justifyContent: 'center',
        marginTop: 6,
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
