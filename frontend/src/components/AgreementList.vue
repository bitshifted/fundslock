<script setup>
import { useAuthStore } from '@/stores/auth.js';
import { onMounted, ref } from 'vue';
import { BACKEND_URL } from '@/config/common.js'
import {statusMap} from '@/assets/abi/enums.js'

const AGREEMENTS_URL = `${BACKEND_URL}/api/v1/agreements`

const authStore = useAuthStore()
const agreementsList = ref([])

async function fetchAgreements() {
    console.log(`auth token: ${authStore.token}`)
    const response = await fetch(AGREEMENTS_URL,{
        headers: {
            'Authorization': `Bearer ${authStore.token}`
        }
    })
    agreementsList.value = await response.json()
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
