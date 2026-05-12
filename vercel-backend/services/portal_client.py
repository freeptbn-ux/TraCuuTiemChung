import os
import requests
from bs4 import BeautifulSoup
import re
from datetime import datetime
from services.auth_service import AuthService

# URLs from constants
LOGIN_URL = "https://tiemchung.vncdc.gov.vn/Account/Login"
INDEX_URL = "https://tiemchung.vncdc.gov.vn/TiemChung/DoiTuong/Index"
SEARCH_URL = "https://tiemchung.vncdc.gov.vn/TiemChung/DoiTuong/TimKiem"
DETAIL_URL = "https://tiemchung.vncdc.gov.vn/TiemChung/DoiTuong/Detail"

class PortalClient:
    def __init__(self):
        self.session = requests.Session()
        self.auth_service = AuthService()
        self.headers = {
            'user-agent': 'Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/144.0.0.0 Safari/537.36 Edg/144.0.0.0',
        }
        self._sync_session()

    def _sync_session(self, force_login=False):
        """Sync requests session with cached cookies."""
        cookies = self.auth_service.get_session_cookies(force_login=force_login)
        self.session.cookies.update(cookies)

    def lookup_patients(self, phone: str):
        """Search for patients by phone number."""
        search_params = {
            'Length': 5, 'LoaiDiaChi': 0, 'VungMienId': '-Khu vực-',
            'ThonApId': '-Thôn/Ấp-', 'NgaySinhTu': '', 'NgaySinhToi': '',
            'GioiTinh': -1, 'LuaTuoi': -1, 'MaDoiTuong': '', 'TenDoiTuong': '',
            'TenMe': '', 'TenBo': '', 'MaDinhDanh': '', 'SoDienThoai': phone,
            'TenNguoiGiamHo': '', 'TinhTrangTheoDoi': -1, 'TinhTrangMangThai': -1,
            'PageNumber': 1, 'PageSize': 20, 'CurrentSystemDate': '',
        }
        
        headers = self.headers.copy()
        headers['X-Requested-With'] = 'XMLHttpRequest'
        headers['Referer'] = INDEX_URL
        
        resp = self.session.get(SEARCH_URL, params=search_params, headers=headers, timeout=30)
        resp.raise_for_status()
        
        if "UserName" in resp.text and "__RequestVerificationToken" in resp.text:
            # Session expired, retry login
            self._sync_session(force_login=True)
            resp = self.session.get(SEARCH_URL, params=search_params, headers=headers, timeout=30)
            resp.raise_for_status()

        return self._parse_search_results(resp.text)

    def get_vaccination_history(self, patient_id: str):
        """Get vaccination history and patient info for a specific patient."""
        detail_params = {'doiTuongId': patient_id}
        headers = self.headers.copy()
        headers['X-Requested-With'] = 'XMLHttpRequest'
        headers['Referer'] = INDEX_URL
        
        resp = self.session.get(DETAIL_URL, params=detail_params, headers=headers, timeout=30)
        resp.raise_for_status()
        
        if "UserName" in resp.text and "__RequestVerificationToken" in resp.text:
            self._sync_session(force_login=True)
            resp = self.session.get(DETAIL_URL, params=detail_params, headers=headers, timeout=30)
            resp.raise_for_status()

        from core.parser import HTMLVaccineParser
        parser = HTMLVaccineParser(resp.text)
        
        patient_info = parser.extract_patient_info()
        history = parser.extract_vaccine_records()
        
        formatted_history = []
        for rec in history:
            formatted_history.append({
                "vaccine_name": rec["raw_name"],
                "dose": str(rec["dose_number"]),
                "date": rec["date"].strftime("%d/%m/%Y")
            })

        return {
            "patient_info": patient_info,
            "history": formatted_history
        }

    def _parse_search_results(self, html_content: str):
        """Parse search results table."""
        soup = BeautifulSoup(html_content, 'html.parser')
        
        # Extract MA_DOI_TUONG via Regex as in old code
        patient_codes = re.findall(r'"MA_DOI_TUONG":"(\d+)"', html_content)
        
        table = soup.find("table", id="doiTuongSearchResult")
        if not table:
            return []

        rows = table.find("tbody").find_all("tr")
        results = []
        for idx, row in enumerate(rows):
            if "không có đối tượng nào" in row.get_text().lower():
                continue
                
            cells = row.find_all("td")
            if len(cells) < 5:
                continue
                
            id_value = row.get('data-id')
            if not id_value:
                onclick = row.get('onclick', '')
                match = re.search(r'OnShowDetail\((\d+)\)', onclick)
                if match:
                    id_value = match.group(1)
            
            if id_value and "," in id_value:
                id_value = id_value.split(",")[0].strip()
                
            name = cells[1].get_text(strip=True)
            birth = cells[4].get_text(strip=True)
            code = patient_codes[idx] if idx < len(patient_codes) else None
            
            results.append({
                "patient_id": id_value,
                "name": name,
                "dob": birth,
                "code": code
            })
        return results

