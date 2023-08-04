import { GetNodes, Node } from "../resolvers/types";

export const parseNodesRes = (res: any): GetNodes => {
  const nodes = res.node.map((node: any) => {
    return parseNodeRes(node);
  });
  return {
    nodes: nodes,
  };
};

export const parseNodeRes = (res: any): Node => {
  return {
    id: res.id,
    name: res.name,
    type: res.type,
    orgId: res.org_id,
    // attached: parseAttachedNodeRes(res.attached),
    status: {
      state: res.status.state,
      connectivity: res.status.connectivity,
    },
  };
};

const parseAttachedNodeRes = (res: any): [Node] => {
  const nodes = res.map((node: any) => {
    return parseNodeRes(node);
  });
  return nodes;
};
