package com.tracuutiemchung.app.data.rules

import kotlinx.serialization.SerialName
import kotlinx.serialization.Serializable
import kotlinx.serialization.json.Json
import kotlinx.serialization.json.JsonObject
import kotlinx.serialization.json.decodeFromJsonElement

@Serializable
data class VaccineRule(
    val vaccineKey: String,
    val displayName: String,
    val type: String,
    val rawNames: List<String> = emptyList(),
    val recognitionKeywords: List<String> = emptyList(),
    val minAgeDays: Int? = null,
    val maxAgeDays: Int? = null,
    val intervalDays: Int? = null,
    val requiredDoses: Int? = null,
    val notes: List<String> = emptyList(),
    val minIntervalDays: List<Int?> = emptyList(),
    val doseSpecificRules: Map<String, DoseSpecificRule> = emptyMap(),
    val boosterIntervalYears: Int? = null,
    val boosterAfterDoseNumber: Int? = null,
    val boosterMaxAgeYears: Int? = null,
    val rulesByAge: List<AgeBasedRegimen> = emptyList(),
    val regimens: List<AgeBasedRegimen> = emptyList(),
    val members: Map<String, MemberConfig> = emptyMap(),
    val courses: List<CourseConfig> = emptyList(),
    val interactions: Map<String, InteractionConfig> = emptyMap(),
    val providesMeaslesProtection: Boolean = false,
    val isLive: Boolean = false,
    val maxAgeMonthsToStart: Int? = null,
    val maxAgeMonthsForCompletion: Int? = null,
    val initialSeriesIntervalDays: Int? = null,
    val minAgeWeeksAtFirstDose: Int? = null,
    val minAgeMonthsAtFirstDose: Int? = null,
    val minAgeMonthsOverall: Int? = null,
    val minAgeYearsAtFirstDose: Int? = null,
    val minAgeDaysAtFirstDose: Int? = null,
    val minAgeMonthsOverallGroup: Int? = null,
)

@Serializable
private data class RawVaccineRule(
    @SerialName("display_name") val displayName: String? = null,
    @SerialName("group_display_name") val groupDisplayName: String? = null,
    val type: String,
    @SerialName("raw_names") val rawNames: List<String> = emptyList(),
    @SerialName("raw_names_members") val rawNamesMembers: Map<String, List<String>> = emptyMap(),
    @SerialName("recognition_keywords") val recognitionKeywords: List<String> = emptyList(),
    @SerialName("doses_required") val dosesRequired: Int? = null,
    @SerialName("min_interval_days") val minIntervalDays: List<Int?> = emptyList(),
    @SerialName("min_age_weeks_at_first_dose") val minAgeWeeksAtFirstDose: Int? = null,
    @SerialName("min_age_weeks_overall") val minAgeWeeksOverall: Int? = null,
    @SerialName("min_age_months_at_first_dose") val minAgeMonthsAtFirstDose: Int? = null,
    @SerialName("min_age_months_overall_group") val minAgeMonthsOverallGroup: Int? = null,
    @SerialName("min_age_years_at_first_dose") val minAgeYearsAtFirstDose: Int? = null,
    @SerialName("dose_specific_rules") val doseSpecificRules: Map<String, JsonObject> = emptyMap(),
    @SerialName("booster_interval_years") val boosterIntervalYears: Int? = null,
    @SerialName("booster_after_dose_number") val boosterAfterDoseNumber: Int? = null,
    @SerialName("booster_max_age_years") val boosterMaxAgeYears: Int? = null,
    @SerialName("rules_by_age") val rulesByAge: List<JsonObject> = emptyList(),
    @SerialName("regimens") val regimens: List<JsonObject> = emptyList(),
    @SerialName("members") val members: Map<String, JsonObject> = emptyMap(),
    @SerialName("courses") val courses: List<JsonObject> = emptyList(),
    @SerialName("interactions") val interactions: Map<String, JsonObject> = emptyMap(),
    @SerialName("provides_measles_protection_group") val providesMeaslesProtectionGroup: Boolean = false,
    @SerialName("max_age_months_to_start_first_dose_group") val maxAgeMonthsToStartGroup: Int? = null,
    @SerialName("max_age_months_for_completion_group") val maxAgeMonthsForCompletionGroup: Int? = null,
    @SerialName("initial_series_interval_days") val initialSeriesIntervalDays: Int? = null,
    @SerialName("is_live") val isLive: Boolean = false,
    @SerialName("min_age_months_overall") val minAgeMonthsOverall: Int? = null,
    @SerialName("min_age_years_overall") val minAgeYearsOverall: Int? = null,
    @SerialName("min_age_days_at_first_dose") val minAgeDaysAtFirstDose: Int? = null,
)

