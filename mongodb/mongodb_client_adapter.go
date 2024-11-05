package mongodb

import (
	"context"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// ClientAdapter is a concrete struct of mongodb client adapter.
type ClientAdapter struct {
	clientOptions *options.ClientOptions
	client        *mongo.Client
}

// NewClientAdapter is a constructor.
func NewClientAdapter(clientOptions *options.ClientOptions) Client {
	return &ClientAdapter{
		clientOptions: clientOptions,
	}
}

// Connect initializes the Client by starting background monitoring goroutines.
// If the Client was created using the NewClient function, this method must be called before a Client can be used.
//
// Connect starts background goroutines to monitor the state of the deployment and does not do any I/O in the main
// goroutine. The Client.Ping method can be used to verify that the connection was created successfully.
func (c *ClientAdapter) Connect(ctx context.Context) (err error) {
	c.client, err = mongo.Connect(ctx, c.clientOptions)
	if err != nil {
		return err
	}

	// Optional: Ping the server to verify the connection
	if err = c.client.Ping(ctx, readpref.Primary()); err != nil {
		logrus.Printf("an error ocurred when connect to mongoDB : %v", err)
		return err
	}
	return
}

// Database returns a handle for a database with the given name configured with the given DatabaseOptions.
func (c *ClientAdapter) Database(name string, opts ...*options.DatabaseOptions) (db Database) {
	mongoDatabase := c.client.Database(name, opts...)
	db = &DatabaseAdapter{db: mongoDatabase}
	return
}

// Disconnect closes sockets to the topology referenced by this Client. It will
// shut down any monitoring goroutines, close the idle connection pool, and will
// wait until all the in use connections have been returned to the connection
// pool and closed before returning. If the context expires via cancellation,
// deadline, or timeout before the in use connections have returned, the in use
// connections will be closed, resulting in the failure of any in flight read
// or write operations. If this method returns with no errors, all connections
// associated with this Client have been closed.
func (c *ClientAdapter) Disconnect(ctx context.Context) (err error) {
	err = c.client.Disconnect(ctx)
	return
}
