/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */
'use client';

/**
 * Generic data table — TanStack Table (headless) rendered with the
 * prototype's `.tbl` markup. Supports sortable headers, column-header
 * dropdown filters, optional row selection and row click, with built-in
 * loading / empty / error states (BUILD-PLAN §3, design finding #9).
 */
import { useState } from 'react';
import Menu from '@mui/material/Menu';
import MenuItem from '@mui/material/MenuItem';
import ArrowDownwardRounded from '@mui/icons-material/ArrowDownwardRounded';
import ArrowUpwardRounded from '@mui/icons-material/ArrowUpwardRounded';
import ExpandMoreRounded from '@mui/icons-material/ExpandMoreRounded';
import {
  flexRender,
  getCoreRowModel,
  getFilteredRowModel,
  getSortedRowModel,
  useReactTable,
} from '@tanstack/react-table';
import type {
  ColumnDef,
  ColumnFiltersState,
  Header,
  OnChangeFn,
  RowData,
  RowSelectionState,
  SortingState,
} from '@tanstack/react-table';
import type { EmptyArtName } from '@/components/EmptyState';
import { EmptyState, ErrorState } from '@/components/EmptyState';
import SkeletonTable from './SkeletonTable';

/* eslint-disable @typescript-eslint/no-unused-vars */
declare module '@tanstack/react-table' {
  interface ColumnMeta<TData extends RowData, TValue> {
    /** right-align + tabular numerals */
    num?: boolean;
    /** renders a column-header dropdown filter with these options */
    filterOptions?: string[];
  }
}
/* eslint-enable @typescript-eslint/no-unused-vars */

export type TableStatus = 'loading' | 'error' | 'ready';

export interface DataTableEmptyProps {
  art?: EmptyArtName;
  title: string;
  sub?: string;
  cta?: string;
  onCta?: () => void;
}

function FilterTh<T>({
  header,
  options,
}: {
  header: Header<T, unknown>;
  options: string[];
}) {
  const [anchor, setAnchor] = useState<HTMLElement | null>(null);
  const col = header.column;
  const value = (col.getFilterValue() as string | undefined) ?? 'all';
  const meta = col.columnDef.meta;

  return (
    <th className={meta?.num ? 'num' : undefined}>
      <button
        type="button"
        className={`colfilter${value !== 'all' ? ' on' : ''}`}
        onClick={(e) => setAnchor(e.currentTarget)}
      >
        {flexRender(col.columnDef.header, header.getContext())}
        <ExpandMoreRounded sx={{ fontSize: 16 }} />
      </button>
      <Menu
        anchorEl={anchor}
        open={!!anchor}
        onClose={() => setAnchor(null)}
        slotProps={{ paper: { sx: { width: 180, mt: 0.5 } } }}
      >
        <MenuItem
          sx={{ fontSize: 13.5 }}
          onClick={() => {
            col.setFilterValue(undefined);
            setAnchor(null);
          }}
        >
          All
        </MenuItem>
        {options.map((o) => (
          <MenuItem
            key={o}
            sx={{ fontSize: 13.5, textTransform: 'capitalize' }}
            selected={o === value}
            onClick={() => {
              col.setFilterValue(o);
              setAnchor(null);
            }}
          >
            {o}
          </MenuItem>
        ))}
      </Menu>
    </th>
  );
}

