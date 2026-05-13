package parity

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
)

type PatientInfo struct {
	Name       string
	Birth      string
	SystemDate string
}

type VaccineRecord struct {
	RawName    string
	DoseNumber int
	Date       time.Time
}

func ParseMinhkhoiHTML(filepath string) (*PatientInfo, []VaccineRecord, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()

	doc, err := goquery.NewDocumentFromReader(file)
	if err != nil {
		return nil, nil, err
	}

	patient := &PatientInfo{}

	// Name
	if val, exists := doc.Find("#txtHoTen").Attr("value"); exists {
		patient.Name = strings.TrimSpace(val)
	} else {
		patient.Name = strings.TrimSpace(doc.Find("#txtHoTen").Text())
	}

	// Birth
	if val, exists := doc.Find("#txtNgaySinh").Attr("value"); exists {
		patient.Birth = strings.TrimSpace(val)
	} else if val, exists := doc.Find("#hfNgaySinhDoiTuong").Attr("value"); exists {
		patient.Birth = strings.TrimSpace(val)
	} else {
		patient.Birth = strings.TrimSpace(doc.Find("#txtNgaySinh").Text())
	}

	// System Date
	if val, exists := doc.Find("#CurrentSystemDate").Attr("value"); exists {
		patient.SystemDate = strings.TrimSpace(val)
	} else if val, exists := doc.Find("#hfNgayHienTai").Attr("value"); exists {
		patient.SystemDate = strings.TrimSpace(val)
	} else {
		patient.SystemDate = strings.TrimSpace(doc.Find("#CurrentSystemDate").Text())
	}

	if patient.SystemDate == "" {
		patient.SystemDate = time.Now().Format("02/01/2006")
	}

	var records []VaccineRecord
	doc.Find("#tblVacxin tr").Each(func(i int, s *goquery.Selection) {
		cols := s.Find("td")
		if cols.Length() >= 5 {
			var rawName string
			
			// Try to get just the text of the main node, ignoring children like spans
			// An easy way is to extract text, but skip inner tags.
			clone := cols.Eq(1).Clone()
			clone.Children().Remove()
			rawName = strings.TrimSpace(clone.Text())
			if rawName == "" {
				rawName = strings.TrimSpace(cols.Eq(1).Text())
			}

			doseStr := strings.TrimSpace(cols.Eq(2).Text())
			doseNum, _ := strconv.Atoi(doseStr)

			dateStr := strings.TrimSpace(cols.Eq(4).Text())
			dateStr = strings.ReplaceAll(dateStr, " ", "")

			dateObj, err := time.Parse("02/01/2006", dateStr)
			if err == nil && rawName != "" {
				records = append(records, VaccineRecord{
					RawName:    rawName,
					DoseNumber: doseNum,
					Date:       dateObj,
				})
			}
		}
	})

	return patient, records, nil
}
