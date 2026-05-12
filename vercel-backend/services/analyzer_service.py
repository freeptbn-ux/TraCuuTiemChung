import json
import os
import re
from datetime import datetime
from collections import defaultdict
from core.utils import normalize_vaccine_name
from core.engine.processor import process_all_vaccine_rules

class AnalyzerService:
    def __init__(self):
        self.rules_path = os.path.join(os.path.dirname(os.path.dirname(__file__)), "assets", "vaccine_rules.json")
        self.rules = self._load_rules()

    def _load_rules(self):
        with open(self.rules_path, "r", encoding="utf-8") as f:
            rules = json.load(f)
        
        # Pre-process rules: normalize all raw_names
        for key in rules:
            rule = rules[key]
            if "raw_names" in rule:
                rule["names_norm"] = [normalize_vaccine_name(n) for n in rule["raw_names"]]
            
            # For group rules
            if "raw_names_members" in rule:
                rule["names_norm_group"] = []
                for member in rule["raw_names_members"]:
                    norms = [normalize_vaccine_name(n) for n in rule["raw_names_members"][member]]
                    rule["names_norm_group"].extend(norms)
            
            # For members in complex groups
            if "members" in rule:
                for member_key in rule["members"]:
                    member = rule["members"][member_key]
                    if "raw_names" in member:
                        member["names_norm"] = [normalize_vaccine_name(n) for n in member["raw_names"]]
            
            # For alternative courses
            if "courses" in rule:
                for course in rule["courses"]:
                    if "raw_names" in course:
                        course["names_norm"] = [normalize_vaccine_name(n) for n in course["raw_names"]]
                        
        return rules

    def analyze(self, patient_info, vaccine_history):
        """
        Analyze vaccination history.
        patient_info: {"name": str, "birth": "dd/mm/yyyy", "system_date": "dd/mm/yyyy"}
        vaccine_history: List of {"vaccine_name": str, "date": "dd/mm/yyyy", "dose": str}
        """
        dob_str = patient_info.get("birth")
        system_date_str = patient_info.get("system_date") or datetime.now().strftime("%d/%m/%Y")
        
        try:
            dob = datetime.strptime(dob_str, "%d/%m/%Y").date()
            analysis_date = datetime.strptime(system_date_str, "%d/%m/%Y").date()
        except Exception as e:
            raise ValueError(f"Invalid date format: {e}")

        # Group administered vaccines by normalized name
        administered_map = defaultdict(list)
        for rec in vaccine_history:
            norm_name = normalize_vaccine_name(rec["vaccine_name"])
            try:
                # Handle dose as int if possible
                dose_str = rec["dose"]
                dose_num = int(re.search(r'\d+', dose_str).group()) if re.search(r'\d+', dose_str) else 0
                
                date_obj = datetime.strptime(rec["date"], "%d/%m/%Y").date()
                administered_map[norm_name].append({
                    "date": date_obj,
                    "dose_number": dose_num,
                    "raw_name": rec["vaccine_name"]
                })
            except:
                continue

        # Sort by date
        for name in administered_map:
            administered_map[name].sort(key=lambda x: x["date"])

        # Run engine
        # In this implementation, other_standard_vaccines can be passed if needed
        # For now, we use a basic list or empty
        other_standard_vaccines = [] 
        
        results = process_all_vaccine_rules(
            administered_map, 
            self.rules, 
            dob, 
            analysis_date, 
            other_standard_vaccines
        )

        # Format results to match Android RecommendationDto
        formatted_recommendations = []
        for item in results:
            tags = item.get("status_tags", [])
            
            # Simple status mapping logic
            status = "DUE_LATER"
            if "due" in tags or "eligible" in tags:
                status = "DUE_NOW"
            elif "overdue" in tags:
                status = "OVERDUE"
            elif "error" in "".join(tags) or "warning" in tags:
                status = "NEEDS_REVIEW"
            elif "completed" in tags:
                status = "COMPLETED"
            
            formatted_recommendations.append({
                "vaccine_name": item.get("vaccine_name_for_popup"),
                "rule_type": "standard", # Placeholder
                "status": status,
                "next_dose": item.get("earliest_next_dose_date").strftime("%d/%m/%Y") if item.get("earliest_next_dose_date") and hasattr(item.get("earliest_next_dose_date"), "strftime") else item.get("earliest_next_dose_date"),
                "message": item.get("description"),
                "status_tags": tags
            })

        return {
            "patient_name": patient_info.get("name"),
            "dob": dob_str,
            "analysis_date": system_date_str,
            "missing_vaccines": formatted_recommendations
        }
