/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Stack, Typography } from '@mui/material';
import { ElementType } from 'react';

interface IEmptyView {
  icon: ElementType;
  title: string;
  description?: string;
  size?: 'small' | 'medium' | 'large';
}
const EmptyView = ({
  title,
  description,
  icon: Icon,
  size = 'medium',
}: IEmptyView) => {
  return (
    <Stack
      spacing={1}
      sx={{
        height: '100%',
        width: '100%',
        display: 'flex',
        overflow: 'auto',
        alignSelf: 'center',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <Icon fontSize={size} color="textPrimary" style={{ opacity: 0.6 }} />
      <Typography variant="body1" fontWeight={500}>
        {title}
      </Typography>
      {description && (
        <Typography
          variant="body2"
          fontWeight={400}
          textAlign={'center'}
          width={{ xs: '60%', md: '25%' }}
        >
          {description}
        </Typography>
      )}
    </Stack>
  );
};

export default EmptyView;