export default function DataTable<T>({
  columns,
  data,
  status = 'ready',
  empty,
  onRetry,
  onRowClick,
  globalFilter,
  enableRowSelection = false,
  rowSelection,
  onRowSelectionChange,
  getRowId,
  initialSorting,
  skeleton,
}: {
  columns: ColumnDef<T, unknown>[];
  data: T[];
  status?: TableStatus;
  empty: DataTableEmptyProps;
  onRetry?: () => void;
  onRowClick?: (row: T) => void;
  globalFilter?: string;
  enableRowSelection?: boolean;
  rowSelection?: RowSelectionState;
  onRowSelectionChange?: OnChangeFn<RowSelectionState>;
  getRowId?: (row: T) => string;
  initialSorting?: SortingState;
  skeleton?: { cols?: number; rows?: number; lead?: boolean };
}) {
  const [sorting, setSorting] = useState<SortingState>(initialSorting ?? []);
  const [columnFilters, setColumnFilters] = useState<ColumnFiltersState>([]);

  const table = useReactTable({
    data,
    columns,
    state: {
      sorting,
      columnFilters,
      globalFilter: globalFilter ?? '',
      ...(rowSelection ? { rowSelection } : {}),
    },
    defaultColumn: { enableSorting: false },
    enableRowSelection,
    onSortingChange: setSorting,
    onColumnFiltersChange: setColumnFilters,
    ...(onRowSelectionChange ? { onRowSelectionChange } : {}),
    ...(getRowId ? { getRowId } : {}),
    globalFilterFn: 'includesString',
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
  });

  if (status === 'loading') {
    return (
      <SkeletonTable
        cols={skeleton?.cols ?? columns.length}
        rows={skeleton?.rows ?? 6}
        lead={skeleton?.lead}
      />
    );
  }
  if (status === 'error') {
    return <ErrorState onRetry={onRetry} />;
  }

  const rows = table.getRowModel().rows;
  if (rows.length === 0) {
    return <EmptyState {...empty} />;
  }

  return (
    <table className="tbl">
      <thead>
        {table.getHeaderGroups().map((hg) => (
          <tr key={hg.id} className="static">
            {hg.headers.map((h) => {
              const col = h.column;
              const meta = col.columnDef.meta;
              if (meta?.filterOptions) {
                return (
                  <FilterTh key={h.id} header={h} options={meta.filterOptions} />
                );
              }
              const sortable = col.getCanSort();
              const sorted = col.getIsSorted();
              const width =
                col.columnDef.size != null ? col.getSize() : undefined;
              return (
                <th
                  key={h.id}
                  className={
                    [sortable ? 'sortable' : '', meta?.num ? 'num' : '']
                      .filter(Boolean)
                      .join(' ') || undefined
                  }
                  style={width != null ? { width } : undefined}
                  onClick={
                    sortable ? col.getToggleSortingHandler() : undefined
                  }
                >
                  <span
                    style={{ display: 'inline-flex', alignItems: 'center', gap: 3 }}
                  >
                    {h.isPlaceholder
                      ? null
                      : flexRender(col.columnDef.header, h.getContext())}
                    {sorted === 'asc' && (
                      <ArrowUpwardRounded sx={{ fontSize: 14 }} />
                    )}
                    {sorted === 'desc' && (
                      <ArrowDownwardRounded sx={{ fontSize: 14 }} />
                    )}
                  </span>
                </th>
              );
            })}
          </tr>
        ))}
      </thead>
      <tbody>
        {rows.map((row) => {
          const clickable = !!onRowClick;
          return (
            <tr
              key={row.id}
              className={clickable ? undefined : 'static'}
              {...(clickable
                ? {
                    role: 'button',
                    tabIndex: 0,
                    onClick: () => onRowClick(row.original),
                    onKeyDown: (e: React.KeyboardEvent) => {
                      if (
                        (e.key === 'Enter' || e.key === ' ') &&
                        e.target === e.currentTarget
                      ) {
                        e.preventDefault();
                        onRowClick(row.original);
                      }
                    },
                  }
                : {})}
            >
              {row.getVisibleCells().map((cell) => {
                const meta = cell.column.columnDef.meta;
                return (
                  <td key={cell.id} className={meta?.num ? 'num' : undefined}>
                    {flexRender(cell.column.columnDef.cell, cell.getContext())}
                  </td>
                );
              })}
            </tr>
          );
        })}
      </tbody>
    </table>
  );
}

/** Leading checkbox column for multi-select tables (agent customers). */
export function selectionColumn<T>(): ColumnDef<T, unknown> {
  return {
    id: 'select',
    size: 34,
    header: ({ table }) => (
      <input
        type="checkbox"
        aria-label="Select all rows"
        checked={table.getIsAllRowsSelected()}
        onChange={table.getToggleAllRowsSelectedHandler()}
      />
    ),
    cell: ({ row }) => (
      <input
        type="checkbox"
        aria-label="Select row"
        checked={row.getIsSelected()}
        onChange={row.getToggleSelectedHandler()}
        onClick={(e) => e.stopPropagation()}
      />
    ),
  };
}
