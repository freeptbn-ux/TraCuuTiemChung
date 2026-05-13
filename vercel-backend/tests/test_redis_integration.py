import os
import sys
import pytest
from unittest.mock import MagicMock, patch

# Add src to path
sys.path.append(os.path.join(os.path.dirname(__file__), '..'))

from services.redis_cache import RedisCache
from services.auth_service import AuthService
from services.portal_client import PortalClient

@pytest.fixture
def mock_env():
    with patch.dict(os.environ, {
        "UPSTASH_REDIS_REST_URL": "http://fake-redis",
        "UPSTASH_REDIS_REST_TOKEN": "fake-token",
        "PORTAL_USERNAME": "test_user",
        "PORTAL_PASSWORD": "test_password"
    }):
        yield

@pytest.fixture
def mock_redis():
    with patch('services.redis_cache.Redis') as mock:
        yield mock

def test_redis_cache_logic(mock_env, mock_redis):
    cache = RedisCache()
    assert cache.client is not None
    cache.set("test_key", "test_value", ex=100)
    # The internal client should be called
    cache.client.set.assert_called_with("test_key", "test_value", ex=100)
    
    cache.client.get.return_value = "cached_val"
    assert cache.get("test_key") == "cached_val"

def test_auth_service_cache_flow(mock_env, mock_redis):
    with patch('services.auth_service.requests.Session') as mock_session:
        # Case 1: Cache hit
        with patch.object(RedisCache, 'get', return_value='{"cookie1": "val1"}'):
            auth = AuthService()
            cookies = auth.get_session_cookies()
            assert cookies == {"cookie1": "val1"}
            # Ensure no login called
            mock_session.return_value.get.assert_not_called()

        # Case 2: Cache miss -> Login
        with patch.object(RedisCache, 'get', return_value=None):
            with patch.object(RedisCache, 'set') as mock_set:
                auth = AuthService()
                # Mock login flow
                mock_resp_get = MagicMock()
                mock_resp_get.text = '<input name="__RequestVerificationToken" value="fake_token">'
                mock_session.return_value.get.return_value = mock_resp_get
                
                mock_resp_post = MagicMock()
                mock_session.return_value.post.return_value = mock_resp_post
                mock_session.return_value.cookies.get_dict.return_value = {"ASPXAUTH": "fake_auth"}
                mock_session.return_value.cookies.__contains__.side_effect = lambda k: k == ".ASPXAUTH"

                cookies = auth.get_session_cookies()
                assert cookies == {"ASPXAUTH": "fake_auth"}
                mock_set.assert_called()

def test_portal_client_integration(mock_env):
    with patch('services.portal_client.AuthService') as mock_auth:
        mock_auth.return_value.get_session_cookies.return_value = {".ASPXAUTH": "valid"}
        
        client = PortalClient()
        assert client.session.cookies.get(".ASPXAUTH") == "valid"
        
        # Test retry logic on session expired
        with patch.object(client.session, 'get') as mock_get:
            # First call returns login page (expired session)
            mock_expired = MagicMock()
            mock_expired.text = "UserName __RequestVerificationToken" # Simulates login redirect
            
            # Second call returns data
            mock_success = MagicMock()
            mock_success.text = "<html>Valid Data</html>"
            
            mock_get.side_effect = [mock_expired, mock_success]
            
            with patch.object(client, '_parse_search_results', return_value=[]):
                client.lookup_patients("0123456789")
                
                # Should have called get_session_cookies with force_login=True
                mock_auth.return_value.get_session_cookies.assert_any_call(force_login=True)
                assert mock_get.call_count == 2

