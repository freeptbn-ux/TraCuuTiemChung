package com.tracuutiemchung.app.data.portal

import com.tracuutiemchung.app.data.model.PatientInfo
import com.tracuutiemchung.app.data.model.VaccinationRecord
import java.time.LocalDate
import java.time.ZoneId
import java.time.format.DateTimeFormatter
import java.time.format.DateTimeParseException
import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import org.jsoup.Jsoup
import org.jsoup.nodes.Document
import org.jsoup.nodes.Element
import org.jsoup.nodes.TextNode

class DefaultPortalParser(
    private val json: Json = Json { ignoreUnknownKeys = true },
) : PortalParser {
    override fun parse(raw: String, sourceType: RawSourceType): Result<PortalLookupResult> = runCatching {
        val normalized = raw.trim()
        if (normalized.isBlank()) throw PortalLookupException.NotFound
        if (containsCaptchaOrOtp(normalized)) throw PortalLookupException.CaptchaRequired
        if (sourceType == RawSourceType.JSON || normalized.startsWith("{")) {
            parseJson(normalized)
        } else {
            parseHtml(normalized, sourceType)
        }
    }.recoverCatching { error ->
        if (error is PortalLookupException) throw error
        throw PortalLookupException.ParseFailed
    }

    override fun parsePatientSummaries(raw: String): Result<List<PortalPatientSummary>> = runCatching {
        val normalized = raw.trim()
        if (normalized.isBlank()) return@runCatching emptyList()
        if (containsCaptchaOrOtp(normalized)) throw PortalLookupException.CaptchaRequired
        val document = Jsoup.parse(normalized)
        document.select("tr").mapNotNull(::parsePatientSummaryRow)
    }.recoverCatching { error ->
        if (error is PortalLookupException) throw error
        throw PortalLookupException.ParseFailed
    }

    private fun parseJson(raw: String): PortalLookupResult {
        val dto = json.decodeFromString(PortalLookupJsonDto.serializer(), raw)
        val records = dto.vaccinations.mapNotNull { record ->
            val vaccineName = record.vaccineName ?: return@mapNotNull null
            val dateText = record.dateText ?: return@mapNotNull null
            VaccinationRecord(
                vaccineName = vaccineName,
                doseNumber = record.doseText?.let(::extractDoseNumber),
                vaccinationDate = dateText,
                provider = record.facilityName,
            )
        }
        if (dto.patient == null && records.isEmpty()) throw PortalLookupException.NotFound
        return PortalLookupResult(
            patient = dto.patient?.toPatientInfo(),
            vaccinations = records,
            rawSourceType = RawSourceType.JSON,
            warnings = buildWarnings(dto.patient?.birthDateText, records),
            systemDate = dto.systemDate,
        )
    }

    private fun parseHtml(raw: String, sourceType: RawSourceType): PortalLookupResult {
        val detailRaw = preferredDetailHtml(raw)
        return runCatching { parseHtmlDocument(detailRaw, sourceType) }
            .recoverCatching { error ->
                if (detailRaw == raw || error !is PortalLookupException.NotFound) throw error
                parseHtmlDocument(raw, sourceType)
            }
            .getOrThrow()
    }

    private fun parseHtmlDocument(raw: String, sourceType: RawSourceType): PortalLookupResult {
        val document = Jsoup.parse(raw)
        val patient = PatientInfo(
            fullName = inputValueById(document, PortalHtmlIds.LEGACY_PATIENT_NAME) ?: "Không rõ",
            phoneNumber = inputValueById(document, "txtDienThoai").orEmpty(),
            dateOfBirth = inputValueById(document, PortalHtmlIds.PATIENT_DOB)
                ?: inputValueById(document, PortalHtmlIds.PATIENT_DOB_HIDDEN),
        )
        val parsedVaccines = parseLegacyVaccineRows(document, raw)
        val records = parsedVaccines.records
        if (records.isEmpty() && patient.fullName == "Không rõ" && patient.phoneNumber.isBlank()) {
            throw PortalLookupException.NotFound
        }
        return PortalLookupResult(
            patient = patient,
            vaccinations = records,
            rawSourceType = sourceType,
            warnings = buildWarnings(patient.dateOfBirth, records, parsedVaccines.warning),
            systemDate = inputValueById(document, PortalHtmlIds.CURRENT_SYSTEM_DATE)
                ?: inputValueById(document, PortalHtmlIds.CURRENT_DATE_HIDDEN)
                ?: vietnamTodayText(),
        )
    }

    private fun parseLegacyVaccineRows(document: Document, raw: String): LegacyVaccineParseResult {
        val table = document.selectFirst("table#${PortalHtmlIds.VACCINE_TABLE}")
            ?: return LegacyVaccineParseResult(
                records = emptyList(),
                warning = "Không tìm thấy bảng vắc-xin (id='${PortalHtmlIds.VACCINE_TABLE}').",
            )
        if (!rawTableHtmlById(raw, PortalHtmlIds.VACCINE_TABLE).orEmpty().contains("<tbody", ignoreCase = true)) {
            return LegacyVaccineParseResult(
                records = emptyList(),
                warning = "Không tìm thấy tbody trong bảng vắc-xin.",
            )
        }
        val tbody = table.selectFirst("tbody")
            ?: return LegacyVaccineParseResult(
                records = emptyList(),
                warning = "Không tìm thấy tbody trong bảng vắc-xin.",
            )
        val rows = tbody.select("tr")
        if (rows.isEmpty()) {
            return LegacyVaccineParseResult(
                records = emptyList(),
                warning = "Không tìm thấy hàng nào (tr) trong tbody.",
            )
        }

        val records = rows.mapNotNull { row ->
            val cols = row.select("td")
            if (cols.size <= LegacyVaccineColumn.DATE) return@mapNotNull null
            val vaccineName = firstDirectText(cols[LegacyVaccineColumn.VACCINE_NAME])
                ?: cols[LegacyVaccineColumn.VACCINE_NAME].text().trim().takeIf { it.isNotBlank() }
                ?: return@mapNotNull null
            val doseText = cols[LegacyVaccineColumn.DOSE_TEXT].text().trim()
            val dateText = cols[LegacyVaccineColumn.DATE].text().trim().takeIf { it.isNotBlank() }
                ?: return@mapNotNull null
            if (!isLegacyValidDate(dateText)) return@mapNotNull null
            VaccinationRecord(
                vaccineName = vaccineName,
                doseNumber = doseText.toIntOrNull() ?: 0,
                vaccinationDate = dateText,
                provider = cols.drop(LegacyVaccineColumn.DATE + 1).firstOrNull { it.text().trim().isNotBlank() }
                    ?.text()
                    ?.trim(),
            )
        }.toList()

        return LegacyVaccineParseResult(
            records = records,
            warning = if (records.isEmpty()) {
                "Không tìm thấy vắc-xin nào trong bảng có định dạng phù hợp."
            } else {
                null
            },
        )
    }

    private fun parsePatientSummaryRow(row: Element): PortalPatientSummary? {
        val cols = row.select("td")
        if (cols.isEmpty()) return null
        val action = row.select("a[onclick],button[onclick],[data-id],[data-doituongid]").firstOrNull() ?: row
        val onclick = action.attr("onclick")
        val patientId = action.attr("data-id").ifBlank { action.attr("data-doituongid") }.ifBlank {
            Regex("OnShowDetail\\((\\d+)\\)", RegexOption.IGNORE_CASE).find(onclick)?.groupValues?.getOrNull(1).orEmpty()
        }
        if (patientId.isBlank()) return null
        fun col(index: Int): String? = cols.getOrNull(index)?.text()?.trim()?.takeIf { it.isNotBlank() }
        val href = action.attr("href").takeIf { it.isNotBlank() && it != "#" }
        return PortalPatientSummary(
            patientId = patientId,
            detailUrl = href,
            detailPayload = onclick.takeIf { it.isNotBlank() },
            fullName = col(1) ?: col(0) ?: "Không rõ",
            birthDateOrYear = col(2),
            gender = col(3),
            phone = col(4),
            address = col(5),
            receivedDate = col(6),
        )
    }

    private fun containsCaptchaOrOtp(value: String): Boolean {
        val lower = value.lowercase()
        // If it looks like a full valid page (contains common dashboard markers), it's not a captcha challenge
        if (lower.contains("quản lý tiêm chủng") || lower.contains("thông tin đối tượng tiêm")) {
            return false
        }
        return lower.contains("captcha") || lower.contains("otp")
    }

    private fun preferredDetailHtml(raw: String): String {
        val markerIndex = raw.lastIndexOf(DETAIL_MARKER, ignoreCase = true)
        if (markerIndex < 0) return raw
        return raw.substring(markerIndex + DETAIL_MARKER.length).trim().takeIf { it.isNotBlank() } ?: raw
    }

    private fun rawTableHtmlById(raw: String, id: String): String? =
        Regex(
            "<table[^>]*id=[\\\"']$id[\\\"'][^>]*>.*?</table>",
            setOf(RegexOption.IGNORE_CASE, RegexOption.DOT_MATCHES_ALL),
        ).find(raw)?.value

    private fun inputValueById(document: Document, id: String): String? =
        document.selectFirst("input#$id")
            ?.attr("value")
            ?.trim()
            ?.takeIf { it.isNotBlank() }

    private fun firstDirectText(element: Element): String? = element.childNodes()
        .filterIsInstance<TextNode>()
        .firstNotNullOfOrNull { textNode -> textNode.text().trim().takeIf { it.isNotBlank() } }

    private fun isLegacyValidDate(value: String): Boolean = try {
        LocalDate.parse(value.replace(" ", ""), LEGACY_DATE_FORMATTER)
        true
    } catch (_: DateTimeParseException) {
        false
    }

    private fun buildWarnings(
        birthDateText: String?,
        records: List<VaccinationRecord>,
        legacyParserWarning: String? = null,
    ): List<String> = buildList {
        if (!legacyParserWarning.isNullOrBlank()) add(legacyParserWarning)
        if (birthDateText.isNullOrBlank()) add("Thiếu ngày sinh, không tự suy đoán.")
        if (records.isEmpty()) add("Không tìm thấy lịch sử tiêm trong phản hồi VNCDC.")
        records.forEachIndexed { index, record ->
            if (record.vaccinationDate.isBlank()) add("Bản ghi ${index + 1} thiếu ngày tiêm.")
        }
    }

    private fun extractDoseNumber(value: String): Int? = Regex("\\d+").find(value)?.value?.toIntOrNull()

    private fun vietnamTodayText(): String = LocalDate.now(VIETNAM_ZONE).format(LEGACY_DATE_FORMATTER)

    companion object {
        private val VIETNAM_ZONE: ZoneId = ZoneId.of("Asia/Ho_Chi_Minh")
        private val LEGACY_DATE_FORMATTER: DateTimeFormatter = DateTimeFormatter.ofPattern("dd/MM/yyyy")

        fun extractSubjectIds(raw: String): List<String> =
            Regex("OnShowDetail\\((\\d+)\\)", RegexOption.IGNORE_CASE)
                .findAll(raw)
                .map { it.groupValues[1] }
                .distinct()
                .toList()
    }
}

