import re
from datetime import datetime, date, timedelta

def normalize_vaccine_name(name: str) -> str:
    """Chuẩn hóa tên vắc xin để so khớp với rules."""
    name = str(name)
    
    # Chuyển "0,5ml" thành "0.5ml"
    name = name.replace(',', '.')
    # Xóa hậu tố liều lượng như " 3mcg/0.5ml"
    name = re.sub(r'\s*\d+mcg/\d+(\.\d+)?ml\s*$', '', name, flags=re.IGNORECASE)

    # Xóa văn bản trong ngoặc đơn
    name = re.sub(r'\s*\(.*?\)\s*', '', name)
    # Xóa hậu tố năm như 2023/2024
    name = re.sub(r'\s+\d{4}/\d{4}\s*$', '', name)
    # Xóa hậu tố 20XX/20XX
    name = re.sub(r'\s+20XX/20XX\s*$', '', name)
    # Xóa hậu tố thể tích như 0.5ml
    name = re.sub(r'\s+\d+(\.\d+)?ml\s*$', '', name)
    
    return name.strip().lower()

def get_age_at_date(dob: date, target_date: date):
    """Tính tuổi tại một thời điểm nhất định."""
    if not dob or not target_date:
        return None, None, None

    total_days = (target_date - dob).days
    if total_days < 0:
        return None, None, None

    years = target_date.year - dob.year
    if (target_date.month, target_date.day) < (dob.month, dob.day):
        years -= 1
    years = max(0, years)

    months_total = (target_date.year - dob.year) * 12 + (target_date.month - dob.month)
    if target_date.day < dob.day:
        months_total -= 1
    months_total = max(0, months_total)
    
    return months_total, total_days, years

def get_age_string(dob_input, current_date_input) -> str:
    """Trả về chuỗi hiển thị tuổi. Chấp nhận string (dd/mm/yyyy) hoặc date object."""
    # Xử lý ngày sinh
    if isinstance(dob_input, (date, datetime)):
        dob = dob_input if isinstance(dob_input, date) else dob_input.date()
    else:
        try:
            dob = datetime.strptime(str(dob_input).strip().replace(" ", ""), "%d/%m/%Y").date()
        except ValueError:
            return "Ngày sinh không hợp lệ"
    
    # Xử lý ngày hiện tại
    if isinstance(current_date_input, (date, datetime)):
        current_dt = current_date_input if isinstance(current_date_input, date) else current_date_input.date()
    else:
        try:
            current_dt = datetime.strptime(str(current_date_input).strip().replace(" ", ""), "%d/%m/%Y").date()
        except ValueError:
            current_dt = date.today()

    if dob > current_dt:
        return "Ngày sinh trong tương lai"

    months, total_days, years = get_age_at_date(dob, current_dt)
    
    if months is None:
         return "Lỗi tính tuổi"

    if months < 1 :
        if total_days < 7:
            return f"{total_days} ngày tuổi"
        else:
            weeks = total_days // 7
            remaining_days = total_days % 7
            if remaining_days == 0:
                return f"{weeks} tuần tuổi"
            else:
                return f"{weeks} tuần {remaining_days} ngày tuổi"
    elif months < 72:
        return f"{months} tháng tuổi"
    else:
        return f"{years} tuổi"
