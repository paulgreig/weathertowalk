package com.paulgreig.weathertowalk

import android.app.Notification
import android.app.NotificationChannel
import android.app.NotificationManager
import android.app.PendingIntent
import android.app.Service
import android.content.Intent
import android.os.Build
import android.os.IBinder
import androidx.core.app.NotificationCompat
import mobile.Mobile
import java.util.concurrent.Executors

/**
 * Foreground service that runs the Go [Mobile.runWalk] implementation (gomobile)
 * off the main thread and broadcasts the result to [MainActivity].
 */
class WeatherWalkService : Service() {

    private val executor = Executors.newSingleThreadExecutor()

    override fun onBind(intent: Intent?): IBinder? = null

    override fun onStartCommand(intent: Intent?, flags: Int, startId: Int): Int {
        val apiKey = intent?.getStringExtra(EXTRA_API_KEY) ?: ""
        val baseURL = intent?.getStringExtra(EXTRA_BASE_URL) ?: ""
        val lat = intent?.getDoubleExtra(EXTRA_LAT, DEFAULT_LAT) ?: DEFAULT_LAT
        val lon = intent?.getDoubleExtra(EXTRA_LON, DEFAULT_LON) ?: DEFAULT_LON

        createChannel()
        val pending = PendingIntent.getActivity(
            this,
            0,
            Intent(this, MainActivity::class.java),
            PendingIntent.FLAG_UPDATE_CURRENT or PendingIntent.FLAG_IMMUTABLE,
        )
        val notification: Notification = NotificationCompat.Builder(this, CHANNEL_ID)
            .setContentTitle(getString(R.string.notification_title))
            .setContentText(getString(R.string.notification_text))
            .setSmallIcon(R.drawable.ic_launcher_foreground)
            .setContentIntent(pending)
            .build()

        startForeground(NOTIFICATION_ID, notification)

        executor.execute {
            val text = try {
                Mobile.runWalk(apiKey, baseURL, lat, lon)
            } catch (e: Exception) {
                "Error: ${e.message ?: e.javaClass.simpleName}"
            }
            sendBroadcast(
                Intent(ACTION_WEATHER_RESULT).setPackage(packageName).putExtra(EXTRA_RESULT_TEXT, text),
            )
            stopForeground(STOP_FOREGROUND_DETACH)
            stopSelf()
        }

        return START_NOT_STICKY
    }

    private fun createChannel() {
        if (Build.VERSION.SDK_INT < Build.VERSION_CODES.O) {
            return
        }
        val channel = NotificationChannel(
            CHANNEL_ID,
            getString(R.string.channel_name),
            NotificationManager.IMPORTANCE_LOW,
        )
        val nm = getSystemService(NotificationManager::class.java)
        nm.createNotificationChannel(channel)
    }

    companion object {
        const val EXTRA_API_KEY = "api_key"
        const val EXTRA_BASE_URL = "base_url"
        const val EXTRA_LAT = "lat"
        const val EXTRA_LON = "lon"

        const val ACTION_WEATHER_RESULT = "com.paulgreig.weathertowalk.ACTION_WEATHER_RESULT"
        const val EXTRA_RESULT_TEXT = "result_text"

        private const val CHANNEL_ID = "weather_walk_fetch"
        private const val NOTIFICATION_ID = 1001

        private const val DEFAULT_LAT = 51.5074
        private const val DEFAULT_LON = -0.1278
    }
}
