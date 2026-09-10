
async function switchToSepolia(ethereum) {
  if (!ethereum) return
  const sepoliaChainId = '0xaa36a7' // Sepolia chain ID in hex
  try {
    await ethereum.request({
      method: 'wallet_switchEthereumChain',
      params: [{ chainId: sepoliaChainId }]
    })
  } catch (error) {
    console.error('Failed to switch to Sepolia:', error)
    alert('Please switch to Sepolia network in your wallet')
    throw error
  }
}

export {switchToSepolia}