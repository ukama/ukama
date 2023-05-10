import colors from '@/styles/theme/colors';
import FormControlLabel from '@mui/material/FormControlLabel';
import Radio, { RadioProps } from '@mui/material/Radio';
import { styled } from '@mui/material/styles';
const BpIcon = styled('span')(({ theme }) => ({
  borderRadius: '50%',
  width: 16,
  height: 16,
  boxShadow:
    theme.palette.mode === 'dark'
      ? colors.dark30Gradient
      : 'inset 0 0 0 1px rgba(16,22,26,.2), inset 0 -1px 0 rgba(16,22,26,.1)',
  backgroundColor:
    theme.palette.mode === 'dark' ? colors.darkBlue : colors.lightChocolate,
  'input:disabled ~ &': {
    boxShadow: 'none',
    background:
      theme.palette.mode === 'dark'
        ? colors.darkGrayGradient
        : colors.darkGradient,
  },
}));

const BpCheckedIcon = styled(BpIcon)({
  backgroundColor: colors.darkPurple,
  backgroundImage:
    'linear-gradient(180deg,hsla(0,0%,100%,.1),hsla(0,0%,100%,0))',
});

const BpRadio = (props: RadioProps) => {
  return (
    <Radio
      sx={{
        '&:hover': {
          bgcolor: 'transparent',
        },
      }}
      disableRipple
      color="default"
      checkedIcon={<BpCheckedIcon />}
      icon={<BpIcon />}
      {...props}
      checked={props.checked}
    />
  );
};

interface ICustomRadioButton {
  value: boolean;
  label: string;
}

const CustomRadioButton = ({ value, label }: ICustomRadioButton) => {
  return (
    <FormControlLabel
      checked={value}
      control={<BpRadio />}
      label={label}
      sx={{
        width: '100%',
        margin: '0px',
        padding: '6px 24px 6px 6px',
        ':hover': {
          backgroundColor: colors.lightPurpleGradient,
        },
      }}
    />
  );
};

export default CustomRadioButton;
