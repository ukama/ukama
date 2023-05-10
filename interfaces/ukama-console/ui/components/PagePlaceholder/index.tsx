import { isDarkmode } from '@/app-recoil';
import { LinkStyle } from '@/styles/global';
import { colors } from '@/styles/theme';
import { Button, Paper, Stack, Typography } from '@mui/material';
import React from 'react';
import { useRecoilValue } from 'recoil';

const PagePlaceholderSvg = React.lazy(() =>
  import('@/public/svg').then((module) => ({
    default: module.PagePlaceholderSvg,
  })),
);

type PagePlaceholderProps = {
  linkText?: string;
  hyperlink?: string;
  description: string;
  buttonTitle?: string;
  handleAction?: Function;
  showActionButton?: boolean;
};

const PagePlaceholder = ({
  hyperlink = '',
  linkText = '',
  description = '',
  buttonTitle = '',
  showActionButton = false,
  handleAction = () => {
    //Default behaviour
  },
}: PagePlaceholderProps) => {
  const _isDarkmode = useRecoilValue(isDarkmode);
  return (
    <Paper sx={{ height: 'inherit' }}>
      <Stack
        spacing={4}
        sx={{
          height: '100%',
          borderRadius: '5px',
          alignItems: 'center',
          justifyContent: 'center',
          p: 10,
        }}
      >
        <PagePlaceholderSvg
          color={_isDarkmode ? colors.white38 : colors.silver}
          color2={_isDarkmode ? colors.nightGrey12 : colors.white}
        />
        <Typography variant="body1" textAlign={'center'}>
          {`${description} `}
          {hyperlink && (
            <LinkStyle
              underline="hover"
              href={hyperlink}
              sx={{
                typography: 'body1',
              }}
            >
              {linkText}
            </LinkStyle>
          )}
        </Typography>

        {showActionButton && (
          <Button
            variant="contained"
            sx={{ width: 190 }}
            onClick={() => handleAction()}
          >
            {buttonTitle}
          </Button>
        )}
      </Stack>
    </Paper>
  );
};

export default PagePlaceholder;
