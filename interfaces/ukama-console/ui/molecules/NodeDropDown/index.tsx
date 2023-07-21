import { Node, NodeStatusEnum } from '@/generated';
import { colors } from '@/styles/theme';
import { hexToRGB } from '@/utils';
import { AddCircleOutlineRounded } from '@mui/icons-material';
import CheckCircleIcon from '@mui/icons-material/CheckCircle';
import CircleIcon from '@mui/icons-material/Circle';
import InfoIcon from '@mui/icons-material/InfoOutlined';
import {
  Button,
  Divider,
  MenuItem,
  Select,
  SelectChangeEvent,
  Stack,
  Typography,
} from '@mui/material';
import LoadingWrapper from '../LoadingWrapper';
import { PaperProps, SelectDisplayProps, useStyles } from './styles';

const getStatus = (status: NodeStatusEnum, time: number) => {
  let str = '';
  switch (status) {
    case NodeStatusEnum.Active:
      str = 'Active';
    case NodeStatusEnum.Maintenance:
      str = 'Maintainance';
    case NodeStatusEnum.Configured:
      str = 'Configured';
    case NodeStatusEnum.Onboarded:
      str = 'Onboarded';
    case NodeStatusEnum.Faulty:
      str = 'Faulty';
    default:
      str = 'Unknown';
  }
  return (
    <Typography variant={'h6'} mr={'6px'}>
      {str}
    </Typography>
  );
};

const getStatusIcon = (status: NodeStatusEnum) => {
  switch (status) {
    case NodeStatusEnum.Active:
      return <CheckCircleIcon htmlColor={colors.green} fontSize={'small'} />;
    case NodeStatusEnum.Maintenance:
      return <InfoIcon htmlColor={colors.yellow} fontSize={'small'} />;
    case NodeStatusEnum.Configured:
      return <InfoIcon htmlColor={colors.black38} fontSize={'small'} />;
    case NodeStatusEnum.Onboarded:
      return <InfoIcon htmlColor={colors.darkGreen05} fontSize={'small'} />;
    case NodeStatusEnum.Faulty:
      return <InfoIcon htmlColor={colors.red} fontSize={'small'} />;
    default:
      return <CircleIcon htmlColor={colors.black38} fontSize={'small'} />;
  }
};

interface INodeDropDown {
  loading: boolean;
  onAddNode: Function;
  nodes: Node[] | [];
  onNodeSelected: Function;
  selectedNode: Node | undefined;
}

const NodeDropDown = ({
  nodes = [],
  onAddNode,
  selectedNode,
  loading = true,
  onNodeSelected,
}: INodeDropDown) => {
  const classes = useStyles();
  const handleChange = (e: SelectChangeEvent<string>) => {
    const { target } = e;
    target.value &&
      onNodeSelected(nodes.find((item: Node) => item.name === target.value));
  };
  return (
    <Stack direction={'row'} spacing={1} alignItems="center">
      {selectedNode && getStatusIcon(selectedNode.status.state)}

      <LoadingWrapper
        height={'fit-content'}
        isLoading={loading}
        width={loading ? '144px' : 'fit-content'}
      >
        <Select
          disableUnderline
          variant="standard"
          onChange={handleChange}
          value={selectedNode?.name}
          SelectDisplayProps={SelectDisplayProps}
          MenuProps={{
            disablePortal: true,
            anchorOrigin: {
              vertical: 'bottom',
              horizontal: 'left',
            },
            transformOrigin: {
              vertical: 'top',
              horizontal: 'left',
            },
            PaperProps: {
              sx: {
                width: '164px',
                ...PaperProps,
              },
            },
          }}
          className={classes.selectStyle}
          renderValue={(selected) => selected}
        >
          {nodes.map(({ id, name }) => (
            <MenuItem
              key={id}
              value={name}
              sx={{
                m: 0,
                p: '6px 16px',
                backgroundColor: `${
                  id === selectedNode?.id
                    ? hexToRGB(colors.secondaryLight, 0.25)
                    : 'inherit'
                } !important`,
                ':hover': {
                  backgroundColor: `${hexToRGB(
                    colors.secondaryLight,
                    0.25,
                  )} !important`,
                },
              }}
            >
              <Typography variant="body1">{name}</Typography>
            </MenuItem>
          ))}
          <Divider />
          <MenuItem
            onClick={(e) => {
              onAddNode();
              e.stopPropagation();
            }}
          >
            <Button
              variant="text"
              sx={{
                p: 0,
                typography: 'body1',
                textTransform: 'none',
              }}
              startIcon={<AddCircleOutlineRounded />}
            >
              Add node
            </Button>
          </MenuItem>
        </Select>
      </LoadingWrapper>

      {/* <LoadingWrapper
        height={38}
        width={'fit-content'}
      >
        {getStatus(nodeStatus.status, nodeStatus.uptime)}
      </LoadingWrapper> */}
    </Stack>
  );
};

export default NodeDropDown;
