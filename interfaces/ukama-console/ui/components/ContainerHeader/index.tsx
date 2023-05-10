import { colors } from '@/styles/theme';
import SearchIcon from '@mui/icons-material/Search';
import { Button, Grid, Stack, Typography } from '@mui/material';
import InputBase from '@mui/material/InputBase';
import { styled } from '@mui/material/styles';

import { useEffect, useState } from 'react';
type ContainerHeaderProps = {
  title?: string;
  stats?: string;
  buttonSize?: 'small' | 'medium' | 'large';
  buttonTitle?: string;
  showButton?: boolean;
  showSearchBox?: boolean;
  handleSearchChange?: Function;
  handleAllNodeUpdate?: Function;
  handleButtonAction?: Function;
};

const StyledInputBase = styled(InputBase)(() => ({
  color: 'inherit',
  '& .MuiInputBase-input': {
    width: '100%',
  },
}));

const ContainerHeader = ({
  title,
  stats,
  buttonTitle,
  handleAllNodeUpdate = () => {
    /* Default empty function */
  },
  showButton = false,
  buttonSize = 'large',
  showSearchBox = false,
  handleSearchChange = () => {
    /* Default empty function */
  },
  handleButtonAction = () => {
    /* Default function */
  },
}: ContainerHeaderProps) => {
  const [currentSearchValue, setCurrentSearchValue] = useState<string>('');

  useEffect(() => {
    handleSearchChange(currentSearchValue.toLowerCase());
  }, [currentSearchValue]);

  return (
    <Grid container spacing={{ xs: 2, md: 0 }} justifyContent="space-between">
      <Grid item xs={6} md={showSearchBox ? 7 : 9}>
        <Stack direction="row" alignItems="center" spacing={{ xs: 1, md: 2 }}>
          <Typography variant="h6">{title}</Typography>
          {stats && (
            <Typography
              variant="subtitle2"
              letterSpacing="4px"
              color={'textSecondary'}
            >
              &#40;{stats}&#41;
            </Typography>
          )}
        </Stack>
      </Grid>

      <Grid
        item
        xs={6}
        md={0}
        justifyContent="flex-end"
        display={{
          xs: showButton ? 'flex' : 'none',
          md: 'none',
        }}
      >
        <Button
          sx={{
            px: { xs: 2, md: 4 },
            width: { xs: '100%', md: 'fit-content' },
          }}
          size={buttonSize}
          variant="contained"
          onClick={() => handleButtonAction()}
        >
          {buttonTitle}
        </Button>
      </Grid>

      <Grid
        item
        xs={12}
        md={3}
        justifyContent={'flex-end'}
        display={showSearchBox ? 'flex' : 'none'}
      >
        <StyledInputBase
          placeholder="Search…"
          value={currentSearchValue}
          onChange={(e: any) => setCurrentSearchValue(e.target.value)}
          sx={{
            height: '42px',
            borderRadius: 2,
            minWidth: { xs: '100%', md: '300px' },
            border: `1px solid ${colors.silver}`,
            padding: '4px 8px 4px 12px !important',
          }}
          endAdornment={<SearchIcon fontSize="small" color="primary" />}
        />
      </Grid>

      <Grid
        item
        xs={0}
        md={showSearchBox ? 2 : 3}
        justifyContent="flex-end"
        display={{
          xs: 'none',
          md: showButton ? 'flex' : 'none',
        }}
      >
        <Button
          sx={{
            px: 2,
            width: { xs: '100%', md: 'fit-content' },
          }}
          size={buttonSize}
          variant="contained"
          onClick={() => handleAllNodeUpdate()}
        >
          {buttonTitle}
        </Button>
      </Grid>
    </Grid>
  );
};

export default ContainerHeader;
