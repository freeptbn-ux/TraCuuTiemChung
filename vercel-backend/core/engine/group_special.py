from datetime import date, timedelta
from ..utils import get_age_at_date
from .checker_utils import (
    get_administered_dose_records, get_age_status_and_earliest_date,
    check_first_dose_age_validity, add_months, add_years
)
from .series import check_single_vaccine_series, check_age_dependent_series

# Constants from config_data.py
MVVAC_TO_MMR_MIN_INTERVAL_DAYS = 84

def check_mmr_equivalent_group(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None):
    """Kiểm tra nhóm MMR/Priorix (có xét đến tương tác với MVVAC)."""
    group_display_name = rule_details.get("group_display_name", "Nhóm MMR")
    
    mvvac_rule = (all_vaccine_rules or {}).get("MVVAC", {})
    mvvac_names_norm = mvvac_rule.get("names_norm", [])
    mvvac_records = get_administered_dose_records(mvvac_names_norm, administered_map)

    if mvvac_records:
        all_mmr_group_names_norm = rule_details.get("names_norm_group", [])
        all_mmr_group_records = get_administered_dose_records(all_mmr_group_names_norm, administered_map)
        num_mmr_group_doses = len(all_mmr_group_records)

        if num_mmr_group_doses == 0:
            last_mvvac_date = mvvac_records[-1]["date"]
            earliest_by_mvvac_interval = last_mvvac_date + timedelta(days=MVVAC_TO_MMR_MIN_INTERVAL_DAYS)
            
            age_status_msg, earliest_by_age, age_tags = get_age_status_and_earliest_date(
                dob, analysis_date, rule_details, group_display_name
            )
            
            earliest_date = earliest_by_mvvac_interval
            if earliest_by_age and earliest_by_age > earliest_date:
                earliest_date = earliest_by_age
            
            if analysis_date >= earliest_date:
                status_tags = ["due"]
                earliest_date = analysis_date
            else:
                status_tags = ["info", "scheduled"]
            
            desc = (
                f"{group_display_name} (Phác đồ MVVAC + MMR: cần tiêm mũi MMR/Priorix đầu tiên "
                f"sau {MVVAC_TO_MMR_MIN_INTERVAL_DAYS} ngày kể từ mũi MVVAC). "
                f"{age_status_msg}."
            )
            
            missing_items_list.append({
                "vaccine_name_for_popup": group_display_name,
                "description": desc,
                "earliest_next_dose_date": earliest_date,
                "status_tags": status_tags
            })
            return

        elif num_mmr_group_doses == 1:
            first_mmr_date = all_mmr_group_records[0]["date"]
            last_mvvac_date = mvvac_records[-1]["date"]
            
            actual_interval = (first_mmr_date - last_mvvac_date).days
            if actual_interval < MVVAC_TO_MMR_MIN_INTERVAL_DAYS:
                missing_items_list.append({
                    "vaccine_name_for_popup": group_display_name,
                    "description": (
                        f"{group_display_name} - ⚠️ Mũi MMR/Priorix đầu tiên chỉ cách MVVAC "
                        f"{actual_interval} ngày (khuyến cáo tối thiểu "
                        f"{MVVAC_TO_MMR_MIN_INTERVAL_DAYS} ngày)."
                    ),
                    "earliest_next_dose_date": None,
                    "status_tags": ["warning", "interval_violation_mvvac_mmr"]
                })

            next_due_date = add_years(first_mmr_date, 3)
            description = f"{group_display_name} - Cần tiêm mũi 2 (phác đồ MVVAC + MMR) sau 3 năm kể từ mũi MMR/Priorix đầu tiên."
            earliest_next_date_for_listing = next_due_date
            status_tags = ["info", "booster_upcoming"] if analysis_date < next_due_date else ["due", "booster_due"]
            if analysis_date >= next_due_date:
                earliest_next_date_for_listing = analysis_date
            
            missing_items_list.append({
                "vaccine_name_for_popup": group_display_name,
                "description": description,
                "earliest_next_dose_date": earliest_next_date_for_listing,
                "status_tags": status_tags
            })
            return
        
        elif num_mmr_group_doses >= 2:
            return

    administered_records = get_administered_dose_records(rule_details.get("names_norm_group", []), administered_map)

    if not administered_records:
        age_status_msg, earliest_date, age_tags = get_age_status_and_earliest_date(dob, analysis_date, rule_details, group_display_name)
        default_doses_str = str(rule_details.get("regimens", [{}])[0].get("doses_required", "một số"))
        desc = f"{group_display_name} (Chưa tiêm - cần {default_doses_str} liều theo phác đồ). {age_status_msg}."
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": desc, "earliest_next_dose_date": earliest_date, "status_tags": age_tags
        })
        return

    if not dob:
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": f"{group_display_name} - Không thể xác định phác đồ do thiếu ngày sinh.",
            "earliest_next_dose_date": None, "status_tags": ["error_dob"]})
        return

    first_dose_date = administered_records[0]["date"]
    age_at_first_dose_months, _, _ = get_age_at_date(dob, first_dose_date)

    if age_at_first_dose_months is None:
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": f"{group_display_name} - Lỗi tính tuổi tại mũi đầu tiên của nhóm.",
            "earliest_next_dose_date": None, "status_tags": ["error_age_calculation"]})
        return

    group_age_check_details = {"min_age_months_at_first_dose": rule_details.get("min_age_months_overall_group")}
    if not check_first_dose_age_validity(dob, first_dose_date, group_age_check_details, group_display_name, missing_items_list):
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": f"{group_display_name} - Mũi đầu tiên của nhóm không hợp lệ về tuổi.",
            "earliest_next_dose_date": None, "status_tags": ["error_age_first_dose", "series_restart_needed"]})
        return
        
    applicable_regimen = None
    if not rule_details.get("regimens"):
         missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": f"{group_display_name} - Lỗi cấu hình: không có 'regimens' cho nhóm.",
            "earliest_next_dose_date": None, "status_tags": ["error_config"]})
         return

    for regimen in rule_details["regimens"]:
        min_m = regimen.get("min_age_at_first_dose_months")
        max_m = regimen.get("max_age_at_first_dose_months")
        condition_met = True
        if min_m is not None and age_at_first_dose_months < min_m: condition_met = False
        if max_m is not None and age_at_first_dose_months > max_m: condition_met = False
        if condition_met:
            applicable_regimen = regimen
            break
            
    if not applicable_regimen:
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": f"{group_display_name} - Không tìm thấy phác đồ phù hợp cho tuổi tiêm mũi đầu ({age_at_first_dose_months} tháng).",
            "earliest_next_dose_date": None, "status_tags": ["error_no_matching_rule"]})
        return
    
    regimen_specific_display_name = f"{group_display_name} ({applicable_regimen.get('regimen_name', 'theo phác đồ đã chọn')})"
    regimen_specific_rule_details = {
        "display_name": regimen_specific_display_name,
        "names_norm": rule_details.get("names_norm_group", []),
        "doses_required": applicable_regimen["doses_required"],
        "min_interval_days": applicable_regimen.get("min_interval_days", [None] * applicable_regimen["doses_required"]),
        "min_age_months_at_first_dose": applicable_regimen.get("min_age_at_first_dose_months"),
        "dose_specific_rules": applicable_regimen.get("dose_specific_rules", {})
    }
    check_single_vaccine_series(f"{rule_key}_regimen", regimen_specific_rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules)


