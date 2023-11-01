/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { HorizontalContainerJustify } from '@/styles/global';
import { Button, Typography } from '@mui/material';

type TableHeaderProps = {
  title: string;
  buttonTitle?: string;
  showSecondaryButton: boolean;
  handleButtonAction?: any;
};

const TableHeader = ({
  title,
  buttonTitle,
  handleButtonAction,
  showSecondaryButton,
}: TableHeaderProps) => {
  return (
    <HorizontalContainerJustify>
      <Typography variant="body2" fontWeight={600}>
        {title}
      </Typography>
      {showSecondaryButton && (
        <Button
          variant="outlined"
          sx={{ width: '144px' }}
          onClick={() => handleButtonAction()}
        >
          {buttonTitle}
        </Button>
      )}
    </HorizontalContainerJustify>
  );
};

export default TableHeader;
