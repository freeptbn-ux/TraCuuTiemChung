import re
from datetime import date, datetime, timedelta
from ..utils import get_age_at_date
from .checker_utils import (
    get_administered_dose_records, get_age_status_and_earliest_date,
    check_first_dose_age_validity, add_months, add_years
)
from .series import check_single_vaccine_series

def check_alternative_courses_group(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None):
    """Kiểm tra nhóm vắc xin có nhiều phác đồ thay thế (ví dụ: Rota)."""
    group_display_name = rule_details["group_display_name"]
    any_course_completed = False
    potential_incomplete_messages = []
    max_age_months_to_start_group = rule_details.get("max_age_months_to_start_first_dose_group")
    all_possible_group_records = get_administered_dose_records(rule_details.get("names_norm",[]), administered_map)
    options_str_parts = []

    current_age_months = None
    if dob and analysis_date:
        current_age_months, _, _ = get_age_at_date(dob, analysis_date)

    if not rule_details.get("courses"):
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": f"{group_display_name} - Lỗi cấu hình: không có 'courses'.",
            "earliest_next_dose_date": None, "status_tags": ["error_config"]})
        return

    for c in rule_details["courses"]:
        course_desc = f"{c['display']} (cần {c['doses_required']} liều"
        min_age_c_months = c.get("min_age_months_at_first_dose")
        min_age_c_weeks = c.get("min_age_weeks_at_first_dose")
        final_min_age_desc = ""
        if min_age_c_months is not None: final_min_age_desc = f", từ {min_age_c_months} tháng"
        elif min_age_c_weeks is not None: final_min_age_desc = f", từ {min_age_c_weeks} tuần"
        elif rule_details.get("min_age_months_at_first_dose") is not None: final_min_age_desc = f", từ {rule_details['min_age_months_at_first_dose']} tháng"
        elif rule_details.get("min_age_weeks_at_first_dose") is not None: final_min_age_desc = f", từ {rule_details['min_age_weeks_at_first_dose']} tuần"
        elif rule_details.get("min_age_months_overall_group") is not None: final_min_age_desc = f", từ {rule_details['min_age_months_overall_group']} tháng"
        elif rule_details.get("min_age_weeks_overall_group") is not None: final_min_age_desc = f", từ {rule_details['min_age_weeks_overall_group']} tuần"
        course_desc += f"{final_min_age_desc})"
        options_str_parts.append(course_desc)
    options_str = "Hoặc ".join(options_str_parts) if options_str_parts else "Không có phác đồ nào được định nghĩa."

    if not all_possible_group_records:
        if dob and max_age_months_to_start_group is not None:
            if current_age_months is not None and current_age_months >= max_age_months_to_start_group:
                missing_items_list.append({
                    "vaccine_name_for_popup": group_display_name, 
                    "description": f"{group_display_name}: Đã qua {max_age_months_to_start_group} tháng tuổi, không còn chỉ định bắt đầu.",
                    "earliest_next_dose_date": None, "status_tags": ["too_old_to_start"]})
                return
        
        age_status_msg, earliest_date, age_tags = get_age_status_and_earliest_date(dob, analysis_date, rule_details, group_display_name)
        desc = f"{group_display_name}: {options_str} (Chưa tiêm). {age_status_msg}."
        missing_items_list.append({
            "vaccine_name_for_popup": group_display_name, 
            "description": desc, "earliest_next_dose_date": earliest_date, "status_tags": age_tags
        })
        return

    any_dose_administered_in_group = bool(all_possible_group_records)
    if any_dose_administered_in_group:
        first_dose_date_overall_group = all_possible_group_records[0]["date"]
        group_first_dose_age_details = {
            "min_age_months_at_first_dose": rule_details.get("min_age_months_overall_group"),
            "min_age_weeks_at_first_dose": rule_details.get("min_age_weeks_overall_group")
        }
        group_first_dose_age_details = {k:v for k,v in group_first_dose_age_details.items() if v is not None}
        if not group_first_dose_age_details:
            group_first_dose_age_details = {"min_age_months_at_first_dose": rule_details.get("min_age_months_at_first_dose"),
                                            "min_age_weeks_at_first_dose": rule_details.get("min_age_weeks_at_first_dose")}
            group_first_dose_age_details = {k:v for k,v in group_first_dose_age_details.items() if v is not None}


        if not check_first_dose_age_validity(dob, first_dose_date_overall_group, group_first_dose_age_details, group_display_name, missing_items_list):
             missing_items_list.append({
                "vaccine_name_for_popup": group_display_name, 
                "description": f"{group_display_name}: Mũi đầu tiên của nhóm ({first_dose_date_overall_group.strftime('%d/%m/%Y')}) không hợp lệ về tuổi.",
                "earliest_next_dose_date": None, "status_tags": ["error_age_first_dose", "series_restart_needed"]
             })
             return

    for course_config in rule_details["courses"]:
        course_specific_display_name = f"{group_display_name} - {course_config['display']}"
        temp_course_rule = {
            "display_name": course_specific_display_name, 
            "names_norm": course_config.get("names_norm", []),
            "doses_required": course_config["doses_required"],
            "min_interval_days": course_config.get("min_interval_days", [None] * course_config["doses_required"]),
            "min_age_months_at_first_dose": course_config.get("min_age_months_at_first_dose", rule_details.get("min_age_months_at_first_dose", rule_details.get("min_age_months_overall_group"))),
            "min_age_weeks_at_first_dose": course_config.get("min_age_weeks_at_first_dose", rule_details.get("min_age_weeks_at_first_dose", rule_details.get("min_age_weeks_overall_group"))),
            "dose_specific_rules": course_config.get("dose_specific_rules", {}),
            "booster_interval_years": course_config.get("booster_interval_years"),
            "booster_after_dose_number": course_config.get("booster_after_dose_number"),
            "booster_max_age_years": course_config.get("booster_max_age_years")
        }
        current_course_messages = []
        course_specific_administered_records = get_administered_dose_records(temp_course_rule["names_norm"], administered_map)
        
        if course_specific_administered_records or not any_dose_administered_in_group: 
            check_single_vaccine_series(f"{rule_key}_{course_config['display']}", temp_course_rule, administered_map, current_course_messages, dob, analysis_date, all_vaccine_rules)
        
        is_rota_rule_key = (rule_key == "Rota")
        max_completion_age_months_rota = rule_details.get("max_age_months_for_completion_group")

        if is_rota_rule_key and current_age_months is not None and \
           max_completion_age_months_rota is not None and \
           current_age_months > max_completion_age_months_rota:
            
            if course_specific_administered_records and current_course_messages:
                transformed_messages = []
                for msg_item in current_course_messages:
                    if "due" in msg_item.get("status_tags", []):
                        transformed_messages.append({
                            "vaccine_name_for_popup": msg_item.get("vaccine_name_for_popup", course_specific_display_name),
                            "description": f"{msg_item.get('vaccine_name_for_popup', course_specific_display_name)}: Đã trên {max_completion_age_months_rota} tháng tuổi, không còn chỉ định hoàn thành phác đồ.",
                            "earliest_next_dose_date": None,
                            "status_tags": ["info", "too_old_to_complete", "rota_too_old_to_complete"] 
                        })
                    else:
                        transformed_messages.append(msg_item)
                current_course_messages = transformed_messages 

        if not current_course_messages and course_specific_administered_records: 
            any_course_completed = True 
            break 
        elif course_specific_administered_records and current_course_messages: 
            for msg_item in current_course_messages: 
                if "vaccine_name_for_popup" not in msg_item: 
                     msg_item["vaccine_name_for_popup"] = course_specific_display_name 
            potential_incomplete_messages.extend(current_course_messages) 
    
    if not any_course_completed: 
        if potential_incomplete_messages: 
            missing_items_list.extend(potential_incomplete_messages) 
        elif any_dose_administered_in_group and not potential_incomplete_messages : 
             missing_items_list.append({ 
                "vaccine_name_for_popup": group_display_name, 
                "description": f"{group_display_name}: {options_str} (Đã có mũi tiêm trong nhóm nhưng không khớp hoàn toàn phác đồ nào hoặc mũi đầu không hợp lệ. Kiểm tra lại).", 
                "earliest_next_dose_date": None, "status_tags": ["error_ambiguous_course"]}) 

