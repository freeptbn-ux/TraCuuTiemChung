package com.tracuutiemchung.app

import org.junit.Assert.assertTrue
import org.junit.Test
import java.io.File

class SecurityConfigTest {

    private fun getFile(path: String): File {
        // Try direct path (for running from app module)
        val file = File(path)
        if (file.exists()) return file
        // Try prepending "app/" (for running from project root)
        val appFile = File("app/$path")
        if (appFile.exists()) return appFile
        return file // Return the original one (which doesn't exist) to trigger the assertion failure with correct message
    }

    @Test
    fun testNetworkSecurityConfigExists() {
        val configFile = getFile("src/main/res/xml/network_security_config.xml")
        assertTrue("Network security config file should exist at ${configFile.absolutePath}", configFile.exists())
        
        val content = configFile.readText()
        assertTrue("Cleartext traffic should be disabled", content.contains("cleartextTrafficPermitted=\"false\""))
        assertTrue("Vercel domain should be present", content.contains("tracuutiemchung.vercel.app"))
    }

    @Test
    fun testAndroidManifestIncludesSecurityConfig() {
        val manifestFile = getFile("src/main/AndroidManifest.xml")
        assertTrue("AndroidManifest.xml should exist at ${manifestFile.absolutePath}", manifestFile.exists())
        
        val content = manifestFile.readText()
        assertTrue("AndroidManifest should reference network_security_config", 
            content.contains("android:networkSecurityConfig=\"@xml/network_security_config\""))
    }

    @Test
    fun testProGuardRulesExist() {
        val proguardFile = getFile("proguard-rules.pro")
        assertTrue("ProGuard rules file should exist at ${proguardFile.absolutePath}", proguardFile.exists())
        
        val content = proguardFile.readText()
        assertTrue("ProGuard should have rules for Retrofit", content.contains("retrofit2"))
        assertTrue("ProGuard should have rules for OkHttp", content.contains("okhttp3"))
        assertTrue("ProGuard should have rules for data models", content.contains("com.tracuutiemchung.app.data.model"))
        assertTrue("ProGuard should dontwarn jspecify", content.contains("org.jspecify.annotations"))
    }
}
