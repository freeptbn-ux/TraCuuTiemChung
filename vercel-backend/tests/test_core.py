import sys
import os
from datetime import date

# Thêm thư mục gốc của backend vào path để import core
sys.path.append(os.path.abspath(os.path.join(os.path.dirname(__file__), '..')))

from core.analyzer_service import AnalyzerService

def test_basic_analysis():
    # Giả lập HTML bảng tiêm chủng
    mock_html = """
    <html>
        <body>
            <span id="txtHoTen">Nguyễn Văn A</span>
            <span id="txtNgaySinh">01/01/2024</span>
            <table id="tblVacxin">
                <tr>
                    <td>1</td>
                    <td>Lao (BCG)</td>
                    <td>1</td>
                    <td>Tiêm chủng mở rộng</td>
                    <td>05/01/2024</td>
                </tr>
                <tr>
                    <td>2</td>
                    <td>Hexaxim</td>
                    <td>1</td>
                    <td>Dịch vụ</td>
                    <td>01/03/2024</td>
                </tr>
            </table>
        </body>
    </html>
    """
    
    print("--- Khởi tạo AnalyzerService ---")
    # Sử dụng đường dẫn tương đối từ thư mục vercel-backend
    service = AnalyzerService(rules_path="assets/vaccine_rules.json")
    
    print("\n--- Chạy phân tích từ HTML giả lập ---")
    results = service.analyze_from_html(mock_html)
    
    print(f"Bệnh nhân: {results['patient_name']}")
    print(f"Ngày sinh: {results['patient_dob']}")
    
    print("\n--- Danh sách đã tiêm ---")
    for adm in results['administered']:
        print(f"- {adm['name']} (Mũi {adm['dose']}): {adm['date']} (Tuổi: {adm['age']})")
        
    print("\n--- Gợi ý tiêm chủng ---")
    for rec in results['recommendations']:
        tags = ", ".join(rec['status_tags'])
        print(f"- [{rec['date']}] {rec['vaccine_name']}: {rec['description']} (Tags: {tags})")

if __name__ == "__main__":
    try:
        test_basic_analysis()
    except Exception as e:
        print(f"Lỗi khi chạy test: {e}")
        import traceback
        traceback.print_exc()
