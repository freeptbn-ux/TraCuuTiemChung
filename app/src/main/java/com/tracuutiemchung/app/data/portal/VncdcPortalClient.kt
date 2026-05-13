package com.tracuutiemchung.app.data.portal

import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.withContext
import okhttp3.FormBody
import okhttp3.HttpUrl.Companion.toHttpUrl
import okhttp3.OkHttpClient
import okhttp3.Request
import okhttp3.RequestBody
import java.io.IOException
import java.util.concurrent.TimeUnit

private const val BASE_URL = "https://tiemchung.vncdc.gov.vn"
private const val LOGIN_URL = "$BASE_URL/Account/Login"
private const val INDEX_URL = "$BASE_URL/TiemChung/DoiTuong/Index"
private const val SEARCH_URL = "$BASE_URL/TiemChung/DoiTuong/TimKiem"
private const val DETAIL_URL = "$BASE_URL/TiemChung/DoiTuong/Detail"

class VncdcPortalClient(
    private val httpClient: OkHttpClient = defaultHttpClient(),
) {
    suspend fun login(username: String, password: String): PortalSession = withContext(Dispatchers.IO) {
        val loginPage = execute(Request.Builder().url(LOGIN_URL).get().build())
        val firstCookies = loginPage.headers.values("Set-Cookie")

        val body = FormBody.Builder()
            .add("username", username)
            .add("password", password)
            .add("remember_me", "true")
            .add("captcha", "")
            .build()

        val loginRequest = baseRequest(LOGIN_URL, firstCookies.cookieHeader())
            .post(body)
            .header("Referer", LOGIN_URL)
            .build()
        val response = execute(loginRequest)
        val allCookies = (firstCookies + response.headers.values("Set-Cookie")).parseCookies()
        val cookieHeader = allCookies.cookieHeader()

        if (!allCookies.containsKey(".ASPXAUTH") && !response.request.url.toString().startsWith(INDEX_URL)) {
            throw LoginInputException("Đăng nhập thất bại. Kiểm tra tài khoản, mật khẩu hoặc captcha VNCDC.")
        }

        val token = runCatching { fetchCsrfToken(cookieHeader) }.getOrNull()
        PortalSession(
            cookieHeader = cookieHeader,
            cookies = allCookies,
            csrfToken = token,
            expiresAtMillis = System.currentTimeMillis() + TimeUnit.MINUTES.toMillis(45),
            source = SessionSource.HTTP_CLIENT,
        )
    }

    suspend fun searchByPhone(phone: String, session: PortalSession): PortalSearchResponse = withContext(Dispatchers.IO) {
        val url = SEARCH_URL.toHttpUrl().newBuilder()
            .addQueryParameter("Length", "5")
            .addQueryParameter("LoaiDiaChi", "0")
            .addQueryParameter("VungMienId", "-Khu vực-")
            .addQueryParameter("ThonApId", "-Thôn/Ấp-")
            .addQueryParameter("NgaySinhTu", "")
            .addQueryParameter("NgaySinhToi", "")
            .addQueryParameter("GioiTinh", "-1")
            .addQueryParameter("LuaTuoi", "-1")
            .addQueryParameter("MaDoiTuong", "")
            .addQueryParameter("TenDoiTuong", "")
            .addQueryParameter("TenMe", "")
            .addQueryParameter("TenBo", "")
            .addQueryParameter("MaDinhDanh", "")
            .addQueryParameter("SoDienThoai", phone)
            .addQueryParameter("TenNguoiGiamHo", "")
            .addQueryParameter("TinhTrangTheoDoi", "-1")
            .addQueryParameter("TinhTrangMangThai", "-1")
            .addQueryParameter("PageNumber", "1")
            .addQueryParameter("PageSize", "20")
            .addQueryParameter("CurrentSystemDate", "")
            .addQueryParameter("X-Requested-With", "XMLHttpRequest")
            .build()

        val response = execute(ajaxRequest(url.toString(), session.cookieHeader).get().build())
        ensureAuthenticated(response.bodyText)
        PortalSearchResponse(
            searchHtml = response.bodyText,
            detailHtmlBySubjectId = emptyMap(),
        )
    }

    suspend fun fetchDetail(patient: PortalPatientSummary, session: PortalSession): String = withContext(Dispatchers.IO) {
        fetchDetail(patient.patientId, session)
    }

    private fun fetchDetail(subjectId: String, session: PortalSession): String {
        val url = DETAIL_URL.toHttpUrl().newBuilder()
            .addQueryParameter("doiTuongId", subjectId)
            .build()
        val response = execute(ajaxRequest(url.toString(), session.cookieHeader).get().build())
        ensureAuthenticated(response.bodyText)
        return response.bodyText
    }

    private fun fetchCsrfToken(cookieHeader: String): String? {
        val response = execute(baseRequest(INDEX_URL, cookieHeader).get().build())
        return Regex("name=[\"']__RequestVerificationToken[\"'][^>]*value=[\"']([^\"']+)[\"']", RegexOption.IGNORE_CASE)
            .find(response.bodyText)
            ?.groupValues
            ?.getOrNull(1)
    }

    private fun baseRequest(url: String, cookieHeader: String? = null): Request.Builder = Request.Builder()
        .url(url)
        .header("User-Agent", "Mozilla/5.0 (Android) AppleWebKit/537.36 Chrome/120 Mobile Safari/537.36")
        .header("Accept-Language", "vi,en;q=0.9")
        .apply { if (!cookieHeader.isNullOrBlank()) header("Cookie", cookieHeader) }

    private fun ajaxRequest(url: String, cookieHeader: String): Request.Builder = baseRequest(url, cookieHeader)
        .header("X-Requested-With", "XMLHttpRequest")
        .header("Referer", INDEX_URL)

    private fun execute(request: Request): PortalHttpResponse {
        httpClient.newCall(request).execute().use { response ->
            val bodyText = response.body?.string().orEmpty()
            if (!response.isSuccessful && !response.isRedirect) {
                throw IOException("VNCDC trả về HTTP ${response.code}.")
            }
            return PortalHttpResponse(request = response.request, headers = response.headers, bodyText = bodyText)
        }
    }

    private fun ensureAuthenticated(bodyText: String) {
        if (bodyText.contains("/Account/Login", ignoreCase = true) || bodyText.contains("name=\"username\"", ignoreCase = true)) {
            throw PortalLookupException.SessionExpired
        }
    }

    companion object {
        private fun defaultHttpClient(): OkHttpClient = OkHttpClient.Builder()
            .connectTimeout(20, TimeUnit.SECONDS)
            .readTimeout(30, TimeUnit.SECONDS)
            .followRedirects(false)
            .build()
    }
}

