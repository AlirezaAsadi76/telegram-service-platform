package productparams

import "telegram-service-platform/entity/smmentity"

type GetSmmServiceByIDRequest struct {
	Id int64
}

type GetSmmServiceByIDResponse struct {
	Smm *smmentity.SMM
}
