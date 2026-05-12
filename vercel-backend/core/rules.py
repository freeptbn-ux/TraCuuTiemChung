import json
import os
from typing import Dict, Any
from .utils import normalize_vaccine_name

class VaccineRulesLoader:
    def __init__(self, rules_path: str = None):
        self.rules_path = rules_path
        self._rules = None

    def _get_default_path(self) -> str:
        current_dir = os.path.dirname(os.path.abspath(__file__))
        rules_path = os.path.join(current_dir, "..", "assets", "vaccine_rules.json")
        if not os.path.exists(rules_path):
            rules_path = os.path.join(os.getcwd(), "assets", "vaccine_rules.json")
        return rules_path

    def load(self) -> Dict[str, Any]:
        path = self.rules_path or self._get_default_path()
        if not os.path.exists(path):
            return {}
            
        with open(path, "r", encoding="utf-8") as f:
            raw_data = json.load(f)
        
        self._rules = self.process_raw_vaccine_rules_data(raw_data)
        return self._rules

    def get_rules(self) -> Dict[str, Any]:
        if self._rules is None:
            self.load()
        return self._rules

    @staticmethod
    def process_raw_vaccine_rules_data(raw_data: Dict[str, Any]) -> Dict[str, Any]:
        """Tiền xử lý rules (chuẩn hóa tên) để tăng tốc độ truy vấn."""
        processed_rules = {}
        for key, details_dict in raw_data.items():
            new_details = details_dict.copy()
            
            # Chuẩn hóa raw_names
            if "raw_names" in new_details:
                new_details["names_norm"] = list(set([normalize_vaccine_name(n) for n in new_details["raw_names"]]))

            # Chuẩn hóa cho MMR_Group
            if "raw_names_members" in new_details:
                all_member_names_norm = []
                for member_key, raw_names_list in new_details["raw_names_members"].items():
                    all_member_names_norm.extend([normalize_vaccine_name(n) for n in raw_names_list])
                new_details["names_norm_group"] = list(set(all_member_names_norm))
            
            # Chuẩn hóa cho các group có courses (Rota, Pneumo, HepA)
            if "courses" in new_details and isinstance(new_details["courses"], list):
                new_courses_list = []
                all_course_names_norm_for_group = []
                for course_dict in new_details["courses"]:
                    new_course = course_dict.copy()
                    if "raw_names" in new_course:
                        normalized_course_names = [normalize_vaccine_name(n) for n in new_course["raw_names"]]
                        new_course["names_norm"] = list(set(normalized_course_names))
                        all_course_names_norm_for_group.extend(normalized_course_names)
                    new_courses_list.append(new_course)
                new_details["courses"] = new_courses_list
                if not new_details.get("names_norm"):
                     new_details["names_norm"] = list(set(all_course_names_norm_for_group))

            # Chuẩn hóa cho các group có "members" (Meningococcal ACYW)
            if "members" in new_details:
                all_member_names_norm = []
                for member_key, member_config in new_details["members"].items():
                    if "raw_names" in member_config:
                        normalized = [normalize_vaccine_name(n) for n in member_config["raw_names"]]
                        member_config["names_norm"] = list(set(normalized))
                        all_member_names_norm.extend(normalized)
                new_details["names_norm_group"] = list(set(all_member_names_norm))

            processed_rules[key] = new_details
        return processed_rules
