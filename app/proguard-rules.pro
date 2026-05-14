# Kotlinx Serialization
-keepattributes *Annotation*, EnclosingMethod, Signature
-keep,allowshrinking class com.tracuutiemchung.app.data.model.** { *; }
-keep,allowshrinking class com.tracuutiemchung.app.data.remote.model.** { *; }
-keepclassmembers class com.tracuutiemchung.app.data.model.** {
    *** Companion;
}
-keepclassmembers class com.tracuutiemchung.app.data.remote.model.** {
    *** Companion;
}

# Retrofit
-keepattributes Signature, InnerClasses, EnclosingMethod
-keepattributes RuntimeVisibleAnnotations, RuntimeVisibleParameterAnnotations
-keepattributes RuntimeInvisibleAnnotations, RuntimeInvisibleParameterAnnotations
-keep class retrofit2.** { *; }
-keep interface retrofit2.** { *; }
-dontwarn retrofit2.**

# OkHttp
-keepattributes Signature
-keepattributes *Annotation*
-keep class okhttp3.** { *; }
-dontwarn okhttp3.**
-dontwarn okio.**
-dontwarn javax.annotation.**
-dontwarn org.conscrypt.**

# Jsoup
-keep class org.jsoup.** { *; }
-dontwarn org.jspecify.annotations.**
