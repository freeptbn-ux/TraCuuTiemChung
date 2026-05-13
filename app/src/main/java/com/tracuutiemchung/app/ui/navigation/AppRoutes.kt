package com.tracuutiemchung.app.ui.navigation

import kotlinx.serialization.Serializable

@Serializable
sealed interface AppRoute {
    @Serializable
    data object Lookup : AppRoute

    @Serializable
    data class Selection(val phoneNumber: String) : AppRoute

    @Serializable
    data class Result(val phoneNumber: String) : AppRoute
}
