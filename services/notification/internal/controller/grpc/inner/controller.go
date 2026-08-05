package inner

import (
	"context"
	"fmt"

	pbinner "github.com/scrynaxx/tr-core/contracts/generated/services/inner/notification"
	"github.com/scrynaxx/tr-core/services/notification/internal/usecase"
	"google.golang.org/protobuf/types/known/emptypb"
)

type Controller struct {
	emailCase usecase.Email
}

func NewController(emailCase usecase.Email) pbinner.NotificationServiceServer {
	return &Controller{
		emailCase: emailCase,
	}
}

func (c *Controller) SendEmail(ctx context.Context, req *pbinner.SendEmailRequest) (*emptypb.Empty, error) {
	err := c.emailCase.Send(ctx, req.Subject, req.Body, req.Receivers)
	if err != nil {
		return nil, fmt.Errorf("send emailCase message: %w", err)
	}

	return new(emptypb.Empty), nil
}
