from bs4 import BeautifulSoup
from collections import defaultdict
from datetime import datetime, date, timezone, timedelta
from .utils import normalize_vaccine_name

# HTML Element IDs
HTML_PATIENT_NAME_ID = 'txtHoTen'
HTML_PATIENT_DOB_ID = 'txtNgaySinh'
HTML_PATIENT_DOB_HF_ID = 'hfNgaySinhDoiTuong'
HTML_SYSTEM_DATE_ID = 'CurrentSystemDate'
HTML_SYSTEM_DATE_HF_ID = 'hfNgayHienTai'
HTML_VACCINE_TABLE_ID = 'tblVacxin'

class HTMLVaccineParser:
    def __init__(self, html_content: str):
        self.soup = BeautifulSoup(html_content, 'html.parser')
        self.normalize_vaccine_name = normalize_vaccine_name

    def extract_patient_info(self):
        patient_name = None
        patient_dob_str = None
        system_date_str = None

        # Tên bệnh nhân
        name_input = self.soup.find('input', id=HTML_PATIENT_NAME_ID)
        if name_input and 'value' in name_input.attrs:
            patient_name = name_input['value'].strip()
        elif self.soup.find('span', id=HTML_PATIENT_NAME_ID):
            patient_name = self.soup.find('span', id=HTML_PATIENT_NAME_ID).text.strip()

        # Ngày sinh
        dob_input = self.soup.find('input', id=HTML_PATIENT_DOB_ID)
        if not dob_input:
            dob_input = self.soup.find('input', id=HTML_PATIENT_DOB_HF_ID)
        
        if dob_input and 'value' in dob_input.attrs:
            patient_dob_str = dob_input['value'].strip()
        elif self.soup.find('span', id=HTML_PATIENT_DOB_ID):
            patient_dob_str = self.soup.find('span', id=HTML_PATIENT_DOB_ID).text.strip()
        
        # Ngày hệ thống
        system_date_input = self.soup.find('input', id=HTML_SYSTEM_DATE_ID)
        if not system_date_input:
            system_date_input = self.soup.find('input', id=HTML_SYSTEM_DATE_HF_ID)
        
        if system_date_input and 'value' in system_date_input.attrs:
            system_date_str = system_date_input['value'].strip()
        elif self.soup.find('span', id=HTML_SYSTEM_DATE_ID):
            system_date_str = self.soup.find('span', id=HTML_SYSTEM_DATE_ID).text.strip()
        else:
            # Lấy ngày hiện tại theo GMT+7 nếu không có trong HTML
            utc_now = datetime.now(timezone.utc)
            gmt7_now = utc_now.astimezone(timezone(timedelta(hours=7)))
            system_date_str = gmt7_now.strftime("%d/%m/%Y")
        
        return {
            "name": patient_name,
            "birth": patient_dob_str,
            "system_date": system_date_str
        }

    def extract_vaccine_records(self):
        vaccine_table = self.soup.find('table', id=HTML_VACCINE_TABLE_ID)
        records = []

        if not vaccine_table:
            return records
        
        # Một số trang portal không có tbody mà tr nằm trực tiếp trong table hoặc dùng thead/tbody
        rows = vaccine_table.find_all('tr')
        
        for row in rows:
            cols = row.find_all('td')
            # Bỏ qua hàng tiêu đề hoặc hàng không đủ cột
            if len(cols) >= 5:
                # Cột 1: STT, Cột 2: Tên vắc xin, Cột 3: Mũi tiêm, Cột 4: Ngày tiêm (portal cũ) hoặc Cột 5: Ngày tiêm
                # Thường cột 2 là tên, cột 3 là mũi, cột 5 (index 4) là ngày tiêm
                # Chỉ lấy text trực tiếp của cột 2, bỏ qua các thẻ con như <span class="sublabel">
                # Điều này giúp tránh việc gộp tên vắc-xin (Varivax) với tên bệnh (Thủy đậu) thành "Varivax Thủy đậu"
                direct_text = cols[1].find(string=True, recursive=False)
                if direct_text:
                    vaccine_name_raw = direct_text.strip()
                else:
                    vaccine_name_raw = cols[1].get_text(separator=" ", strip=True)
                dose_text_raw = cols[2].get_text(strip=True)
                date_text_raw = cols[4].get_text(strip=True)
                
                if vaccine_name_raw and date_text_raw:
                    try:
                        dose_number_int = int(dose_text_raw)
                    except ValueError:
                        dose_number_int = 0
                    
                    try:
                        # Chuẩn hóa ngày tiêm
                        date_obj = datetime.strptime(date_text_raw.replace(" ",""), "%d/%m/%Y").date()
                        records.append({
                            "raw_name": vaccine_name_raw,
                            "dose_number": dose_number_int,
                            "date": date_obj
                        })
                    except ValueError:
                        pass
        
        return records
