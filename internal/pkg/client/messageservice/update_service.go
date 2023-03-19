package messageservice

import "github.com/sudal-91/monitoring/pkg/message"

type updateService struct {
}

func (u *updateService) router(data message.Update) {
	switch {
	case data.NewConnection != nil:
		return
	}

}
