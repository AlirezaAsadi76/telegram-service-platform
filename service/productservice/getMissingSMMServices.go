package productservice

import (
	"context"
	"fmt"
	"telegram-service-platform/entity/smmentity"
	"telegram-service-platform/pkg/richerror"
)

func (s Service) GetMissingSMMServices(ctx context.Context) ([]smmentity.SMM, error) {
	const op = "productservice.GetMissingSMMServices"

	adapterServices, err := s.adapter.AllServices(ctx)
	fmt.Println(-1)
	if err != nil {
		return nil, richerror.New(op, err).WithKind(richerror.KindExternalAPI)
	}
	fmt.Println(0)
	dbServices, sErr := s.repository.SMMServiceGetAll(ctx)
	if sErr != nil {
		return nil, richerror.New(op, sErr).WithKind(richerror.KindQueryFailure)
	}
	fmt.Println(1)
	adapterMap := make(map[int64]bool)
	for _, svc := range adapterServices.Services {
		adapterMap[svc.Service] = true
	}

	var missingServices []smmentity.SMM
	for _, dbSvc := range dbServices {
		if !adapterMap[dbSvc.Service] {
			missingServices = append(missingServices, dbSvc)
		}
	}

	return missingServices, nil
}
