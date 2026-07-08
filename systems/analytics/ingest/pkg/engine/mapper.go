/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package engine

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
	"text/template"
)

// GetPath extracts a value from a decoded JSON item using a "$.a.b" path.
// "$" returns the item itself.
func GetPath(item interface{}, path string) (interface{}, bool) {
	if path == "$" {
		return item, true
	}

	path = strings.TrimPrefix(path, "$.")
	cur := item

	for _, seg := range strings.Split(path, ".") {
		m, ok := cur.(map[string]interface{})
		if !ok {
			return nil, false
		}

		cur, ok = m[seg]
		if !ok {
			return nil, false
		}
	}

	return cur, true
}

// MapItem projects one source item through the spec field map, then merges
// the for_each binds (lineage propagation: child rows carry ancestor ids).
func MapItem(item interface{}, fieldMap map[string]string, binds map[string]string) map[string]interface{} {
	out := make(map[string]interface{}, len(fieldMap)+len(binds))

	for field, path := range fieldMap {
		if v, ok := GetPath(item, path); ok {
			out[field] = v
		}
	}

	for k, v := range binds {
		if _, exists := out[k]; !exists {
			out[k] = v
		}
	}

	return out
}

// HashFields produces a deterministic content hash of mapped fields (the
// change-log dedup key component).
func HashFields(fields map[string]interface{}) string {
	keys := make([]string, 0, len(fields))
	for k := range fields {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var b bytes.Buffer
	for _, k := range keys {
		v, _ := json.Marshal(fields[k])
		b.WriteString(k)
		b.WriteByte('=')
		b.Write(v)
		b.WriteByte(';')
	}

	sum := sha256.Sum256(b.Bytes())

	return hex.EncodeToString(sum[:16])
}

// RenderTemplate expands {{.key}} references (window bounds + for_each binds)
// in endpoints and params.
func RenderTemplate(tmpl string, ctx map[string]string) (string, error) {
	if !strings.Contains(tmpl, "{{") {
		return tmpl, nil
	}

	t, err := template.New("t").Option("missingkey=error").Parse(tmpl)
	if err != nil {
		return "", fmt.Errorf("parsing template %q: %w", tmpl, err)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		return "", fmt.Errorf("rendering template %q: %w", tmpl, err)
	}

	return buf.String(), nil
}

// ItemsAt extracts the result array from a response body using the spec's
// items path ("$" for a bare JSON array).
func ItemsAt(body []byte, itemsPath string) ([]interface{}, error) {
	var doc interface{}
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil, fmt.Errorf("decoding response: %w", err)
	}

	node, ok := GetPath(doc, itemsPath)
	if !ok {
		// A missing items key is an ERROR, not an empty set: treating it as
		// empty would tombstone every entity of a snapshot dataset when a
		// source renames its wrapper key or returns an error object.
		return nil, fmt.Errorf("items path %q not found in response", itemsPath)
	}

	if node == nil {
		return []interface{}{}, nil // key present, null value = empty set
	}

	arr, ok := node.([]interface{})
	if !ok {
		return nil, fmt.Errorf("items path %q does not point at an array", itemsPath)
	}

	return arr, nil
}
