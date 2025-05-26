/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { HorizontalContainerJustify } from '@/styles/global';
import { colors } from '@/theme';
import { Search } from '@mui/icons-material';
import { Button, Grid, TextField, Typography } from '@mui/material';

interface IPageContainerHeader {
  title: string;
  search?: string;
  buttonId?: string;
  subtitle?: string;
  showSearch?: boolean;
  buttonTitle?: string;
  onSearchChange?: (value: string) => void;
  handleButtonAction?: () => void;
}

const PageContainerHeader = ({
  title,
  search,
  subtitle,
  buttonId,
  onSearchChange,
  buttonTitle = '',
  showSearch = true,
  handleButtonAction,
}: IPageContainerHeader) => {
  return (
    <HorizontalContainerJustify>
      <Grid container justifyContent={'space-between'} spacing={1}>
        <Grid container item xs={12} md="auto" alignItems={'center'}>
          <Grid item xs={'auto'}>
            <Typography variant="h6" mr={1}>
              {title}
            </Typography>
          </Grid>
          {subtitle && (
            <Grid item xs={'auto'}>
              <Typography variant="subtitle2">({subtitle})</Typography>
            </Grid>
          )}
          {showSearch && (
            <Grid
              item
              xs={12}
              md={'auto'}
              ml={{ xs: 0, md: 1.4 }}
              mt={{ xs: 1, md: 0 }}
            >
              <TextField
                id="subscriber-search"
                variant="outlined"
                size="small"
                placeholder="Search"
                value={search}
                onChange={(e) =>
                  onSearchChange && onSearchChange(e.target.value)
                }
                sx={{ width: { xs: '100%', lg: '250px' } }}
                InputLabelProps={{
                  shrink: false,
                }}
                InputProps={{
                  endAdornment: <Search htmlColor={colors.black54} />,
                }}
              />
            </Grid>
          )}
        </Grid>
        {buttonTitle && (
          <Grid item xs={12} md={'auto'}>
            <Button
              id={buttonId}
              variant="contained"
              color="primary"
              size="medium"
              sx={{ width: { xs: '100%', md: 'fit-content' } }}
              onClick={() => handleButtonAction && handleButtonAction()}
            >
              {buttonTitle}
            </Button>
          </Grid>
        )}
      </Grid>
    </HorizontalContainerJustify>
  );
};
export default PageContainerHeader;
