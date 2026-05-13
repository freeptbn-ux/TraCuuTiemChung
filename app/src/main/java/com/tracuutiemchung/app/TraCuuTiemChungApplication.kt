package com.tracuutiemchung.app

import android.app.Application
import com.tracuutiemchung.app.BuildConfig
import timber.log.Timber

class TraCuuTiemChungApplication : Application() {
    override fun onCreate() {
        super.onCreate()
        if (BuildConfig.DEBUG) {
            Timber.plant(Timber.DebugTree())
        }
    }
}
