<script setup>
import { RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue';

import { useAppKitAccount } from "@reown/appkit/vue";
import { onMounted } from 'vue';

const { isConnected, address } = useAppKitAccount();

</script>

<template>
  <div class="d-flex flex-column min-vh-100">
    <NavBar />
    <main class="container">
      <RouterView />
    </main>
  </div>
</template>

<style scoped>

</style>
