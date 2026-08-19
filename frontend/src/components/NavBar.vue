<script setup>
import { onMounted, onUnmounted } from 'vue'
import {appKitModal} from '@/auth/appkit.js'
import {refreshToken} from '@/auth/authenticate.js'

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
// import { createAppKit } from '@reown/appkit'
// import { mainnet, sepolia } from '@reown/appkit/networks'
// import { BrowserProvider } from 'ethers'
// import { EthersAdapter } from "@reown/appkit-adapter-ethers";

// const PROJECT_ID = 'e103ded1feccd94ad1f759a8a13dce0f' // Replace with your Reown dashboard project ID
// const BACKEND_URL = 'http://localhost:3000'
// const ACCESS_TOKEN_KEY = 'access_token'

// // Initialize WalletConnect Modal
// const modal = createAppKit({
//   adapters: [new EthersAdapter()],
//   networks: [mainnet, sepolia],
//   projectId: PROJECT_ID,
//   metadata: {
//     name: 'Go Web3 Test App',
//     description: 'Testing Go Auth with WalletConnect',
//     url: window.location.origin,
//     icons: ['https://assets.reown.com/reown-profile-pic.png']
//   },
//   features: {
//     analytics: true,
//     email: true
//   }
// })

// modal.subscribeEvents((event) => {
//   const { type, event: name, properties } = event.data
//   console.log('Event name: ' + name)
//   if (name === 'INITIALIZE') {
//     console.log('AppKit initialized!')
//   }
//   if (name === 'CONNECT_SUCCESS') {
//     console.log('Wallet connected!', properties)
//   }
//   if (name === 'DISCONNECT_SUCCESS') {
//     console.log('Wallet disconnected!')
//   }
// })



// import { createAppKit, useAppKit, useAppKitEvents, useAppKitProvider } from '@reown/appkit/vue';
// import { EthersAdapter } from '@reown/appkit-adapter-ethers';
// import { mainnet, sepolia } from '@reown/appkit/networks'
// import { watch } from 'vue';
// import { useAppKitAccount } from "@reown/appkit/vue";

// const PROJECT_ID = 'e103ded1feccd94ad1f759a8a13dce0f' // AppKit project ID


// createAppKit({
//   projectId: PROJECT_ID,
//   adapters: [new EthersAdapter()],
//   networks: [mainnet, sepolia],
//   metadata: {
//     name: 'FundsLock',
//     description: 'Web3 escrow application for locking funds until a certain date or condition is met.',
//     icon: 'https://assets.reown.com/reown-profile-pic.png',
//   },
//   features: {
//    analytics: true,
//    email: true,
//   },
// });

// const events = useAppKitEvents();
// const accountData = useAppKitAccount();
// const walletProvider = useAppKitProvider("eip155");

// watch(events, (event) => {
//   const { type, event: name, properties } = event.data;
//   console.log("Event name: " + name);
//   if (name == "INITIALIZE") {
//     console.log("AppKit initialized!");
//     console.log("Is connected: " + accountData.value.isConnected);
//     console.log( walletProvider.signer);
//   }
//   if (name === "CONNECT_SUCCESS") {
//     console.log("Wallet connected!", properties);
//   }

//   if (name === "DISCONNECT_SUCCESS") {
//     console.log("Wallet disconnected!");
//   }
// });

//  const modal = useAppKit();
 
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
          <a class="nav-link active" aria-current="page" href="#">Home</a>
        </li>
        <li class="nav-item">
          <a class="nav-link" href="#">Link</a>
        </li>
        <li class="nav-item dropdown">
          <a class="nav-link dropdown-toggle" href="#" id="navbarDropdown" role="button" data-bs-toggle="dropdown" aria-expanded="false">
            Dropdown
          </a>
          <ul class="dropdown-menu" aria-labelledby="navbarDropdown">
            <li><a class="dropdown-item" href="#">Action</a></li>
            <li><a class="dropdown-item" href="#">Another action</a></li>
            <li><hr class="dropdown-divider"></li>
            <li><a class="dropdown-item" href="#">Something else here</a></li>
          </ul>
        </li>
        <li class="nav-item">
          <a class="nav-link disabled" href="#" tabindex="-1" aria-disabled="true">Disabled</a>
        </li>
      </ul>
      <appkit-button balance="show" />
    </div>
  </div>
</nav>
</template>