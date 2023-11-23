/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { isDarkmode } from '@/app-recoil';
import { Brightness3, Brightness7 } from '@mui/icons-material';
import { IconButton } from '@mui/material';
import { useRecoilState } from 'recoil';

const DarkModToggle = () => {
  const [_isDarkMod, _setIsDarkMod] = useRecoilState(isDarkmode);
  const icon = _isDarkMod ? <Brightness7 /> : <Brightness3 />;
  const handleToggle = () => _setIsDarkMod(!_isDarkMod);

  return (
    <IconButton
      size="small"
      color="inherit"
      onClick={handleToggle}
      sx={{ p: '8px' }}
      aria-label="darkmode-btn"
    >
      {icon}
    </IconButton>
  );
};

export default DarkModToggle;
