import pytest
from fastapi.testclient import TestClient
from api.index import app, API_KEY, get_portal_client, verify_api_key
from unittest.mock import MagicMock, patch

client = TestClient(app, raise_server_exceptions=False)

def test_global_exception_handler():
    # Mock verify_api_key to raise an unexpected exception (unhandled by route try-except)
    # Also mock get_portal_client to avoid real initialization errors
    def crash():
        raise RuntimeError("Auth system crash")
        
    app.dependency_overrides[get_portal_client] = lambda: MagicMock()
    app.dependency_overrides[verify_api_key] = crash
    
    try:
        response = client.post(
            "/api/lookup",
            json={"phone": "0987654321"},
            headers={"X-API-KEY": API_KEY}
        )
        
        assert response.status_code == 500
        data = response.json()
        assert data["status"] == "error"
        assert data["message"] == "Internal server error"
        assert "Auth system crash" in data["detail"]
    finally:
        app.dependency_overrides.clear()

def test_portal_error_handling():
    mock_portal = MagicMock()
    mock_portal.lookup_patients.side_effect = Exception("Portal down")
    
    app.dependency_overrides[get_portal_client] = lambda: mock_portal
    
    try:
        response = client.post(
            "/api/lookup",
            json={"phone": "0987654321"},
            headers={"X-API-KEY": API_KEY}
        )
        
        assert response.status_code == 500
        # In api/index.py, the lookup route has:
        # except Exception as e:
        #     raise HTTPException(status_code=500, detail=f"Portal error: {str(e)}")
        assert "Portal error: Portal down" in response.json()["detail"]
    finally:
        app.dependency_overrides.clear()
