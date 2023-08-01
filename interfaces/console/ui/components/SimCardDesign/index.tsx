import { colors } from '@/styles/theme';
import ShareIcon from '@mui/icons-material/Share';
import {
  Button,
  Card,
  Divider,
  IconButton,
  Stack,
  Theme,
  Typography,
} from '@mui/material';
import { makeStyles } from '@mui/styles';
type StyleProps = {
  isSelected?: boolean;
};

const useStyles = makeStyles<Theme, StyleProps>(() => ({
  cardStyle: {
    marginBottom: 16,
    cursor: 'pointer',
    padding: '6px 10px',
    border: `1px solid ${colors.darkGradient}`,
  },
}));

type SimCardDesignProps = {
  id: number;
  title: string;
  serial: string;
  isActivate?: boolean;
  isSelected: boolean;
  handleItemClick: Function;
};

const SimCardDesign = ({
  id,
  title,
  serial,
  isSelected,
  isActivate,
  handleItemClick,
}: SimCardDesignProps) => {
  const classes = useStyles({ isSelected });
  return (
    <Card className={classes.cardStyle} onClick={() => handleItemClick(id)}>
      <Stack
        direction="row"
        spacing={6}
        justifyContent="space-between"
        sx={{ alignItems: 'center', width: '100%' }}
      >
        <Stack
          direction="row"
          spacing={1}
          sx={{ alignItems: 'center' }}
          divider={<Divider orientation="vertical" flexItem />}
        >
          <Typography
            variant="body1"
            sx={{
              fontSize: '14px',
              fontWeight: 'bold',
            }}
          >
            {title}
          </Typography>
          <Typography variant="body1">{serial}</Typography>
        </Stack>

        {isActivate && (
          <Button sx={{ color: colors.black70 }}>AWAITING ACTIVATION</Button>
        )}
        <IconButton>
          <ShareIcon sx={{ color: colors.black70 }} />
        </IconButton>
      </Stack>
    </Card>
  );
};

export default SimCardDesign;
