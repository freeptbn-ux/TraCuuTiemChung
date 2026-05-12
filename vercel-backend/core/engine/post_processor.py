from datetime import date, timedelta
from ..utils import normalize_vaccine_name

def get_vaccine_live_status_by_norm_name(norm_name, vaccine_rules):
    """Kiểm tra xem vắc-xin có phải là vắc-xin sống hay không dựa trên tên chuẩn hóa."""
    for rule in vaccine_rules.values():
        if norm_name in rule.get("names_norm", []) and rule.get("is_live"):
            return True
        if norm_name in rule.get("names_norm_group", []) and rule.get("is_live"):
            return True
        if "courses" in rule:
            for course in rule.get("courses", []):
                if norm_name in course.get("names_norm", []) and course.get("is_live"):
                    return True
    return False

def is_missing_item_live(item, vaccine_rules):
    """Kiểm tra xem mũi tiêm đang thiếu có phải là vắc-xin sống không."""
    popup_name = item.get("vaccine_name_for_popup", "")
    description = item.get("description", "").lower()
    if not popup_name:
        return False

    for rule_key, rule in vaccine_rules.items():
        rule_display_name = rule.get("group_display_name", rule.get("display_name"))

        if rule_display_name and rule_display_name in popup_name:
            if rule_key == "JE_Group":
                live_mentioned = "imojev" in description or "jeev" in description
                inactivated_mentioned = "jevax" in description or "vnnb" in description
                if live_mentioned and not inactivated_mentioned: return True
                if inactivated_mentioned and not live_mentioned: return False

            if rule.get("is_live"): return True 
            
            if "courses" in rule:
                for course in rule.get("courses", []):
                    if course.get("display") and course.get("display").lower() in description:
                        return course.get("is_live", False)
                if any(course.get("is_live") for course in rule.get("courses", [])):
                    if "jevax" in popup_name.lower() or "vnnb" in popup_name.lower():
                        return False 
                    return True
    return False

def apply_spacing_and_sort(missing_items, administered_map, vaccine_rules, analysis_date):
    """Áp dụng quy tắc khoảng cách giữa các loại vắc-xin (sống vs sống, sống vs chết) và sắp xếp kết quả."""
    def sort_key(item):
        return (item.get("earliest_next_dose_date") or date.max, item.get("description", "").lower())

    if not administered_map or not missing_items:
        missing_items.sort(key=sort_key)
        return missing_items

    all_records = [record for records_list in administered_map.values() for record in records_list]
    if not all_records:
        missing_items.sort(key=sort_key)
        return missing_items
        
    all_records.sort(key=lambda x: x["date"])

    last_live_vaccine_date = None
    last_overall_vaccine_date = all_records[-1]["date"]

    for record in reversed(all_records):
        norm_name = normalize_vaccine_name(record["raw_name"])
        if get_vaccine_live_status_by_norm_name(norm_name, vaccine_rules):
            last_live_vaccine_date = record["date"]
            break

    adjusted_items = []
    for item in missing_items:
        new_item = item.copy()
        original_date = new_item.get("earliest_next_dose_date")

        if original_date is None:
            adjusted_items.append(new_item)
            continue
            
        # Quy tắc: cách mũi tiêm gần nhất 14 ngày (chung cho các loại vắc xin)
        potential_dates = [original_date, last_overall_vaccine_date + timedelta(days=14)]

        # Quy tắc: vắc-xin sống cách vắc-xin sống 28 ngày
        if is_missing_item_live(new_item, vaccine_rules) and last_live_vaccine_date:
            potential_dates.append(last_live_vaccine_date + timedelta(days=28))
        
        final_date = max(potential_dates)
        new_item["earliest_next_dose_date"] = max(final_date, analysis_date)
        adjusted_items.append(new_item)

    adjusted_items.sort(key=sort_key)
    return adjusted_items
