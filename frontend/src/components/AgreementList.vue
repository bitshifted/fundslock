<script setup>
import { useAuthStore } from '@/stores/auth.js';
import { onMounted } from 'vue';
import { BACKEND_URL } from '@/config/common.js'

const AGREEMENTS_URL = `${BACKEND_URL}/api/v1/agreements`

const authStore = useAuthStore()
async function fetchAgreements() {
    console.log(`auth token: ${authStore.token}`)
    const response = await fetch(AGREEMENTS_URL,{
        headers: {
            'Authorization': `Bearer ${authStore.token}`
        }
    })
    console.log(response.body)
}

onMounted(() => {
    fetchAgreements()
})

</script>

<template>
    <h1>Your agreements</h1>
</template>