/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Box, Button, Card, Stack, Typography } from '@mui/material';

const Page = () => {
  const handleGoToLogin = () => {
    typeof window !== 'undefined' &&
      window.location.replace(process.env.NEXT_PUBLIC_AUTH_APP_URL || '');
  };
  return (
    <Box
      sx={{
        width: '100%',
        height: 'calc(100vh - 10vh)',
        display: 'flex',
        overflow: 'hidden',
        alignItems: 'center',
        justifyContent: 'center',
      }}
    >
      <Card sx={{ width: 'fit-content', height: 'fit-content', p: 3 }}>
        <Stack direction={'column'} spacing={2}>
          <Typography variant="subtitle1">
            Require data not found, Please re-login.
          </Typography>
          <Button variant="contained" onClick={handleGoToLogin}>
            Re-Login
          </Button>
        </Stack>
      </Card>
    </Box>
  );
};
export default Page;
