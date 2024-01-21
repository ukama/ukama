/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import * as React from 'react';
import {
  Select,
  MenuItem,
  Stack,
  Divider,
  Typography,
  Button,
  Grid,
} from '@mui/material';
import { FiberManualRecord, Add } from '@mui/icons-material';
import { colors } from '@/styles/theme';
import { SiteDto } from '@/generated';

type SiteHeaderProps = {
  sites?: SiteDto[];
  addSiteAction: () => void;
  restartSiteAction: () => void;
  onSiteSelect: (siteId: string) => void; // New prop for passing selected site ID to parent
};

const SiteHeader: React.FC<SiteHeaderProps> = ({
  sites = [],
  addSiteAction,
  onSiteSelect,
  restartSiteAction,
}) => {
  const [selectedSite, setSelectedSite] = React.useState('');

  const handleSiteChange = (e: any) => {
    e.stopPropagation();
    const selectedValue = e.target.value;
    if (selectedValue === undefined) {
      return;
    } else {
      const selectedSite = sites.find((site) => site.id === selectedValue);
      if (selectedSite) {
        onSiteSelect(selectedSite.id);
      }
    }
    setSelectedSite(selectedValue);
    onSiteSelect(selectedSite);
  };

  return (
    <>
      <Grid container spacing={2}>
        <Grid item xs={6}>
          <Stack direction="row" spacing={2} alignItems={'center'}>
            <Select
              disableUnderline
              variant="standard"
              value={selectedSite}
              onChange={handleSiteChange}
              MenuProps={{
                disablePortal: true,
                anchorOrigin: {
                  vertical: 'bottom',
                  horizontal: 'left',
                },
                transformOrigin: {
                  vertical: 'top',
                  horizontal: 'left',
                },
              }}
              displayEmpty
            >
              {sites.map(({ name, isDeactivated }) => (
                <MenuItem
                  key={name}
                  value={name}
                  sx={{
                    m: 0,
                    p: '6px 16px',
                    justifyContent: 'space-between',
                    backgroundColor: `${'inherit'} !important`,
                    ':hover': {
                      backgroundColor: `
                      colors.secondaryLight,
                      0.25,
                    )} !important`,
                    },
                  }}
                >
                  <Stack direction={'row'} alignItems={'center'} spacing={1}>
                    {isDeactivated ? (
                      <FiberManualRecord
                        htmlColor={colors.green}
                        fontSize="small"
                      />
                    ) : (
                      <FiberManualRecord
                        htmlColor={colors.red}
                        fontSize="small"
                      />
                    )}
                    <Typography variant="h6">{name}</Typography>
                  </Stack>
                </MenuItem>
              ))}
              <Divider />

              <Button
                variant="text"
                startIcon={<Add />}
                onClick={addSiteAction}
              >
                Add site
              </Button>
            </Select>
          </Stack>
        </Grid>

        <Grid item xs={6} container justifyContent={'flex-end'}>
          <Button
            variant="contained"
            color="primary"
            onClick={restartSiteAction}
          >
            RESTART SITE
          </Button>
        </Grid>
      </Grid>
    </>
  );
};

export default SiteHeader;
