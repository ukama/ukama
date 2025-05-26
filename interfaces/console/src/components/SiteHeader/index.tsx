/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { Add, FiberManualRecord } from '@mui/icons-material';
import {
  Button,
  Divider,
  Grid,
  MenuItem,
  Select,
  SelectChangeEvent,
  Stack,
  Typography,
} from '@mui/material';
import * as React from 'react';

type Site = {
  name: string;
  health: 'online' | 'offline';
  duration: string;
};

type SiteHeaderProps = {
  sites: Site[];
  sitesAction: (site: Site) => void;
  addSiteAction: () => void;
  restartSiteAction: () => void;
};

const SiteHeader: React.FC<SiteHeaderProps> = ({
  sites,
  sitesAction,
  addSiteAction,
  restartSiteAction,
}) => {
  const [selectedSite, setSelectedSite] = React.useState('site1');
  const [siteHealth] = React.useState('online');

  const handleSiteChange = (e: SelectChangeEvent<string>) => {
    e.stopPropagation();
    const selectedValue = e.target.value;
    if (selectedValue === undefined) {
      return;
    } else {
      const selectedSite = sites.find((site) => site.name === selectedValue);
      if (selectedSite) {
        sitesAction(selectedSite);
      }
    }
    setSelectedSite(selectedValue);
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
              {sites.map(({ name, health }) => (
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
                    {health === 'online' ? (
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
                    <Typography variant="body1">{name}</Typography>
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
            <Typography variant="body1">
              {`is ${siteHealth} for ${
                sites.find((s) => s.name === selectedSite)?.duration
              }`}
            </Typography>
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
