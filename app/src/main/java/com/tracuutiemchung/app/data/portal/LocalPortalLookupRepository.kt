package com.tracuutiemchung.app.data.portal

class LocalPortalLookupRepository(
    private val sessionStore: SessionStore,
    private val parser: PortalParser = DefaultPortalParser(),
    private val searchProvider: suspend (phone: String, session: PortalSession) -> String = { phone, session ->
        VncdcPortalClient().searchByPhone(phone, session).searchHtml
    },
    private val detailProvider: suspend (patient: PortalPatientSummary, session: PortalSession) -> Pair<String, RawSourceType> = { patient, session ->
        VncdcPortalClient().fetchDetail(patient, session) to RawSourceType.HTML
    },
    private val responseProvider: (suspend (phone: String, session: PortalSession) -> Pair<String, RawSourceType>)? = null,
) : PortalLookupRepository {
    override suspend fun searchPatientsByPhone(phone: String): Result<List<PortalPatientSummary>> {
        val normalizedPhone = phone.trim()
        val session = validateLookup(normalizedPhone).getOrElse { return Result.failure(it) }
        return runCatching {
            if (responseProvider != null) {
                val (raw, _) = responseProvider.invoke(normalizedPhone, session)
                listOf(compatibilitySummary(normalizedPhone, raw))
            } else {
                parser.parsePatientSummaries(searchProvider(normalizedPhone, session)).getOrThrow()
            }
        }.recoverCatching { throw PortalLookupException.NetworkFailed(it.message ?: "VNCDC lỗi hoặc timeout.") }
    }

    override suspend fun lookupVaccinations(patient: PortalPatientSummary): Result<PortalLookupResult> {
        val session = sessionStore.get() ?: return Result.failure(PortalLookupException.MissingSession)
        if (session.isExpired()) return Result.failure(PortalLookupException.SessionExpired)
        return runCatching {
            if (responseProvider != null) {
                responseProvider.invoke(patient.phone.orEmpty(), session)
            } else {
                detailProvider(patient, session)
            }
        }.recoverCatching { throw PortalLookupException.NetworkFailed(it.message ?: "VNCDC lỗi hoặc timeout.") }
            .mapCatching { (raw, sourceType) ->
                val parsed = parser.parse(raw, sourceType).getOrThrow()
                parsed.copy(
                    patient = parsed.patient?.copy(phoneNumber = parsed.patient.phoneNumber.ifBlank { patient.phone.orEmpty() }),
                )
            }
    }

    override suspend fun lookupByPhone(phone: String): Result<PortalLookupResult> =
        searchPatientsByPhone(phone).mapCatching { patients ->
            val patient = patients.firstOrNull() ?: throw PortalLookupException.NotFound
            lookupVaccinations(patient).getOrThrow()
        }

    private fun validateLookup(phone: String): Result<PortalSession> {
        if (!isValidVietnamPhone(phone)) return Result.failure(PortalLookupException.InvalidPhone)
        val session = sessionStore.get() ?: return Result.failure(PortalLookupException.MissingSession)
        if (session.isExpired()) return Result.failure(PortalLookupException.SessionExpired)
        return Result.success(session)
    }

    private fun compatibilitySummary(phone: String, raw: String): PortalPatientSummary {
        val result = parser.parse(raw, RawSourceType.JSON).getOrNull()
        return PortalPatientSummary(
            patientId = "compat-$phone",
            fullName = result?.patient?.fullName ?: "Không rõ",
            birthDateOrYear = result?.patient?.dateOfBirth,
            phone = phone,
            address = result?.patient?.address,
        )
    }

    private fun isValidVietnamPhone(phone: String): Boolean = Regex("^0\\d{9}$").matches(phone)

    private fun PortalSession.isExpired(): Boolean = expiresAtMillis?.let { it <= System.currentTimeMillis() } == true
}

private fun demoResponse(phone: String, session: PortalSession): Pair<String, RawSourceType> {
    val sourceType = when (session.source) {
        SessionSource.HTTP_CLIENT -> RawSourceType.JSON
        SessionSource.WEBVIEW -> RawSourceType.WEBVIEW_DOM
    }
    val raw = """
        {
          "patient": {
            "name": "Nguyễn Văn A",
            "birthDateText": "01/01/2020",
            "phone": "$phone"
          },
          "vaccinations": [
            {
              "vaccineName": "Sởi - Quai bị - Rubella",
              "doseText": "Mũi 1",
              "dateText": "15/03/2021",
              "facilityName": "VNCDC"
            },
            {
              "vaccineName": "Viêm gan B",
              "doseText": "Mũi 3",
              "dateText": "20/08/2021",
              "facilityName": "VNCDC"
            }
          ]
        }
    """.trimIndent()
    return raw to sourceType
}
