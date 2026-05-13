package com.tracuutiemchung.app.data.remote

import com.tracuutiemchung.app.BuildConfig
import okhttp3.mockwebserver.MockResponse
import okhttp3.mockwebserver.MockWebServer
import org.junit.After
import org.junit.Assert.assertEquals
import org.junit.Before
import org.junit.Test
import kotlinx.coroutines.runBlocking
import java.lang.reflect.Field

class RetrofitClientTest {
    private lateinit var server: MockWebServer

    @Before
    fun setUp() {
        server = MockWebServer()
        server.start()
    }

    @After
    fun tearDown() {
        server.shutdown()
    }

    @Test
    fun `test retrofit client adds X-API-KEY header`() = runBlocking {
        // We need to use reflection because RetrofitClient is an object with private members
        // and BASE_URL is hardcoded to BuildConfig.BASE_URL
        
        // Since we can't easily change the baseUrl of the singleton, 
        // we'll verify the interceptor logic by creating a similar OkHttpClient
        
        val okHttpClientField: Field = RetrofitClient::class.java.getDeclaredField("okHttpClient")
        okHttpClientField.isAccessible = true
        val client = okHttpClientField.get(RetrofitClient) as okhttp3.OkHttpClient
        
        // Mock a response
        server.enqueue(MockResponse().setBody("{}"))
        
        // We can't easily point RetrofitClient to our mock server because BASE_URL is fixed
        // So we test the client directly
        val request = okhttp3.Request.Builder()
            .url(server.url("/"))
            .build()
            
        client.newCall(request).execute().use { response ->
            val recordedRequest = server.takeRequest()
            assertEquals(BuildConfig.X_API_KEY, recordedRequest.getHeader("X-API-KEY"))
        }
    }
}
