/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
import { IconStyle } from '@/styles/global';
import AccountCircleIcon from '@mui/icons-material/AccountCircle';
import {
  Divider,
  IconButton,
  Paper,
  Popover,
  Stack,
  Typography,
} from '@mui/material';
import Link from 'next/link';
import { useState } from 'react';

const AccountPopover = () => {
  const [anchorEl, setAnchorEl] = useState<null | HTMLElement>(null);

  const handleClick = (event: React.MouseEvent<HTMLElement>) => {
    setAnchorEl(event.currentTarget);
  };

  const handleClose = () => {
    setAnchorEl(null);
  };

  const open = Boolean(anchorEl);
  const id = open ? 'alert-popover' : undefined;

  return (
    <>
      <IconButton
        sx={{
          ...IconStyle,
        }}
        onClick={handleClick}
      >
        <AccountCircleIcon />
      </IconButton>
      <Popover
        id={id}
        open={open}
        sx={{ mt: 1.5 }}
        anchorEl={anchorEl}
        onClose={handleClose}
        anchorOrigin={{
          vertical: 'bottom',
          horizontal: 'right',
        }}
        transformOrigin={{
          vertical: 'top',
          horizontal: 'center',
        }}
      >
        <Paper sx={{ px: 3, py: 2 }}>
          <Stack spacing={1}>
            <Link
              href={`${process.env.NEXT_PUBLIC_AUTH_APP_URL}/user/account-settings`}
              prefetch={true}
              style={{
                borderRadius: 4,
                textDecoration: 'none',
              }}
            >
              <Typography
                variant="body1"
                color={'text.primary'}
                sx={{
                  ':hover': {
                    color: 'primary.main',
                  },
                }}
              >
                Account Settings
              </Typography>
            </Link>
            <Divider />
            <Link
              href={`${process.env.NEXT_PUBLIC_AUTH_APP_URL}/user/logout`}
              prefetch={true}
              style={{
                borderRadius: 4,
                textDecoration: 'none',
              }}
            >
              <Typography
                variant="body1"
                color={'text.primary'}
                sx={{
                  ':hover': {
                    color: 'primary.main',
                  },
                }}
              >
                Logout of account
              </Typography>
            </Link>
          </Stack>
        </Paper>
      </Popover>
    </>
  );
};

export default AccountPopover;
