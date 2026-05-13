import os
import requests
from bs4 import BeautifulSoup
import json
from services.redis_cache import RedisCache

LOGIN_URL = "https://tiemchung.vncdc.gov.vn/Account/Login"
SESSION_KEY = "portal_session_cookies"

class AuthService:
    def __init__(self):
        self.username = os.getenv("PORTAL_USERNAME")
        self.password = os.getenv("PORTAL_PASSWORD")
        self.cache = RedisCache()
        self.headers = {
            'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0',
        }

    def get_session_cookies(self, force_login=False):
        """Get cookies from cache or perform login."""
        if not force_login:
            cached_cookies = self.cache.get(SESSION_KEY)
            if cached_cookies:
                print("DEBUG: Using session from Redis")
                return json.loads(cached_cookies)

        print("DEBUG: Đang login và lưu cache")
        return self.login()

    def login(self):
        """Perform login and cache cookies."""
        if not self.username or not self.password:
            raise Exception("PORTAL_USERNAME or PORTAL_PASSWORD not set")

        session = requests.Session()
        
        # Get login page for CSRF
        resp = session.get(LOGIN_URL, headers=self.headers, timeout=30)
        resp.raise_for_status()
        
        soup = BeautifulSoup(resp.text, 'html.parser')
        token_tag = soup.find("input", {"name": "__RequestVerificationToken"})
        if not token_tag:
            raise Exception("Could not find __RequestVerificationToken")
        
        token = token_tag['value']
        
        # Form data (using files to match PortalClient implementation if needed, 
        # but regular data is usually fine for these fields)
        form_data = {
            '__RequestVerificationToken': (None, token),
            'UserName': (None, self.username),
            'password': (None, self.password),
            'remember_me': (None, 'false'),
        }
        
        resp = session.post(LOGIN_URL, files=form_data, headers=self.headers, timeout=30)
        resp.raise_for_status()
        
        # Verification
        if ".ASPXAUTH" not in session.cookies:
            raise Exception("Login failed: ASPXAUTH cookie not found")

        cookies_dict = session.cookies.get_dict()
        # Store in Redis with 30 mins TTL
        self.cache.set(SESSION_KEY, json.dumps(cookies_dict), ex=1800)
        
        return cookies_dict
