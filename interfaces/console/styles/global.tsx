import { hexToRGB } from '@/utils';
import {
  Box,
  Link,
  Paper,
  Skeleton,
  Tooltip,
  TooltipProps,
  styled,
  tooltipClasses,
} from '@mui/material';
import { makeStyles } from '@mui/styles';
import { colors } from './theme';

interface IRoundedCard {
  radius?: string;
}

const globalUseStyles = makeStyles(() => ({
  inputFieldStyle: {
    height: '24px',
    padding: '12px 14px',
  },
  disableInputFieldStyle: {
    padding: '4px 0px',
    '-webkit-text-fill-color': `${colors.black} !important`,
  },
  backToNodeGroupButtonStyle: {
    position: 'fixed',
    left: '50%',
    bottom: '20px',
    transform: 'translate(-50%, -50%)',
    margin: '0 auto',
    pointer: 'cursor',
  },
  GridContainer: {
    padding: '1em',
  },
}));

const HorizontalContainerJustify = styled(Box)((props) => ({
  width: '100%',
  height: 'auto',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'space-between',
  [props.theme.breakpoints.down('md')]: {
    flexDirection: 'column',
  },
}));

const HorizontalContainer = styled(Box)({
  width: '100%',
  height: 'auto',
  display: 'flex',
  alignItems: 'center',
  flexDirection: 'row',
});

const VerticalContainer = styled(Box)({
  width: '100%',
  height: '100%',
  display: 'flex',
  alignItems: 'center',
  overflow: 'scroll',
  flexDirection: 'column',
});

const CenterContainer = styled(Box)({
  width: '100%',
  height: '100%',
  display: 'flex',
  padding: '18px',
  alignItems: 'center',
  flexDirection: 'column',
  justifyContent: 'center',
});

const LinkStyle = styled(Link)({
  fontSize: '14px',
  width: 'fit-content',
  alignSelf: 'flex-end',
  color: colors.primaryMain,
  letterSpacing: '0.4px',
  textDecoration: 'none',
  '&:hover': {
    textDecoration: 'underline',
  },
});

const MessageContainer = styled(Box)({
  paddingBottom: '5%',
});

const ContainerJustifySpaceBtw = styled(Box)({
  width: '100%',
  display: 'flex',
  paddingBottom: 10,
  flexDirection: 'row',
  justifyContent: 'space-between',
  textAlign: 'center',
});

const RoundedCard = styled(Paper, {
  shouldForwardProp: (prop) => prop !== 'radius',
})<IRoundedCard>(({ theme, radius = '10px' }) => ({
  width: '100%',
  padding: '18px 28px',
  height: '100%',
  borderRadius: radius,
  display: 'inline-block',
  boxShadow: '2px 2px 6px rgba(0, 0, 0, 0.05)',
  [theme.breakpoints.down('sm')]: {
    padding: '18px',
  },
}));

const SkeletonRoundedCard = styled(Skeleton)(() => ({
  width: '100%',
  height: '100%',
  borderRadius: '10px',
  display: 'inline-block',
  backgroundColor: 'rgba(255, 255, 255, 0.25)',
}));

const ContainerMax = styled(Box)(() => ({
  width: '100%',
  height: '100%',
}));

const SimpleCardWithBorder = styled(Box)((props) => ({
  borderRadius: '4px',
  border: `1px solid ${hexToRGB(props.theme.palette.text.primary, 0.1)}`,
}));

const PageContainer = styled(Paper, {
  shouldForwardProp: (prop) => prop !== 'radius',
})<IRoundedCard>(({ theme, radius = '5px' }) => ({
  marginTop: '18px',
  borderRadius: radius,
  padding: '24px 32px',
  overflow: 'scroll',
  height: 'calc(100vh - 210px)',
  [theme.breakpoints.down('sm')]: {
    marginTop: '12px',
    padding: '12px 18px',
  },
}));

const DarkTooltip = styled(({ className, ...props }: TooltipProps) => (
  <Tooltip {...props} classes={{ popper: className }} />
))(({ theme }) => ({
  [`& .${tooltipClasses.tooltip}`]: {
    backgroundColor: colors.secondaryDark,
    color: 'rgba(255, 255, 255, 0.87)',
    boxShadow: theme.shadows[1],
    fontSize: '16px',
    fontWeight: 600,
  },
}));
export {
  CenterContainer,
  ContainerJustifySpaceBtw,
  ContainerMax,
  DarkTooltip,
  HorizontalContainer,
  HorizontalContainerJustify,
  LinkStyle,
  MessageContainer,
  PageContainer,
  RoundedCard,
  SimpleCardWithBorder,
  SkeletonRoundedCard,
  VerticalContainer,
  globalUseStyles,
};
