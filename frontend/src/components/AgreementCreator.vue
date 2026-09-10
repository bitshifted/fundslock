<script setup>
import contractAbi from '@/assets/abi/FundsLock.json'
import { CONTRACT_ADDRESS } from '@/config/common';
import { onMounted, ref } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useAppKitProvider } from '@reown/appkit/vue';
import { BrowserProvider } from 'ethers';
import { Contract } from 'ethers';

const appKitProvider = useAppKitProvider('eip155')


function getAbi() {
    return contractAbi.abi ? contractAbi.abi  : contractAbi
}

const counterparty = ref('')
const amount = ref(0)
const currency = ref('ETH')
const network = ref('sepolia')
const side = ref('buyer')

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

const  createAgreement = async () => {
  const ethereum = appKitProvider?.walletProvider
  if (!ethereum) {
    alert('Please connect your wallet first')
    return
  }
  try {
    await switchToSepolia(ethereum)
    const provider = new BrowserProvider(ethereum)
    const signer = await provider.getSigner()
    const contract = new Contract(CONTRACT_ADDRESS, getAbi(), signer)
    const address = await signer.getAddress()
    const calcAmount = amount.value * 1000000000000000000
    const buyerAddress = side.value === 'buyer' ? address : counterparty.value
    const sellerAddress = side.value === 'seller' ? address : counterparty.value

    console.log(`Creating agreement with address: ${address}, counterparty: ${counterparty.value}, amount: ${calcAmount}`)
    const tx = await contract.createAgreement(sellerAddress, buyerAddress, calcAmount)
    console.log("tansaction: " + tx)
  } catch(err) {
    console.log(err)
  }
  
}

</script>

<template>
  <div class="container">
     <h1>Create new agreement</h1>
   <form @submit.prevent="createAgreement">
  <div class="form-group row mb-3">
    <label for="side" class="col-sm-2 col-form-label">Side: </label>
    <div class="col-sm-1">
      <select id="side" class="form-control" v-model="side">
        <option value="buyer">Buyer</option>
        <option value="seller">Seller</option>
      </select>
    </div>
    <label for="side" class="col-sm-2 col-form-label">Network: </label>
    <div class="col-sm-1">
      <select id="network" class="form-control" v-model="network">
        <option value="ethereum">Ethereum</option>
        <option value="sepolia">Sepolia</option>
      </select>
    </div>
  </div>
  <div class="form-group row">
    <label for="counterparty" class="col-sm-2 col-form-label">Counterparty address</label>
    <div class="col-sm-10">
      <input type="text" v-model="counterparty"  class="form-control" id="counterparty" placeholder="Counterparty address">
    </div>
  </div>
   <div class="form-group row">
    <label for="amount" class="col-sm-2 col-form-label">Agreement amount:</label>
    <div class="input-group col-sm-10">
      <input type="text" v-model="amount" class="form-control" id="amount" placeholder="Agreement amount">
      <select v-model="currency">
        <option value="ETH">ETH</option>
        <option value="USDC">USDC</option>
        <option value="USDT">USDT</option>
      </select>
    </div>
  </div>
  <div class="form-group row">
    <div class="col-sm-10">
      <button type="submit" class="btn btn-primary">Create</button>
    </div>
  </div>
</form>
  </div>
   

</template>