
async function switchChain(ethereum, chainId) {
  if (!ethereum) return
  try {
    await ethereum.request({
      method: 'wallet_switchEthereumChain',
      params: [{ chainId }]
    })
  } catch (error) {
    console.error('Failed to switch chain:', error)
    alert('Please switch to the desired network in your wallet')
    throw error
  }
}

export {switchChain as switchToSepolia}