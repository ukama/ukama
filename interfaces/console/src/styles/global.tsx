/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */

import { colors } from '@/theme';
import { hexToRGB } from '@/utils';
import {
  Box,
  Container,
  IconButton,
  Link,
  Paper,
  Skeleton,
  TextField,
  Tooltip,
  TooltipProps,
  styled,
  tooltipClasses,
} from '@mui/material';

interface IRoundedCard {
  radius?: string;
}

const GlobalInput = styled(TextField)({
  height: '24px',
  padding: '12px 14px',
});

const RootContainer = styled(Container)((props) => ({
  height: 'auto',
  padding: '0px !important',
  background:
    props.theme.palette.mode === 'dark' ? colors.darkGreen05 : colors.white,
  boxShadow:
    '-4px 0px 4px 4px rgba(0, 0, 0, 0.05), 4px 4px 4px 4px rgba(0, 0, 0, 0.05)',
  borderRadius: '5px',
}));

const GradiantBar = styled(Box)({
  width: '100%',
  height: '12px',
  background:
    'linear-gradient(90deg, #00D3EB 0%, #2190F6 14.06%, #6974F8 44.27%, #6974F8 58.85%, #271452 100%)',
  borderRadius: '4px 4px 0px 0px',
});

const GradiantBarNoRadius = styled(Box)({
  width: '100%',
  height: '12px',
  background:
    'linear-gradient(90deg, #00D3EB 0%, #2190F6 14.06%, #6974F8 44.27%, #6974F8 58.85%, #271452 100%)',
});

const HorizontalContainerJustify = styled(Box)(() => ({
  width: '100%',
  height: 'auto',
  display: 'flex',
  flexDirection: 'row',
  alignItems: 'center',
  justifyContent: 'space-between',
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
  // backgroundColor: 'rgba(55, 57, 62, 0.11)',
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
})<IRoundedCard>(({ theme, radius = '10px' }) => ({
  marginTop: '18px',
  borderRadius: radius,
  padding: '24px 32px',
  overflow: 'scroll',
  height: 'calc(100vh - 216px)',
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

const IconStyle = {
  '.MuiSvgIcon-root': {
    width: '24px',
    height: '24px',
    fill: colors.white,
  },
  '.MuiBadge-root': {
    '.MuiSvgIcon-root': {
      width: '24px',
      height: '24px',
      fill: colors.white,
    },
  },
};

const ComponentContainer = {
  width: 'auto',
  height: 'auto',
  overflow: 'hidden',
};

const ScrollContainer = styled(Box)(() => ({
  position: 'relative',
  width: '100%',
  overflow: 'hidden',
}));

const ScrollableContent = styled(Box)(({ theme }) => ({
  display: 'flex',
  gap: theme.spacing(2),
  overflowX: 'auto',
  scrollBehavior: 'smooth',
  msOverflowStyle: 'none',
  scrollbarWidth: 'none',
  '&::-webkit-scrollbar': {
    display: 'none',
  },
  padding: theme.spacing(1),
}));

const NavigationWrapper = styled(Box)(({ theme }) => ({
  display: 'flex',
  gap: theme.spacing(1),
  position: 'absolute',
  top: 160,
  right: 50,
}));

const NavigationButton = styled(IconButton)(({ theme }) => ({
  backgroundColor: theme.palette.background.paper,
  boxShadow: theme.shadows[1],
  width: 30,
  height: 30,
  padding: 6,
  '&:hover': {
    backgroundColor: theme.palette.grey[100],
  },
  '&.Mui-disabled': {
    backgroundColor: theme.palette.grey[100],
    color: theme.palette.grey[400],
  },
  border: `1px solid ${colors.black38}`,
}));

const CardWrapper = styled(Box)(() => ({
  width: 'calc(25% - 12px)',
  minWidth: 200,
  flexShrink: 0,
}));
const DataPlanEmptyView = styled(Box)(({ theme }) => ({
  display: 'flex',
  flexDirection: 'column',
  alignItems: 'center',
  justifyContent: 'center',
  padding: theme.spacing(4),
  width: '100%',
  color: theme.palette.text.secondary,
}));
export {
  CardWrapper,
  CenterContainer,
  ComponentContainer,
  ContainerJustifySpaceBtw,
  ContainerMax,
  DarkTooltip,
  DataPlanEmptyView,
  GlobalInput,
  GradiantBar,
  GradiantBarNoRadius,
  HorizontalContainer,
  HorizontalContainerJustify,
  IconStyle,
  LinkStyle,
  MessageContainer,
  NavigationButton,
  NavigationWrapper,
  PageContainer,
  RootContainer,
  RoundedCard,
  ScrollContainer,
  ScrollableContent,
  SimpleCardWithBorder,
  SkeletonRoundedCard,
  VerticalContainer,
};
