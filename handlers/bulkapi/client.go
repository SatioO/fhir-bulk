package bulkapi

type client struct{}

func NewBulkFHIRClient() *client {
	return &client{}
}

func (c *client) CreateNewJob() error {
	return nil
}
