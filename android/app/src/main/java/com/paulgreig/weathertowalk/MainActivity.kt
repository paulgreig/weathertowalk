package com.paulgreig.weathertowalk

import android.Manifest
import android.content.BroadcastReceiver
import android.content.Context
import android.content.Intent
import android.content.IntentFilter
import android.content.pm.PackageManager
import android.os.Build
import android.os.Bundle
import android.widget.Toast
import androidx.activity.result.contract.ActivityResultContracts
import androidx.appcompat.app.AppCompatActivity
import androidx.core.content.ContextCompat
import com.paulgreig.weathertowalk.databinding.ActivityMainBinding

class MainActivity : AppCompatActivity() {

    private lateinit var binding: ActivityMainBinding

    private val resultReceiver = object : BroadcastReceiver() {
        override fun onReceive(context: Context?, intent: Intent?) {
            if (intent?.action != WeatherWalkService.ACTION_WEATHER_RESULT) {
                return
            }
            val text = intent.getStringExtra(WeatherWalkService.EXTRA_RESULT_TEXT) ?: return
            binding.output.text = text
        }
    }

    private val requestNotifications = registerForActivityResult(
        ActivityResultContracts.RequestPermission(),
    ) { granted ->
        if (!granted && Build.VERSION.SDK_INT >= 33) {
            Toast.makeText(this, "Notifications disabled; foreground service may be less visible.", Toast.LENGTH_LONG).show()
        }
    }

    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        binding = ActivityMainBinding.inflate(layoutInflater)
        setContentView(binding.root)

        if (Build.VERSION.SDK_INT >= 33) {
            when {
                ContextCompat.checkSelfPermission(this, Manifest.permission.POST_NOTIFICATIONS) ==
                    PackageManager.PERMISSION_GRANTED -> { /* ok */ }
                else -> requestNotifications.launch(Manifest.permission.POST_NOTIFICATIONS)
            }
        }

        binding.buttonFetch.setOnClickListener {
            val apiKey = binding.inputApiKey.text?.toString()?.trim().orEmpty()
            if (apiKey.length != 32) {
                Toast.makeText(this, "API key must be 32 hex characters.", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }
            val lat = binding.inputLat.text?.toString()?.toDoubleOrNull()
            val lon = binding.inputLon.text?.toString()?.toDoubleOrNull()
            if (lat == null || lon == null) {
                Toast.makeText(this, "Enter valid latitude and longitude.", Toast.LENGTH_SHORT).show()
                return@setOnClickListener
            }
            val base = binding.inputBaseUrl.text?.toString()?.trim().orEmpty()
            val intent = Intent(this, WeatherWalkService::class.java).apply {
                putExtra(WeatherWalkService.EXTRA_API_KEY, apiKey)
                putExtra(WeatherWalkService.EXTRA_BASE_URL, base)
                putExtra(WeatherWalkService.EXTRA_LAT, lat)
                putExtra(WeatherWalkService.EXTRA_LON, lon)
            }
            ContextCompat.startForegroundService(this, intent)
            binding.output.text = getString(R.string.notification_text)
        }
    }

    override fun onStart() {
        super.onStart()
        val filter = IntentFilter(WeatherWalkService.ACTION_WEATHER_RESULT)
        if (Build.VERSION.SDK_INT >= Build.VERSION_CODES.TIRAMISU) {
            registerReceiver(resultReceiver, filter, Context.RECEIVER_NOT_EXPORTED)
        } else {
            @Suppress("DEPRECATION")
            registerReceiver(resultReceiver, filter)
        }
    }

    override fun onStop() {
        super.onStop()
        unregisterReceiver(resultReceiver)
    }
}
