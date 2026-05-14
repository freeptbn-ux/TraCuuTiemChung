package com.tracuutiemchung.app.data.portal

import org.junit.Assert.assertEquals
import org.junit.Assert.assertTrue
import org.junit.Test

class DefaultPortalParserPatientSummaryTest {
    private val parser = DefaultPortalParser()

    @Test
    fun emptySearchHtmlReturnsEmptyList() {
        val html = """
            <table><tbody></tbody></table>
        """.trimIndent()

        val result = parser.parsePatientSummaries(html).getOrThrow()

        assertTrue(result.isEmpty())
    }

    @Test
    fun searchHtmlWithTwoRowsReturnsPatientSummaries() {
        val html = """
            <table><tbody>
              <tr>
                <td>1</td><td>Nguyễn Văn A</td><td>2019</td><td>Nam</td>
                <td>0912345678</td><td>Hà Nội</td><td>01/05/2026</td>
                <td><a onclick="OnShowDetail(101)">Chi tiết</a></td>
              </tr>
              <tr>
                <td>2</td><td>Trần Thị B</td><td>02/02/2020</td><td>Nữ</td>
                <td>0912345678</td><td>Đà Nẵng</td><td>02/05/2026</td>
                <td><button onclick="OnShowDetail(202)">Chi tiết</button></td>
              </tr>
            </tbody></table>
        """.trimIndent()

        val result = parser.parsePatientSummaries(html).getOrThrow()

        assertEquals(2, result.size)
        assertEquals("101", result[0].patientId)
        assertEquals("Nguyễn Văn A", result[0].fullName)
        assertEquals("2019", result[0].birthDateOrYear)
        assertEquals("Nam", result[0].gender)
        assertEquals("0912345678", result[0].phone)
        assertEquals("Hà Nội", result[0].address)
        assertEquals("01/05/2026", result[0].receivedDate)
        assertEquals("202", result[1].patientId)
        assertEquals("Trần Thị B", result[1].fullName)
    }
}
