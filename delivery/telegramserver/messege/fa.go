package message

func ProductStarsSelect() Message {

	return New(
		"⭐ انتخاب Telegram Stars\n\n" +
			"لطفا تعداد مورد نظر خود را انتخاب کنید:",
	)
}

func ProductPremiumSelect() Message {

	return New(
		"👑 انتخاب Telegram Premium\n\n" +
			"مدت زمان اشتراک را انتخاب کنید:",
	)
}

func ProductConfirm() Message {

	return New(
		"آیا از خرید خود مطمئن هستید؟",
	)
}

func BuySuccess() Message {

	return New(
		"✅ سفارش شما با موفقیت ثبت شد.",
	)
}

func BuyFailed() Message {

	return New(
		"⚠️ ثبت سفارش انجام نشد.",
	)
}
