import json
import os

def load_json(path):
    with open(path, 'r', encoding='utf-8') as f:
        return json.load(f)

def main():
    py_data = load_json('testdata/python_giahan_output.json')
    go_data = load_json('testdata/go_giahan_output.json')
    rules = load_json('vercel-backend/assets/vaccine_rules.json')

    # Map Python results by display name or normalized name
    py_map = {v['vaccine_name']: v for v in py_data['missing_vaccines']}
    
    # Map Go results
    go_map = {v['vaccine_name_for_popup']: v for v in go_data['missing_vaccines']}

    # Get all vaccine display names from rules
    all_vaccines = sorted(list(rules.keys()))

    report = []
    report.append("# Full Comparison Report: Gia-Han.html")
    report.append("| Vaccine (Rule Key) | Python Status | Go Status | Match? | Note |")
    report.append("| :--- | :--- | :--- | :--- | :--- |")

    for v_key in all_vaccines:
        display_name = rules[v_key].get('display_name', v_key)
        
        # Find in Python (Python output uses display_name or similar)
        py_res = None
        for k, v in py_map.items():
            if display_name in k or k in display_name:
                py_res = v
                break
        
        # Find in Go (Go output uses vaccine_name_for_popup which is display_name)
        go_res = go_map.get(display_name)

        py_status = "COMPLETED / NOT LISTED"
        if py_res:
            py_status = f"{py_res['status']} ({py_res.get('next_dose', 'N/A')})"
        
        go_status = "COMPLETED / NOT LISTED"
        if go_res:
            tags = ", ".join(go_res.get('status_tags', []))
            next_date = go_res.get('earliest_next_dose_date', 'N/A')
            if next_date != 'N/A':
                # Convert ISO to DD/MM/YYYY
                try:
                    from datetime import datetime
                    dt = datetime.fromisoformat(next_date.replace('Z', '+00:00'))
                    next_date = dt.strftime("%d/%m/%Y")
                except:
                    pass
            go_status = f"{tags.upper()} ({next_date})"

        match = "✅"
        if (py_res is None) != (go_res is None):
            match = "❌"
        elif py_res and go_res:
            # Check if statuses are roughly the same
            if "DUE_NOW" in py_status and "DUE" in go_status:
                match = "✅"
            elif "DUE_LATER" in py_status and "TOO_YOUNG" in go_status:
                match = "✅"
            else:
                # If they are different types of due/too young
                match = "⚠️"

        report.append(f"| {display_name} | {py_status} | {go_status} | {match} | |")

    with open('testdata/full_comparison_giahan.md', 'w', encoding='utf-8') as f:
        f.write("\n".join(report))

    print("Generated testdata/full_comparison_giahan.md")

if __name__ == '__main__':
    main()
