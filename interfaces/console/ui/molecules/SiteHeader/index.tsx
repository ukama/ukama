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
  const [siteHealth, setSiteHealth] = React.useState('online');
  const [isSiteSelected, setIsSiteSelected] = React.useState(false);

  const handleSiteChange = (e: any) => {
    e.stopPropagation();
    const selectedValue = e.target.value;
    if (selectedValue === undefined) {
      return;
    } else {
      setIsSiteSelected(true);
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
