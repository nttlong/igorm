package contextservice

import "context"

type ContextService struct {
	onGetContext func() context.Context
}

func NewContextService() *ContextService {
	return &ContextService{}
}
func (s *ContextService) OnGetContext(onGetContext func() context.Context) {
	s.onGetContext = onGetContext
}
func (s *ContextService) GetContext() context.Context {
	if s.onGetContext == nil {
		return context.Background()
	}

	return s.onGetContext()

}
