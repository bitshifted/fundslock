<script>
import contractAbi from '@/assets/abi/FundsLock.json'
import { CONTRACT_ADDRESS } from '@/config/common';
import { onMounted } from 'vue';


function getAbi() {
    return contractAbi.abi ? contractAbi.abi  : contractAbi
}

onMounted(() => {
    console.log("Contract address: " + CONTRACT_ADDRESS)
})

</script>

<template>
    <h1>Create agreement</h1>

</template>