def check_cumulative_group_doses(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None):
    """Kiểm tra tổng số liều tích lũy của một nhóm vắc xin."""
    rule_display_name = rule_details.get("group_display_name", rule_details.get("display_name", rule_key))
    total_doses_needed = rule_details.get("required_total_unique_doses", 0)
    all_group_records = get_administered_dose_records(rule_details.get("names_norm", []), administered_map)

    if total_doses_needed > 0 and len(all_group_records) >= total_doses_needed:
        return

    if not all_group_records:
        if total_doses_needed > 0:
            age_status_msg, earliest_date, age_tags = get_age_status_and_earliest_date(dob, analysis_date, rule_details, rule_display_name)
            desc = f"{rule_display_name} (Chưa tiêm - cần {total_doses_needed} liều). {age_status_msg}."
            missing_items_list.append({
                "vaccine_name_for_popup": rule_display_name, 
                "description": desc, "earliest_next_dose_date": earliest_date, "status_tags": age_tags
            })
        return

    if not check_first_dose_age_validity(dob, all_group_records[0]["date"], rule_details, rule_display_name, missing_items_list):
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name, 
            "description": f"{rule_display_name} - Cần {total_doses_needed} liều (do mũi đầu của nhóm không hợp lệ về tuổi).",
            "earliest_next_dose_date": None, "status_tags": ["error_age_first_dose", "series_restart_needed"]})
        return
        
    num_doses_taken = len(all_group_records)

    if num_doses_taken < total_doses_needed:
        remaining = total_doses_needed - num_doses_taken
        desc = f"{rule_display_name} - Cần thêm {remaining} liều (đã tiêm {num_doses_taken}/{total_doses_needed} liều)."
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name, 
            "description": desc,
            "earliest_next_dose_date": analysis_date,
            "status_tags": ["due"]
        })

