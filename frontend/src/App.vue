<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'

import NavBar from './views/NavBar.vue'

const route = useRoute()
const showNavBar = computed(() => route.name !== 'auth')
</script>

<template>
  <div class="common-layout" :class="{ 'auth-layout': !showNavBar }">
    <el-container>
      <el-header v-if="showNavBar" class="app-header">
        <NavBar />
      </el-header>
      <el-main class="app-main">
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<style scoped>
.common-layout {
  min-height: 100vh;
}

.auth-layout {
  min-height: 100vh;
}

.app-header {
  padding: 0;
  height: auto;
}

.app-main {
  padding: 0;
}
</style>
