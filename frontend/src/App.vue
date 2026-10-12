<template>
  <div class="app-layout">
    <!-- Top Navigation Bar -->
    <Navbar />

    <!-- Main View Content -->
    <div class="content-container">
      <router-view />
    </div>

    <!-- Modals & Drawers -->
    <FilterShieldModal />
    <FilterSettingsModal />
    <QuarantineDrawer />

    <!-- Application Footer -->
    <footer class="app-footer">
      <div class="footer-content">
        <div class="footer-left">
          <div class="footer-logo">
            <ShieldCheck :size="18" class="text-emerald-400" />
            <span class="footer-title">AI-Free YouTube</span>
          </div>
          <p class="footer-tagline">
            Preserving human voice, craft, and authentic intellect on the web.
          </p>
        </div>

        <div class="footer-center">
          <div class="status-indicator">
            <div class="status-dot"></div>
            <span>Backend Engine: Connected & Protecting</span>
          </div>
        </div>

        <div class="footer-right">
          <button @click="openSettings" class="footer-link">
            Filter Rules & Thresholds
          </button>
          <span class="footer-sep">•</span>
          <button @click="openQuarantine" class="footer-link">
            Quarantine Drawer ({{ quarantinedCount }})
          </button>
        </div>
      </div>
    </footer>
  </div>
</template>

<script setup>
import { onMounted, computed } from 'vue';
import { ShieldCheck } from 'lucide-vue-next';
import Navbar from './components/Navbar.vue';
import FilterShieldModal from './components/FilterShieldModal.vue';
import FilterSettingsModal from './components/FilterSettingsModal.vue';
import QuarantineDrawer from './components/QuarantineDrawer.vue';
import { useFilterStore } from './stores/filterStore';

const filterStore = useFilterStore();
const quarantinedCount = computed(() => filterStore.quarantinedCount.value);

function openSettings() {
  filterStore.toggleSettingsModal(true);
}

function openQuarantine() {
  filterStore.toggleQuarantineDrawer(true);
}

onMounted(() => {
  filterStore.fetchFilterSettings();
  filterStore.fetchGlobalStats();
});
</script>

<style scoped>
.app-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.content-container {
  flex: 1;
}

.app-footer {
  border-top: 1px solid var(--border-subtle);
  background-color: #080b12;
  padding: 1.5rem;
  margin-top: auto;
}

.footer-content {
  max-width: 1560px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1.25rem;
  font-size: 0.8rem;
}

.footer-left {
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 0.4rem;
}

.footer-title {
  font-weight: 700;
  color: #ffffff;
}

.footer-tagline {
  color: var(--text-muted);
}

.footer-center {
  display: flex;
  align-items: center;
}

.status-indicator {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.35rem 0.75rem;
  background-color: rgba(16, 185, 129, 0.08);
  border: 1px solid rgba(16, 185, 129, 0.25);
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
  color: var(--accent-shield);
}

.status-dot {
  width: 6px;
  height: 6px;
  border-radius: 50%;
  background-color: var(--accent-shield);
  box-shadow: 0 0 6px var(--accent-shield);
}

.footer-right {
  display: flex;
  align-items: center;
  gap: 0.6rem;
}

.footer-link {
  background: none;
  border: none;
  color: var(--text-secondary);
  font-size: 0.8rem;
  cursor: pointer;
  transition: color 0.15s ease;
}
.footer-link:hover {
  color: #38bdf8;
}

.footer-sep {
  color: var(--text-muted);
}
</style>
