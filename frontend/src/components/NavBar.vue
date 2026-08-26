<script setup>
import { onMounted, onUnmounted } from 'vue'
import {appKitModal} from '@/auth/appkit.js'
import {refreshToken} from '@/auth/authenticate.js'
import { RouterLink } from 'vue-router'

let _unsubscribeProviders

onMounted(async () => {
  // Subscribe to provider changes
  _unsubscribeProviders = appKitModal.subscribeProviders((state) => {
    console.log('Provider state changed:', state)
    const provider = state['eip155']
    if (provider) {
      console.log('Wallet connected')
      refreshToken()
    } else {
      console.log('Wallet disconnected')
    }
  })

  // Also attempt to log the current provider state immediately (some AppKit versions
  // expose a `getProviders` method). This ensures we see 'Wallet disconnected' on reload.
  try {
    if (typeof appKitModal.getProviders === 'function') {
      const current = await appKitModal.getProviders()
      console.log('Initial provider state:', current)
      const p = current['eip155']
      if (p) console.log('Wallet connected')
      else console.log('Wallet disconnected')
    }
  } catch (e) {
    console.warn('Could not read initial providers from appKitModal', e)
  }
})

onUnmounted(() => {
  if (typeof _unsubscribeProviders === 'function') {
    _unsubscribeProviders()
  }
})
 
</script>

<template>

<nav class="navbar navbar-expand-md navbar-light bg-light">
  <div class="container">
    <a class="navbar-brand" href="#">FundsLock</a>
    <button class="navbar-toggler" type="button" data-bs-toggle="collapse" data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent" aria-expanded="false" aria-label="Toggle navigation">
      <span class="navbar-toggler-icon"></span>
    </button>
    <div class="collapse navbar-collapse" id="navbarSupportedContent">
      <ul class="navbar-nav me-auto mb-2 mb-lg-0">
        <li class="nav-item">
          <RouterLink to="/" class="nav-link active" aria-current="page">Home</RouterLink>
          <!-- <a class="nav-link active" aria-current="page" href="/">Home</a> -->
        </li>
        <li class="nav-item">
          <RouterLink to="/account" class="nav-link">Account</RouterLink>
          <!-- <a class="nav-link" href="/account">Account</a> -->
        </li>
      </ul>
      <appkit-button balance="show" />
    </div>
  </div>
</nav>
</template>