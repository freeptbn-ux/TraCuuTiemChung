from datetime import date, timedelta
from ..utils import normalize_vaccine_name, get_age_at_date
from .checker_utils import get_administered_dose_records
from .series import check_single_vaccine_series, check_age_dependent_series
from .group_alternative import check_alternative_courses_group, check_alternative_courses_age_range_group
from .group_special import check_mmr_equivalent_group, check_cumulative_group_doses, check_flu_group, check_meningococcal_acyw_group

# Rule Type Constants
RULE_TYPE_SINGLE_SERIES = "single_series"
RULE_TYPE_SINGLE_DOSE_MIN_AGE = "single_dose_min_age"
RULE_TYPE_SINGLE_SERIES_MIN_AGE = "single_series_min_age"
RULE_TYPE_AGE_DEPENDENT = "age_dependent_series"
RULE_TYPE_GROUP_CUMULATIVE_UNIQUE = "group_cumulative_unique_doses"
RULE_TYPE_GROUP_CUMULATIVE_UNIQUE_MIN_AGE = "group_cumulative_unique_doses_min_age"
RULE_TYPE_GROUP_ALTERNATIVE = "group_alternative_courses"
RULE_TYPE_GROUP_ALTERNATIVE_MIN_AGE = "group_alternative_courses_min_age"
RULE_TYPE_GROUP_ALTERNATIVE_AGE_RANGE = "group_alternative_courses_age_range"
RULE_TYPE_FLU_GROUP = "flu_group"
RULE_TYPE_MMR_EQUIVALENT_GROUP = "mmr_equivalent_group"
RULE_TYPE_MENINGOCOCCAL_ACYW_GROUP = "meningococcal_acyw_group"

