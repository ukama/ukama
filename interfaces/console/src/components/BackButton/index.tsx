/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { ArrowBack } from '@mui/icons-material';
import { IconButton, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';

interface IBackButton {
  title: string;
}

const BackButton = ({ title }: IBackButton) => {
  const router = useRouter();
  return (
    <Stack
      direction={'row'}
      alignItems={'center'}
      spacing={1.5}
      sx={{
        ':hover': {
          p: {
            color: colors.primaryMain,
            cursor: 'pointer',
          },
          '.MuiButtonBase-root': {
            color: colors.primaryMain,
          },
        },
      }}
      onClick={() => router.push('/console/home')}
    >
      <IconButton size="small" sx={{ p: 0 }}>
        <ArrowBack />
      </IconButton>
      <Typography variant={'body2'} fontWeight={500}>
        {title}
      </Typography>
    </Stack>
  );
};
export default BackButton;