class VaccineRuleRepository(
    private val assetReader: (String) -> String,
    private val json: Json = Json { ignoreUnknownKeys = true },
) {
    fun loadRules(): List<VaccineRule> {
        val content = assetReader(VACCINE_RULES_ASSET)
        val rawRules = json.decodeFromString<Map<String, RawVaccineRule>>(content)
        return rawRules.map { (key, raw) -> raw.toVaccineRule(key) }
    }

    fun loadStandardVaccines(): List<String> {
        val content = assetReader(STANDARD_VACCINES_ASSET)
        return json.decodeFromString<List<String>>(content)
    }

    private fun RawVaccineRule.toVaccineRule(key: String): VaccineRule {
        val display = displayName ?: groupDisplayName ?: key
        val interval = minIntervalDays.filterNotNull().firstOrNull()
        val minAge = minAgeWeeksAtFirstDose?.weeksToDays()
            ?: minAgeWeeksOverall?.weeksToDays()
            ?: minAgeMonthsAtFirstDose?.monthsToDays()
            ?: minAgeMonthsOverallGroup?.monthsToDays()
            ?: minAgeYearsAtFirstDose?.yearsToDays()
            ?: minAgeMonthsOverall?.monthsToDays()
            ?: minAgeYearsOverall?.yearsToDays()
            ?: minAgeDaysAtFirstDose

        return VaccineRule(
            vaccineKey = key,
            displayName = display,
            type = type,
            rawNames = rawNames + rawNamesMembers.values.flatten(),
            recognitionKeywords = recognitionKeywords,
            minAgeDays = minAge,
            intervalDays = interval,
            requiredDoses = dosesRequired,
            minIntervalDays = minIntervalDays,
            doseSpecificRules = doseSpecificRules.mapValues { (_, v) -> json.decodeFromJsonElement<DoseSpecificRule>(v) },
            boosterIntervalYears = boosterIntervalYears,
            boosterAfterDoseNumber = boosterAfterDoseNumber,
            boosterMaxAgeYears = boosterMaxAgeYears,
            rulesByAge = rulesByAge.map { json.decodeFromJsonElement<AgeBasedRegimen>(it) },
            regimens = regimens.map { json.decodeFromJsonElement<AgeBasedRegimen>(it) },
            members = members.mapValues { (_, v) -> json.decodeFromJsonElement<MemberConfig>(v) },
            courses = courses.map { json.decodeFromJsonElement<CourseConfig>(it) },
            interactions = interactions.mapValues { (_, v) -> json.decodeFromJsonElement<InteractionConfig>(v) },
            providesMeaslesProtection = providesMeaslesProtectionGroup,
            isLive = isLive,
            maxAgeMonthsToStart = maxAgeMonthsToStartGroup,
            maxAgeMonthsForCompletion = maxAgeMonthsForCompletionGroup,
            initialSeriesIntervalDays = initialSeriesIntervalDays,
            minAgeWeeksAtFirstDose = minAgeWeeksAtFirstDose,
            minAgeMonthsAtFirstDose = minAgeMonthsAtFirstDose,
            minAgeMonthsOverall = minAgeMonthsOverall,
            minAgeYearsAtFirstDose = minAgeYearsAtFirstDose,
            minAgeDaysAtFirstDose = minAgeDaysAtFirstDose,
            minAgeMonthsOverallGroup = minAgeMonthsOverallGroup,
        )
    }

    companion object {
        const val VACCINE_RULES_ASSET = "vaccine_rules.json"
        const val STANDARD_VACCINES_ASSET = "standard_vaccines.json"
    }
}

class VaccineNameNormalizer(rules: List<VaccineRule>) {
    private val aliases: List<Pair<String, VaccineRule>> = rules.flatMap { rule ->
        val ruleAliases = (rule.rawNames + rule.recognitionKeywords + rule.displayName + rule.vaccineKey)
        val courseAliases = rule.courses.flatMap { it.rawNames }
        val memberAliases = rule.members.values.flatMap { it.rawNames }
        
        (ruleAliases + courseAliases + memberAliases)
            .filter { it.isNotBlank() }
            .map { alias -> alias.normalizeForMatching() to rule }
    }.sortedByDescending { (alias, _) -> alias.length }

    fun normalize(vncdcName: String): VaccineRule? {
        val normalizedName = vncdcName.normalizeForMatching()
        return aliases.firstOrNull { (alias, _) ->
            normalizedName == alias || normalizedName.contains(alias) || alias.contains(normalizedName)
        }?.second
    }
}

private fun Int.weeksToDays(): Int = this * 7
private fun Int.monthsToDays(): Int = this * 30
private fun Int.yearsToDays(): Int = this * 365

private fun String.normalizeForMatching(): String = lowercase()
    .replace('đ', 'd')
    .replace(Regex("[^a-z0-9]+"), " ")
    .trim()
    .replace(Regex("\\s+"), " ")
