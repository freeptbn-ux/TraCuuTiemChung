package com.tracuutiemchung.app.data.portal

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class PortalClientParserIntegrationTest {
    private val parser = DefaultPortalParser()

    @Test
    fun searchResponseAsParseRaw_keepsDetailInputsAndTableAfterMarker() {
        val searchHtml = """
            <html><body>
                <a onclick="OnShowDetail(101)">Chi tiết</a>
                <input id="txtHoTen" value="Tên từ danh sách" />
            </body></html>
        """.trimIndent()
        val detailHtml = detailHtml(
            name = "Nguyen Van Detail",
            dobInput = "<input id=\"txtNgaySinh\" value=\"01/02/2020\" />",
            systemDateInput = "<input id=\"CurrentSystemDate\" value=\"12/05/2026\" />",
        )
        val response = PortalSearchResponse(
            searchHtml = searchHtml,
            detailHtmlBySubjectId = mapOf("101" to detailHtml),
        )

        val raw = response.asParseRaw()
        val result = parser.parse(raw, RawSourceType.HTML).getOrThrow()

        assertTrue(raw.contains(DETAIL_MARKER))
        assertTrue(raw.substringAfter(DETAIL_MARKER).contains("id=\"txtHoTen\""))
        assertTrue(raw.substringAfter(DETAIL_MARKER).contains("id=\"tblVacxin\""))
        assertEquals("Nguyen Van Detail", result.patient?.fullName)
        assertEquals("01/02/2020", result.patient?.dateOfBirth)
        assertEquals("12/05/2026", result.systemDate)
        assertEquals(1, result.vaccinations.size)
        assertEquals("BCG", result.vaccinations.single().vaccineName)
    }

    @Test
    fun parserFallsBackToFullRawWhenDetailMarkerSectionHasNoUsableData() {
        val raw = """
            ${detailHtml(name = "Fallback Patient")}
            $DETAIL_MARKER
            <div>Partial response without patient fields or vaccine table.</div>
        """.trimIndent()

        val result = parser.parse(raw, RawSourceType.HTML).getOrThrow()

        assertEquals("Fallback Patient", result.patient?.fullName)
        assertEquals(1, result.vaccinations.size)
        assertEquals("BCG", result.vaccinations.single().vaccineName)
    }

    @Test
    fun parserHandlesPartialDetailHtmlAfterMarkerWithHiddenDateFields() {
        val raw = """
            <table><tr><td><a onclick="OnShowDetail(202)">Chi tiết</a></td></tr></table>
            $DETAIL_MARKER
            ${detailHtml(
                name = "Partial Detail",
                dobInput = "<input id=\"hfNgaySinhDoiTuong\" value=\"05/06/2021\" />",
                systemDateInput = "<input id=\"hfNgayHienTai\" value=\"13/05/2026\" />",
            )}
        """.trimIndent()

        val result = parser.parse(raw, RawSourceType.HTML).getOrThrow()

        assertEquals("Partial Detail", result.patient?.fullName)
        assertEquals("05/06/2021", result.patient?.dateOfBirth)
        assertEquals("13/05/2026", result.systemDate)
        assertEquals(1, result.vaccinations.size)
    }

    @Test
    fun extractSubjectIdsStillReadsListHtml() {
        val searchHtml = """
            <table>
                <tr><td><a onclick="OnShowDetail(101)">A</a></td></tr>
                <tr><td><a onclick="OnShowDetail(202)">B</a></td></tr>
                <tr><td><a onclick="OnShowDetail(101)">A duplicate</a></td></tr>
            </table>
        """.trimIndent()

        assertEquals(listOf("101", "202"), DefaultPortalParser.extractSubjectIds(searchHtml))
    }

    @Test
    fun captchaOrOtpResponseStillReturnsCaptchaRequired() {
        val result = parser.parse(
            """
                <html><body>
                    <form><input name="otp" /></form>
                </body></html>
            """.trimIndent(),
            RawSourceType.HTML,
        )

        assertTrue(result.exceptionOrNull() is PortalLookupException.CaptchaRequired)
    }

    @Test
    fun notFoundAndParseFailedRemainResultFailuresWithoutThrowingAtCallSite() {
        val notFound = parser.parse("<html><body>Không có dữ liệu phù hợp</body></html>", RawSourceType.HTML)
        val parseFailed = parser.parse("{ invalid json }", RawSourceType.JSON)

        assertTrue(notFound.isFailure)
        assertTrue(notFound.exceptionOrNull() is PortalLookupException.NotFound)
        assertTrue(parseFailed.isFailure)
        assertTrue(parseFailed.exceptionOrNull() is PortalLookupException.ParseFailed)
    }

    private fun detailHtml(
        name: String,
        dobInput: String = "<input id=\"txtNgaySinh\" value=\"01/02/2020\" />",
        systemDateInput: String = "<input id=\"CurrentSystemDate\" value=\"12/05/2026\" />",
    ): String = """
        <div class="detail-fragment">
            <input id="txtHoTen" value="$name" />
            $dobInput
            $systemDateInput
            <table id="tblVacxin">
                <tbody>
                    <tr>
                        <td>1</td><td>BCG</td><td>1</td><td>Lô</td><td>03/04/2020</td><td>VNCDC</td>
                    </tr>
                </tbody>
            </table>
        </div>
    """.trimIndent()
}
