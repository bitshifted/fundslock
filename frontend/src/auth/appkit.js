import { createAppKit } from '@reown/appkit'
import { mainnet, sepolia } from '@reown/appkit/networks'
import { EthersAdapter } from "@reown/appkit-adapter-ethers";
import {siweConfig } from '@/auth/siwe.js'


const PROJECT_ID = 'e103ded1feccd94ad1f759a8a13dce0f' // Replace with your Reown dashboard project ID


// Initialize WalletConnect Modal
const appKitModal = createAppKit({
  adapters: [new EthersAdapter()],
  networks: [mainnet, sepolia],
  defaultNetwork: sepolia,
  projectId: PROJECT_ID,
  metadata: {
    name: 'FundsLock',
    description: 'FundsLock distributed escrow system',
    url: window.location.origin,
    icons: ['https://assets.reown.com/reown-profile-pic.png']
  },
  features: {
    analytics: true,
    email: true
  },
  siweConfig: siweConfig
})


export {appKitModal}
