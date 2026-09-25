package mainhandler

import "fmt"

func serviceSelectedMessage(
	serviceName string,
	platformName string,
	categoryIcon string,
	categoryName string,
	minQuantity int64,
	maxQuantity int64,
) string {
	return fmt.Sprintf(
		"✅ سرویس «%s» با موفقیت انتخاب شد.\n\n"+
			"📱 پلتفرم: %s\n"+
			"%s دسته‌بندی: %s\n\n"+
			"• حداقل تعداد سفارش: %d\n"+
			"• حداکثر تعداد سفارش: %d\n\n"+
			"🔢 لطفاً تعداد مورد نظر خود را فقط به صورت عدد ارسال کنید:",
		serviceName,
		platformName,
		categoryIcon,
		categoryName,
		minQuantity,
		maxQuantity,
	)
}
