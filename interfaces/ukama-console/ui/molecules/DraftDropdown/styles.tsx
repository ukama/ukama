import { colors } from '@/styles/theme';
import { Theme } from '@mui/material';
import { makeStyles } from '@mui/styles';

const useStyles = makeStyles<Theme>(() => ({
  selectStyle: () => ({
    marginTop: '8px',
    width: 'fit-content',
    color: colors.primaryMain,
    '.MuiSelect-icon': {
      paddingLeft: '16px',
    },
  }),
}));

const SelectDisplayProps = {
  style: {
    fontWeight: 400,
    display: 'flex',
    fontSize: '20px',
    marginLeft: '4px',
    alignItems: 'center',
    minWidth: 'fit-content',
  },
};

const PaperProps = {
  width: 240,
  boxShadow:
    '0px 5px 5px -3px rgba(0, 0, 0, 0.2), 0px 8px 10px 1px rgba(0, 0, 0, 0.14), 0px 3px 14px 2px rgba(0, 0, 0, 0.12)',
  borderRadius: '4px',
};

const ICON_STYLE = {
  fontSize: '18px',
  mr: 0.8,
};

export { useStyles, SelectDisplayProps, PaperProps, ICON_STYLE };