data class PortalSearchResponse(
    val searchHtml: String,
    val detailHtmlBySubjectId: Map<String, String>,
) {
    fun asParseRaw(): String = buildString {
        appendLine(searchHtml)
        detailHtmlBySubjectId.values.forEach { detailHtml ->
            appendLine("\n$DETAIL_MARKER")
            appendLine(detailHtml)
        }
    }
}

private data class PortalHttpResponse(
    val request: Request,
    val headers: okhttp3.Headers,
    val bodyText: String,
)

private fun List<String>.parseCookies(): Map<String, String> = buildMap {
    forEach { raw: String ->
        val firstSegmentEnd = raw.indexOf(';').let { index: Int -> if (index >= 0) index else raw.length }
        val pair = raw.substring(0, firstSegmentEnd)
        val separatorIndex = pair.indexOf('=')
        if (separatorIndex > 0) {
            val name = pair.substring(0, separatorIndex).trim()
            val value = pair.substring(separatorIndex + 1).trim()
            if (name.isNotBlank() && value.isNotBlank()) put(name, value)
        }
    }
}

private fun Map<String, String>.cookieHeader(): String = entries.joinToString("; ") { (name, value) -> "$name=$value" }

private fun List<String>.cookieHeader(): String = parseCookies().cookieHeader()
