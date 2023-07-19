import { HorizontalContainerJustify } from '@/styles/global';
import { colors } from '@/styles/theme';
import { Search } from '@mui/icons-material';
import {
  Button,
  Grid,
  TextField,
  Stack,
  Tooltip,
  Chip,
  Typography,
} from '@mui/material';
import ErrorOutlineIcon from '@mui/icons-material/ErrorOutline';
interface IPageContainerHeader {
  title: string;
  search?: string;
  subtitle?: string;
  showSearch?: boolean;
  buttonTitle: string;
  onSearchChange?: Function;
  handleButtonAction: Function;
  warningIcon?: boolean;
}

const PageContainerHeader = ({
  title,
  search,
  warningIcon = false,
  subtitle,
  buttonTitle,
  onSearchChange,
  showSearch = true,
  handleButtonAction,
}: IPageContainerHeader) => {
  const tooltipContentStyles = {
    marginBottom: '8px',
    listStyleType: 'disc',
    paddingLeft: '16px',
  };
  return (
    <HorizontalContainerJustify>
      <Grid container justifyContent={'space-between'} spacing={1}>
        <Grid container item xs={12} md="auto" alignItems={'center'}>
          <Grid item xs={'auto'}>
            <Stack direction="row" spacing={2} alignItems={'center'}>
              <Typography variant="h6" mr={1}>
                {title}
              </Typography>
              {warningIcon && (
                <Tooltip
                  title={
                    <>
                      {' '}
                      <h4>All data plans adhere to the following rules:</h4>
                      <ul style={tooltipContentStyles}>
                        <li>
                          All fixed Ukama data plans are honored at every
                          network
                        </li>
                        <li>
                          Unused data from billing cycle does not carry over
                        </li>
                        <li>
                          No changes can be made to the data plan mid-cycle
                        </li>
                      </ul>
                    </>
                  }
                  placement="bottom"
                  sx={{ backgroundColor: colors.vulcan }}
                >
                  <ErrorOutlineIcon sx={{ color: colors.black38 }} />
                </Tooltip>
              )}
            </Stack>
          </Grid>
          {subtitle && (
            <Grid item xs={'auto'}>
              <Typography variant="subtitle2">({subtitle})</Typography>
            </Grid>
          )}
          {showSearch && (
            <Grid item xs={12} md={'auto'} ml={{ xs: 0, md: 1.4 }}>
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
        <Grid item xs={12} md={'auto'}>
          <Button
            variant="contained"
            color="primary"
            size="medium"
            sx={{ width: { xs: '100%', md: '220px' } }}
            onClick={() => handleButtonAction()}
          >
            {buttonTitle}
          </Button>
        </Grid>
      </Grid>
    </HorizontalContainerJustify>
  );
};
export default PageContainerHeader;
