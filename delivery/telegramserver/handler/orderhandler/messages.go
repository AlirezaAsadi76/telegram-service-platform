package orderhandler

import "fmt"

const (
	orderFlowExpiredMessage = "⚠️ این سفارش دیگر معتبر نیست. لطفاً سرویس مورد نظر خود را دوباره انتخاب کنید."

	noActiveOrderFlowMessage = "⚠️ لطفاً ابتدا از منوی اصلی یک سرویس را انتخاب کنید."

	invalidQuantityMessage = "❌ تعداد نامعتبر است. لطفاً فقط عدد وارد کنید."

	priceCalculationFailedMessage = "❌ خطایی در محاسبه قیمت سفارش رخ داد. لطفاً چند لحظه بعد دوباره تلاش کنید."

	genericOrderProcessingErrorMessage = "❌ خطایی در پردازش درخواست شما رخ داد. لطفاً دوباره تلاش کنید یا با پشتیبانی تماس بگیرید."

	insufficientBalanceMessage = "❌ موجودی کیف پول شما کافی نیست.\n" +
		"لطفاً ابتدا کیف پول خود را شارژ کنید.\n\n" +
		"سفارش شما تا ۱۵ دقیقه دیگر در سیستم باقی می‌ماند تا پس از شارژ، پرداخت را انجام دهید."

	cancelOrderMessage = "❌ سفارش شما با موفقیت لغو شد.\n" +
		"می‌توانید از منوی اصلی سرویس جدیدی انتخاب کنید."

	paymentStatusMessage = "⏳ <b>وضعیت پرداخت:</b>\n\n" +
		"لطفاً چند لحظه صبر کنید تا سیستم پرداخت را بررسی کند.\n\n" +
		"<i>(در نسخه نهایی، این دکمه مستقیماً وضعیت را از OrderService استعلام می‌گیرد)</i>"
)

func validationErrorMessage(
	message string,
) string {
	return "❌ " + message
}

func quantityResultMessage(
	quantity int64,
	usd string,
	toman string,
) string {
	return fmt.Sprintf(
		"✅ تعداد <b>%d</b> با موفقیت ثبت شد.\n\n"+
			"💰 <b>قیمت برآورد شده برای این سفارش:</b>\n"+
			"• به دلار: <code>%s $</code>\n"+
			"• به تومان: <code>%s تومان</code>\n\n"+
			"🔗 حالا لطفاً لینک کانال یا گروه خود را ارسال کنید.\n"+
			"(مثال: <code>https://t.me/YourChannel</code>)\n\n"+
			"⚠️ <b>توجه:</b> لینک باید عمومی (Public) باشد تا سرویس قابل انجام باشد.",
		quantity,
		usd,
		toman,
	)
}

func confirmOrderMessage(
	platform string,
	service string,
	quantity int64,
	link string,
	amount string,
) string {
	return fmt.Sprintf(
		"📋 <b>خلاصه سفارش شما:</b>\n\n"+
			"📱 پلتفرم: %s\n"+
			"📦 سرویس: %s\n"+
			"🔢 تعداد: <code>%d</code> عدد\n"+
			"🔗 لینک: %s\n\n"+
			"💰 <b>مبلغ قابل پرداخت:</b> <code>%s تومان</code>\n\n"+
			"آیا اطلاعات فوق صحیح است؟",
		platform,
		service,
		quantity,
		link,
		amount,
	)
}

func directPaymentCreatedMessage(
	methodName string,
	orderID uint64,
	amount string,
) string {
	return fmt.Sprintf(
		"💳 <b>درخواست %s ثبت شد</b>\n\n"+
			"📋 شماره سفارش: <code>%d</code>\n"+
			"💰 مبلغ: <code>%s</code> تومان\n\n"+
			"برای تکمیل پرداخت، روی دکمه زیر کلیک کنید:\n\n"+
			"<i>⚠️ توجه: در حال حاضر این بخش در محیط شبیه‌سازی (Stub) قرار دارد.</i>",
		methodName,
		orderID,
		amount,
	)
}

func walletPurchaseSuccessMessage(
	amount string,
) string {
	return fmt.Sprintf(
		"✅ <b>پرداخت با موفقیت انجام شد!</b>\n\n"+
			"🎉 سفارش شما ثبت گردید و در حال پردازش است.\n\n"+
			"💰 مبلغ کسر شده: <code>%s</code> تومان\n\n"+
			"📌 می‌توانید وضعیت سفارش خود را از بخش «💳 تراکنش‌ها» پیگیری کنید.",
		amount,
	)
}
