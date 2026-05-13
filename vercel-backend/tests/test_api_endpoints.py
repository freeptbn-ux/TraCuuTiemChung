import pytest
from fastapi.testclient import TestClient
from api.index import app, X_API_KEY, get_portal_client
from unittest.mock import MagicMock, patch

client = TestClient(app)

@pytest.fixture
def mock_portal_client():
    mock = MagicMock()
    app.dependency_overrides[get_portal_client] = lambda: mock
    yield mock
    app.dependency_overrides.clear()

def test_health_check():
    response = client.get("/api/health")
    assert response.status_code == 200
    assert response.json()["status"] == "ok"
    assert "environment" in response.json()

def test_root():
    response = client.get("/")
    assert response.status_code == 200
    assert "Welcome" in response.json()["message"]

def test_lookup_no_key(mock_portal_client):
    response = client.post("/api/lookup", json={"phone": "0987654321"})
    assert response.status_code == 422 # FastAPI returns 422 for missing required Header

def test_lookup_invalid_key(mock_portal_client):
    response = client.post(
        "/api/lookup", 
        json={"phone": "0987654321"},
        headers={"X-API-KEY": "wrong_key"}
    )
    assert response.status_code == 403

def test_lookup_success(mock_portal_client):
    mock_portal_client.lookup_patients.return_value = [
        {"patient_id": "123", "name": "Nguyen Van A", "dob": "01/01/2020", "code": "PAT123"}
    ]
    
    response = client.post(
        "/api/lookup",
        json={"phone": "0987654321"},
        headers={"X-API-KEY": X_API_KEY}
    )
    
    assert response.status_code == 200
    assert response.json()["status"] == "success"
    assert len(response.json()["data"]) == 1
    assert response.json()["data"][0]["name"] == "Nguyen Van A"

def test_analyze_success(mock_portal_client):
    # Mock portal.get_vaccination_history
    mock_portal_client.get_vaccination_history.return_value = {
        "patient_info": {"name": "Nguyen Van A", "birth": "01/01/2020", "system_date": "12/05/2026"},
        "history": [
            {"vaccine_name": "BCG", "dose": "1", "date": "02/01/2020"}
        ]
    }
    
    # AnalyzerService will be instantiated and called. 
    # We can mock it too if we want to isolate, but here we test the integration of endpoints.
    
    response = client.post(
        "/api/analyze",
        json={"patient_id": "123", "phone": "0987654321"},
        headers={"X-API-KEY": X_API_KEY}
    )
    
    assert response.status_code == 200
    assert response.json()["status"] == "success"
    assert "missing_vaccines" in response.json()["data"]
    # BCG was given, so it shouldn't be in missing_vaccines (if rule says 1 dose)
    # This verifies the whole flow from endpoint -> portal mock -> analyzer
