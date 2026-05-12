# Phase 03: Android Environment Configuration
Status: ⬜ Pending
Dependencies: Phase 02 (cần URL production)

## Objective
Cấu hình Android app để tự động chuyển đổi giữa môi trường Development và Production.

## Implementation Steps
1. [ ] **Local Properties Setup:** Chuyển `X-API-KEY` vào `local.properties`.
2. [ ] **Gradle BuildConfig:** Cấu hình `build.gradle.kts` để đọc key từ `local.properties` và inject vào code qua `BuildConfig`.
3. [ ] **Build Types Configuration:** 
    - Định nghĩa `BASE_URL` khác nhau cho `debug` và `release`.
    - Setup `versionName` và `versionCode` chuẩn.
4. [ ] **Retrofit Update:** Đảm bảo `RetrofitClient` sử dụng `BuildConfig.BASE_URL` và tự động thêm header `X-API-KEY`.

## Files to Create/Modify
- `app/build.gradle.kts`
- `app/src/main/java/.../RetrofitClient.kt`
- `local.properties`

## Test Criteria
- [ ] Build Debug: Gọi tới localhost/dev API.
- [ ] Build Release: Gọi tới Vercel production API.
