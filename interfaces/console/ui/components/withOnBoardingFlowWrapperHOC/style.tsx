import colors from '@/styles/theme/colors';
import { Box, Container, styled } from '@mui/material';

const RootContainer = styled(Container)((props) => ({
  background:
    props.theme.palette.mode === 'dark' ? colors.darkGreen05 : colors.white,
  boxShadow:
    '-4px 0px 4px 4px rgba(0, 0, 0, 0.05), 4px 4px 4px 4px rgba(0, 0, 0, 0.05)',
  borderRadius: '5px',
  position: 'absolute',
  left: '50%',
  top: '50%',
  transform: 'translate(-50%, -50%)',
}));

const GradiantBar = styled(Box)({
  width: '100%',
  height: '12px',
  background:
    'linear-gradient(90deg, #00D3EB 0%, #2190F6 14.06%, #6974F8 44.27%, #6974F8 58.85%, #271452 100%)',
  borderRadius: '4px 4px 0px 0px',
});

const ComponentContainer = {
  width: 'auto',
  height: 'auto',
  margin: '20px',
};

export { RootContainer, GradiantBar, ComponentContainer };
