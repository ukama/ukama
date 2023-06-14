import { colors } from '@/styles/theme';
import SearchIcon from '@mui/icons-material/Search';
import { InputBase, styled } from '@mui/material';

interface ISearchBar {
  value: string;
  placeholderText?: string;
  handleOnChange: Function;
}

const StyledInputBase = styled(InputBase)(() => ({
  color: 'inherit',
  '& .MuiInputBase-input': {
    width: '100%',
  },
}));

const SearchBar = ({
  value,
  placeholderText = 'Search',
  handleOnChange,
}: ISearchBar) => (
  <StyledInputBase
    fullWidth
    value={value}
    placeholder={placeholderText}
    onChange={(e: any) => handleOnChange(e.target.value)}
    sx={{
      zIndex: 400,
      height: '42px',
      borderRadius: '4px',
      fontSize: '14px !important',
      minWidth: { xs: '100%', md: '300px' },
      border: `1px solid ${colors.silver}`,
      padding: '4px 8px 4px 12px !important',
      backgroundColor: colors.white,
    }}
    endAdornment={<SearchIcon fontSize="small" color="primary" />}
  />
);

export default SearchBar;
