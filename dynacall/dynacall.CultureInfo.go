package dynacall

import "time"

type CultureInfo struct {
	// Ngôn ngữ (ví dụ: "en", "vi", "fr").
	// Thường là mã ngôn ngữ ISO 639-1.
	Language string `json:"language"`

	// Quốc gia/Khu vực (ví dụ: "US", "VN", "FR").
	// Thường là mã quốc gia ISO 3166-1 alpha-2.
	// Kết hợp với Language tạo thành mã văn hóa đầy đủ (ví dụ: "en-US", "vi-VN").
	Region string `json:"region"`

	// Tên hiển thị đầy đủ của văn hóa (ví dụ: "English (United States)", "Tiếng Việt (Việt Nam)").
	DisplayName string `json:"displayName"`

	// Tên văn hóa chuẩn (ví dụ: "en-US", "vi-VN").
	Name string `json:"name"`

	// Ngôn ngữ mẹ đẻ (ví dụ: "English", "Tiếng Việt").
	NativeLanguage string `json:"nativeLanguage"`

	// Tên quốc gia/khu vực mẹ đẻ (ví dụ: "United States", "Việt Nam").
	NativeRegion string `json:"nativeRegion"`

	// Định hướng đọc văn bản (ví dụ: "ltr" - left-to-right, "rtl" - right-to-left).
	TextDirection string `json:"textDirection"` // "ltr" hoặc "rtl"

	// --- Định dạng Ngày và Thời gian ---
	DateTimeFormat DateTimeFormatInfo `json:"dateTimeFormat"`

	// --- Định dạng Số ---
	NumberFormat NumberFormatInfo `json:"numberFormat"`

	// --- Định dạng Tiền tệ ---
	CurrencyFormat CurrencyFormatInfo `json:"currencyFormat"`

	// Các thông tin khác có thể có tùy theo yêu cầu cụ thể:
	// Calendar (lịch sử dụng, ví dụ: GregorianCalendar, JapaneseCalendar)
	// CompareInfo (quy tắc so sánh chuỗi, phân biệt chữ hoa/thường, dấu phụ)
	// Parent (văn hóa gốc mà văn hóa này kế thừa từ đó)
}

// DateTimeFormatInfo chứa thông tin định dạng ngày và thời gian.
type DateTimeFormatInfo struct {
	// Định dạng ngày ngắn (ví dụ: "M/d/yyyy" -> 1/15/2023)
	ShortDatePattern string `json:"shortDatePattern"`
	// Định dạng ngày dài (ví dụ: "dddd, MMMM dd, yyyy" -> Monday, January 15, 2023)
	LongDatePattern string `json:"longDatePattern"`
	// Định dạng thời gian ngắn (ví dụ: "h:mm tt" -> 3:04 PM)
	ShortTimePattern string `json:"shortTimePattern"`
	// Định dạng thời gian dài (ví dụ: "h:mm:ss tt" -> 3:04:05 PM)
	LongTimePattern string `json:"longTimePattern"`
	// Định dạng ngày và thời gian đầy đủ (ví dụ: "dddd, MMMM dd, yyyy h:mm:ss tt")
	FullDateTimePattern string `json:"fullDateTimePattern"`

	// Tên các ngày trong tuần (ví dụ: "Sunday", "Monday", ...)
	DayNames []string `json:"dayNames"`
	// Tên viết tắt các ngày trong tuần (ví dụ: "Sun", "Mon", ...)
	AbbreviatedDayNames []string `json:"abbreviatedDayNames"`
	// Tên các tháng (ví dụ: "January", "February", ...)
	MonthNames []string `json:"monthNames"`
	// Tên viết tắt các tháng (ví dụ: "Jan", "Feb", ...)
	AbbreviatedMonthNames []string `json:"abbreviatedMonthNames"`

	// Ký hiệu AM (ví dụ: "AM")
	AMDesignator string `json:"amDesignator"`
	// Ký hiệu PM (ví dụ: "PM")
	PMDesignator string `json:"pmDesignator"`

	// Ký tự phân tách ngày (ví dụ: "/")
	DateSeparator string `json:"dateSeparator"`
	// Ký tự phân tách thời gian (ví dụ: ":")
	TimeSeparator string `json:"timeSeparator"`

	// Ngày đầu tiên của tuần (ví dụ: time.Sunday)
	FirstDayOfWeek time.Weekday `json:"firstDayOfWeek"`
}

// NumberFormatInfo chứa thông tin định dạng số.
type NumberFormatInfo struct {
	// Ký tự phân tách thập phân (ví dụ: "." trong "123.45")
	DecimalSeparator string `json:"decimalSeparator"`
	// Ký tự phân tách hàng nghìn (ví dụ: "," trong "1,234,567")
	GroupSeparator string `json:"groupSeparator"`
	// Kích thước nhóm chữ số (ví dụ: [3] cho 1,234,567; [3,2] cho 12,34,567 - Ấn Độ)
	GroupSizes []int `json:"groupSizes"`
	// Số chữ số thập phân mặc định (ví dụ: 2 cho 123.45)
	DecimalDigits int `json:"decimalDigits"`

	// Ký hiệu dương (ví dụ: "+")
	PositiveSign string `json:"positiveSign"`
	// Ký hiệu âm (ví dụ: "-")
	NegativeSign string `json:"negativeSign"`
	// Dạng số âm (ví dụ: "(1.0)", "-1.0")
	NegativePattern int `json:"negativePattern"` // 0: (n), 1: -n, 2: - n, 3: n-, 4: n -

	// Dạng số dương (ít dùng, thường là 0)
	PositivePattern int `json:"positivePattern"` // 0: n, 1: +n, 2: + n, 3: n+, 4: n +
}

// CurrencyFormatInfo chứa thông tin định dạng tiền tệ.
type CurrencyFormatInfo struct {
	// Ký hiệu tiền tệ (ví dụ: "$", "₫")
	CurrencySymbol string `json:"currencySymbol"`
	// Mã ISO 4217 của tiền tệ (ví dụ: "USD", "VND")
	ISOCurrencySymbol string `json:"isoCurrencySymbol"`
	// Ký tự phân tách thập phân tiền tệ (thường giống DecimalSeparator)
	CurrencyDecimalSeparator string `json:"currencyDecimalSeparator"`
	// Ký tự phân tách hàng nghìn tiền tệ (thường giống GroupSeparator)
	CurrencyGroupSeparator string `json:"currencyGroupSeparator"`
	// Kích thước nhóm chữ số tiền tệ
	CurrencyGroupSizes []int `json:"currencyGroupSizes"`
	// Số chữ số thập phân tiền tệ (ví dụ: 2 cho USD, 0 cho VND)
	CurrencyDecimalDigits int `json:"currencyDecimalDigits"`

	// Vị trí ký hiệu tiền tệ cho số dương (ví dụ: "$1.0", "1.0 $")
	PositiveCurrencyPattern int `json:"positiveCurrencyPattern"` // 0: $n, 1: n$, 2: $ n, 3: n $
	// Vị trí ký hiệu tiền tệ cho số âm (ví dụ: "($1.0)", "-$1.0")
	NegativeCurrencyPattern int `json:"negativeCurrencyPattern"` // 0: ($n), 1: -$n, 2: $-n, 3: $n-, 4: (n$), 5: -n$, 6: n-$, 7: n$-

}
