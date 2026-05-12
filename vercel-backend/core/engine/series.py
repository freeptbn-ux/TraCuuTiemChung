from datetime import date, timedelta
from ..utils import get_age_at_date
from .checker_utils import (
    GRACE_PERIOD_DAYS, add_months, add_years,
    get_administered_dose_records, get_age_status_and_earliest_date,
    check_first_dose_age_validity
)

# Constants for Rule Types (moved from config_data)
RULE_TYPE_MMR_EQUIVALENT_GROUP = "mmr_equivalent_group"

def check_single_vaccine_series(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None):
    """Kiểm tra phác đồ đơn lẻ (hoặc chuỗi liều đơn)."""
    rule_display_name = rule_details.get("display_name", rule_key)
    doses_required = rule_details.get("doses_required", 0)
    names_norm_list = rule_details.get("names_norm", [])

    administered_records = get_administered_dose_records(names_norm_list, administered_map)

    # Đếm số mũi đã tiêm
    valid_doses_count = len(administered_records)
    last_valid_dose_date = administered_records[-1]["date"] if administered_records else None
    
    # Nếu đủ mũi, kiểm tra mũi nhắc (booster)
    if doses_required > 0 and valid_doses_count >= doses_required:
        booster_interval_years = rule_details.get("booster_interval_years")
        booster_applies_after_dose_num = rule_details.get("booster_after_dose_number", doses_required)
        booster_max_age_years = rule_details.get("booster_max_age_years")

        if booster_interval_years and valid_doses_count >= booster_applies_after_dose_num:
            date_for_booster_base = administered_records[-1]["date"] 

            if date_for_booster_base:
                next_booster_due_date = add_years(date_for_booster_base, booster_interval_years)

                if dob and booster_max_age_years is not None:
                    _, _, age_at_booster_years = get_age_at_date(dob, next_booster_due_date)
                    if age_at_booster_years is not None and age_at_booster_years >= booster_max_age_years:
                        return 
                
                description_suffix = None
                current_status_tags = []
                actual_next_booster_date_for_listing = None

                if analysis_date >= next_booster_due_date:
                    actual_next_booster_date_for_listing = analysis_date
                    current_status_tags = ["due", "booster_due"]
                    description_suffix = f"Cần tiêm mũi nhắc lại định kỳ {booster_interval_years} năm (đã đến hạn)."
                else:
                    actual_next_booster_date_for_listing = next_booster_due_date
                    current_status_tags = ["info", "booster_upcoming"]
                    description_suffix = f"Cần tiêm mũi nhắc lại định kỳ {booster_interval_years} năm."

                if description_suffix:
                    missing_items_list.append({
                        "vaccine_name_for_popup": rule_display_name,
                        "description": f"{rule_display_name} - {description_suffix}",
                        "earliest_next_dose_date": actual_next_booster_date_for_listing,
                        "status_tags": current_status_tags
                    })
        return

    # MVVAC vs MMR interaction
    if rule_key == "MVVAC" and all_vaccine_rules and dob:
        measles_covered_by_other_vaccine_info = None
        for other_key, other_rule_details_item in all_vaccine_rules.items():
            if other_key == rule_key: continue
            provides_protection = False
            names_to_check_for_coverage = []
            group_display_name_for_msg = other_rule_details_item.get("display_name", other_rule_details_item.get("group_display_name", other_key))

            if other_rule_details_item.get("type") == RULE_TYPE_MMR_EQUIVALENT_GROUP:
                if other_rule_details_item.get("provides_measles_protection_group"):
                    provides_protection = True
                    names_to_check_for_coverage = other_rule_details_item.get("names_norm_group", [])
            elif other_rule_details_item.get("provides_measles_protection"):
                provides_protection = True
                names_to_check_for_coverage = other_rule_details_item.get("names_norm", [])
            
            if provides_protection and names_to_check_for_coverage:
                other_records = get_administered_dose_records(names_to_check_for_coverage, administered_map)
                if other_records:
                    for record in other_records:
                        age_at_mmr_months, _, _ = get_age_at_date(dob, record["date"])
                        if age_at_mmr_months is not None and age_at_mmr_months >= 9:
                            measles_covered_by_other_vaccine_info = group_display_name_for_msg
                            break
                if measles_covered_by_other_vaccine_info:
                    break
        
        if measles_covered_by_other_vaccine_info:
            administered_mvvac_records = get_administered_dose_records(names_norm_list, administered_map)
            if not administered_mvvac_records:
                missing_items_list.append({
                    "vaccine_name_for_popup": rule_display_name,
                    "description": f"{rule_display_name}: Miễn dịch sởi có thể đã được bao phủ bởi {measles_covered_by_other_vaccine_info}.",
                    "earliest_next_dose_date": None,
                    "status_tags": ["info", "coverage_by_other"]
                })
            return

    # VA-MENGOC-BC reverse interaction
    if rule_key == "VA-MENGOC-BC" and all_vaccine_rules and dob:
        age_months, _, _ = get_age_at_date(dob, analysis_date)
        if age_months is not None and age_months >= 24:
            acyw_rule = all_vaccine_rules.get("MeningococcalACYW_Group", {})
            acyw_names = []
            if acyw_rule:
                for member in acyw_rule.get("members", {}).values():
                    acyw_names.extend(member.get("names_norm", []))
            
            acyw_records = get_administered_dose_records(acyw_names, administered_map)
            if acyw_records:
                missing_items_list.append({
                    "vaccine_name_for_popup": rule_display_name,
                    "description": f"Lưu ý: Trẻ từ 2 tuổi trở lên nên tiêm VA-Mengoc BC SAU khi đã tiêm MenQuadfi/Menactra để đạt hiệu quả bảo vệ tốt nhất.",
                    "earliest_next_dose_date": None,
                    "status_tags": ["warning", "interaction_reverse_mengoc"]
                })

    if not administered_records:
        if doses_required > 0:
            age_status_msg, earliest_date_for_first_dose, age_tags = get_age_status_and_earliest_date(dob, analysis_date, rule_details, rule_display_name)
            desc = f"{rule_display_name} (Chưa tiêm - cần {doses_required} liều). {age_status_msg}"
            missing_items_list.append({
                "vaccine_name_for_popup": rule_display_name,
                "description": desc,
                "earliest_next_dose_date": earliest_date_for_first_dose,
                "status_tags": age_tags if "error" not in "".join(age_tags) else ["error_initial_check"] + age_tags
            })
        return
        
    if not check_first_dose_age_validity(dob, administered_records[0]["date"], rule_details, rule_display_name, missing_items_list):
         missing_items_list.append({ 
            "vaccine_name_for_popup": rule_display_name,
            "description": f"{rule_display_name} - Cần {doses_required} liều hợp lệ (do mũi đầu tiên không hợp lệ về tuổi). Cân nhắc tiêm lại theo đúng phác đồ.",
            "earliest_next_dose_date": None, 
            "status_tags": ["error_series_restart_needed", "error_age_first_dose"]
        })
         return

    if valid_doses_count < doses_required:
        remaining_doses = doses_required - valid_doses_count
        next_dose_number_display = valid_doses_count + 1
        
        description_parts = [f"{rule_display_name} - Cần thêm {remaining_doses} liều."]
        condition_parts = []
        earliest_next_dose_date = None
        status_tags_for_missing = ["due"]
        
        next_dose_rule_idx = valid_doses_count
        
        min_intervals_for_next_dose_series = rule_details.get("min_interval_days", [])
        date_by_interval = None

        if last_valid_dose_date and min_intervals_for_next_dose_series and \
           next_dose_rule_idx < len(min_intervals_for_next_dose_series) and \
           min_intervals_for_next_dose_series[next_dose_rule_idx] is not None:
            
            interval_days_for_next = min_intervals_for_next_dose_series[next_dose_rule_idx]
            date_by_interval = last_valid_dose_date + timedelta(days=interval_days_for_next)
            
            unit, val = "ngày", interval_days_for_next
            if interval_days_for_next >= 365 * 2: unit, val = "năm", interval_days_for_next // 365
            elif interval_days_for_next >= 30 * 2: unit, val = "tháng", interval_days_for_next // 30
            elif interval_days_for_next >= 7 * 2: unit, val = "tuần", interval_days_for_next // 7
            condition_parts.append(f"Mũi {next_dose_number_display} cách mũi {valid_doses_count} tối thiểu {val} {unit}")
        elif not last_valid_dose_date and doses_required > 0:
            age_status_msg, earliest_date_calc, age_tags_calc = get_age_status_and_earliest_date(dob, analysis_date, rule_details, rule_display_name)
            condition_parts.append(age_status_msg)
            status_tags_for_missing = age_tags_calc
            earliest_next_dose_date = earliest_date_calc

        dose_specific_options_for_next_scheduled_dose = rule_details.get("dose_specific_rules", {}).get(str(next_dose_number_display))
        date_by_alt_age_years = None
        date_by_abs_min_age_months = None

        if dose_specific_options_for_next_scheduled_dose and dob:
            alt_min_y = dose_specific_options_for_next_scheduled_dose.get("alternative_min_age_years")
            alt_max_y = dose_specific_options_for_next_scheduled_dose.get("alternative_max_age_years")
            if alt_min_y is not None and alt_max_y is not None:
                date_by_alt_age_years = add_years(dob, alt_min_y)
                condition_parts.append(f"HOẶC khi trẻ {alt_min_y}-{alt_max_y-1} tuổi")

            abs_min_m = dose_specific_options_for_next_scheduled_dose.get("min_absolute_age_months")
            if abs_min_m is not None:
                date_by_abs_min_age_months = add_months(dob, abs_min_m)
                condition_parts.append(f"và trẻ ít nhất {abs_min_m} tháng tuổi")
        
        if earliest_next_dose_date is None and date_by_interval:
            earliest_next_dose_date = date_by_interval
        
        if date_by_abs_min_age_months:
            if earliest_next_dose_date is None or date_by_abs_min_age_months > earliest_next_dose_date:
                earliest_next_dose_date = date_by_abs_min_age_months
        
        if date_by_alt_age_years:
            effective_alt_date = date_by_alt_age_years
            if date_by_abs_min_age_months and effective_alt_date < date_by_abs_min_age_months: 
                effective_alt_date = date_by_abs_min_age_months
            
            if earliest_next_dose_date is None or effective_alt_date < earliest_next_dose_date:
                 earliest_next_dose_date = effective_alt_date
        
        if earliest_next_dose_date and earliest_next_dose_date < analysis_date:
            earliest_next_dose_date = analysis_date
        
        if not condition_parts and doses_required > valid_doses_count:
            condition_parts.append(f"Mũi {next_dose_number_display}")
        
        description_parts.append(" ".join(condition_parts) + ".")
        
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name,
            "description": " ".join(description_parts),
            "earliest_next_dose_date": earliest_next_dose_date,
            "status_tags": status_tags_for_missing
        })