def process_all_vaccine_rules(administered_vaccine_details_map, vaccine_rules, dob, analysis_date, other_standard_vaccines):
    """
    Xử lý tất cả các quy tắc vắc-xin dựa trên hồ sơ tiêm chủng để tạo danh sách các mũi tiêm còn thiếu/cần tiêm.
    """
    current_missing_items = []
    processed_rule_names = set()
    sorted_rule_keys = list(vaccine_rules.keys())

    # --- Start: Special Pneumococcal Vaccine Logic (Phế cầu) ---
    prevenar13_key, vaxneuvance_key, synflorix_key, pneumovax23_key = "Prevenar13", "Vaxneuvance", "Synflorix", "Pneumovax23"
    pneumo_rules = {k: vaccine_rules.get(k, {}) for k in [prevenar13_key, vaxneuvance_key, synflorix_key, pneumovax23_key]}
    
    pneumo_records = {k: get_administered_dose_records(v.get("names_norm", []), administered_vaccine_details_map) for k, v in pneumo_rules.items()}
    num_doses = {k: len(v) for k, v in pneumo_records.items()}

    patient_age_years = None
    if dob and analysis_date:
        _, _, patient_age_years = get_age_at_date(dob, analysis_date)

    pneumo_rules_to_skip = set()
    active_series_keys = [k for k in [prevenar13_key, vaxneuvance_key, synflorix_key] if num_doses[k] > 0]

    if num_doses[pneumovax23_key] > 0:
        pneumo_rules_to_skip.update(pneumo_rules.keys())
    elif len(active_series_keys) > 1:
        mixed_names = [pneumo_rules[k].get('display_name') for k in active_series_keys]
        current_missing_items.append({
            "description": f"Cảnh báo: Đã ghi nhận tiêm xen kẽ các loại phế cầu ({' và '.join(mixed_names)}). Không nên sử dụng xen kẽ.",
            "earliest_next_dose_date": None, "status_tags": ["error_interchange", "pneumo_mixed"],
            "vaccine_name_for_popup": "Phế cầu (nhiều loại)"
        })
        pneumo_rules_to_skip.update(pneumo_rules.keys())
    elif any(num_doses[k] >= 4 for k in active_series_keys):
        pneumo_rules_to_skip.update(pneumo_rules.keys())
    elif active_series_keys:
        primary_key = active_series_keys[0]
        primary_name = pneumo_rules[primary_key].get("display_name", "")
        
        if patient_age_years is not None and patient_age_years >= 2:
            if num_doses[primary_key] < 3:
                current_missing_items.append({
                    "description": f"{pneumo_rules[pneumovax23_key].get('display_name')}: Có thể tiêm 1 mũi để hoàn thành phác đồ phế cầu (do đã trên 2 tuổi và đã tiêm < 3 mũi {primary_name}).",
                    "earliest_next_dose_date": analysis_date, "status_tags": ["info", "alternative_completion"],
                    "vaccine_name_for_popup": pneumo_rules[pneumovax23_key].get('display_name')
                })
                pneumo_rules_to_skip.update(pneumo_rules.keys())
            elif num_doses[primary_key] == 3:
                temp_missing = []
                check_age_dependent_series(primary_key, pneumo_rules[primary_key], administered_vaccine_details_map, temp_missing, dob, analysis_date, vaccine_rules)
                if temp_missing:
                    current_missing_items.append({
                        "description": f"{pneumo_rules[pneumovax23_key].get('display_name')}: Có thể tiêm 1 mũi thay thế cho mũi 4 của {primary_name} (do đã trên 2 tuổi).",
                        "earliest_next_dose_date": analysis_date, "status_tags": ["info", "alternative_booster"],
                        "vaccine_name_for_popup": pneumo_rules[pneumovax23_key].get('display_name')
                    })
                    pneumo_rules_to_skip.update(pneumo_rules.keys())
        
        other_keys = {prevenar13_key, vaxneuvance_key, synflorix_key} - {primary_key}
        pneumo_rules_to_skip.update(other_keys)
    # --- End: Special Pneumococcal Vaccine Logic ---

    for rule_key in sorted_rule_keys:
        if rule_key in pneumo_rules_to_skip:
            processed_rule_names.update(vaccine_rules[rule_key].get("names_norm", []))
            continue

        rule_details = vaccine_rules[rule_key]
        rule_type = rule_details.get("type")
        
        current_rule_names = set(rule_details.get("names_norm", []))
        current_rule_names.update(rule_details.get("names_norm_group", []))
        if "courses" in rule_details:
            for course in rule_details.get("courses", []):
                current_rule_names.update(course.get("names_norm", []))
        processed_rule_names.update(current_rule_names)
        
        checker_args = (rule_key, rule_details, administered_vaccine_details_map,
                        current_missing_items, dob, analysis_date, vaccine_rules)

        checker_map = {
            RULE_TYPE_SINGLE_SERIES: check_single_vaccine_series,
            RULE_TYPE_SINGLE_DOSE_MIN_AGE: check_single_vaccine_series,
            RULE_TYPE_SINGLE_SERIES_MIN_AGE: check_single_vaccine_series,
            RULE_TYPE_AGE_DEPENDENT: check_age_dependent_series,
            RULE_TYPE_MMR_EQUIVALENT_GROUP: check_mmr_equivalent_group,
            RULE_TYPE_GROUP_CUMULATIVE_UNIQUE: check_cumulative_group_doses,
            RULE_TYPE_GROUP_CUMULATIVE_UNIQUE_MIN_AGE: check_cumulative_group_doses,
            RULE_TYPE_GROUP_ALTERNATIVE: check_alternative_courses_group,
            RULE_TYPE_GROUP_ALTERNATIVE_MIN_AGE: check_alternative_courses_group,
            RULE_TYPE_GROUP_ALTERNATIVE_AGE_RANGE: check_alternative_courses_age_range_group,
            RULE_TYPE_FLU_GROUP: check_flu_group,
            RULE_TYPE_MENINGOCOCCAL_ACYW_GROUP: check_meningococcal_acyw_group
        }
        if rule_type in checker_map:
            checker_map[rule_type](*checker_args)
    
    if other_standard_vaccines:
        for vaccine_display_name in other_standard_vaccines:
            norm_name = normalize_vaccine_name(vaccine_display_name)
            if norm_name not in processed_rule_names and not administered_vaccine_details_map.get(norm_name):
                current_missing_items.append({
                    "description": f"{vaccine_display_name} (Chưa tiêm/uống)",
                    "earliest_next_dose_date": analysis_date,
                    "status_tags": ["due", "standard_unadministered"],
                    "vaccine_name_for_popup": vaccine_display_name
                })
                
    return current_missing_items
