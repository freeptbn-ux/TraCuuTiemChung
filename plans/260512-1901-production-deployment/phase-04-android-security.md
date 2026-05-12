# Phase 04: Android Security & Hardening
Status: ⬜ Pending
Dependencies: Phase 03

## Objective
Bảo mật mã nguồn và đường truyền dữ liệu của ứng dụng.

## Implementation Steps
1. [ ] **Enable R8/ProGuard:** Bật `isMinifyEnabled = true` trong build type `release`.
2. [ ] **ProGuard Rules:** Cấu hình `proguard-rules.pro` để không làm hỏng các data class của Retrofit/Gson.
3. [ ] **Network Security Config:**
    - Tạo `res/xml/network_security_config.xml`.
    - Chỉ cho phép HTTPS (Cleartext traffic = false).
4. [ ] **Certificate Pinning (Optional):** Pin certificate của Vercel để chống MITM attack.

## Files to Create/Modify
- `app/build.gradle.kts`
- `app/proguard-rules.pro`
- `app/src/main/res/xml/network_security_config.xml`
- `app/src/main/AndroidManifest.xml`

## Test Criteria
- [ ] App Release không bị crash do làm rối mã nguồn.
- [ ] App từ chối kết nối nếu dùng proxy/intercept mà không có cert chuẩn.
