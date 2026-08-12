/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine

import (
	"encoding/json"
	"fmt"
	"net/url"
	"sort"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/ukama/ukama/systems/analytics/ingest/pkg/db"
	"github.com/ukama/ukama/systems/analytics/ingest/pkg/resolver"
	"github.com/ukama/ukama/systems/analytics/schema"
	"github.com/ukama/ukama/systems/common/rest/client"
)

// Puller executes a single (dataset, window) pull: resolve address, fan out
// for_each iterations, fetch, map, and persist per the pull strategy.
type Puller struct {
	resolver resolver.Resolver
	raw      db.RawRepo
	rest     *client.Resty
	retry    int
	org      string
}

func NewPuller(res resolver.Resolver, raw db.RawRepo, org string, retry int, debug bool) *Puller {
	return &Puller{
		resolver: res,
		raw:      raw,
		rest:     client.NewResty(client.WithDebug(debug)),
		retry:    retry,
		org:      org,
	}
}

// Execute runs the pull for one window. Returns the number of records
// written (for snapshots: change rows).
func (p *Puller) Execute(pull schema.PullSpec, win schema.Window) (int, error) {
	iterations, err := p.iterations(pull, win)
	if err != nil {
		return 0, err
	}

	records := make([]schema.RawRecord, 0)
	now := time.Now().UTC()

	for _, binds := range iterations {
		items, err := p.fetchWithRetry(pull, win, binds)
		if err != nil {
			// on_error: record — for health-probe style pulls where an
			// unreachable target IS the signal: write a synthetic row with
			// the binds + unreachable:true instead of failing the window.
			if pull.OnError == "record" {
				log.Warnf("dataset %s window %d iteration %v unreachable, recording: %v",
					pull.Key, win.ID, binds, err)

				items = []interface{}{map[string]interface{}{"unreachable": true}}
			} else {
				return 0, fmt.Errorf("dataset %s window %d iteration %v: %w",
					pull.Key, win.ID, binds, err)
			}
		}

		for _, item := range items {
			fields := MapItem(item, pull.Map, binds)

			// Propagate the unreachable marker from on_error:record rows.
			if u, ok := item.(map[string]interface{}); ok {
				if unreachable, ok := u["unreachable"].(bool); ok && unreachable {
					fields["unreachable"] = true
				}
			}

			// Entity key: single mapped field, or a composite ("a,b,c") whose
			// component values are joined with "|" in spec order.
			entity := ""

			if entityFields := pull.EntityFields(); len(entityFields) > 0 {
				parts := make([]string, 0, len(entityFields))

				for _, ef := range entityFields {
					value := ""
					if v, ok := fields[ef.Name]; ok && v != nil {
						value = fmt.Sprintf("%v", v)
					}

					if value == "" && !ef.Optional {
						// A snapshot row without its entity key would collapse
						// the change-log — fail loudly instead of storing junk.
						return 0, fmt.Errorf("dataset %s window %d: entity field %q missing/empty in mapped item (check the spec's map paths against the source response)",
							pull.Key, win.ID, ef.Name)
					}

					parts = append(parts, value)
				}

				entity = strings.Join(parts, "|")
			}

			hash := HashFields(fields)

			payloadJSON, err := json.Marshal(item)
			if err != nil {
				return 0, fmt.Errorf("dataset %s: marshaling payload: %w", pull.Key, err)
			}

			fieldsJSON, err := json.Marshal(fields)
			if err != nil {
				return 0, fmt.Errorf("dataset %s: marshaling fields: %w", pull.Key, err)
			}

			dedup := hash
			if entity != "" {
				dedup = entity + ":" + hash
			}

			records = append(records, schema.RawRecord{
				OrgID:       p.org,
				DatasetKey:  pull.Key,
				WindowID:    win.ID,
				EntityKey:   entity,
				ContentHash: hash,
				DedupKey:    dedup,
				EventTime:   win.Start,
				Payload:     string(payloadJSON),
				Fields:      string(fieldsJSON),
				IngestedAt:  now,
			})
		}
	}

	if pull.Strategy == schema.StrategyFullSnapshot {
		return p.raw.UpsertSnapshot(p.org, pull.Key, win.ID, records)
	}

	if err := p.raw.InsertWindowed(records); err != nil {
		return 0, err
	}

	return len(records), nil
}