def check_age_dependent_series(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None):
    """Kiểm tra phác đồ phụ thuộc tuổi bắt đầu tiêm."""
    rule_display_name = rule_details["display_name"]
    
    current_rule_administered_records = get_administered_dose_records(rule_details.get("names_norm", []), administered_map)

    if not current_rule_administered_records:
        age_status_msg, earliest_date, age_tags = get_age_status_and_earliest_date(dob, analysis_date, rule_details, rule_display_name)
        default_doses = "một số"
        if rule_details.get("rules_by_age") and rule_details["rules_by_age"]:
            default_doses = rule_details["rules_by_age"][0].get("doses_required", "một số")
        
        desc = f"{rule_display_name} (Chưa tiêm - cần {default_doses} liều theo phác đồ). {age_status_msg}."
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name,
            "description": desc, "earliest_next_dose_date": earliest_date, "status_tags": age_tags
        })
        return

    if not dob:
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name,
            "description": f"{rule_display_name} - Không thể xác định phác đồ do thiếu ngày sinh.",
            "earliest_next_dose_date": None, "status_tags": ["error_dob"]
        })
        return
        
    first_dose_date = current_rule_administered_records[0]["date"]
    age_at_first_dose_months, _, _ = get_age_at_date(dob, first_dose_date)

    if age_at_first_dose_months is None:
         missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name,
            "description": f"{rule_display_name} - Lỗi tính tuổi tại mũi đầu tiên.",
            "earliest_next_dose_date": None, "status_tags": ["error_age_calculation"]
         })
         return
    
    overall_age_check_details = {
        "min_age_months_at_first_dose": rule_details.get("min_age_months_overall"),
        "min_age_weeks_at_first_dose": rule_details.get("min_age_weeks_overall"),
        "min_age_years_at_first_dose": rule_details.get("min_age_years_overall"),
        "min_age_days_at_first_dose": rule_details.get("min_age_days_overall")
    }
    overall_age_check_details = {k:v for k,v in overall_age_check_details.items() if v is not None}

    is_first_dose_valid_for_group_initiation = check_first_dose_age_validity(dob, first_dose_date, overall_age_check_details, rule_display_name, [])
    
    applicable_age_rule = None
    if not rule_details.get("rules_by_age"):
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name,
            "description": f"{rule_display_name} - Lỗi cấu hình: không có 'rules_by_age'.",
            "earliest_next_dose_date": None, "status_tags": ["error_config"]})
        return

    for age_rule in rule_details["rules_by_age"]:
        min_m = age_rule.get("min_age_at_first_dose_months")
        max_m = age_rule.get("max_age_at_first_dose_months")  
        condition_met = True
        if min_m is not None and age_at_first_dose_months < min_m: condition_met = False
        if max_m is not None and age_at_first_dose_months > max_m: condition_met = False
        if condition_met:
            applicable_age_rule = age_rule
            break
    
    if not applicable_age_rule:
        if not is_first_dose_valid_for_group_initiation:
             check_first_dose_age_validity(dob, first_dose_date, overall_age_check_details, rule_display_name, missing_items_list)
             default_doses_needed = "một số"
             if rule_details.get("rules_by_age") and rule_details["rules_by_age"]:
                default_doses_needed = rule_details["rules_by_age"][0].get("doses_required", "một số")
             missing_items_list.append({
                "vaccine_name_for_popup": rule_display_name,
                "description": f"{rule_display_name} - Cần {default_doses_needed} liều (do mũi đầu không hợp lệ về tuổi chung của nhóm). Cân nhắc tiêm lại theo đúng phác đồ.",
                "earliest_next_dose_date": None, 
                "status_tags": ["error_age_first_dose", "series_restart_needed"]
            })
        else: 
            missing_items_list.append({
                "vaccine_name_for_popup": rule_display_name,
                "description": f"{rule_display_name} - Không tìm thấy quy tắc tuổi phù hợp cho mũi đầu tiên (trẻ {age_at_first_dose_months} tháng tuổi khi tiêm).",
                "earliest_next_dose_date": None, "status_tags": ["error_no_matching_rule"]})
        return

    temp_rule_for_check_name = f"{rule_display_name} ({applicable_age_rule.get('regimen_name', f'phác đồ cho mũi 1 lúc {age_at_first_dose_months} tháng')})"
    temp_rule_for_check = {
        "display_name": temp_rule_for_check_name, 
        "names_norm": rule_details.get("names_norm", []), 
        "doses_required": applicable_age_rule["doses_required"], 
        "min_interval_days": applicable_age_rule.get("min_interval_days", [None] * applicable_age_rule["doses_required"]), 
        "min_age_months_at_first_dose": applicable_age_rule.get("min_age_months_at_first_dose"), 
        "min_age_weeks_at_first_dose": applicable_age_rule.get("min_age_weeks_at_first_dose"),
        "dose_specific_rules": applicable_age_rule.get("dose_specific_rules", {}),
        "booster_interval_years": applicable_age_rule.get("booster_interval_years"),
        "booster_after_dose_number": applicable_age_rule.get("booster_after_dose_number"),
        "booster_max_age_years": applicable_age_rule.get("booster_max_age_years")
    }
    
    check_single_vaccine_series(f"{rule_key}_age_specific", temp_rule_for_check, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules)
