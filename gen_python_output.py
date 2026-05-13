import json
import sys
import os

# add vercel-backend to path
sys.path.append(os.path.join(os.path.dirname(__file__), 'vercel-backend'))

from core.parser import HTMLVaccineParser
from services.analyzer_service import AnalyzerService

def main():
    with open('test/Minhkhoi.html', 'r', encoding='utf-8') as f:
        content = f.read()

    parser = HTMLVaccineParser(content)
    patient_info = parser.extract_patient_info()
    records = parser.extract_vaccine_records()

    vaccine_history = []
    for r in records:
        vaccine_history.append({
            "vaccine_name": r["raw_name"],
            "date": r["date"].strftime("%d/%m/%Y"),
            "dose": str(r["dose_number"])
        })

    service = AnalyzerService()
    result = service.analyze(patient_info, vaccine_history)

    # ensure testdata exists
    os.makedirs('testdata', exist_ok=True)
    with open('testdata/python_minhkhoi_output.json', 'w', encoding='utf-8') as f:
        json.dump(result, f, ensure_ascii=False, indent=2)

if __name__ == '__main__':
    main()
