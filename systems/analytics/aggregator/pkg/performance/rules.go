/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain it at http://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2026-present, Ukama Inc.
 */

package performance

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/ukama/ukama/systems/analytics/schema"
)

// Status rules: minimal "<name> <op> <literal>" comparisons over report
// columns and entity attributes. Deliberately NOT a general expression
// language — anything fancier becomes a registered Go resolver.

// EvaluateStatus returns the label of the first matching rule ("" when no
// rule matches and there is no default).
func EvaluateStatus(rules []schema.StatusRule, values map[string]interface{}) (string, error) {
	for _, rule := range rules {
		if rule.Default {
			return rule.Label, nil
		}

		match, err := evalCondition(rule.When, values)
		if err != nil {
			return "", err
		}

		if match {
			return rule.Label, nil
		}
	}

	return "", nil
}

func evalCondition(cond string, values map[string]interface{}) (bool, error) {
	parts := strings.Fields(cond)
	if len(parts) != 3 {
		return false, fmt.Errorf("status rule %q: expected '<name> <op> <literal>'", cond)
	}

	name, op, lit := parts[0], parts[1], parts[2]

	actual, ok := values[name]
	if !ok {
		return false, fmt.Errorf("status rule %q: unknown column/attribute %q", cond, name)
	}

	// Numeric comparison when both sides parse as numbers.
	if litNum, err := strconv.ParseFloat(lit, 64); err == nil {
		actNum, ok := toNumber(actual)
		if !ok {
			return false, nil
		}

		switch op {
		case "==":
			return actNum == litNum, nil
		case "!=":
			return actNum != litNum, nil
		case "<":
			return actNum < litNum, nil
		case "<=":
			return actNum <= litNum, nil
		case ">":
			return actNum > litNum, nil
		case ">=":
			return actNum >= litNum, nil
		default:
			return false, fmt.Errorf("status rule %q: unknown operator %q", cond, op)
		}
	}

	// Bool/string equality otherwise.
	actStr := strings.ToLower(fmt.Sprintf("%v", actual))
	litStr := strings.ToLower(strings.Trim(lit, `"'`))

	switch op {
	case "==":
		return actStr == litStr, nil
	case "!=":
		return actStr != litStr, nil
	default:
		return false, fmt.Errorf("status rule %q: operator %q needs numeric operands", cond, op)
	}
}

func toNumber(v interface{}) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case bool:
		if t {
			return 1, true
		}

		return 0, true
	case string:
		f, err := strconv.ParseFloat(t, 64)

		return f, err == nil
	default:
		return 0, false
	}
}
