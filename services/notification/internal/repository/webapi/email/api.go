package email

import (
	"context"
	"fmt"

	"gopkg.in/gomail.v2"
)

type Client struct {
	host     string
	port     uint16
	username string
	password string
	sender   string
	name     string
}

func NewAPI(host string, port uint16, username, password, sender, name string) (*Client, error) {
	if host == "" {
		return nil, fmt.Errorf("host is empty")
	}
	if port == 0 {
		return nil, fmt.Errorf("port is empty")
	}
	if username == "" {
		return nil, fmt.Errorf("username is empty")
	}
	if password == "" {
		return nil, fmt.Errorf("password is empty")
	}
	if sender == "" {
		return nil, fmt.Errorf("sender is empty")
	}
	if name == "" {
		return nil, fmt.Errorf("name is empty")
	}

	return &Client{
		host:     host,
		port:     port,
		username: username,
		password: password,
		sender:   sender,
		name:     name,
	}, nil
}

func (c *Client) Send(ctx context.Context, subject, body string, receivers []string) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	dialer := gomail.NewDialer(c.host, int(c.port), c.username, c.password)
	dialer.SSL = true

	conn, err := dialer.Dial()
	if err != nil {
		return fmt.Errorf("dial connection: %w", err)
	}
	defer conn.Close()

	msg := gomail.NewMessage()
	msg.SetAddressHeader("From", c.sender, c.name)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", body)

	err = conn.Send(c.sender, receivers, msg)
	if err != nil {
		return fmt.Errorf("send email: %w", err)
	}

	return nil
}
