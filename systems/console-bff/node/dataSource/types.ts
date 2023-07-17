interface NodeId {
  id: string;
}

interface Node {
  allocated: boolean;
  attached: string[];
  name: string;
  network: string;
  node: string;
  state: string;
  type: string;
}

interface GetNode {
  node: Node;
}

interface GetNodes {
  nodes: Node[];
}

interface DeleteNode {
  node: string;
}

interface AttachNodeArgs {
  anodel: string;
  anoder: string;
  parentNode: string;
}

interface AddNodeArgs {
  id: string;
  state: string;
}

interface AddNodeToNetworkArgs {
  nodeId: string;
  networkId: string;
}

interface UpdateNodeStateArgs {
  id: string;
  state: string;
}

interface UpdateNodeState {
  id: string;
  state: string;
}

interface UpdateNodeArgs {
  id: string;
  name: string;
}
