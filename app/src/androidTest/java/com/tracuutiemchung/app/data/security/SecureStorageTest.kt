package com.tracuutiemchung.app.data.security

import android.content.Context
import androidx.test.core.app.ApplicationProvider
import androidx.test.ext.junit.runners.AndroidJUnit4
import org.junit.Assert.assertEquals
import org.junit.Assert.assertNull
import org.junit.Before
import org.junit.Test
import org.junit.runner.RunWith

@RunWith(AndroidJUnit4::class)
class SecureStorageTest {
    private lateinit var secureStorage: SecureStorage

    @Before
    fun setup() {
        val context = ApplicationProvider.getApplicationContext<Context>()
        secureStorage = SecureStorage(context)
        secureStorage.clear()
    }

    @Test
    fun testSaveAndGetString() {
        val key = "test_key"
        val value = "test_value"
        secureStorage.saveString(key, value)
        assertEquals(value, secureStorage.getString(key))
    }

    @Test
    fun testClear() {
        val key = "test_key"
        val value = "test_value"
        secureStorage.saveString(key, value)
        secureStorage.clear()
        assertNull(secureStorage.getString(key))
    }
}
