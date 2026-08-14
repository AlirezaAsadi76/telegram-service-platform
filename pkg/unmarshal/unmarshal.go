package unmarshal

import (
	"encoding/json"
	"log"
	"telegram-service-platform/entity/notificationentity"
	"telegram-service-platform/pkg/msgerror"
	"telegram-service-platform/pkg/richerror"
)

func UnmarshalToUint64(rec string) (uint64, error) {

	const Op = "unmarshal.UnmarshalToUint64"
	var data uint64
	if err := json.Unmarshal([]byte(rec), &data); err != nil {
		log.Printf("unmarshal data id failed: %v", err)
		return 0, richerror.New(Op, err).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.UnmarshalFailed)
	}

	return data, nil
}

func UnmarshalToNotification(rec string) (notificationentity.Notification, error) {

	const Op = "unmarshal.UnmarshalToNotification"
	var data notificationentity.Notification
	if err := json.Unmarshal([]byte(rec), &data); err != nil {
		log.Printf("unmarshal data id failed: %v", err)
		return notificationentity.Notification{}, richerror.New(Op, err).WithKind(richerror.KindSerializationFailure).WithMessage(msgerror.UnmarshalFailed)
	}

	return data, nil
}
