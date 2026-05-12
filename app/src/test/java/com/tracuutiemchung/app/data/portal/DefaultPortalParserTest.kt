package com.tracuutiemchung.app.data.portal

import org.junit.Assert.assertEquals
import org.junit.Assert.assertFalse
import org.junit.Assert.assertTrue
import org.junit.Test

class DefaultPortalParserTest {
    private val parser = DefaultPortalParser()

    @Test
    fun parseJsonReturnsPatientAndVaccinations() {
        val raw = """
            {
              "patient": {"name":"Trần Bé B","birthDateText":"02/02/2021","phone":"0912345678"},
              "vaccinations": [
                {"vaccineName":"Viêm gan B","doseText":"Mũi 2","dateText":"03/03/2021","facilityName":"Trạm y tế A"}
              ]
            }
        """.trimIndent()

        val result = parser.parse(raw, RawSourceType.JSON).getOrThrow()

        assertEquals(RawSourceType.JSON, result.rawSourceType)
        assertEquals("Trần Bé B", result.patient?.fullName)
        assertEquals("0912345678", result.patient?.phoneNumber)
        assertEquals(1, result.vaccinations.size)
        assertEquals("Viêm gan B", result.vaccinations.first().vaccineName)
        assertEquals(2, result.vaccinations.first().doseNumber)
        assertEquals("03/03/2021", result.vaccinations.first().vaccinationDate)
        assertTrue(result.warnings.isEmpty())
    }

    @Test
    fun parseHtmlReturnsRowsAndMissingBirthDateWarning() {
        val raw = """
            <html><body>
              <input id="txtHoTen" value="Lê Bé C" />
              <input id="txtDienThoai" value="0987654321" />
              <table id="tblVacxin">
                <tbody>
                  <tr><td>1</td><td>Sởi</td><td>1</td><td>Lô</td><td>04/04/2022</td><td>VNCDC</td></tr>
                </tbody>
              </table>
            </body></html>
        """.trimIndent()

        val result = parser.parse(raw, RawSourceType.HTML).getOrThrow()

        assertEquals("Lê Bé C", result.patient?.fullName)
        assertEquals("0987654321", result.patient?.phoneNumber)
        assertEquals(1, result.vaccinations.size)
        assertEquals("Sởi", result.vaccinations.first().vaccineName)
        assertTrue(result.warnings.contains("Thiếu ngày sinh, không tự suy đoán."))
    }

    @Test
    fun captchaResponseReturnsCaptchaRequiredFailure() {
        val result = parser.parse("<html>captcha required</html>", RawSourceType.HTML)

        assertTrue(result.isFailure)
        assertTrue(result.exceptionOrNull() is PortalLookupException.CaptchaRequired)
    }

    @Test
    fun blankOrNoDataResponseReturnsNotFound() {
        val blank = parser.parse("   ", RawSourceType.HTML)
        val noData = parser.parse("<html><body>Không có dữ liệu phù hợp</body></html>", RawSourceType.HTML)

        assertTrue(blank.exceptionOrNull() is PortalLookupException.NotFound)
        assertTrue(noData.exceptionOrNull() is PortalLookupException.NotFound)
    }

    @Test
    fun changedJsonFormatReturnsClearParseFailure() {
        val result = parser.parse("{ this is not valid json }", RawSourceType.JSON)

        assertTrue(result.isFailure)
        assertTrue(result.exceptionOrNull() is PortalLookupException.ParseFailed)
        assertFalse(result.exceptionOrNull() is IllegalArgumentException)
    }
}