// iterations builds the bind sets: one empty set for plain pulls, one per
// parent row for for_each pulls (parent state as of the same window).
func (p *Puller) iterations(pull schema.PullSpec, win schema.Window) ([]map[string]string, error) {
	if pull.ForEach == nil {
		return []map[string]string{{}}, nil
	}

	parents, err := p.raw.StateAsOf(p.org, pull.ForEach.Dataset, win.ID)
	if err != nil {
		return nil, fmt.Errorf("loading parent dataset %s: %w", pull.ForEach.Dataset, err)
	}

	iterations := make([]map[string]string, 0, len(parents))

	for _, parent := range parents {
		fields := map[string]interface{}{}
		if err := json.Unmarshal([]byte(parent.Fields), &fields); err != nil {
			return nil, fmt.Errorf("parent %s fields: %w", parent.EntityKey, err)
		}

		// Optional parent-row filter (e.g. only tnode/anode for health).
		if f := pull.ForEach.Filter; f != nil {
			value := fmt.Sprintf("%v", fields[f.Field])

			keep := false
			for _, want := range f.In {
				if strings.EqualFold(value, want) {
					keep = true

					break
				}
			}

			if !keep {
				continue
			}
		}

		// A parent row missing a bind field cannot produce a meaningful child
		// pull — e.g. a node the registry serializes without its `site` block
		// because it is not attached to a site yet, so site_id/network_id are
		// absent from the mapped fields. Skip that row instead of failing the
		// whole dataset window: one unassigned node must not blank a KPI for
		// the entire org.
		binds := map[string]string{}
		skip := false

		for _, b := range pull.ForEach.Bind {
			v, ok := fields[b]
			if !ok {
				log.Warnf("dataset %s: skipping parent %s row %s, missing bind field %q",
					pull.Key, pull.ForEach.Dataset, parent.EntityKey, b)

				skip = true

				break
			}

			binds[b] = fmt.Sprintf("%v", v)
		}

		if skip {
			continue
		}

		iterations = append(iterations, binds)
	}

	// Deterministic execution order.
	sort.Slice(iterations, func(i, j int) bool {
		return bindKey(iterations[i]) < bindKey(iterations[j])
	})

	return iterations, nil
}

func (p *Puller) fetchWithRetry(pull schema.PullSpec, win schema.Window, binds map[string]string) ([]interface{}, error) {
	var lastErr error

	for attempt := 0; attempt < p.retry; attempt++ {
		items, err := p.fetch(pull, win, binds)
		if err == nil {
			return items, nil
		}

		lastErr = err

		// Address may be stale: re-resolve before the next attempt.
		p.resolver.Invalidate(p.org, pull.System)

		log.Warnf("pull %s attempt %d failed: %v", pull.Key, attempt+1, err)
		time.Sleep(time.Duration(attempt+1) * time.Second)
	}

	return nil, lastErr
}

func (p *Puller) fetch(pull schema.PullSpec, win schema.Window, binds map[string]string) ([]interface{}, error) {
	ctx := map[string]string{
		"WindowStart": win.Start.Format(time.RFC3339),
		"WindowEnd":   win.End.Format(time.RFC3339),
	}
	for k, v := range binds {
		ctx[k] = v
	}

	base := pull.BaseURL
	if base == "" {
		var (
			resolved string
			err      error
		)

		if pull.Gateway == "node" {
			resolved, err = p.resolver.ResolveNodeGw(p.org, pull.System)
		} else {
			resolved, err = p.resolver.Resolve(p.org, pull.System)
		}

		if err != nil {
			return nil, err
		}

		base = resolved
	}

	endpoint, err := RenderTemplate(pull.Endpoint, ctx)
	if err != nil {
		return nil, err
	}

	fullURL := strings.TrimSuffix(base, "/") + endpoint

	query := url.Values{}
	for k, tmpl := range pull.Params {
		v, err := RenderTemplate(tmpl, ctx)
		if err != nil {
			return nil, err
		}
		query.Set(k, v)
	}

	if len(query) > 0 {
		fullURL = fullURL + "?" + query.Encode()
	}

	resp, err := p.rest.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("GET %s: %w", fullURL, err)
	}

	return ItemsAt(resp.Body(), pull.Items)
}

func bindKey(binds map[string]string) string {
	keys := make([]string, 0, len(binds))
	for k := range binds {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var sb strings.Builder
	for _, k := range keys {
		sb.WriteString(k)
		sb.WriteByte('=')
		sb.WriteString(binds[k])
		sb.WriteByte(';')
	}

	return sb.String()
}
