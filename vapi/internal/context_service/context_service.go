package contextservice

import "context"

type ContextService struct {
	Ctx context.Context
}

func (s *ContextService) GetContext() context.Context {
	return s.Ctx
}
func NewContextService(context context.Context) *ContextService {
	return &ContextService{
		Ctx: context,
	}
}