def check_alternative_courses_age_range_group(rule_key, rule_details, administered_map, missing_items_list, dob, analysis_date, all_vaccine_rules=None): 
    """Kiểm tra nhóm vắc xin có phác đồ thay thế phụ thuộc vào dải tuổi (ví dụ: JE, HepA)."""
    group_display_name = rule_details["group_display_name"] 
    courses_to_skip = set()

    if rule_key == "JE_Group":
        jevax_course = next((c for c in rule_details.get("courses", []) if "Jevax" in c.get("display", "")), None)
        imojev_course = next((c for c in rule_details.get("courses", []) if "Imojev" in c.get("display", "")), None)
        jeev_course = next((c for c in rule_details.get("courses", []) if "JEEV" in c.get("display", "")), None)

        if jevax_course:
            jevax_records = get_administered_dose_records(jevax_course.get("names_norm", []), administered_map)
            imojev_records = get_administered_dose_records(imojev_course.get("names_norm", []), administered_map) if imojev_course else []
            jeev_records = get_administered_dose_records(jeev_course.get("names_norm", []), administered_map) if jeev_course else []

            if len(jevax_records) >= 3 and len(imojev_records) >= 1:
                return 

            if jevax_records and jeev_records:
                last_jevax_date = jevax_records[-1]["date"]
                first_jeev_date = jeev_records[0]["date"]

                if first_jeev_date > last_jevax_date:
                    missing_items_list.append({
                        "vaccine_name_for_popup": group_display_name,
                        "description": f"{group_display_name} - Cảnh báo: Đã tiêm VNNB (Jevax) sau đó chuyển sang JEEV. Sẽ ưu tiên nhắc lịch theo phác đồ JEEV.",
                        "earliest_next_dose_date": None, "status_tags": ["info", "je_mixed_warning", "switched_to_jeev"]
                    })
                    courses_to_skip.add(jevax_course['display'])

            if imojev_records and jevax_records:
                if imojev_records[0]["date"] < jevax_records[0]["date"]:
                    missing_items_list.append({
                        "vaccine_name_for_popup": group_display_name,
                        "description": f"{group_display_name} - Cảnh báo: Đã bắt đầu với Imojev, không nên tiêm Jevax/VNNB. Cần hoàn thành phác đồ Imojev.",
                        "earliest_next_dose_date": None, "status_tags": ["error_interchange", "je_mixed_warning"]
                    })
                    courses_to_skip.add(jevax_course['display'])
                elif imojev_records[0]["date"] > jevax_records[-1]["date"]:
                    missing_items_list.append({
                        "vaccine_name_for_popup": group_display_name,
                        "description": f"{group_display_name} - Cảnh báo: Đã tiêm VNNB (Jevax) sau đó chuyển sang Imojev. Sẽ ưu tiên nhắc lịch theo phác đồ Imojev.",
                        "earliest_next_dose_date": None, "status_tags": ["info", "je_mixed_warning", "switched_to_imojev"]
                    })
                    courses_to_skip.add(jevax_course['display'])
            
            if jevax_course['display'] not in courses_to_skip and len(jevax_records) >= 3 and not imojev_records and not jeev_records:
                booster_interval = jevax_course.get("booster_interval_years")
                if booster_interval:
                    base_date = jevax_records[-1]["date"]
                    next_due = add_years(base_date, booster_interval)
                    
                    is_too_old = False
                    max_age = jevax_course.get("booster_max_age_years")
                    if dob and max_age is not None:
                        _, _, age_at_due = get_age_at_date(dob, next_due)
                        if age_at_due is not None and age_at_due >= max_age:
                            is_too_old = True

                    if not is_too_old:
                        desc = (f"{group_display_name} - Cần tiêm mũi nhắc lại. Tùy chọn: dùng Jevax (3 năm/lần đến 15 tuổi) "
                                f"hoặc tiêm 1 mũi Imojev để hoàn thành.")
                        final_date = next_due if analysis_date < next_due else analysis_date
                        status_tags = ["info", "booster_upcoming"] if analysis_date < next_due else ["due", "booster_due"]
                        
                        missing_items_list.append({
                            "vaccine_name_for_popup": group_display_name,
                            "description": desc,
                            "earliest_next_dose_date": final_date,
                            "status_tags": status_tags
                        })
                
                courses_to_skip.add(jevax_course['display'])
                if imojev_course:
                    courses_to_skip.add(imojev_course['display'])
            
            elif jevax_course['display'] not in courses_to_skip and 0 < len(jevax_records) < 3 and not imojev_records and not jeev_records and (not imojev_course or len(imojev_records) < imojev_course.get("doses_required", 2)):
                if len(jevax_records) == 1:
                    next_date = jevax_records[0]["date"] + timedelta(days=30)
                    if next_date < analysis_date: next_date = analysis_date
                    missing_items_list.append({
                        "vaccine_name_for_popup": "Tùy chọn VNNB -> Imojev",
                        "description": f"{group_display_name} - Tùy chọn: Sau 1 mũi Jevax, có thể chuyển sang tiêm 2 mũi Imojev (mũi Imojev đầu tiên cách mũi Jevax 1 tháng).",
                        "earliest_next_dose_date": next_date, "status_tags": ["info", "alternative_course"]
                    })
                elif len(jevax_records) == 2:
                    next_date = add_years(jevax_records[1]["date"], 1)
                    if next_date < analysis_date: next_date = analysis_date
                    missing_items_list.append({
                        "vaccine_name_for_popup": "Tùy chọn VNNB -> Imojev",
                        "description": f"{group_display_name} - Tùy chọn: Sau 2 mũi Jevax, có thể tiêm 1 mũi Imojev sau 1 năm để hoàn thành.",
                        "earliest_next_dose_date": next_date, "status_tags": ["info", "alternative_course"]
                    })

    any_course_completed_in_group = False 
    incomplete_started_valid_courses_msgs = [] 
    possible_options_texts = [] 
    current_age_months, _, current_age_years = get_age_at_date(dob, analysis_date) if dob and analysis_date else (None, None, None) 
    has_any_administered_dose_in_group = any( 
        get_administered_dose_records(course.get("names_norm",[]), administered_map) 
        for course in rule_details.get("courses", []) 
    )
    if not rule_details.get("courses"): 
        missing_items_list.append({ 
            "vaccine_name_for_popup": group_display_name, 
            "description": f"{group_display_name} - Lỗi cấu hình: không có 'courses'.", 
            "earliest_next_dose_date": None, "status_tags": ["error_config"]}) 
        return 

    at_least_one_course_eligible_to_start_now = False 

    for course_config in rule_details["courses"]: 
        if course_config['display'] in courses_to_skip:
            continue
            
        course_display_name_full = f"{group_display_name} - {course_config['display']}" 
        course_administered_records = get_administered_dose_records(course_config.get("names_norm",[]), administered_map) 
        temp_course_rule_details = { 
            "display_name": course_display_name_full, 
            "names_norm": course_config.get("names_norm", []), 
            "doses_required": course_config["doses_required"], 
            "min_interval_days": course_config.get("min_interval_days", [None] * course_config["doses_required"]), 
            "min_age_months_at_first_dose": course_config.get("min_age_months_at_first_dose"), 
            "max_age_years_at_first_dose_course": course_config.get("max_age_years_at_first_dose"), 
            "dose_specific_rules": course_config.get("dose_specific_rules", {}),
            "booster_interval_years": course_config.get("booster_interval_years"),
            "booster_after_dose_number": course_config.get("booster_after_dose_number"),
            "booster_max_age_years": course_config.get("booster_max_age_years")
        }
        course_min_age_months = course_config.get("min_age_months_at_first_dose") 
        course_max_age_years_for_starting = course_config.get("max_age_years_at_first_dose")  
        is_eligible_now_to_start_this_course = False 
        earliest_start_for_this_course_option = None 

        if dob and current_age_months is not None and current_age_years is not None: 
            age_status_msg_course, earliest_date_course, age_tags_course = get_age_status_and_earliest_date( 
                dob, analysis_date, {"min_age_months_at_first_dose": course_min_age_months} 
            )
            eligible_by_min_age = "eligible" in age_tags_course 
            if "too_young" in age_tags_course and earliest_date_course: 
                earliest_start_for_this_course_option = earliest_date_course 
            elif eligible_by_min_age: 
                 earliest_start_for_this_course_option = analysis_date 
            eligible_by_max_age = True 
            if course_max_age_years_for_starting is not None and current_age_years >= course_max_age_years_for_starting: 
                eligible_by_max_age = False 
            if eligible_by_min_age and eligible_by_max_age: 
                is_eligible_now_to_start_this_course = True 
                at_least_one_course_eligible_to_start_now = True 
        
        option_desc_for_summary = "" 
        if is_eligible_now_to_start_this_course: 
            age_range_str = f"(từ {course_min_age_months or '?'} tháng đến <{course_max_age_years_for_starting or '?'} tuổi)" 
            date_part = f"Sớm nhất: {earliest_start_for_this_course_option.strftime('%d/%m/%Y')}" if earliest_start_for_this_course_option else "" 
            option_desc_for_summary = f"{course_config['display']} ({course_config['doses_required']} liều {age_range_str}. {date_part})" 

        if not course_administered_records: 
            if option_desc_for_summary: 
                possible_options_texts.append(option_desc_for_summary) 
            continue 

        first_dose_this_course_date = course_administered_records[0]["date"] 
        age_at_first_dose_months_course, _, age_at_first_dose_years_course = get_age_at_date(dob, first_dose_this_course_date) if dob else (None,None,None) 
        first_dose_age_valid_for_this_course = True 
        if dob and age_at_first_dose_months_course is not None: 
            temp_min_age_check_details = {"min_age_months_at_first_dose": course_min_age_months} 
            if not check_first_dose_age_validity(dob, first_dose_this_course_date, temp_min_age_check_details, course_display_name_full, incomplete_started_valid_courses_msgs): 
                first_dose_age_valid_for_this_course = False 
            if course_max_age_years_for_starting is not None and age_at_first_dose_years_course is not None and \
               age_at_first_dose_years_course >= course_max_age_years_for_starting: 
                first_dose_age_valid_for_this_course = False 
                incomplete_started_valid_courses_msgs.append({ 
                    "vaccine_name_for_popup": course_display_name_full, 
                    "description":f"{course_display_name_full} - Mũi 1 tiêm sau tuổi tối đa cho phép bắt đầu phác đồ này ({course_max_age_years_for_starting} tuổi).", 
                    "earliest_next_dose_date": None, "status_tags": ["error_age_first_dose", "too_old_at_first_dose"]}) 
        elif not dob: 
            first_dose_age_valid_for_this_course = False 
            incomplete_started_valid_courses_msgs.append({ 
                "vaccine_name_for_popup": course_display_name_full, 
                "description":f"{course_display_name_full} - Không thể xác thực tuổi tiêm mũi 1 do thiếu ngày sinh.", 
                "earliest_next_dose_date": None, "status_tags": ["error_dob"]}) 

        if not first_dose_age_valid_for_this_course: 
            if option_desc_for_summary: 
                 possible_options_texts.append(f"{option_desc_for_summary} - mũi đã tiêm không hợp lệ, cân nhắc tiêm lại.") 
            continue 

        course_specific_missing_items = [] 
        check_single_vaccine_series(f"{rule_key}_{course_config['display']}", temp_course_rule_details, administered_map, course_specific_missing_items, dob, analysis_date, all_vaccine_rules) 
        
        if not course_specific_missing_items: 
            any_course_completed_in_group = True 
            break 
        else: 
            for msg_item in course_specific_missing_items: 
                if "vaccine_name_for_popup" not in msg_item: 
                     msg_item["vaccine_name_for_popup"] = course_display_name_full 
            incomplete_started_valid_courses_msgs.extend(course_specific_missing_items) 

    if not any_course_completed_in_group: 
        if incomplete_started_valid_courses_msgs: 
            missing_items_list.extend(incomplete_started_valid_courses_msgs) 
        elif has_any_administered_dose_in_group : 
             missing_items_list.append({ 
                 "vaccine_name_for_popup": group_display_name, 
                 "description":f"{group_display_name} (Các mũi đã tiêm không phù hợp với phác đồ nào trong khoảng tuổi cho phép hoặc đã qua tuổi. Kiểm tra lại.)", 
                 "earliest_next_dose_date": None, "status_tags": ["error_ambiguous_course"]}) 
        elif possible_options_texts: 
            final_desc = ""
            earliest_start_date = None
            status_tags = ["due"]
            
            if rule_key in ["HepA", "JE_Group"]:
                final_desc = f"{group_display_name}: Chưa tiêm. "
                options = []
                min_age_for_any_course = float('inf')
                
                if dob and current_age_months is not None:
                    for course in rule_details.get("courses", []):
                        min_age = course.get("min_age_months_at_first_dose")
                        if min_age is not None and min_age < min_age_for_any_course:
                            min_age_for_any_course = min_age
                        
                        if min_age is not None and current_age_months >= min_age:
                            options.append(course.get("display", "Không rõ"))
                
                    if options:
                        final_desc += f"Lựa chọn có thể tiêm: {', '.join(sorted(list(set(options))))}."
                        earliest_start_date = analysis_date
                    else: 
                        final_desc += f"Chưa đến tuổi tiêm (bắt đầu từ {min_age_for_any_course} tháng)."
                        earliest_start_date = add_months(dob, min_age_for_any_course)
                        status_tags = ["too_young"]
                else: 
                    final_desc = f"{group_display_name}: Chưa tiêm. Cần có ngày sinh để gợi ý phác đồ."
                    status_tags = ["error_dob"]

                missing_items_list.append({
                    "vaccine_name_for_popup": group_display_name,
                    "description": final_desc,
                    "earliest_next_dose_date": earliest_start_date,
                    "status_tags": status_tags
                })
            else: 
                desc = f"{group_display_name}: Chưa hoàn thành. Lựa chọn có thể tiêm: {' Hoặc '.join(sorted(list(set(possible_options_texts))))}"
                earliest_overall = analysis_date 
                actual_earliest_date_from_options = None
                for opt_text in possible_options_texts:
                    match_date = re.search(r"Sớm nhất: (\d{2}/\d{2}/\d{4})", opt_text)
                    if match_date:
                        try:
                            dt_obj = datetime.strptime(match_date.group(1), "%d/%m/%Y").date()
                            if actual_earliest_date_from_options is None or dt_obj < actual_earliest_date_from_options:
                                actual_earliest_date_from_options = dt_obj
                        except ValueError:
                            pass
                if actual_earliest_date_from_options and actual_earliest_date_from_options > analysis_date:
                     earliest_overall = actual_earliest_date_from_options
                elif actual_earliest_date_from_options and actual_earliest_date_from_options <= analysis_date:
                     earliest_overall = analysis_date

                missing_items_list.append({ 
                    "vaccine_name_for_popup": group_display_name, 
                    "description": desc, "earliest_next_dose_date": earliest_overall, 
                    "status_tags": ["due", "multiple_options"] 
                })
        elif not at_least_one_course_eligible_to_start_now : 
            missing_items_list.append({ 
                "vaccine_name_for_popup": group_display_name, 
                "description":f"{group_display_name} (Chưa tiêm và hiện không có lựa chọn phác đồ nào phù hợp tuổi, hoặc đã qua tuổi cho tất cả phác đồ).", 
                "earliest_next_dose_date": None, "status_tags": ["info_too_old_or_no_option"]})
