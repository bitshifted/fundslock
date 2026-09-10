<script setup>
import contractAbi from '@/assets/abi/FundsLock.json'
import { useAuthStore } from '@/stores/auth.js';
import { onMounted, ref,computed } from 'vue';
import { BACKEND_URL,CONTRACT_ADDRESS } from '@/config/common.js'
import {STATUS_FUNDED, STATUS_RELEASED, STATUS_SELLER_ACCEPTED, STATUS_SELLER_REQUESTED_RELEASE, statusMap} from '@/assets/abi/enums.js'
import { useAppKitAccount, useAppKitProvider } from '@reown/appkit/vue';
import { switchToSepolia } from '@/eth/index.js'
import { BrowserProvider, parseEther } from 'ethers';
import { Contract } from 'ethers';

const AGREEMENTS_URL = `${BACKEND_URL}/api/v1/agreements`

const eip155Account = useAppKitAccount({ namespace: "eip155" }); 
const appKitProvider = useAppKitProvider('eip155')

function getAbi() {
    return contractAbi.abi ? contractAbi.abi  : contractAbi
}

const authStore = useAuthStore()
const agreementsList = ref([])
const opRunning = ref([])
const opSuccess = ref([])


const isOpRunning = (agreementId) => {
    return opRunning.value.includes(agreementId)
}

const showSuccessAlert = (agreementId) => {
    return opSuccess.value.includes(agreementId)
}

const resetSuccess = (agreementId) => {
    opSuccess.value = opSuccess.value.filter(id => id !== agreementId)
}

async function fetchAgreements() {
    console.log(`auth token: ${authStore.token}`)
    const response = await fetch(AGREEMENTS_URL,{
        headers: {
            'Authorization': `Bearer ${authStore.token}`
        }
    })
    agreementsList.value = await response.json()
    const curAddress = eip155Account.value.address
    console.log(`current address: ${curAddress}`)
}

function statusLabelColor(status) {
    switch (status) {
        case 0:
            return 'badge text-bg-primary'
        case 1:
            return 'badge text-bg-secondary'
        case 2:
            return 'badge text-bg-success'
        case 3:
            return 'badge text-bg-danger'
        default:
            return 'badge text-bg-light'
    }
}

function isSeller(agreement) {
    const curAddress = eip155Account.value.address
    return agreement.seller.toLowerCase() === curAddress.toLowerCase()
}

function isBuyer(agreement) {
    const curAddress = eip155Account.value.address
    return agreement.buyer.toLowerCase() === curAddress.toLowerCase()
}

function isAgreementAccepted(agreement) {
    return agreement.statusChanges.some(statusChange => statusChange.status === STATUS_SELLER_ACCEPTED) // Check if status 2 (Accepted) exists in statusChanges
}

function canBeFunded(agreement) {
    return !agreement.statusChanges.some(statusChange => statusChange.status === STATUS_FUNDED || statusChange.status === STATUS_RELEASED) // Check if status 1 (Funded) or 5 (Released) exists in statusChanges
}

function canRequestRelease(agreement) {
    return !agreement.statusChanges.some(statusChange => statusChange.status === STATUS_RELEASED || statusChange.status === STATUS_SELLER_REQUESTED_RELEASE) // Check if status 5 (Released) or 1 (Funded) exists in statusChanges
}

function isReleaseRequested(agreement) {
    return agreement.statusChanges.some(statusChange => statusChange.status === STATUS_SELLER_REQUESTED_RELEASE) && !agreement.statusChanges.some(statusChange => statusChange.status === STATUS_RELEASED) // Check if status 1 (Funded) or 5 (Released) exists in statusChanges
}

const acceptAgreement = async (agreementId) => {
    opRunning.value.push(agreementId)
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
    
    const tx = await contract.sellerAcceptAgreement(agreementId)
    console.log("tansaction: " + tx)
  } catch(err) {
    console.log(err)
  }
  opRunning.value = opRunning.value.filter(id => id !== agreementId)
}

