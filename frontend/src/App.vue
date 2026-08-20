<script setup>
import { RouterView } from 'vue-router'
import NavBar from './components/NavBar.vue';


</script>

<template>
  <div class="d-flex flex-column min-vh-100">
    <NavBar />
    <main>
      <RouterView />
    </main>
  </div>
</template>

<style scoped>

</style>
