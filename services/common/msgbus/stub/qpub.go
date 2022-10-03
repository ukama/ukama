package stub

type QPubStub struct {
}

func (q QPubStub) Publish(payload any, routingKey string) error {
	return nil
}

func (q QPubStub) PublishToQueue(queueName string, payload any) error {
	return nil
}
func (q QPubStub) Close() error {
	return nil
}
