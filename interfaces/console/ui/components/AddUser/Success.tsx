/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { Typography } from '@mui/material';

interface ISuccess {
  description: string;
}

const Success = ({ description }: ISuccess) => {
  return (
    <Typography variant="body1" sx={{ mb: 2 }}>
      {description}
    </Typography>
  );
};

export default Success;