def check_flu_group(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None):
    """Kiểm tra tiêm phòng Cúm (dựa trên từ khóa và lịch hàng năm)."""
    rule_display_name = rule_details["group_display_name"]
    
    recognition_keywords = rule_details.get("recognition_keywords", [])
    administered_records = []
    
    if recognition_keywords:
        for norm_name, dose_list in administered_map.items():
            if not dose_list: continue
            raw_name = dose_list[0]["raw_name"].lower() 
            for keyword in recognition_keywords:
                if keyword.lower() in raw_name:
                    administered_records.extend(dose_list)
                    break 
    
    administered_records.sort(key=lambda x: x["date"])

    min_age_months_flu = rule_details.get("min_age_months_at_first_dose", 6)

    if not administered_records:
        age_status_msg, earliest_date, age_tags = get_age_status_and_earliest_date(dob, analysis_date, {"min_age_months_at_first_dose": min_age_months_flu}, rule_display_name)
        desc = f"{rule_display_name} (Chưa tiêm. Lần đầu (nếu <9 tuổi) có thể cần 2 mũi cách nhau ~1 tháng, sau đó nhắc lại hàng năm). {age_status_msg}."
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name, 
            "description": desc, "earliest_next_dose_date": earliest_date, "status_tags": age_tags
        })
        return

    if not check_first_dose_age_validity(dob, administered_records[0]["date"], {"min_age_months_at_first_dose": min_age_months_flu}, rule_display_name, missing_items_list):
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name, 
            "description": f"{rule_display_name} - Cần tiêm lại đúng độ tuổi (lần đầu (nếu <9 tuổi) 2 mũi, sau đó hàng năm).",
            "earliest_next_dose_date": None, "status_tags": ["error_age_first_dose", "series_restart_needed"]})
        return

    num_doses_recorded = len(administered_records)
    age_at_first_flu_shot_years = None
    if dob:
        _, _, age_at_first_flu_shot_years = get_age_at_date(dob, administered_records[0]["date"])
    
    is_first_vaccination_under_9 = age_at_first_flu_shot_years is not None and age_at_first_flu_shot_years < 9

    if is_first_vaccination_under_9 and num_doses_recorded == 1:
        required_initial_interval_days = rule_details.get("initial_series_interval_days", 28)
        earliest_next_dose2_date = administered_records[0]["date"] + timedelta(days=required_initial_interval_days)
        if earliest_next_dose2_date < analysis_date: earliest_next_dose2_date = analysis_date
        desc_second_dose = f"{rule_display_name} - Cần mũi 2 (do <9 tuổi lần đầu tiêm) cách mũi 1 khoảng {required_initial_interval_days // 7} tuần. Sau đó nhắc lại hàng năm."
        
        already_missing_second_dose = any(
            item.get("vaccine_name_for_popup") == rule_display_name and "flu_second_dose" in item.get("status_tags", [])
            for item in missing_items_list
        )
        if not already_missing_second_dose:
            missing_items_list.append({
                "vaccine_name_for_popup": rule_display_name,
                "description": desc_second_dose, 
                "earliest_next_dose_date": earliest_next_dose2_date,
                "status_tags": ["due", "flu_second_dose"]
            })
    
    last_dose_date = administered_records[-1]["date"]
    ideal_next_annual_date = add_years(last_dose_date, 1)

    if is_first_vaccination_under_9 and num_doses_recorded == 1:
        return

    if analysis_date >= ideal_next_annual_date:
        earliest_booster_date = analysis_date
        status_tags = ["due", "flu_annual"]
        desc_annual = f"{rule_display_name} - Cần tiêm nhắc lại hàng năm (đã đến lịch)."
    else:
        earliest_booster_date = ideal_next_annual_date
        status_tags = ["info", "booster_upcoming"]
        desc_annual = f"{rule_display_name} - Lịch tiêm nhắc lại hàng năm tiếp theo."

    already_missing_annual_booster = any(
        item.get("vaccine_name_for_popup") == rule_display_name and 
        ("flu_annual" in item.get("status_tags", []) or "booster_upcoming" in item.get("status_tags", []))
        for item in missing_items_list
    )
    if not already_missing_annual_booster:
        missing_items_list.append({
            "vaccine_name_for_popup": rule_display_name, 
            "description": desc_annual, "earliest_next_dose_date": earliest_booster_date,
            "status_tags": status_tags
        })

