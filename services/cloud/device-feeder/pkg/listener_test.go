package pkg

import (
	"testing"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/stretchr/testify/assert"
	"github.com/ukama/ukama/services/common/msgbus"
	"github.com/wagslane/go-rabbitmq"
)

func Test_isRetryLimitReached_ValidHeader(t *testing.T) {

	q := &QueueListener{
		maxRetryCount: 1,
	}

	delivery := rabbitmq.Delivery{
		Delivery: amqp.Delivery{
			Headers: amqp.Table{
				"x-death": []interface{}{
					amqp.Table{
						"count":    int64(2),
						"reason":   "expired",
						"time":     "2018-12-06T15:00:00.000000000Z",
						"routing":  msgbus.DeviceFeederRequestRoutingKey,
						"exchange": deadLetterExchangeName,
					},
				},
			},
		},
	}

	t.Run("happyPath", func(tt *testing.T) {
		// Act
		ret := q.isRetryLimitReached(delivery)

		// Assert
		assert.Equal(tt, true, ret)
	})

	t.Run("LessThenLimit", func(tt *testing.T) {
		q.maxRetryCount = 3
		ret := q.isRetryLimitReached(delivery)

		// Assert
		assert.Equal(tt, false, ret)

	})

}

func Test_isRetryLimitReached_NoHeader(t *testing.T) {
	// Arrange
	q := &QueueListener{
		maxRetryCount: 1,
	}
	delivery := rabbitmq.Delivery{
		Delivery: amqp.Delivery{
			Headers: amqp.Table{},
		},
	}

	// Act
	ret := q.isRetryLimitReached(delivery)

	// Assert
	assert.Equal(t, false, ret)
}
