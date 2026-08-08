package keyboard

import "github.com/go-telegram/bot/models"

type Builder struct {
	rows [][]models.InlineKeyboardButton
}

func NewBuilder() *Builder {
	return &Builder{
		rows: make([][]models.InlineKeyboardButton, 0),
	}
}

func (b *Builder) AddRow(buttons ...Button) *Builder {

	row := make([]models.InlineKeyboardButton, 0, len(buttons))

	for _, btn := range buttons {
		row = append(row, btn.ToTelegramInlineButton())
	}

	b.rows = append(b.rows, row)

	return b
}

func (b *Builder) AddButtonsPerRow(buttons []Button, size int) *Builder {

	row := make([]Button, 0, size)

	for _, btn := range buttons {

		row = append(row, btn)

		if len(row) == size {

			b.AddRow(row...)
			row = row[:0]
		}
	}

	if len(row) > 0 {
		b.AddRow(row...)
	}

	return b
}

func (b *Builder) Build() *models.InlineKeyboardMarkup {

	return &models.InlineKeyboardMarkup{
		InlineKeyboard: b.rows,
	}
}
