from datetime import datetime, date, timedelta, timezone
from collections import defaultdict
import traceback

from .utils import normalize_vaccine_name, get_age_string
from .parser import HTMLVaccineParser
from .rules import VaccineRulesLoader
from .engine.processor import process_all_vaccine_rules
from .engine.post_processor import apply_spacing_and_sort

# Default standard vaccines for analysis if not provided in records
STANDARD_VACCINES_STRING = "MMR-II;Varivax;Influvac Tetra 20XX/20XX;Vaxigrip Tetra 0.5ml;VA - MENGOC - BC;MVVAC;Rota Teq;Rotarix 1.5ml;ROTARIXTM;Rotavin;Rotavin-M1;Typhim Vi (Lọ 1 liều/0.5ml);Morcvax (Lọ 1 liều - 1.5ml);Avaxim 80U;HAVAX;Priorix;JEEV 3mcg/0,5ml"

class AnalyzerService:
    """
    Dịch vụ bao bọc (Wrapper) kết nối Parser, Rules Loader và Engine để thực hiện phân tích đầy đủ.
    """
    def __init__(self, rules_path="assets/vaccine_rules.json"):
        self.rules_loader = VaccineRulesLoader(rules_path)
        self.vaccine_rules = self.rules_loader.get_rules()
        self.standard_vaccines = [v.strip() for v in STANDARD_VACCINES_STRING.split(';') if v.strip()]

    def analyze_from_html(self, html_content):
        """
        Phân tích từ mã nguồn HTML của trang Cổng thông tin tiêm chủng.
        """
        parser = HTMLVaccineParser(html_content)
        patient_info = parser.extract_patient_info()
        administered_list = parser.extract_vaccine_records()
        
        return self.analyze(patient_info, administered_list)

    def analyze(self, patient_info, administered_list):
        """
        Phân tích dựa trên thông tin bệnh nhân và danh sách các mũi tiêm đã thực hiện.
        """
        results = {
            "patient_name": patient_info.get("name"),
            "patient_dob": patient_info.get("birth"),
            "administered": [],
            "recommendations": [],
            "error": None
        }

        try:
            # 1. Chuẩn bị ngày tháng
            dob_str = patient_info.get("birth", "")
            patient_dob_obj = None
            if dob_str:
                try:
                    # Giả định format DD/MM/YYYY từ parser
                    patient_dob_obj = datetime.strptime(dob_str.strip(), "%d/%m/%Y").date()
                except ValueError:
                    pass

            # Sử dụng múi giờ Việt Nam (GMT+7) cho ngày phân tích
            gmt7_now = datetime.now(timezone(timedelta(hours=7)))
            analysis_date = gmt7_now.date()

            # 2. Phân loại dữ liệu đã tiêm
            administered_map = defaultdict(list)
            
            for record in administered_list:
                name = record.get("raw_name", "")
                date_obj = record.get("date")
                dose_int = record.get("dose_number", 0)

                if not name or not date_obj:
                    continue

                norm_name = normalize_vaccine_name(name)
                
                # Lưu vào map phục vụ engine
                administered_map[norm_name].append({
                    "dose_number": dose_int,
                    "date": date_obj,
                    "raw_name": name
                })
                
                # Lưu vào results để trả về frontend
                age_at_vaccination = ""
                if patient_dob_obj:
                    age_at_vaccination = get_age_string(patient_dob_obj, date_obj)
                
                results["administered"].append({
                    "name": name,
                    "dose": dose_int,
                    "date": date_obj.strftime("%d/%m/%Y"),
                    "age": age_at_vaccination
                })

            # Sắp xếp hồ sơ tiêm
            results["administered"].sort(key=lambda x: datetime.strptime(x["date"], "%d/%m/%Y"))
            for key in administered_map:
                administered_map[key].sort(key=lambda x: x["date"])

            # 3. Thực hiện logic phân tích (Engine)
            if not self.vaccine_rules:
                results["error"] = "Dữ liệu quy tắc vắc-xin không khả dụng."
                return results

            # Giai đoạn 1: Chạy các checker
            missing_raw = process_all_vaccine_rules(
                administered_map, 
                self.vaccine_rules, 
                patient_dob_obj, 
                analysis_date, 
                self.standard_vaccines
            )

            # Giai đoạn 2: Hậu xử lý (khoảng cách và sắp xếp)
            missing_final = apply_spacing_and_sort(
                missing_raw, 
                administered_map, 
                self.vaccine_rules, 
                analysis_date
            )

            # 4. Định dạng kết quả gợi ý
            for item in missing_final:
                date_obj = item.get("earliest_next_dose_date")
                date_str = date_obj.strftime("%d/%m/%Y") if date_obj else ""
                
                results["recommendations"].append({
                    "vaccine_name": item.get("vaccine_name_for_popup", ""),
                    "description": item.get("description", ""),
                    "date": date_str,
                    "status_tags": item.get("status_tags", [])
                })

        except Exception as e:
            traceback.print_exc()
            results["error"] = f"Lỗi hệ thống trong quá trình phân tích: {str(e)}"

        return results
