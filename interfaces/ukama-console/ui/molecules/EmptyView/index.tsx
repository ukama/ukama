import { Stack, Typography } from '@mui/material';
import { ElementType } from 'react';

interface IEmptyView {
  icon: ElementType;
  title: string;
  size?: 'small' | 'medium' | 'large';
}
const EmptyView = ({ title, icon: Icon, size = 'medium' }: IEmptyView) => {
  return (
    <Stack
      spacing={1}
      sx={{
        height: '100%',
        width: '100%',
        display: 'flex',
        alignSelf: 'center',
        justifyContent: 'center',
        alignItems: 'center',
      }}
    >
      <Icon fontSize={size} color="textPrimary" style={{ opacity: 0.6 }} />
      <Typography variant="body2">{title}</Typography>
    </Stack>
  );
};

export default EmptyView;
