import CloseIcon from '@mui/icons-material/Close';
import {
  Button,
  Checkbox,
  Dialog,
  DialogActions,
  DialogContent,
  DialogTitle,
  IconButton,
  List,
  ListItem,
  ListItemIcon,
  ListItemText,
  Stack,
  Typography,
} from '@mui/material';

type AddNodeDialogProps = {
  title: string;
  isOpen: boolean;
  description: string;
  handleCloseAction: any;
  labelSuccessBtn?: string;
  handleSuccessAction?: any;
  labelNegativeBtn?: string;
  handleNodeCheck: Function;
  data: Record<string, string | boolean>[];
};

const AddNodeDialog = ({
  title,
  isOpen,
  data = [],
  description,
  labelSuccessBtn,
  handleNodeCheck,
  labelNegativeBtn,
  handleCloseAction,
  handleSuccessAction,
}: AddNodeDialogProps) => {
  return (
    <Dialog fullWidth open={isOpen} maxWidth="sm" onClose={handleCloseAction}>
      <Stack direction="row" alignItems="center" justifyContent="space-between">
        <DialogTitle>{title}</DialogTitle>
        <IconButton
          onClick={handleCloseAction}
          sx={{ position: 'relative', right: 8 }}
        >
          <CloseIcon fontSize="small" />
        </IconButton>
      </Stack>

      <DialogContent>
        <Typography variant="body2">{description}</Typography>
        <List
          sx={{ width: '100%', mt: 2 }}
          subheader={
            <Typography variant="body1" fontWeight={600} pb={1}>
              Nodes Available
            </Typography>
          }
        >
          {data.map(({ id, name, isChecked }) => {
            const labelId = `node-checkbox-${id}`;
            return (
              <ListItem key={id.toString()} disablePadding>
                <ListItemIcon sx={{ ml: 1 }}>
                  <Checkbox
                    edge="start"
                    disableRipple
                    sx={{ p: 0 }}
                    checked={isChecked as boolean}
                    onChange={(e) => handleNodeCheck(id, e.target.checked)}
                    inputProps={{ 'aria-labelledby': labelId }}
                  />
                </ListItemIcon>
                <ListItemText id={id.toString()} primary={name} />
              </ListItem>
            );
          })}
        </List>
      </DialogContent>

      <DialogActions>
        <Stack
          width={'100%'}
          spacing={2}
          direction={'row'}
          alignItems="center"
          justifyContent={'space-between'}
        >
          {labelNegativeBtn && (
            <Button
              sx={{ p: 0 }}
              variant="text"
              color={'primary'}
              onClick={handleCloseAction}
            >
              {labelNegativeBtn}
            </Button>
          )}
          {labelSuccessBtn && (
            <Button
              size="small"
              variant="contained"
              onClick={() =>
                handleSuccessAction
                  ? handleSuccessAction()
                  : handleCloseAction()
              }
            >
              {labelSuccessBtn}
            </Button>
          )}
        </Stack>
      </DialogActions>
    </Dialog>
  );
};

export default AddNodeDialog;
