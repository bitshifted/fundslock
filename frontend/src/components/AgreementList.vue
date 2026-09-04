<script setup>
import { useAuthStore } from '@/stores/auth.js';
import { onMounted, ref } from 'vue';
import { BACKEND_URL } from '@/config/common.js'

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

onMounted(() => {
    fetchAgreements()
})

</script>

<template>
    <h1>Your agreements</h1>
    <div class="accordion mb-5" id="agreements">
        <div class="accordion-item" v-for="(obj,index) in agreementsList" :key="index" >
            <div class="accordion-header">
                <button class="accordion-button" type="button" data-bs-toggle="collapse" :data-bs-target="`#collapse${index}`" aria-expanded="true" :aria-controls="`collapse${index}`">
                Agreement #{{  index }}
            </button>
            </div>
            <div :id="`collapse${index}`" class="accordion-collapse collapse show" data-bs-parent="#agreements">
            <div class="accordion-body">
                <ul>
                    <li v-for="data  in obj" :key="obj.agreement_id">
                        time: {{ data.timestamp }} amount: {{ data.amount }} status: {{  data.status }}
                    </li>
                </ul>
            </div>
            </div>
        </div>
    </div>
</template>
