import networks from '@/assets/networks.json'
import { SEPOLIA_RPC_URL, ARBITRUM_SEPOLIA_RPC_URL } from '@/config/common.js'

const supportedNetworks = Object.keys(networks)

const contractAddressForNetwork = (networkId) => {
  return networks[networkId]?.FundsLock?.address ?? null
}

const networkInfo = {
    "sepolia": {
        displayName: "Ethereum Sepolia",
        chainId: "0xaa36a7",
        rpcUrl: SEPOLIA_RPC_URL,
    },
    "arbitrum-sepolia": {
        displayName: "Arbitrum Sepolia",
        chainId: "0x66eee",
        rpcUrl: ARBITRUM_SEPOLIA_RPC_URL,
    },
}

export { supportedNetworks, contractAddressForNetwork, networkInfo }