def check_meningococcal_acyw_group(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None):
    """
    Logic check cho nhóm Não mô cầu ACYW-135 (Menactra và MenQuadfi).
    Tích hợp cảnh báo tương tác với VA-Mengoc BC và 6in1.
    """
    group_display_name = rule_details.get("group_display_name", "Não mô cầu ACYW-135")
    members = rule_details.get("members", {})
    menactra_config = members.get("MENACTRA", {})
    menquadfi_config = members.get("MENQUADFI", {})
    interactions = rule_details.get("interactions", {})

    menactra_records = get_administered_dose_records(menactra_config.get("names_norm", []), administered_map)
    menquadfi_records = get_administered_dose_records(menquadfi_config.get("names_norm", []), administered_map)

    def apply_interactions_and_append(suggestion_item):
        if not suggestion_item: return
        
        base_date = suggestion_item.get("earliest_next_dose_date")
        if base_date is None: 
            missing_items_list.append(suggestion_item)
            return
            
        final_date = base_date
        item_display = suggestion_item.get("vaccine_name_for_popup", group_display_name)
        
        mengoc_int = interactions.get("VA-MENGOC-BC")
        if mengoc_int and dob:
            age_m, _, _ = get_age_at_date(dob, analysis_date)
            min_age_m = mengoc_int.get("applies_when_age_months_gte", 0)
            if age_m is not None and age_m >= min_age_m:
                mengoc_rule = (all_vaccine_rules or {}).get("VA-MENGOC-BC", {})
                m_records = get_administered_dose_records(mengoc_rule.get("names_norm", []), administered_map)
                if m_records:
                    last_m_date = m_records[-1]["date"]
                    if mengoc_int["min_interval_days"] == 60:
                        earliest_after = add_months(last_m_date, 2)
                    else:
                        earliest_after = last_m_date + timedelta(days=mengoc_int["min_interval_days"])
                    if analysis_date < earliest_after:
                        missing_items_list.append({
                            "vaccine_name_for_popup": item_display,
                            "description": mengoc_int["message"],
                            "earliest_next_dose_date": None,
                            "status_tags": ["warning", "interaction_mengoc_bc"]
                        })
                    if final_date < earliest_after:
                        final_date = earliest_after
                        
        six_int = interactions.get("Six_In_One_Combined")
        if six_int:
            six_rule = (all_vaccine_rules or {}).get("Six_In_One_Combined", {})
            s_records = get_administered_dose_records(six_rule.get("names_norm", []), administered_map)
            if s_records:
                last_s_date = s_records[-1]["date"]
                earliest_after = last_s_date + timedelta(days=six_int["min_interval_days"])
                if analysis_date < earliest_after:
                    missing_items_list.append({
                        "vaccine_name_for_popup": item_display,
                        "description": six_int["message"],
                        "earliest_next_dose_date": None,
                        "status_tags": ["warning", "interaction_6in1"]
                    })
                if final_date < earliest_after:
                    final_date = earliest_after
        
        suggestion_item["earliest_next_dose_date"] = max(final_date, analysis_date)
        missing_items_list.append(suggestion_item)

    if menactra_records:
        first_menactra_date = menactra_records[0]["date"]
        age_at_first_menactra_months, _, _ = get_age_at_date(dob, first_menactra_date)
        
        applicable_menactra_rule = None
        if age_at_first_menactra_months is not None:
            for r in menactra_config.get("rules_by_age", []):
                min_m = r.get("min_age_at_first_dose_months")
                max_m = r.get("max_age_at_first_dose_months")
                if (min_m is None or age_at_first_menactra_months >= min_m) and \
                   (max_m is None or age_at_first_menactra_months <= max_m):
                    applicable_menactra_rule = r
                    break
        
        if applicable_menactra_rule:
            doses_required = applicable_menactra_rule["doses_required"]
            if len(menactra_records) >= doses_required:
                return 

            temp_rule = menactra_config.copy()
            temp_rule["display_name"] = menactra_config.get("display", "Menactra")
            
            temp_list = []
            check_age_dependent_series(
                "MENACTRA", temp_rule, administered_map, 
                temp_list, dob, analysis_date, all_vaccine_rules
            )
            for item in temp_list:
                apply_interactions_and_append(item)
            return

    if menquadfi_records:
        first_mq_date = menquadfi_records[0]["date"]
        age_at_first_mq_months, _, _ = get_age_at_date(dob, first_mq_date)
        
        if age_at_first_mq_months is not None:
            applicable_mq_rule = None
            for r in menquadfi_config.get("rules_by_age", []):
                min_m = r.get("min_age_at_first_dose_months")
                max_m = r.get("max_age_at_first_dose_months")
                
                if r.get("min_age_weeks_at_first_dose") == 6 and age_at_first_mq_months < 6:
                    applicable_mq_rule = r
                    break
                if (min_m is None or age_at_first_mq_months >= min_m) and \
                   (max_m is None or age_at_first_mq_months <= max_m):
                    applicable_mq_rule = r
                    break
        else:
            mq_display = menquadfi_config.get("display", "MenQuadfi")
            missing_items_list.append({
                "vaccine_name_for_popup": mq_display,
                "description": f"{mq_display} - Lỗi: Mũi tiêm trước ngày sinh.",
                "earliest_next_dose_date": None,
                "status_tags": ["error_age_calculation"]
            })
            return

        if applicable_mq_rule:
            doses_required = applicable_mq_rule["doses_required"]
            mq_display = menquadfi_config.get("display", "MenQuadfi")
            
            if len(menquadfi_records) < doses_required:
                next_dose_num = len(menquadfi_records) + 1
                min_interval = applicable_mq_rule["min_interval_days"][next_dose_num - 1]
                last_dose_date = menquadfi_records[-1]["date"]
                
                if min_interval == 60:
                    earliest_date = add_months(last_dose_date, 2)
                else:
                    earliest_date = last_dose_date + timedelta(days=min_interval)
                
                apply_interactions_and_append({
                    "vaccine_name_for_popup": mq_display,
                    "description": f"{mq_display} - Cần tiêm mũi {next_dose_num} (phác đồ {doses_required} mũi cơ bản).",
                    "earliest_next_dose_date": earliest_date,
                    "status_tags": ["due"]
                })
                return
            else:
                booster_config = applicable_mq_rule.get("booster")
                if booster_config:
                    min_booster_age_m = booster_config["min_age_months"]
                    min_interval_days = booster_config["min_interval_days_from_last"]
                    
                    last_dose_date = menquadfi_records[-1]["date"]
                    earliest_by_age = add_months(dob, min_booster_age_m)
                    earliest_by_interval = last_dose_date + timedelta(days=min_interval_days)
                    
                    apply_interactions_and_append({
                        "vaccine_name_for_popup": mq_display,
                        "description": f"{mq_display} - {booster_config['description']}.",
                        "earliest_next_dose_date": max(earliest_by_age, earliest_by_interval),
                        "status_tags": ["due", "booster_due"] if analysis_date >= max(earliest_by_age, earliest_by_interval) else ["info", "booster_upcoming"]
                    })
                return

    if not dob:
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name,
            "description": f"{group_display_name} - Không có ngày sinh.",
            "earliest_next_dose_date": None,
            "status_tags": ["error_dob"]
        })
        return

    mq_display = menquadfi_config.get("display", "MenQuadfi")
    min_mq_weeks = menquadfi_config.get("min_age_weeks_overall", 6)
    earliest_mq_date = dob + timedelta(weeks=min_mq_weeks)
    
    if analysis_date >= earliest_mq_date:
        current_age_m, _, _ = get_age_at_date(dob, analysis_date)
        desc_suffix = "(Chưa tiêm - Ưu tiên do phổ tuổi rộng từ 6 tuần)"
        if current_age_m is not None:
            if current_age_m >= 12:
                desc_suffix = "(Gợi ý 1 mũi duy nhất)"
            elif current_age_m >= 6:
                desc_suffix = "(Gợi ý mũi 1 - phác đồ 1 mũi + 1 nhắc)"
            else:
                desc_suffix = "(Gợi ý mũi 1 - phác đồ 3 mũi + 1 nhắc)"

        apply_interactions_and_append({
            "vaccine_name_for_popup": mq_display,
            "description": f"{mq_display} {desc_suffix}.",
            "earliest_next_dose_date": analysis_date,
            "status_tags": ["due"]
        })
    else:
        apply_interactions_and_append({
            "vaccine_name_for_popup": mq_display,
            "description": f"{mq_display} (Chưa đủ tuổi, có thể tiêm từ 6 tuần tuổi).",
            "earliest_next_dose_date": earliest_mq_date,
            "status_tags": ["too_young"]
        })

    current_age_months, _, _ = get_age_at_date(dob, analysis_date)
    min_menactra_m = menactra_config.get("min_age_months_overall", 9)
    if current_age_months is not None and current_age_months >= min_menactra_m:
        menactra_display = menactra_config.get("display", "Menactra")
        apply_interactions_and_append({
            "vaccine_name_for_popup": menactra_display,
            "description": f"{menactra_display} (Hoặc lựa chọn thay thế cho MenQuadfi).",
            "earliest_next_dose_date": analysis_date,
            "status_tags": ["due"]
        })
