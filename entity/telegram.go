package entity

import "strconv"

type TelegramId int64

func (id TelegramId) String() string {
	return strconv.FormatInt(int64(id), 10)
}

func (id TelegramId) Int64() int64 {
	return int64(id)
}
