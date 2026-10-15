/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 * Offline views of observed samples. No network requests or verdict changes.
 */
(() => {
  'use strict';
  const payload = JSON.parse(document.getElementById('report-data').textContent);
  const report = payload.report;
  const services = report.services || [];
  const containers = report.containers || [];
  const nodes = report.nodes || [];
  const scenarios = report.scenarios || [];
  const colors = ['var(--blue)', 'var(--orange)', 'var(--green)',
    'var(--purple)', 'var(--red)', 'var(--muted)'];
  const $ = id => document.getElementById(id);
  const finite = value => typeof value === 'number' && Number.isFinite(value);
  const fmt = (value, digits = 1) => finite(value)
    ? value.toLocaleString('en-US', {maximumFractionDigits: digits}) : 'N/A';
  const iso = value => finite(value)
    ? new Date(value * 1000).toISOString().replace('.000', '') : 'N/A';
  const stat = (item, metric, phase = 'during', key = 'peak') =>
    item?.statistics?.[metric]?.[phase]?.[key] ?? null;
  const points = (item, metric) => {
    const values = item?.series?.[metric];
    return Array.isArray(values) ? values.filter(p =>
      Array.isArray(p) && finite(p[0]) && finite(p[1])) : [];
  };
  const mesh = item => (item.service || '').split('/').slice(2).join('/')
    .startsWith('ukama-mesh-node-');
  const core = services.filter(item => !mesh(item));
  const meshes = services.filter(mesh);
  const step = report.config?.step || 30;
  const charts = new Map();
  let selectedService = null;
  let pinned = false;

  function el(tag, text, className) {
    const node = document.createElement(tag);
    if (text !== undefined) node.textContent = text;
    if (className) node.className = className;
    return node;
  }

  function button(text, action, className) {
    const node = el('button', text, className);
    node.type = 'button';
    node.addEventListener('click', action);
    return node;
  }

  function option(select, label, value) {
    const node = el('option', label);
    node.value = value;
    select.append(node);
  }

  function aggregate(items, metric, mode = 'sum', resource = false) {
    const values = new Map();
    for (const item of items) {
      for (const [t, v] of points(resource ? item.resources : item, metric)) {
        if (!values.has(t)) values.set(t, v);
        else values.set(t, mode === 'max'
          ? Math.max(values.get(t), v) : values.get(t) + v);
      }
    }
    return [...values].sort((a, b) => a[0] - b[0]);
  }

  function sampleCount(items, metric) {
    const values = new Map();
    for (const item of items) {
      for (const [t] of points(item, metric)) {
        values.set(t, (values.get(t) || 0) + 1);
      }
    }
    return [...values].sort((a, b) => a[0] - b[0]);
  }

  function meshPodCount() {
    const owned = new Set(meshes.map(s => s.service));
    const values = new Map();
    for (const item of containers.filter(c => owned.has(c.service))) {
      for (const [t] of points(item, 'memory_working_mib')) {
        if (!values.has(t)) values.set(t, new Set());
        values.get(t).add(item.namespace + '/' + item.pod);
      }
    }
    return [...values].map(([t, pods]) => [t, pods.size])
      .sort((a, b) => a[0] - b[0]);
  }

  function activity(data) {
    const bounds = data.windows;
    return [[bounds.before[0], 0], ...(data.activity || []),
      [bounds.after[1], (data.activity || []).at(-1)?.[1] || 0]];
  }

  function line(label, values, color = colors[0], extra = {}) {
    return {label, values, color, gap: step * 1.5, ...extra};
  }

  function bounds() {
    const phase = $('time-window').value;
    return phase === 'all'
      ? [report.windows.before[0], report.windows.after[1]]
      : report.windows[phase];
  }

  function svgNode(tag, attrs = {}, text) {
    const node = document.createElementNS('http://www.w3.org/2000/svg', tag);
    for (const [key, value] of Object.entries(attrs)) {
      node.setAttribute(key, String(value));
    }
    if (text !== undefined) node.textContent = text;
    return node;
  }

  function nearest(values, t, hold = false) {
    let lo = 0, hi = values.length;
    while (lo < hi) {
      const mid = (lo + hi) >> 1;
      if (values[mid][0] <= t) lo = mid + 1;
      else hi = mid;
    }
    if (hold) return values[lo - 1];
    const a = values[lo - 1], b = values[lo];
    if (!a) return b;
    if (!b) return a;
    return t - a[0] <= b[0] - t ? a : b;
  }

  function makeChart(parent, title, unit, lines, options = {}) {
    const figure = el('figure', undefined, 'chart');
    figure.append(el('h3', title));
    const legend = el('div', undefined, 'legend');
    for (const item of lines) {
      const label = el('span', item.label);
      const swatch = el('i', undefined, 'swatch');
      swatch.style.backgroundColor = item.color;
      label.prepend(swatch);
      legend.append(label);
    }
    const svg = svgNode('svg', {role: 'img', 'aria-label': title + ', ' + unit});
    const readout = el('div', 'No sample selected.', 'readout');
    figure.append(legend, svg, readout);
    parent.append(figure);
    const chart = {figure, svg, readout, unit, lines, ...options};
    charts.set(figure, chart);
    chart.observer = new ResizeObserver(() => drawChart(chart));
    chart.observer.observe(figure);
    drawChart(chart);
    return chart;
  }

  function clearCharts(id) {
    const parent = $(id);
    for (const figure of parent.querySelectorAll('.chart')) {
      charts.get(figure)?.observer.disconnect();
      charts.delete(figure);
    }
    parent.replaceChildren();
    return parent;
  }

  function drawChart(chart) {
    if (!chart.figure.getClientRects().length) return;
    const svg = chart.svg;
    const width = Math.max(240, Math.floor(svg.getBoundingClientRect().width));
    const height = 190, left = 62, right = width - 14;
    const top = 18, bottom = height - 39;
    const [from, until] = chart.domain || bounds();
    const to = Math.max(from + 1, until);
    const x = t => left + (t - from) / (to - from) * (right - left);
    const samples = chart.lines.flatMap(l => l.values
      .filter(p => p[0] >= from && p[0] <= to).map(p => p[1]));
    const maxValue = Math.max(1, chart.reference || 0,
      samples.reduce((a, v) => Math.max(a, v), 0)) * 1.08;
    const y = value => bottom - value / maxValue * (bottom - top);
    chart.geometry = {from, to, x, y, left, right, top, bottom};
    svg.replaceChildren();
    svg.setAttribute('viewBox', `0 0 ${width} ${height}`);
    svg.setAttribute('height', height);
    if (!chart.elapsed) {
      const [start, end] = report.windows.during;
      const a = Math.max(from, start), b = Math.min(to, end);
      if (b > a) svg.append(svgNode('rect', {x: x(a), y: top,
        width: x(b) - x(a), height: bottom - top, fill: 'var(--shade)'}));
    }
    for (let i = 0; i <= 3; i++) {
      const value = maxValue * i / 3, yy = y(value);
      svg.append(svgNode('line', {x1: left, x2: right, y1: yy, y2: yy,
        stroke: 'var(--line)', 'stroke-width': 0.6}));
      svg.append(svgNode('text', {x: left - 8, y: yy + 4,
        'text-anchor': 'end'}, fmt(value, maxValue < 10 ? 1 : 0)));
    }
    const tickCount = width < 400 ? 2 : 4;
    for (let i = 0; i <= tickCount; i++) {
      const t = from + (to - from) * i / tickCount;
      const label = chart.elapsed ? fmt(t / 60, 1) : iso(t).slice(11, 19);
      svg.append(svgNode('text', {x: x(t), y: bottom + 18,
        'text-anchor': i === 0 ? 'start' : i === tickCount ? 'end' : 'middle'}, label));
    }
    svg.append(svgNode('text', {x: left, y: 11}, chart.unit));
    svg.append(svgNode('text', {x: right, y: height - 2,
      'text-anchor': 'end'}, chart.elapsed ? 'Minutes from execution start' : 'UTC'));
    if (chart.reference) svg.append(svgNode('line', {
      x1: left, x2: right, y1: y(chart.reference), y2: y(chart.reference),
      stroke: 'var(--red)', 'stroke-dasharray': '4 4'}));
    for (const item of chart.lines) {
      let path = '', previous = null, count = 0;
      const visible = item.values.filter(p => p[0] >= from && p[0] <= to);
      if (item.hold) {
        const first = nearest(item.values, from, true);
        const last = nearest(item.values, to, true);
        if (first && (!visible.length || visible[0][0] > from)) {
          visible.unshift([from, first[1]]);
        }
        if (last && to <= item.values.at(-1)[0] &&
            (!visible.length || visible.at(-1)[0] < to)) {
          visible.push([to, last[1]]);
        }
      }
      for (const [t, value] of visible) {
        const joined = previous !== null && (item.hold || t - previous <= item.gap);
        const xx = x(t).toFixed(2), yy = y(value).toFixed(2);
        path += joined ? (item.hold ? `H${xx}V${yy}` : `L${xx},${yy}`)
          : `M${xx},${yy}`;
        // A marker keeps isolated samples visible without filling missing gaps.
        svg.append(svgNode('circle', {cx: xx, cy: yy, r: 1.5, fill: item.color}));
        previous = t;
        count++;
      }
      if (count) svg.append(svgNode('path', {d: path, fill: 'none',
        stroke: item.color, 'stroke-width': 1.8}));
    }
    if (!samples.length) svg.append(svgNode('text', {
      x: (left + right) / 2, y: (top + bottom) / 2,
      'text-anchor': 'middle'}, 'N/A — no samples in this window'));
    if (!chart.elapsed) for (const s of scenarios) {
      if (s.outcome !== 'FAIL' || !finite(s.ended_at) ||
          s.ended_at < from || s.ended_at > to) continue;
      svg.append(svgNode('line', {x1: x(s.ended_at), x2: x(s.ended_at),
        y1: top, y2: top + 7, stroke: 'var(--red)'}));
    }
    const guide = svgNode('line', {class: 'guide', x1: left, x2: left,
      y1: top, y2: bottom, stroke: 'var(--muted)', visibility: 'hidden'});
    svg.append(guide);
    const hit = svgNode('rect', {x: left, y: top, width: right - left,
      height: bottom - top, fill: 'transparent'});
    function move(event, pin = false) {
      if (pinned && !pin) return;
      const rect = svg.getBoundingClientRect();
      const px = (event.clientX - rect.left) * width / rect.width;
      const t = from + Math.max(0, Math.min(1,
        (px - left) / (right - left))) * (to - from);
      if (pin) pinned = !pinned;
      for (const other of charts.values()) {
        if (Boolean(other.elapsed) === Boolean(chart.elapsed)) cursor(other, t);
      }
      const active = scenarios.filter(s => s.outcome !== 'SKIP' &&
        finite(s.started_at) && s.started_at <= t && s.ended_at > t);
      $('cursor-info').textContent = chart.elapsed
        ? `${fmt(t / 60)} minutes from execution start`
        : `${iso(t)} · ${active.length} active scenarios` +
          (pinned ? ' · pinned (tap chart to release)' : '');
    }
    hit.addEventListener('pointermove', event => move(event));
    hit.addEventListener('click', event => move(event, true));
    svg.append(hit);
  }

  function cursor(chart, t) {
    const g = chart.geometry;
    if (!g || t < g.from || t > g.to) return;
    const guide = chart.svg.querySelector('.guide');
    if (!guide) return;
    guide.setAttribute('x1', g.x(t));
    guide.setAttribute('x2', g.x(t));
    guide.setAttribute('visibility', 'visible');
    const parts = chart.lines.map(item => {
      const p = nearest(item.values, t, item.hold);
      const known = p && (item.hold
        ? t <= item.values.at(-1)[0]
        : Math.abs(t - p[0]) <= item.gap / 2);
      const sampleTime = chart.elapsed ? `${fmt((p?.[0] || 0) / 60)}m`
        : iso(p?.[0]).slice(11, 19);
      return `${item.label}: ${known ? fmt(p[1]) : 'N/A'}` +
        (known && !item.hold ? ` (sample ${sampleTime})` : '');
    });
    chart.readout.textContent = parts.join(' · ');
  }

  function table(id, headers, rows, pageSize = 20) {
    const host = $(id);
    host.replaceChildren();
    let page = 0;
    function render() {
      host.replaceChildren();
      const wrap = el('div', undefined, 'table-wrap');
      const grid = el('table'), head = el('thead'), tr = el('tr');
      headers.forEach(h => tr.append(el('th', h)));
      head.append(tr);
      grid.append(head);
      const body = el('tbody');
      for (const row of rows.slice(page * pageSize, (page + 1) * pageSize)) {
        const tr = el('tr');
        for (const value of row) {
          const td = el('td');
          if (typeof value === 'function') td.append(value());
          else if (value && typeof value === 'object') {
            td.textContent = value.text;
            td.className = value.className || '';
          } else td.textContent = value ?? 'N/A';
          tr.append(td);
        }
        body.append(tr);
      }
      grid.append(body);
      wrap.append(grid);
      host.append(wrap);
      if (!rows.length) host.append(el('p', 'No matching observations.', 'muted'));
      if (rows.length > pageSize) {
        const pages = Math.ceil(rows.length / pageSize);
        const controls = el('div', undefined, 'pagination');
        const prev = button('Previous', () => { page--; render(); });
        const next = button('Next', () => { page++; render(); });
        prev.disabled = page === 0;
        next.disabled = page + 1 === pages;
        controls.append(prev, el('span',
          `${page * pageSize + 1}–${Math.min((page + 1) * pageSize, rows.length)} of ${rows.length}`), next);
        host.append(controls);
      }
    }
    render();
  }

  function heat(value, warning = 70, danger = 90) {
    return {text: fmt(value), className: !finite(value) ? 'unknown'
      : value >= danger ? 'hot' : value >= warning ? 'warm' : ''};
  }

  function showTab(name) {
    document.querySelectorAll('.view').forEach(v => { v.hidden = v.id !== name; });
    document.querySelectorAll('[data-tab]').forEach(b => {
      b.setAttribute('aria-pressed', String(b.dataset.tab === name));
    });
    $('time-toolbar').hidden = name === 'compare' || name === 'coverage';
    pinned = false;
    $('cursor-info').textContent = '';
    for (const chart of charts.values()) drawChart(chart);
  }

  function overview() {
    const result = report.results || {};
    const info = [
      ['Executed pass rate', `${fmt(report.pass_rate)}%`,
        `${result.PASS || 0} pass / ${result.FAIL || 0} fail / ${result.SKIP || 0} skip`],
      ['Workers launched', fmt(report.run?.workers_launched, 0),
        `${fmt(report.peak_concurrency, 0)} peak concurrent scenarios`],
      ['Execution', `${fmt(report.wall_seconds / 60)} min`,
        `${report.selected_count ?? scenarios.length} scenarios selected`],
      ['Cleanup warnings / failures', String(report.cleanup_warnings || 0),
        report.recovery_complete ? 'Recovery window complete' : 'Recovery window partial']
    ];
    for (const [label, value, detail] of info) {
      const block = el('div', undefined, 'stat');
      block.append(el('span', label), el('strong', value), el('span', detail));
      $('summary').append(block);
    }
    $('batch-label').textContent = `${report.batch} · ${report.run?.suite || 'suite unknown'} · ` +
      `${iso(report.windows.during[0])} — ${iso(report.windows.during[1])}`;
    const c = payload.coverage;
    const q = Object.values(report.queries || {});
    $('coverage-banner').textContent =
      `${q.filter(x => x.status === 'OK').length}/${q.length} queries have samples. ` +
      `Usable CPU limits: ${c.cpu_limit_pct}/${c.containers} observed containers; ` +
      `memory limits: ${c.memory_limit_pct}/${c.containers}. ` +
      `${core.length} other workload groups + ${meshes.length} mesh-node groups. ` +
      'Unknown coverage is not healthy coverage.';
    const allFindings = payload.findings || [];
    const findings = allFindings.filter(row => !mesh(row));
    findings.slice(0, 4).forEach(row => {
      const item = el('div', undefined, 'finding');
      item.append(button(row.service, () => selectService(row.service), 'link'),
        el('p', row.message), el('small', row.context));
      $('findings').append(item);
    });
    if (!findings.length) $('findings').append(el('p',
      'No configured markers triggered in observed data. Review missing coverage.'));
    if (findings.length > 4) $('findings').append(button(
      `${findings.length - 4} more findings — browse services`, () => showTab('services')));
    const meshFindings = allFindings.length - findings.length;
    if (meshFindings) $('findings').append(button(
      `${meshFindings} findings in mesh workloads — inspect group`,
      () => showTab('mesh')));
    const parent = $('overview-charts');
    makeChart(parent, 'Scenario concurrency', 'Active scenarios',
      [line('Active', activity(report), colors[0], {hold: true})]);
    makeChart(parent, 'CPU by workload group', 'mCPU', [
      line('Other workloads', aggregate(core, 'cpu_mcpu')),
      line('Mesh nodes', aggregate(meshes, 'cpu_mcpu'), colors[1])]);
    makeChart(parent, 'Memory by workload group', 'MiB', [
      line('Other workloads', aggregate(core, 'memory_mib')),
      line('Mesh nodes', aggregate(meshes, 'memory_mib'), colors[1])]);
    makeChart(parent, 'Worst container throttling', 'Throttled periods (%)',
      [line('Maximum observed', aggregate(services, 'cpu_throttle_pct', 'max', true))],
      {reference: 100});
    makeChart(parent, 'Node CPU pressure', '%', [
      line('Highest CPU busy', aggregate(nodes, 'cpu_busy_pct', 'max')),
      line('Highest CPU PSI some', aggregate(nodes, 'cpu_psi_pct', 'max'), colors[1])],
      {reference: 100});
    makeChart(parent, 'Node memory pressure', '%', [
      line('Highest memory unavailable', aggregate(nodes, 'memory_used_pct', 'max')),
      line('Highest memory PSI some', aggregate(nodes, 'memory_psi_pct', 'max'), colors[1])],
      {reference: 100});
    makeChart(parent, 'Workload sample coverage', 'Groups with memory samples', [
      line('Other workloads', sampleCount(core, 'memory_mib')),
      line('Mesh nodes', sampleCount(meshes, 'memory_mib'), colors[1])]);
    const changes = new Map();
    for (const event of report.events || []) {
      const match = /^restart increase ([\d.]+)/.exec(event.event);
      if (match) changes.set(event.time, (changes.get(event.time) || 0) + Number(match[1]));
    }
    let count = 0;
    const restarts = [[report.windows.during[0], 0],
      ...[...changes].sort((a,b) => a[0]-b[0]).map(([t,v]) => [t, count += v]),
      [report.windows.during[1], count]];
    const known = services.some(s => finite(s.restart_increases));
    makeChart(parent, 'Observed restart increases', 'Cumulative lower bound',
      [line('During execution', known ? restarts : [], colors[4], {hold: true})]);
  }

  function renderServices() {
    const query = $('service-search').value.toLowerCase();
    const source = $('include-mesh').checked ? services : core;
    const rows = source.filter(s => s.service.toLowerCase().includes(query));
    const scores = new Map();
    for (const f of payload.findings || []) {
      scores.set(f.service, Math.max(scores.get(f.service) || 0, 10-f.priority));
    }
    const key = $('service-sort').value;
    const score = s => ({attention: scores.get(s.service) || 0,
      cpu: stat(s, 'cpu_mcpu'), memory: stat(s, 'memory_mib'),
      throttle: stat(s.resources, 'cpu_throttle_pct'),
      memoryLimit: stat(s.resources, 'memory_limit_pct')})[key] ?? -Infinity;
    rows.sort((a,b) => key === 'name' ? a.service.localeCompare(b.service)
      : score(b)-score(a) || a.service.localeCompare(b.service));
    table('service-table', ['Workload', 'CPU mCPU', 'Memory MiB',
      'CPU / limit %', 'Throttled periods %', 'Memory / limit %',
      'Restarts ≥', 'OOM events / kills ≥', 'Replica deficit s'], rows.map(s => [
      () => button(s.service, () => selectService(s.service), 'link'),
      fmt(stat(s, 'cpu_mcpu')), fmt(stat(s, 'memory_mib')),
      heat(stat(s.resources, 'cpu_limit_pct')),
      heat(stat(s.resources, 'cpu_throttle_pct'), 10, 50),
      heat(stat(s.resources, 'memory_limit_pct')),
      heat(s.restart_increases, 1, 3),
      `${fmt(s.resources?.totals?.oom_events?.during, 0)} / ` +
        fmt(s.resources?.totals?.oom_kills?.during, 0),
      heat(s.degraded_sample_seconds, 1, 300)]));
  }

  function selectService(identity) {
    selectedService = services.find(s => s.service === identity);
    if (!selectedService) return;
    $('service-detail').hidden = false;
    $('service-title').textContent = identity;
    const owned = containers.filter(c => c.service === identity);
    $('service-description').textContent =
      `${selectedService.pods?.length || 0} observed pods; ${owned.length} observed containers. ` +
      `Observed restarts ≥ ${fmt(selectedService.restart_increases, 0)}. ` +
      'Service totals can change with replica count. Percentages show the worst container.';
    $('container-select').replaceChildren();
    option($('container-select'), 'Whole workload', '');
    owned.forEach(c => option($('container-select'), c.identity, c.identity));
    renderServiceCharts();
    showTab('services');
    $('service-detail').scrollIntoView({block: 'start', behavior: 'instant'});
  }

  function renderServiceCharts() {
    if (!selectedService) return;
    const selected = $('container-select').value;
    const container = containers.find(c => c.identity === selected);
    const resource = container || selectedService.resources || {};
    const parent = clearCharts('service-charts');
    const cpu = points(container || selectedService, 'cpu_mcpu');
    const memory = points(container || selectedService,
      container ? 'memory_working_mib' : 'memory_mib');
    makeChart(parent, 'CPU usage and configuration', 'mCPU', [
      line('Usage', cpu), line('Request', points(resource, 'cpu_request_mcpu'), colors[2]),
      line('Limit', points(resource, 'cpu_limit_mcpu'), colors[4])]);
    makeChart(parent, 'Memory usage and configuration', 'MiB', [
      line('Working set', memory),
      line('Total incl. cache', points(resource, 'memory_usage_mib'), colors[1]),
      line('Limit', points(resource, 'memory_limit_mib'), colors[4])]);
    makeChart(parent, 'CPU throttled periods', '% (not wall time lost)',
      [line('Throttled periods', points(resource, 'cpu_throttle_pct'))], {reference: 100});
    makeChart(parent, 'Memory against limit', '%', [
      line('Total / limit', points(resource, 'memory_limit_pct')),
      line('Working set / limit', points(resource, 'memory_working_pct'), colors[2])],
      {reference: 100});
    makeChart(parent, 'CPU against limit', '%',
      [line('Usage / limit', points(resource, 'cpu_limit_pct'))], {reference: 100});
    makeChart(parent, 'Container pressure', 'PSI (%)', [
      line('CPU some', points(resource, 'cpu_psi_pct')),
      line('Memory some', points(resource, 'memory_psi_pct'), colors[1])], {reference: 100});
    if (!container) {
      makeChart(parent, 'Available and desired replicas', 'Replicas', [
        line('Available', points(selectedService, 'replicas')),
        line('Desired', points(selectedService, 'desired'), colors[2])]);
      makeChart(parent, 'Request activity', 'Requests / second', [
        line('gRPC', points(selectedService, 'grpc_rps')),
        line('DB', points(selectedService, 'db_rps'), colors[1])]);
      makeChart(parent, 'Network activity', 'Bytes / second', [
        line('Receive', points(selectedService, 'rx_bps')),
        line('Transmit', points(selectedService, 'tx_bps'), colors[1])]);
    }
  }

  function renderMesh() {
    $('mesh-summary').textContent = `${meshes.length} mesh-node workload groups observed ` +
      'across the full window. Mesh membership uses the ukama-mesh-node- name prefix.';
    const parent = $('mesh-charts');
    makeChart(parent, 'Sampled mesh pods', 'Pods with memory samples',
      [line('Observed', meshPodCount())]);
    makeChart(parent, 'Mesh CPU', 'mCPU',
      [line('Sum at each timestamp', aggregate(meshes, 'cpu_mcpu'))]);
    makeChart(parent, 'Mesh memory working set', 'MiB',
      [line('Sum at each timestamp', aggregate(meshes, 'memory_mib'))]);
    makeChart(parent, 'Scenario concurrency', 'Active scenarios',
      [line('Active', activity(report), colors[1], {hold: true})]);
  }

  function renderNodes() {
    const exporters = nodes.filter(n => points(n, 'cpu_busy_pct').length ||
      points(n, 'memory_used_pct').length);
    exporters.sort((a,b) => (stat(b, 'cpu_psi_pct') || 0) -
      (stat(a, 'cpu_psi_pct') || 0));
    table('node-table', ['Node / exporter', 'Peak CPU busy %',
      'Peak memory unavailable %', 'Peak CPU PSI %', 'Peak memory PSI %'],
      exporters.map(n => [() => button(n.node || n.identity, () => {
        $('node-select').value = n.identity; renderNodeCharts();
      }, 'link'), heat(stat(n, 'cpu_busy_pct')), heat(stat(n, 'memory_used_pct')),
      heat(stat(n, 'cpu_psi_pct'), 10, 50), heat(stat(n, 'memory_psi_pct'), 1, 10)]));
    nodes.forEach(n => option($('node-select'), `${n.node || n.identity} — ${n.identity}`, n.identity));
    if (exporters.length) $('node-select').value = exporters[0].identity;
    table('pressure-table', ['Kubernetes node', 'MemoryPressure peak (0/1)'],
      nodes.filter(n => stat(n, 'memory_pressure') !== null)
        .map(n => [n.node || n.identity, heat(stat(n, 'memory_pressure'), 1, 1)]));
    renderNodeCharts();
  }

  function renderNodeCharts() {
    const n = nodes.find(n => n.identity === $('node-select').value);
    const parent = clearCharts('node-charts');
    for (const [name, label, unit, ref] of [
      ['cpu_busy_pct', 'CPU busy', '%', 100],
      ['memory_used_pct', 'Memory unavailable', '%', 100],
      ['memory_available_mib', 'Memory available', 'MiB', null],
      ['cpu_psi_pct', 'CPU pressure', 'PSI some (%)', 100],
      ['memory_psi_pct', 'Memory pressure', 'PSI some (%)', 100],
      ['memory_pressure', 'Kubernetes MemoryPressure', 'Condition (0/1)', 1]]) {
      makeChart(parent, label, unit, [line(label, points(n, name))], {reference: ref});
    }
  }

  function renderScenarios() {
    const q = $('scenario-search').value.toLowerCase();
    const outcome = $('scenario-outcome').value;
    const rows = scenarios.filter(s => (!outcome || s.outcome === outcome) &&
      [s.scenario, s.worker, s.failure_reason, s.run_id].join(' ').toLowerCase().includes(q));
    table('scenario-table', ['Scenario', 'Worker', 'Outcome', 'Start UTC',
      'End UTC', 'Cleanup', 'Failure detail'], rows.map(s => [
      s.scenario || s.run_id, s.worker, {text: s.outcome, className: s.outcome?.toLowerCase()},
      iso(s.started_at), iso(s.ended_at), s.cleanup || 'N/A', s.failure_reason || '—']));
  }

  function renderEvents() {
    const query = $('event-search').value.toLowerCase();
    const rows = (report.events || []).filter(e =>
      [e.service, e.event].join(' ').toLowerCase().includes(query));
    table('event-table', ['UTC', 'Workload', 'Event', 'Overlapping scenarios'],
      rows.map(e => [iso(e.time), e.service, e.event, () => {
        const detail = el('details');
        const runs = e.overlapping_runs || [];
        detail.append(el('summary', `${runs.length} overlapping runs`));
        const list = el('ul');
        runs.forEach(r => list.append(el('li', r)));
        detail.append(list);
        return detail;
      }]));
  }

  const comparison = [report];
  function validReport(data) {
    return data && typeof data.batch === 'string' &&
      data.windows && ['before', 'during', 'after'].every(name =>
        Array.isArray(data.windows[name]) && data.windows[name].length === 2 &&
        data.windows[name].every(finite) && data.windows[name][1] >= data.windows[name][0]) &&
      Array.isArray(data.services) && data.services.every(s =>
        s && typeof s.service === 'string' && s.series &&
        typeof s.series === 'object') &&
      (data.activity === undefined || (Array.isArray(data.activity) &&
        data.activity.every(p => Array.isArray(p) && p.length === 2 &&
          p.every(finite)))) &&
      data.results && typeof data.results === 'object';
  }

  function renderComparison() {
    $('compare-list').replaceChildren();
    comparison.forEach((r,i) => {
      const item = el('span', `${r.batch} `);
      if (i) item.append(button('Remove', () => {
        comparison.splice(i,1); renderComparison();
      }));
      $('compare-list').append(item);
    });
    const measurements = [
      ['Workers requested / launched', r => `${r.run?.workers_requested ?? 'N/A'} / ${r.run?.workers_launched ?? 'N/A'}`],
      ['Peak scenario concurrency', r => fmt(r.peak_concurrency, 0)],
      ['PASS / FAIL / SKIP', r => ['PASS','FAIL','SKIP'].map(k => r.results[k] || 0).join(' / ')],
      ['Executed pass rate (%)', r => fmt(r.pass_rate)],
      ['Wall time (min)', r => fmt(r.wall_seconds / 60)],
      ['Cleanup warnings / failures', r => r.cleanup_warnings ?? 'N/A'],
      ['Selected scenarios', r => r.selected_count ?? 'N/A'],
      ['Namespace / service filters', r => `${r.config?.namespace || 'N/A'} / ${r.config?.service || 'N/A'}`],
      ['Sample step (s)', r => r.config?.step || 'N/A'],
      ['Recovery complete', r => r.recovery_complete ? 'Yes' : 'No'],
      ['Source / binary SHA256', r => r.run?.source_sha256 || r.run?.binary_sha256 || 'N/A']
    ];
    table('compare-table', ['Measurement', ...comparison.map(r => r.batch)],
      measurements.map(([name, fn]) => [name, ...comparison.map(fn)]));
    const chosen = $('compare-focus').value;
    $('compare-focus').replaceChildren();
    option($('compare-focus'), 'Other workloads (mesh excluded)', 'core');
    option($('compare-focus'), 'Mesh workloads', 'mesh');
    const names = [...new Set(comparison.flatMap(r => r.services.map(s => s.service)))].sort();
    names.filter(n => !n.split('/').at(-1).startsWith('ukama-mesh-node-'))
      .forEach(n => option($('compare-focus'), n, n));
    if (['core','mesh',...names].includes(chosen)) $('compare-focus').value = chosen;
    renderCompareCharts();
  }

  function renderCompareCharts() {
    const focus = $('compare-focus').value || 'core';
    const metric = $('compare-metric').value;
    const resource = metric.endsWith('_pct');
    const lines = comparison.map((r,i) => {
      const rows = r.services.filter(s => focus === 'core' ? !mesh(s)
        : focus === 'mesh' ? mesh(s) : s.service === focus);
      const values = aggregate(rows, metric, resource ? 'max' : 'sum', resource);
      return line(r.batch, values.map(([t,v]) => [t-r.windows.during[0], v]),
        colors[i], {gap: (r.config?.step || 30)*1.5});
    });
    const domain = [Math.min(...comparison.map(r => r.windows.before[0]-r.windows.during[0])),
      Math.max(...comparison.map(r => r.windows.after[1]-r.windows.during[0]))];
    const parent = clearCharts('compare-charts');
    makeChart(parent, $('compare-metric').selectedOptions[0].textContent,
      resource ? '%' : metric === 'cpu_mcpu' ? 'mCPU' : 'MiB', lines,
      {domain, elapsed: true, reference: resource ? 100 : null});
    makeChart(parent, 'Scenario concurrency', 'Active scenarios',
      comparison.map((r,i) => line(r.batch,
        activity(r).map(([t,v]) => [t-r.windows.during[0],v]), colors[i], {hold: true})),
      {domain, elapsed: true});
  }

  function coverage() {
    const c = payload.coverage;
    $('coverage-detail').textContent =
      `${c.containers} containers with CPU or memory samples during execution. ` +
      `CPU / limit: ${c.cpu_limit_pct}; memory / limit: ${c.memory_limit_pct}; ` +
      `CPU throttled periods: ${c.cpu_throttle_pct}; memory PSI: ${c.memory_psi_pct}.`;
    table('coverage-table', ['Query', 'Status', 'Series', 'Samples', 'Detail'],
      Object.entries(report.queries || {}).map(([name,q]) => [name,
        {text:q.status, className:q.status === 'OK' ? '' : 'unknown'},
        q.series, q.samples, q.error || [...(q.warnings || []), ...(q.infos || [])].join('; ')]), 50);
    (report.notes || []).forEach(note => $('notes').append(el('li', note)));
  }

  document.querySelectorAll('[data-tab]').forEach(b => {
    b.addEventListener('click', () => showTab(b.dataset.tab));
  });
  $('time-window').addEventListener('change', () => {
    pinned = false;
    $('cursor-info').textContent = '';
    for (const chart of charts.values()) drawChart(chart);
  });
  ['service-search', 'service-sort', 'include-mesh'].forEach(id =>
    $(id).addEventListener('input', renderServices));
  $('container-select').addEventListener('change', renderServiceCharts);
  $('node-select').addEventListener('change', renderNodeCharts);
  ['scenario-search', 'scenario-outcome'].forEach(id =>
    $(id).addEventListener('input', renderScenarios));
  $('event-search').addEventListener('input', renderEvents);
  ['compare-focus','compare-metric'].forEach(id =>
    $(id).addEventListener('change', renderCompareCharts));
  $('browse-mesh').addEventListener('click', () => {
    $('include-mesh').checked = true;
    $('service-search').value = 'ukama-mesh-node-';
    renderServices(); showTab('services');
  });
  $('compare-files').addEventListener('change', async event => {
    $('compare-error').textContent = '';
    for (const file of event.target.files) {
      if (comparison.length >= 6) {
        $('compare-error').textContent = 'Remove a run before adding more (six maximum).';
        break;
      }
      try {
        const data = JSON.parse(await file.text());
        if (!validReport(data)) throw Error('expected a metrics/report.json file');
        if (comparison.some(r => r.batch === data.batch)) {
          throw Error('this batch is already loaded');
        }
        comparison.push(data);
      } catch (error) {
        $('compare-error').textContent = `${file.name}: ${error.message}`;
      }
    }
    event.target.value = '';
    renderComparison();
  });
  overview();
  renderServices();
  renderMesh();
  renderNodes();
  makeChart($('scenario-charts'), 'Scenario concurrency', 'Active scenarios',
    [line('Active', activity(report), colors[0], {hold:true})]);
  renderScenarios();
  renderEvents();
  renderComparison();
  coverage();
})();
