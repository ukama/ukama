/* SPDX-License-Identifier: MPL-2.0
 * Copyright (c) 2026-present, Ukama Inc.
 * Read the rendered chart/rail. Never inspect React props or use a service
 * response as an acceptance oracle. Expectations belong to the scenario.
 */
import type {Locator, Page} from 'playwright';
import {Budget, WorkerError, normalize} from './contract.js';

export const ANALYTICS_LABELS = ['Chart tooltip','Chart x axis','Chart y axis','Chart legend','Chart segments','Chart state','Chart range','Metric value','App resource','App identity'];
export const ANALYTICS_ACTIONS = ['chart_hover','chart_range','metric_select'];
export const analyticsView = (view: string) => ['network_node_detail','network_site_detail','business_revenue','business_packages','business_home'].includes(view);
const hovered = new WeakMap<Page,{path:string,title:string,at:number}>();
async function unique(target: Locator): Promise<Locator | null> {
  const visible = target.filter({visible:true}), n = await visible.count();
  if (n > 1) throw new WorkerError('AMBIGUOUS_LOCATOR','Analytics observation matches multiple visible elements');
  return n === 1 ? visible : null;
}
async function value(target: Locator) {
  const found = await unique(target);
  return found ? normalize(await found.innerText({timeout:500})) : null;
}
export class Analytics {
  constructor(private page: Page) {}
  private main() { return this.page.locator('main.main'); }
  private async chart(title: string) {
    // :text-is matches a text-only sec-title as well as a child span.
    const card = this.main().locator('.card').filter({has:this.page.locator('.sec-title').filter({hasText:new RegExp('^'+title.replace(/[.*+?^${}()|[\]\\]/g,'\\$&')+'$')})});
    const found = await unique(card);
    return found;
  }
  async run(action: string, label: string, input: string, budget: Budget) {
    hovered.delete(this.page);
    const click = (target: Locator) => target.click({timeout:budget.remaining()});
    if (action === 'metric_select') {
      if (label === 'tab' && ['Overview','Resources','Software','Connectivity','Radio'].includes(input))
        await click(this.main().getByRole('tab',{name:input,exact:true}));
      else if (label === 'section')
        await click(this.main().locator('.card-selectable[role="button"]').filter({has:this.page.locator('.sec-title').filter({hasText:new RegExp('^'+input.replace(/[.*+?^${}()|[\]\\]/g,'\\$&')+'$')})}));
      else if (label === 'component' && ['Node','Switch','Charge controller','Backhaul','Solar panels','Batteries'].includes(input))
        await click(this.main().locator('button.comp-tile').filter({has:this.page.getByText(input,{exact:true})}));
      else if (label === 'app' && /^[a-zA-Z0-9_.-]{1,80}$/.test(input))
        await click(this.main().getByRole('button',{name:`View ${input} resources`,exact:true}));
      else if (label === 'dialog' && input === 'Close')
        await click(this.page.getByRole('dialog').getByRole('button',{name:'Close',exact:true}));
      else throw new WorkerError('INVALID_INPUT','Unsupported metric selection');
      return;
    }
    const card = await budget.poll(()=>this.chart(label),v=>v!==null,'Chart card is not visible');
    if (!card) throw new WorkerError('UNSUPPORTED_LOCATOR','Missing chart');
    if (action === 'chart_range') {
      if (!['Day','Week','Month'].includes(input)) throw new WorkerError('INVALID_INPUT','Chart range must be Day, Week or Month');
      await click(card.getByRole('group',{name:'Time range',exact:true}).getByRole('button',{name:input,exact:true}));
      return;
    }
    if (action !== 'chart_hover' || !/^(?:0|[1-9]\d?)(?:\.\d+)?$|^100$/.test(input))
      throw new WorkerError('INVALID_INPUT','Chart hover is a percentage from 0 to 100 of the rendered plot');
    const surface = await unique(card.locator('svg.recharts-surface'));
    if (!surface) throw new WorkerError('UNSUPPORTED_LOCATOR','Chart has no unique rendered SVG');
    // Recharts publishes its clipping rectangle in SVG coordinates. This is
    // visible plot geometry, not hidden application data or a tooltip oracle.
    const plot = await surface.evaluate(svg => {
      const rects = svg.querySelectorAll('defs clipPath rect');
      if (rects.length !== 1) return null;
      const rect = rects[0] as SVGRectElement, matrix = rect.getScreenCTM();
      if (!matrix) return null;
      const a = new DOMPoint(rect.x.baseVal.value,rect.y.baseVal.value).matrixTransform(matrix);
      const z = new DOMPoint(rect.x.baseVal.value+rect.width.baseVal.value,rect.y.baseVal.value+rect.height.baseVal.value).matrixTransform(matrix);
      return {x:a.x,y:a.y,width:z.x-a.x,height:z.y-a.y};
    });
    if (!plot || plot.width<=2 || plot.height<=2) throw new WorkerError('UNSUPPORTED_LOCATOR','Chart plot bounds are absent or ambiguous');
    // Move out first: an old tooltip must not satisfy the next sample check.
    await this.page.mouse.move(0,0);
    await this.page.mouse.move(plot.x+Math.max(1,Math.min(plot.width-1,plot.width*Number(input)/100)),plot.y+plot.height/2);
    hovered.set(this.page,{path:this.page.url(),title:label,at:Date.now()});
  }
  async observe(label: string, subject: string): Promise<string | null> {
    if (label === 'Metric value') {
      const row = this.main().locator('.kv-row').filter({has:this.page.getByText(subject,{exact:true})});
      const target = await unique(row);
      return target ? value(target.locator(':scope > span > .tnum')) : null;
    }
    if (label === 'App identity') return value(this.page.getByRole('dialog').getByRole('heading'));
    if (label === 'App resource') {
      const dialog = this.page.getByRole('dialog').filter({visible:true});
      // Anchor to the label's direct parent to avoid including dialog wrappers.
      const key = await unique(dialog.getByText(subject,{exact:true}));
      return key ? value(key.locator('..').locator(':scope > span.tnum')) : null;
    }
    const card = await this.chart(subject);
    if (!card || await card.locator('.MuiSkeleton-root:visible').count()) return null;
    if (label === 'Chart state') {
      const empty = await card.getByText('No data',{exact:true}).isVisible();
      const error = await card.getByText("Couldn't load metric",{exact:true}).isVisible();
      const plot = await card.locator('svg.recharts-surface').isVisible();
      if (Number(empty)+Number(error)+Number(plot)!==1) return null;
      return error?'error':empty?'empty':'data';
    }
    if (label === 'Chart range') return value(card.getByRole('group',{name:'Time range',exact:true}).locator('button[aria-pressed="true"]'));
    if (label === 'Chart tooltip') {
      const sample = hovered.get(this.page);
      if (!sample || sample.path!==this.page.url() || sample.title!==subject)
        throw new WorkerError('MISSING_OBSERVER','Tooltip requires a hover on this chart and page');
      if (Date.now()-sample.at<150) return null;
      const tooltip = await unique(card.locator('.recharts-tooltip-wrapper'));
      return tooltip ? normalize(await tooltip.innerText({timeout:500})) : 'absent';
    }
    if (label === 'Chart x axis' || label === 'Chart y axis') {
      const axis = await unique(card.locator(label==='Chart x axis'?'.recharts-xAxis':'.recharts-yAxis'));
      if (!axis) return null;
      const ticks = axis.locator('.recharts-cartesian-axis-tick-value').filter({visible:true});
      return (await ticks.allTextContents()).map(normalize).join('|') || null;
    }
    if (label === 'Chart legend') {
      // MetricChartCard's final body child contains the threshold/Trend legend.
      const body = card.locator(':scope > div:last-child');
      const legend = body.locator(':scope > div:last-child > span').filter({visible:true});
      return (await legend.allTextContents()).map(normalize).join('|') || null;
    }
    if (label === 'Chart segments') {
      const path = await unique(card.locator('path.recharts-line-curve'));
      if (!path) return null;
      const d = await path.getAttribute('d');
      if (!d || /NaN|Infinity/.test(d)) return null;
      return String((d.match(/[Mm]/g) || []).length);
    }
    throw new WorkerError('INVALID_INPUT','Unknown analytics assertion');
  }
}
