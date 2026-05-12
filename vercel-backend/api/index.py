from fastapi import FastAPI, Header, HTTPException, Depends
from pydantic import BaseModel
from typing import List, Optional
import os
from services.portal_client import PortalClient
from services.analyzer_service import AnalyzerService
from dotenv import load_dotenv

# Load environment variables
load_dotenv()

app = FastAPI(title="Vaccine Analyzer API")

# Security key
API_KEY = os.getenv("API_KEY", "default_secret_key")

def verify_api_key(x_api_key: str = Header(...)):
    if x_api_key != API_KEY:
        raise HTTPException(status_code=403, detail="Invalid API Key")
    return x_api_key

# Request Models
class LookupRequest(BaseModel):
    phone: str

class AnalyzeRequest(BaseModel):
    patient_id: str
    phone: str # Might be used for re-login if session expires

# Services (Singleton-like or instantiated per request)
# For now, per request is safer for session handling if not using Redis yet
def get_portal_client():
    return PortalClient()

def get_analyzer_service():
    return AnalyzerService()

@app.get("/health")
async def health_check():
    return {"status": "ok", "message": "Backend is running"}

@app.get("/")
async def root():
    return {"message": "Welcome to Vaccine Analyzer API"}

@app.post("/api/lookup")
async def lookup(
    request: LookupRequest, 
    x_api_key: str = Depends(verify_api_key),
    portal: PortalClient = Depends(get_portal_client)
):
    try:
        results = portal.lookup_patients(request.phone)
        return {"status": "success", "data": results}
    except Exception as e:
        raise HTTPException(status_code=500, detail=f"Portal error: {str(e)}")

@app.post("/api/analyze")
async def analyze(
    request: AnalyzeRequest,
    x_api_key: str = Depends(verify_api_key),
    portal: PortalClient = Depends(get_portal_client),
    analyzer: AnalyzerService = Depends(get_analyzer_service)
):
    try:
        # 1. Get history from portal
        # We might need patient basic info (name, dob) as well for analysis
        # Search again to get patient info if not provided in request, 
        # but usually patient_id is enough for detail page.
        # Actually, the Detail page has the name and DOB.
        
        # We need a way to get HTML content or parsed info from Detail page.
        # Let's modify PortalClient to return both info and history.
        
        # For now, let's assume we need to fetch detail page.
        # I'll update PortalClient to have a method that returns the detail page info.
        
        # [Wait] PortalClient.get_vaccination_history currently only returns history list.
        # I need patient name and DOB too.
        
        # I'll modify PortalClient.get_vaccination_history to return both.
        
        history_data = portal.get_vaccination_history(request.patient_id)
        # Assuming history_data is {"patient_info": {...}, "history": [...]}
        
        if not history_data or "patient_info" not in history_data:
             # Fallback: if I can't get info from detail page directly, I might need to lookup by ID?
             # But the detail page MUST have it.
             raise HTTPException(status_code=404, detail="Patient details not found")

        # 2. Run analysis
        analysis_result = analyzer.analyze(history_data["patient_info"], history_data["history"])
        
        return {"status": "success", "data": analysis_result}
        
    except Exception as e:
        import traceback
        print(traceback.format_exc())
        raise HTTPException(status_code=500, detail=f"Analysis error: {str(e)}")
