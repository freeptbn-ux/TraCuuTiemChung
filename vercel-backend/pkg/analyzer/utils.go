package analyzer

import (
	"regexp"
	"strings"
	"time"
)

// NormalizeVaccineName chuẩn hóa tên vắc xin để so khớp với rules.
func NormalizeVaccineName(name string) string {
	// Chuyển "0,5ml" thành "0.5ml"
	name = strings.ReplaceAll(name, ",", ".")

	// Xóa hậu tố liều lượng như " 3mcg/0.5ml"
	reDose := regexp.MustCompile(`(?i)\s*\d+mcg/\d+(\.\d+)?ml\s*$`)
	name = reDose.ReplaceAllString(name, "")

	// Xóa văn bản trong ngoặc đơn
	reParen := regexp.MustCompile(`\s*\(.*?\)\s*`)
	name = reParen.ReplaceAllString(name, " ")

	// Xóa hậu tố năm như 2023/2024
	reYear := regexp.MustCompile(`\s+\d{4}/\d{4}\s*$`)
	name = reYear.ReplaceAllString(name, "")

	// Xóa hậu tố 20XX/20XX
	reXX := regexp.MustCompile(`\s+20XX/20XX\s*$`)
	name = reXX.ReplaceAllString(name, "")

	// Xóa hậu tố thể tích như 0.5ml
	reVol := regexp.MustCompile(`\s+\d+(\.\d+)?ml\s*$`)
	name = reVol.ReplaceAllString(name, "")

	// Xóa hậu tố -TCDV hoặc -TCMR
	reSuffix := regexp.MustCompile(`(?i)-(TCDV|TCMR)$`)
	name = reSuffix.ReplaceAllString(name, "")

	return strings.TrimSpace(strings.ToLower(name))
}

// ParseDateDDMMYYYY chuyển đổi chuỗi ngày DD/MM/YYYY sang time.Time.
func ParseDateDDMMYYYY(dateStr string) (time.Time, error) {
	dateStr = strings.TrimSpace(dateStr)
	// Remove extra spaces inside date if any
	dateStr = strings.ReplaceAll(dateStr, " ", "")
	return time.Parse("02/01/2006", dateStr)
}

// GetAgeAtDate tính tuổi tại một thời điểm nhất định.
func GetAgeAtDate(dob, targetDate time.Time) (months, days, years int) {
	if targetDate.Before(dob) {
		return 0, 0, 0
	}

	days = int(targetDate.Sub(dob).Hours() / 24)

	years = targetDate.Year() - dob.Year()
	if targetDate.Month() < dob.Month() || (targetDate.Month() == dob.Month() && targetDate.Day() < dob.Day()) {
		years--
	}
	if years < 0 {
		years = 0
	}

	months = (targetDate.Year()-dob.Year())*12 + int(targetDate.Month()-dob.Month())
	if targetDate.Day() < dob.Day() {
		months--
	}
	if months < 0 {
		months = 0
	}

	return
}

// AddMonths cộng thêm tháng vào một ngày, xử lý trường hợp ngày cuối tháng.
func AddMonths(sourceDate time.Time, months int) time.Time {
	year := sourceDate.Year()
	month := int(sourceDate.Month()) + months
	day := sourceDate.Day()

	year += (month - 1) / 12
	month = (month-1)%12 + 1
	if month <= 0 {
		month += 12
		year--
	}

	lastDay := lastDayOfMonth(year, time.Month(month))
	if day > lastDay {
		day = lastDay
	}

	return time.Date(year, time.Month(month), day, sourceDate.Hour(), sourceDate.Minute(), sourceDate.Second(), sourceDate.Nanosecond(), sourceDate.Location())
}

func lastDayOfMonth(year int, month time.Month) int {
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}

// AddYears cộng thêm năm vào một ngày.
func AddYears(sourceDate time.Time, years int) time.Time {
	year := sourceDate.Year() + years
	month := sourceDate.Month()
	day := sourceDate.Day()

	lastDay := lastDayOfMonth(year, month)
	if day > lastDay {
		day = lastDay
	}

	return time.Date(year, month, day, sourceDate.Hour(), sourceDate.Minute(), sourceDate.Second(), sourceDate.Nanosecond(), sourceDate.Location())
}
