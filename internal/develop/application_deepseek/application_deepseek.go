package application_deepseek

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/mahdi-cpp/messages-api/internal/collection_manager_gemini"
)

// Post implements collectionItem and holds app-specific data
type Post struct {
	ID      uuid.UUID     `json:"id"`
	Name    string        `json:"name"`
	Timeout time.Duration `json:"timeout"`
	// Add other fields as needed
}

func (a *Post) SetID(id uuid.UUID) { a.ID = id }
func (a *Post) GetID() uuid.UUID   { return a.ID }

type command struct {
	action    string
	app       *Post
	id        uuid.UUID
	replyChan chan reply
}

type reply struct {
	app *Post
	err error
}

// ApplicationManager manages applications with timeout controls
type ApplicationManager struct {
	cm         *collection_manager_gemini.Manager[*Post]
	cmdChannel chan command
	timeout    time.Duration // Default timeout for operations
}

func NewApplicationManager(dataDir string, defaultTimeout time.Duration) (*ApplicationManager, error) {

	cm, err := collection_manager_gemini.New[*Post](dataDir)
	if err != nil {
		return nil, err
	}

	am := &ApplicationManager{
		cm:         cm,
		cmdChannel: make(chan command),
		timeout:    defaultTimeout,
	}

	go am.run() // Start command processor
	return am, nil
}

func (am *ApplicationManager) run() {
	for cmd := range am.cmdChannel {
		ctx, cancel := context.WithTimeout(context.Background(), am.timeout)
		go func(cmd command) {
			defer cancel()
			am.processCommand(ctx, cmd)
		}(cmd)
	}
}

func (am *ApplicationManager) processCommand(ctx context.Context, cmd command) {
	select {
	case <-ctx.Done():
		cmd.replyChan <- reply{err: ctx.Err()}
	default:
		switch cmd.action {
		case "create":
			app, err := am.cm.Create(cmd.app)
			cmd.replyChan <- reply{app: app, err: err}
		case "read":
			app, err := am.cm.Read(cmd.id)
			cmd.replyChan <- reply{app: app, err: err}
			// Add cases for update, delete, etc.
		}
	}
}

func (am *ApplicationManager) CreateApp(app *Post) (*Post, error) {

	replyChan := make(chan reply)
	am.cmdChannel <- command{action: "create", app: app, replyChan: replyChan}

	select {
	case r := <-replyChan:
		return r.app, r.err
	case <-time.After(am.timeout):
		return nil, context.DeadlineExceeded
	}
}

func (am *ApplicationManager) ReadApp(id uuid.UUID) (*Post, error) {
	replyChan := make(chan reply)
	am.cmdChannel <- command{action: "read", id: id, replyChan: replyChan}
	select {
	case r := <-replyChan:
		return r.app, r.err
	case <-time.After(am.timeout):
		return nil, context.DeadlineExceeded
	}
}

// Similarly implement Update, Delete, etc.
