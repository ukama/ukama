/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
'use client';
import { Button, Paper, Stack, Typography } from '@mui/material';
import { useRouter } from 'next/navigation';

const NodeNotFoundPage = () => {
  const router = useRouter();
  return (
    <Paper elevation={0} sx={{ px: { xs: 2, md: 4 }, py: { xs: 1, md: 2 } }}>
      <Stack direction={'column'} spacing={2}>
        <Typography variant="h6">No new node found!</Typography>
        <Typography variant="body1">
          Please check that your node is On. If it&#39;s On you&#39;ll get
          notification when it get online and ready to configure.
        </Typography>
        <Button
          variant="contained"
          sx={{ width: 'fit-content', alignSelf: 'flex-end' }}
          onClick={() => router.push('/')}
        >
          Back to Home
        </Button>
      </Stack>
    </Paper>
  );
};

export default NodeNotFoundPage;
