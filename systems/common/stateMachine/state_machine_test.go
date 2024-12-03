/*
 * This Source Code Form is subject to the terms of the Mozilla Public
 * License, v. 2.0. If a copy of the MPL was not distributed with this
 * file, You can obtain one at https://mozilla.org/MPL/2.0/.
 *
 * Copyright (c) 2023-present, Ukama Inc.
 */
package statemachine

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNodeStateMachine(t *testing.T) {
	var capturedEvents []Event
	mockHandler := func(event Event) {
		capturedEvents = append(capturedEvents, event)
	}

	sm := NewStateMachine(mockHandler)

	configFile := createTempConfigFile(t, nodeStateConfig)

	t.Run("Configured State Transitions", func(t *testing.T) {
		t.Run("ready event", func(t *testing.T) {
			instance, err := sm.NewInstance(configFile, "test-node-ready", "Configured")
			require.NoError(t, err)

			err = instance.Transition("ready")
			assert.NoError(t, err)
			assert.Equal(t, "Operational", instance.CurrentState)
		})

		t.Run("config event", func(t *testing.T) {
			instance, err := sm.NewInstance(configFile, "test-node-config", "Configured")
			require.NoError(t, err)

			err = instance.Transition("config")
			assert.NoError(t, err)
			assert.Equal(t, "Configured", instance.CurrentState)
		})

		t.Run("fault event", func(t *testing.T) {
			instance, err := sm.NewInstance(configFile, "test-node-fault", "Configured")
			require.NoError(t, err)

			err = instance.Transition("fault")
			assert.NoError(t, err)
			assert.Equal(t, "Faulty", instance.CurrentState)
		})

		substateCases := []struct {
			name          string
			event         string
			expectedState string
			expectedEvents []string
		}{
			{"update substate", "update", "Configured", []string{ "online", "ready"}},
			{"upgrade substate", "upgrade", "Configured", []string{ "online", "ready"}},
			{"downgrade substate", "downgrade", "Configured", []string{ "online", "ready"}},
		}

		for _, tc := range substateCases {
			t.Run(tc.name, func(t *testing.T) {
				instance, err := sm.NewInstance(configFile, "test-node-"+tc.name, "Configured")
				require.NoError(t, err)

				err = instance.Transition(tc.event)
				assert.NoError(t, err)
				assert.Equal(t, tc.expectedState, instance.CurrentState)
				assert.Equal(t, tc.expectedEvents, instance.ExpectedEvents)
				assert.Equal(t, tc.event, instance.CurrentSubstate)
			})
		}
	})

	
}

func createTempConfigFile(t *testing.T, content string) string {
	t.Helper()
	tempFile, err := os.CreateTemp("", "node-state-config-*.json")
	require.NoError(t, err)
	defer tempFile.Close()

	_, err = tempFile.Write([]byte(content))
	require.NoError(t, err)

	return tempFile.Name()
}


const nodeStateConfig = `{
  "version": "0.0.0",
  "entity": "node",
  "file": "nodeState.json",
  "states": [
    {
      "name": "Unknown",
      "description": "Initial state when node comes online for the first time, or state after offBoarding",
      "events": ["online", "offline", "reboot", "onboarding"],
      "transition": [
        {
          "to_state": "Configured",
          "trigger": ["onboarding"]
        }
      ],
      "substate": {
        "events": ["online", "offline", "reboot"],
        "transition": [
          {
            "to_state": "on",
            "trigger": ["online"]
          },
          {
            "to_state": "off",
            "trigger": ["offline"]
          },
          {
            "to_state": "reboot",
            "trigger": ["reboot"],
            "expectedEvents": ["offline", "online"],
            "timeout": 180
          }
        ]
      }
    },
    {
      "name": "Configured",
      "description": "Node is in configuration state",
      "events": [
        "online",
        "offline",
        "reboot",
        "config",
        "ready",
        "upgrade",
        "update",
        "downgrade",
        "fault"
      ],
      "transition": [
        {
          "to_state": "Faulty",
          "trigger": ["fault"]
        },
        {
          "to_state": "Configured",
          "trigger": ["config"]
        },
        {
          "to_state": "Operational",
          "trigger": ["ready"]
        }
      ],
      "substate": {
        "events": [
          "online",
          "offline",
          "reboot",
          "upgrade",
          "update",
          "downgrade"
        ],
        "transition": [
          {
            "to_state": "on",
            "trigger": ["online"]
          },
          {
            "to_state": "off",
            "trigger": ["offline"]
          },
          {
            "to_state": "reboot",
            "trigger": ["reboot"],
            "expectedEvents": ["offline", "online"],
            "timeout": 180
          },
          {
            "to_state": "update",
            "trigger":["ready"],
            "expectedEvents": ["offline", "online", "ready"],
            "timeout": 600
          },
          {
            "to_state": "update",
            "trigger": ["update"],
            "expectedEvents": ["offline", "online", "ready"],
            "timeout": 600
          },
          {
            "to_state": "upgrade",
            "trigger": ["upgrade"],
            "expectedEvents": ["offline", "online", "ready"],
            "timeout": 180
          },
          {
            "to_state": "downgrade",
            "trigger": ["downgrade"],
            "expectedEvents": ["offline", "online", "ready"],
            "timeout": 180
          }
        ]
      }
    },
    {
      "name": "Operational",
      "description": "Node is fully operational and part of a site",
      "events": [
        "online",
        "offline",
        "reboot",
        "fault",
        "offboarding",
        "reset"
      ],
      "transition": [
        {
          "to_state": "Faulty",
          "trigger": ["fault"]
        },
        {
          "to_state": "Configured",
          "trigger": ["reset"]
        },
        {
          "to_state": "Unknown",
          "trigger": ["offboarding"]
        }
      ],
      "substate": {
        "events": ["online", "offline", "reboot"],
        "transition": [
          {
            "to_state": "on",
            "trigger": ["online"]
          },
          {
            "to_state": "off",
            "trigger": ["offline"]
          },
          {
            "to_state": "reboot",
            "trigger": ["reboot"],
            "expectedEvents": ["offline", "online"],
            "timeout": 180
          }
        ]
      }
    },
    {
      "name": "Faulty",
      "description": "Node is in a faulty state",
      "events": [
        "offline",
        "reboot",
        "fault",
        "online",
        "config"
      ],
      "transition": [
        {
          "to_state": "Faulty",
          "trigger": ["fault"]
        },
        {
          "to_state": "Configured",
          "trigger": ["config"]
        }
      ],
      "substate": {
        "events": ["online", "offline", "reboot"],
        "transition": [
          {
            "to_state": "on",
            "trigger": ["online"]
          },
          {
            "to_state": "off",
            "trigger": ["offline"]
          },
          {
            "to_state": "reboot",
            "trigger": ["reboot"],
            "expectedEvents": ["offline", "online"],
            "timeout": 180
          }
        ]
      }
    }
  ]
}`