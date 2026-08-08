package keyboard

func Back(back string) Button {

	return Button{
		Text:  "⬅️ Back",
		Data:  back,
		Style: Primary,
	}
}
