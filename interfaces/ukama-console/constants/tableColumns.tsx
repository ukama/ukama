import { ColumnsWithOptions } from '@/types';

const DataTableWithOptionColumns: ColumnsWithOptions[] = [
  { id: 'name', label: 'Name', minWidth: 150 },
  { id: 'dataUsage', label: 'Usage', minWidth: 100 },
  {
    label: '',
    minWidth: 50,
    id: 'actions',
    align: 'right',
  },
];

const CurrentBillColumns: ColumnsWithOptions[] = [
  {
    id: 'name',
    label: 'Name',
  },
  {
    id: 'rate',
    label: 'Data used',
  },
  {
    id: 'subtotal',
    label: 'Subtotal',
  },
];
const HistoryBillingColumns: ColumnsWithOptions[] = [
  {
    id: 'date',
    label: 'Date',
  },
  {
    id: 'usage',
    label: 'Data Usage',
  },
  {
    id: 'total',
    label: 'Total',
  },
  {
    id: 'pdf',
    label: '',
  },
];
const NodeAppsColumns = [
  {
    id: 'version',
    label: 'Version',
    minWidth: 200,
  },
  {
    id: 'date',
    label: 'Date',
    minWidth: 200,
  },
  {
    id: 'notes',
    label: 'Notes',
    minWidth: 600,
  },
];

export {
  CurrentBillColumns,
  DataTableWithOptionColumns,
  HistoryBillingColumns,
  NodeAppsColumns,
};
