/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

import { Grid2, Link, Typography } from "@mui/material";
import CardUI from "../../_components/CardUI";

export interface StatusItem {
  label: string;
  value: string;
}

export interface StatusGridCardData {
  id: string;
  title: string;
  icon: string;
  items: StatusItem[];
  impact: string;
  viewLabel: string;
  viewHref?: string;
  onViewClick?: () => void;
}

const DEFAULT_CARDS: StatusGridCardData[] = [
  {
    id: "health",
    title: "HEALTH",
    icon: "🔧",
    items: [
      { label: "Status", value: "🟢 Healthy" },
      { label: "Compute", value: "🟢 Normal" },
      { label: "Storage", value: "🟢 Normal" },
      { label: "Power", value: "🟢 Normal" },
      { label: "Thermal", value: "🟢 Normal" },
    ],
    impact: "No unit constraints",
    viewLabel: "View Health",
  },
  {
    id: "network",
    title: "NETWORK",
    icon: "🌐",
    items: [
      { label: "Status", value: "🟡 Backhaul Limited" },
      { label: "Cellular", value: "🟢 Healthy" },
      { label: "Backhaul", value: "🟡 Degraded" },
    ],
    impact: "Throughput capped at ~20 Mbps",
    viewLabel: "View Network",
  },
  {
    id: "rf",
    title: "RF",
    icon: "📡",
    items: [
      { label: "Status", value: "🟢 Healthy" },
      { label: "Coverage", value: "operating as expected" },
    ],
    impact: "No RF-related degradation",
    viewLabel: "View RF",
  },
  {
    id: "software",
    title: "SOFTWARE",
    icon: "⚙️",
    items: [
      { label: "Status", value: "🟢 Healthy" },
      { label: "Critical apps", value: "All running" },
    ],
    impact: "No software impact",
    viewLabel: "View Software",
  },
];

const MAX_ITEMS_PER_COLUMN = 3;

export interface OverviewStatusGridProps {
  cards?: StatusGridCardData[];
  /** Max items per column before adding a new column. Default: 3 */
  maxItemsPerColumn?: number;
  selectedMetric: string;
}

function chunkItems<T>(items: T[], chunkSize: number): T[][] {
  const chunks: T[][] = [];
  for (let i = 0; i < items.length; i += chunkSize) {
    chunks.push(items.slice(i, i + chunkSize));
  }
  return chunks;
}

function StatusGridCard({
  card,
  maxItemsPerColumn = MAX_ITEMS_PER_COLUMN,
}: Readonly<{
  card: StatusGridCardData;
  maxItemsPerColumn?: number;
}>) {
  const handleClick = (e: React.MouseEvent) => {
    if (card.onViewClick) {
      e.preventDefault();
      card.onViewClick();
    }
  };

  const itemColumns = chunkItems(card.items, maxItemsPerColumn);

  return (
    <CardUI>
      <Typography
        variant="overline"
        sx={{
          display: "flex",
          fontWeight: 700,
          alignItems: "center",
        }}
      >
        {card.title}
      </Typography>
      <Grid2 container spacing={2} rowSpacing={2} columnSpacing={2}> 
        {itemColumns.map((columnItems) => (
          <Grid2 size={{ xs: 12, sm: 6 }} key={columnItems.map((i) => i.label).join("-")}>
            {columnItems.map((item) => (
              <Typography
                key={item.label}
                variant="body2"
                sx={{
                  fontFamily: "monospace",
                  fontSize: "0.8125rem",
                  mb: 0.5,
                  "&:last-of-type": { mb: 0 },
                }}
              >
                {item.label}: {item.value}  
              </Typography>
            ))}
          </Grid2>
        ))}
      </Grid2>
      <Typography
        variant="body2"
        sx={{
          fontFamily: "monospace",
          fontSize: "0.75rem",
          mt: 1.5,
          mb: 1,
        }}
      >
        Impact: {card.impact}
      </Typography>
      <Link
        href={card.viewHref ?? "#"}
        onClick={handleClick}
        sx={{
          fontFamily: "monospace",
          fontSize: "0.75rem",
          textDecoration: "none",
          "&:hover": { textDecoration: "underline" },
          mt: "auto",
          alignSelf: "flex-start",
        }}
      >
        [ {card.viewLabel} → ]
      </Link>
    </CardUI>
  );
}

export default function OverviewStatusGrid({
  selectedMetric,
  cards = DEFAULT_CARDS,
  maxItemsPerColumn = MAX_ITEMS_PER_COLUMN,
}: Readonly<OverviewStatusGridProps>) {
  return (
    <Grid2
      container
      rowSpacing={2}
      columnSpacing={2}
    >
      {cards.map((card) => (
        card.id === selectedMetric && <Grid2 size={12} key={card.id} >
          <StatusGridCard
            card={card}
            maxItemsPerColumn={maxItemsPerColumn}
          />
        </Grid2>
      ))}
    </Grid2>
  );
}