const fundAgreement = async (agreementId, amount) => {
    opRunning.value.push(agreementId)
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
    const val = parseEther(String(amount).substring(0, 18))
    console.log(`Funding agreement with id: ${agreementId} and amount: ${amount}`)
    const tx = await contract.fundAgreement(agreementId, { value: val })
    console.log("tansaction: " + tx)
  } catch(err) {
    console.log(err)
  }
  opRunning.value = opRunning.value.filter(id => id !== agreementId)
}

const requestRelease = async (agreementId) => {
    opRunning.value.push(agreementId)
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
    
    const tx = await contract.requestRelease(agreementId)
    console.log("tansaction: " + tx)
    opSuccess.value.push(agreementId)
  } catch(err) {
    console.log(err)
  }
  opRunning.value = opRunning.value.filter(id => id !== agreementId)
}

const releaseFunds = async (agreementId) => {
    opRunning.value.push(agreementId)
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
    
    const tx = await contract.releaseFunds(agreementId)
    console.log("tansaction: " + tx)
  } catch(err) {
    console.log(err)
  }
    opRunning.value = opRunning.value.filter(id => id !== agreementId)
}

onMounted(() => {
    fetchAgreements()
})

</script>

<template>
    <h1>Your agreements</h1>
    <div class="accordion mb-5" id="agreements">
        <div class="accordion-item" v-for="agreement in agreementsList" >
            <div class="accordion-header">
                <button class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="`#collapse${agreement.agreementId}`" aria-expanded="false" :aria-controls="`collapse${index}`">
                Agreement #{{  agreement.agreementId }} &nbsp; <span :class="statusLabelColor(agreement.status)">{{ statusMap.get(agreement.status) }}</span>
            </button>
            </div>
            <div :id="`collapse${agreement.agreementId}`" class="accordion-collapse collapse" data-bs-parent="#agreements">
            <div class="accordion-body">
                <p>Seller: {{ agreement.seller }}</p>
                <p>Buyer: {{ agreement.buyer }}</p>
                <p>Amount: {{ agreement.amount }}</p>
                <div class="d-flex gap-2">
                    <button type="button" class="btn btn-primary" v-if="isSeller(agreement) && !isAgreementAccepted(agreement)" :disabled="isOpRunning(agreement.agreementId)" @click="acceptAgreement(agreement.agreementId)">Accept agreement</button>
                    <button type="button" class="btn btn-primary" v-if="isBuyer(agreement) && canBeFunded(agreement)" :disabled="isOpRunning(agreement.agreementId)" @click="fundAgreement(agreement.agreementId, agreement.amount)">Fund agreement</button>
                    <button type="button" class="btn btn-primary" v-if="isSeller(agreement) && canRequestRelease(agreement)" :disabled="isOpRunning(agreement.agreementId)" @click="requestRelease(agreement.agreementId)">Request Release</button>
                    <button type="button" class="btn btn-primary" v-if="isBuyer(agreement) && isReleaseRequested(agreement)" :disabled="isOpRunning(agreement.agreementId)" @click="releaseFunds(agreement.agreementId)">Release funds</button>
                    <div class="spinner-border" role="status" v-if="isOpRunning(agreement.agreementId)">
                    <span class="visually-hidden">Loading...</span>
                    </div>
                </div>
                <div class="d-flex p-3">
                    <div class="alert alert-success alert-dismissible" role="alert" v-if="showSuccessAlert(agreement.agreementId)">
                    <div>Operation completed successfully!</div>
                    <button type="button" class="btn-close" data-bs-dismiss="alert" aria-label="Close" @click="resetSuccess(agreement.agreementId)"></button>
                    </div>
                </div>
                <ul>
                    <li v-for="statusChange in agreement.statusChanges" :key="statusChange.timestamp">
                        time: {{ statusChange.timestamp }} status: <span :class="statusLabelColor(statusChange.status)">{{ statusMap.get(statusChange.status) }}</span>
                    </li>
                </ul>
            </div>
            </div>
        </div>
    </div>
</template>
