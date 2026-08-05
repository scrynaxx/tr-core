package customerclient

type Client interface {
}

type client struct {
}

func New() Client {
	return &client{}
}
