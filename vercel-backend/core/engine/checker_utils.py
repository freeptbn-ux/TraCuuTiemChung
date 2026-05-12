import calendar
from datetime import date, timedelta
from ..utils import get_age_at_date

GRACE_PERIOD_DAYS = 0

def add_months(source_date: date, months: int) -> date:
    """Cộng thêm tháng vào một ngày, xử lý trường hợp ngày cuối tháng."""
    month = source_date.month - 1 + months
    year = source_date.year + month // 12
    month = month % 12 + 1
    day = source_date.day
    try:
        return date(year, month, day)
    except ValueError:
        # Trường hợp ngày 31 cộng 1 tháng (tháng sau chỉ có 30 ngày)
        last_day_of_month = calendar.monthrange(year, month)[1]
        return date(year, month, last_day_of_month)

def add_years(source_date: date, years_to_add: int) -> date:
    """Cộng thêm năm vào một ngày."""
    new_year = source_date.year + years_to_add
    try:
        return date(new_year, source_date.month, source_date.day)
    except ValueError:
        # Xử lý nhuận 29/02
        return date(new_year, source_date.month, 28)

def get_administered_dose_records(names_norm_list, administered_map):
    """Lấy tất cả các bản ghi cho các vắc xin trong list, sắp xếp theo ngày."""
    records = []
    if not names_norm_list:
        return records
    for norm_name in names_norm_list:
        records.extend(administered_map.get(norm_name, []))
    # Bản ghi có dạng dict: {"dose_number": int, "date": date, ...}
    records.sort(key=lambda x: x["date"])
    return records

def get_age_status_and_earliest_date(dob, analysis_date, rule_details, for_group_display_name=""):
    """Xác định trạng thái tuổi và ngày tiêm sớm nhất có thể."""
    display_name_prefix = f"{for_group_display_name} - " if for_group_display_name else ""

    if not dob:
        return (f"{display_name_prefix}Không có ngày sinh để kiểm tra tuổi", None, ["error_dob"])

    age_at_analysis_months, age_at_analysis_total_days, age_at_analysis_years = get_age_at_date(dob, analysis_date)

    if age_at_analysis_months is None:
        return (f"{display_name_prefix}Ngày phân tích trước ngày sinh", None, ["error_date"])

    min_age_months = rule_details.get("min_age_months_overall_group") or \
                     rule_details.get("min_age_months_at_first_dose") or \
                     rule_details.get("min_age_months_overall")
    min_age_weeks = rule_details.get("min_age_weeks_overall_group") or \
                    rule_details.get("min_age_weeks_at_first_dose")
    min_age_years = rule_details.get("min_age_years_overall_group") or \
                    rule_details.get("min_age_years_at_first_dose")
    min_age_days_val = rule_details.get("min_age_days_overall_group") or \
                       rule_details.get("min_age_days_at_first_dose")

    earliest_acceptable_date = None
    status_message = ""
    status_tags = ["eligible"]

    if min_age_months is not None:
        target_min_age_date = add_months(dob, min_age_months)
        earliest_acceptable_date = target_min_age_date - timedelta(days=GRACE_PERIOD_DAYS)
        if analysis_date < earliest_acceptable_date:
            status_message = f"cần {min_age_months} tháng tuổi"
            status_tags = ["too_young"]
    elif min_age_years is not None:
        target_min_age_date = add_years(dob, min_age_years)
        earliest_acceptable_date = target_min_age_date - timedelta(days=GRACE_PERIOD_DAYS)
        if analysis_date < earliest_acceptable_date:
            status_message = f"cần {min_age_years} tuổi"
            status_tags = ["too_young"]
    elif min_age_weeks is not None:
        required_total_days_for_weeks = min_age_weeks * 7
        effective_min_total_days = required_total_days_for_weeks - GRACE_PERIOD_DAYS
        if age_at_analysis_total_days < effective_min_total_days:
            earliest_acceptable_date = dob + timedelta(days=effective_min_total_days)
            status_message = f"cần {min_age_weeks} tuần tuổi"
            status_tags = ["too_young"]
    elif min_age_days_val is not None:
        effective_min_total_days = min_age_days_val - GRACE_PERIOD_DAYS
        if age_at_analysis_total_days < effective_min_total_days:
            earliest_acceptable_date = dob + timedelta(days=effective_min_total_days)
            display_req = f"{min_age_days_val // 30} tháng" if min_age_days_val >= 60 else f"{min_age_days_val} ngày"
            status_message = f"cần >{display_req} tuổi"
            status_tags = ["too_young"]

    if "too_young" in status_tags:
        return status_message, earliest_acceptable_date, status_tags
    return "đủ điều kiện tuổi", analysis_date if dob else None, status_tags

def check_first_dose_age_validity(dob, first_dose_date, rule_details, rule_display_name, missing_items_list):
    """Kiểm tra tuổi mũi 1."""
    if not dob: return True

    age_at_first_dose_months, age_at_first_dose_total_days, age_at_first_dose_years = get_age_at_date(dob, first_dose_date)

    if age_at_first_dose_months is None:
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name,
            "description": f"{rule_display_name} - Lỗi tính tuổi cho mũi đầu (ngày tiêm có thể trước ngày sinh).",
            "earliest_next_dose_date": None,
            "status_tags": ["error_age_calculation"]
        })
        return False

    min_age_months = rule_details.get("min_age_months_overall_group") or \
                     rule_details.get("min_age_months_at_first_dose") or \
                     rule_details.get("min_age_months_overall")
    min_age_weeks = rule_details.get("min_age_weeks_overall_group") or \
                    rule_details.get("min_age_weeks_at_first_dose")
    min_age_years = rule_details.get("min_age_years_overall_group") or \
                    rule_details.get("min_age_years_at_first_dose")
    min_age_days_val = rule_details.get("min_age_days_overall_group") or \
                       rule_details.get("min_age_days_at_first_dose")

    error_detail = None
    if min_age_days_val is not None:
        effective_min_days = min_age_days_val - GRACE_PERIOD_DAYS
        if age_at_first_dose_total_days < effective_min_days:
            display_age = f"{min_age_days_val // 30} tháng" if min_age_days_val >= 60 else f"{min_age_days_val} ngày"
            error_detail = f"Mũi 1 tiêm quá sớm (cần >{display_age}, thực tế {age_at_first_dose_total_days} ngày tuổi)."
    elif min_age_weeks is not None:
        effective_min_days = (min_age_weeks * 7) - GRACE_PERIOD_DAYS
        if age_at_first_dose_total_days < effective_min_days:
            error_detail = f"Mũi 1 tiêm quá sớm (cần {min_age_weeks} tuần, thực tế {age_at_first_dose_total_days} ngày tuổi)."
    elif min_age_months is not None:
        earliest_allowed = add_months(dob, min_age_months) - timedelta(days=GRACE_PERIOD_DAYS)
        if first_dose_date < earliest_allowed:
            error_detail = f"Mũi 1 tiêm quá sớm (cần {min_age_months} tháng tuổi)."
    elif min_age_years is not None:
        earliest_allowed = add_years(dob, min_age_years) - timedelta(days=GRACE_PERIOD_DAYS)
        if first_dose_date < earliest_allowed:
            error_detail = f"Mũi 1 tiêm quá sớm (cần {min_age_years} tuổi)."

    if error_detail:
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name,
            "description": f"{rule_display_name} - {error_detail}",
            "earliest_next_dose_date": None,
            "status_tags": ["error_age_first_dose", "too_early"]
        })
        return False
    return True
