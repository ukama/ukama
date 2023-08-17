import ChevronRightIcon from '@mui/icons-material/ChevronRight';
import ExpandMoreIcon from '@mui/icons-material/ExpandMore';
import { TreeItem, TreeView } from '@mui/lab';
import {
  Box,
  Checkbox,
  FormControlLabel,
  FormGroup,
  Stack,
  Typography,
} from '@mui/material';

export const LabelOverlayUI = ({ name }: { name: string }) => {
  return (
    <Box
      sx={{
        top: 24,
        left: 24,
        zIndex: 400,
        display: 'flex',
        padding: '4px 12px',
        borderRadius: '4px',
        width: 'fit-content',
        position: 'absolute',
        boxShadow: '2px 2px 6px 0px rgba(0, 0, 0, 0.05)',
        background: (theme) => theme.palette.background.paper,
      }}
    >
      <Typography variant="h6" fontWeight={500}>
        {name}
      </Typography>
    </Box>
  );
};

export const SitesTree = () => {
  return (
    <Box
      sx={{
        top: 24,
        right: 24,
        zIndex: 400,
        display: 'flex',
        padding: '16px 24px',
        borderRadius: '4px',
        width: '220px',
        position: 'absolute',
        boxShadow: '2px 2px 6px 0px rgba(0, 0, 0, 0.05)',
        background: (theme) => theme.palette.background.paper,
      }}
    >
      <Stack spacing={0.5}>
        <Typography variant="body1" fontWeight={600}>
          Network (2)
        </Typography>
        <TreeView
          aria-label="sites-tree"
          defaultCollapseIcon={<ExpandMoreIcon />}
          defaultExpandIcon={<ChevronRightIcon />}
          sx={{
            flexGrow: 1,
            overflowY: 'auto',
            height: 'fit-content',
            maxHeight: '400px',
          }}
        >
          <TreeItem nodeId="1" label="Site 1">
            <TreeItem nodeId="2" label="Tower Node 1" />
            <TreeItem nodeId="3" label="Tower Node 2" />
            <TreeItem nodeId="4" label="Amplifier unit 1" />
          </TreeItem>
          <TreeItem nodeId="5" label="Site 2">
            <TreeItem nodeId="6" label="Tower Node 1" />
            <TreeItem nodeId="7" label="Amplifier unit 1" />
            <TreeItem nodeId="8" label="Amplifier unit 2" />
          </TreeItem>
        </TreeView>
      </Stack>
    </Box>
  );
};

export const SitesSelection = () => {
  return (
    <Box
      sx={{
        bottom: 12,
        left: 24,
        zIndex: 400,
        display: 'flex',
        padding: '16px 24px',
        borderRadius: '4px',
        width: 'fit-content',
        position: 'absolute',
        boxShadow: '2px 2px 6px 0px rgba(0, 0, 0, 0.05)',
        background: (theme) => theme.palette.background.paper,
      }}
    >
      <FormGroup>
        <FormControlLabel
          control={<Checkbox defaultChecked sx={{ p: '4px' }} />}
          label="All"
        />
        <FormControlLabel
          control={<Checkbox sx={{ p: '4px' }} />}
          label="Online (1)"
        />
        <FormControlLabel
          control={<Checkbox sx={{ p: '4px' }} />}
          label="Offline (1)"
        />
        <FormControlLabel
          control={<Checkbox sx={{ p: '4px' }} />}
          label="Needs attention (0)"
        />
      </FormGroup>
    </Box>
  );
};
