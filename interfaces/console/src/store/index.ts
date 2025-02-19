import { create } from 'zustand';

type TActionState = 'none' | 'in_progress' | 'success' | 'failed';

type TRestartNode = {
  id: string;
  nodeState: string;
  actionState: TActionState;
};
interface AppState {
  restartNode: TRestartNode;
  setRestartNode: () => void;
}

const appStore = create<AppState>((set) => ({
  restartNode: {
    id: '',
    nodeState: 'unknown',
    actionState: 'none',
  },
  setRestartNode: () =>
    set((state) => ({
      restartNode: {
        id: state.restartNode.id,
        nodeState: state.restartNode.nodeState,
        actionState: state.restartNode.actionState,
      },
    })),
}));
