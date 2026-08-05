package main

import (
	"fmt"
	"telegram-service-platform/adapter/fragmentadapter"
)

func main() {

	fragment := fragmentadapter.New(
		fragmentadapter.Config{},
	)

	session, err := fragment.GetSession()

	if err != nil {

		panic(err)

	}

	fmt.Println(
		"hash:",
		session.Hash,
	)

	price, err := fragment.GetStarsPrice(
		1909,
	)

	if err != nil {
		panic(err)
	}
	fmt.Printf(
		"%+v\n",
		price,
	)
	stars, derr := fragment.GetStarsBuyState()

	if derr != nil {
		panic(derr)
	}
	fmt.Printf(
		"%+v\n",
		stars,
	)

}