private data class LegacyVaccineParseResult(
    val records: List<VaccinationRecord>,
    val warning: String? = null,
)

internal const val DETAIL_MARKER = "<!-- VNCDC_DETAIL -->"

object PortalHtmlIds {
    const val LEGACY_PATIENT_NAME = "txtHoTen"
    const val PATIENT_NAME = "txtTenDoiTuong"
    const val PATIENT_DOB = "txtNgaySinh"
    const val PATIENT_DOB_HIDDEN = "hfNgaySinhDoiTuong"
    const val CURRENT_SYSTEM_DATE = "CurrentSystemDate"
    const val CURRENT_DATE_HIDDEN = "hfNgayHienTai"
    const val VACCINE_TABLE = "tblVacxin"
}

object LegacyVaccineColumn {
    const val VACCINE_NAME = 1
    const val DOSE_TEXT = 2
    const val DATE = 4
}

@Serializable
private data class PortalLookupJsonDto(
    val patient: PatientJsonDto? = null,
    val vaccinations: List<VaccinationJsonDto> = emptyList(),
    val systemDate: String? = null,
)

@Serializable
private data class PatientJsonDto(
    val name: String? = null,
    @SerialName("birthDateText") val birthDateText: String? = null,
    val phone: String? = null,
) {
    fun toPatientInfo(): PatientInfo = PatientInfo(
        fullName = name ?: "Không rõ",
        phoneNumber = phone.orEmpty(),
        dateOfBirth = birthDateText,
    )
}

@Serializable
private data class VaccinationJsonDto(
    val vaccineName: String? = null,
    val doseText: String? = null,
    val dateText: String? = null,
    val facilityName: String? = null,
    val lotNumber: String? = null,
)
