import { colors } from '@/styles/theme';
import SearchIcon from '@mui/icons-material/Search';
import {
  InputBase,
  List,
  ListItem,
  Paper,
  Typography,
  styled,
} from '@mui/material';
import { LatLngLiteral } from 'leaflet';
import useOnclickOutside from 'react-cool-onclickoutside';

import usePlacesAutocomplete, {
  getGeocode,
  getLatLng,
} from 'use-places-autocomplete';

interface ISearchBar {
  placeholderText?: string;
  handleLocationSelected: (loc: LatLngLiteral) => void;
}

const StyledInputBase = styled(InputBase)(() => ({
  color: 'inherit',
  '& .MuiInputBase-input': {
    width: '100%',
  },
}));

const SearchBar = ({
  handleLocationSelected,
  placeholderText = 'Search',
}: ISearchBar) => {
  const {
    ready,
    value,
    suggestions: { status, data },
    setValue,
    clearSuggestions,
  } = usePlacesAutocomplete({
    callbackName: 'ukama_search_cb',
    requestOptions: {},
    cache: 24 * 60 * 60,
    debounce: 300,
  });
  const ref = useOnclickOutside(() => {
    clearSuggestions();
  });

  const handleInput = (e: any) => {
    setValue(e.target.value);
  };

  const handleSelect =
    ({ description }: any) =>
    () => {
      setValue(description, false);
      clearSuggestions();

      getGeocode({ address: description }).then((results) => {
        const { lat, lng } = getLatLng(results[0]);
        handleLocationSelected({ lat, lng });
      });
    };

  const renderSuggestions = () =>
    data.map((suggestion) => {
      const {
        place_id,
        structured_formatting: { main_text, secondary_text },
      } = suggestion;

      return (
        <ListItem key={place_id} onClick={handleSelect(suggestion)}>
          <Typography variant="body2">
            <strong>{main_text}</strong>
            {secondary_text && <small> - {secondary_text}</small>}
          </Typography>
        </ListItem>
      );
    });
  return (
    <div ref={ref} style={{ zIndex: 600, width: 'inherit' }}>
      <StyledInputBase
        fullWidth
        value={value}
        disabled={!ready}
        onChange={handleInput}
        placeholder={placeholderText}
        sx={{
          height: '42px',
          borderRadius: '4px',
          fontSize: '14px !important',
          minWidth: { xs: '100%', md: '300px' },
          border: `1px solid ${colors.silver}`,
          padding: '4px 8px 4px 12px !important',
          backgroundColor: colors.white,
        }}
        endAdornment={<SearchIcon fontSize="small" color="disabled" />}
      />
      {status === 'OK' && (
        <Paper sx={{ cursor: 'pointer' }}>
          <List>{renderSuggestions()}</List>
        </Paper>
      )}
    </div>
  );
};

export default SearchBar;
