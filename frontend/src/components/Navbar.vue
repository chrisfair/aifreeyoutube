<template>
  <header class="navbar-wrapper">
    <div class="navbar-content">
      <!-- Brand Logo -->
      <router-link to="/" class="brand-logo">
        <div class="logo-shield">
          <ShieldCheck :size="24" class="text-emerald-400" />
        </div>
        <div class="brand-text">
          <span class="brand-title">AI-Free <span class="brand-sub">YouTube</span></span>
          <span class="brand-tag">Human Content Shield</span>
        </div>
      </router-link>

      <!-- Search Bar -->
      <div class="search-container">
        <form @submit.prevent="executeSearch" class="search-form">
          <div class="search-input-wrapper">
            <Search :size="18" class="search-icon" />
            <input
              v-model="searchQuery"
              type="text"
              placeholder="Search human-verified videos, tutorials, creators..."
              class="search-input"
            />
            <button
              v-if="searchQuery"
              type="button"
              @click="searchQuery = ''"
              class="clear-btn"
            >
              <X :size="16" />
            </button>
          </div>
          <button type="submit" class="search-submit-btn">
            Search
          </button>
        </form>
      </div>

      <!-- Action Buttons & Shield Status -->
      <div class="navbar-actions">
        <!-- Shield Status Badge Button -->
        <button
          @click="toggleQuarantine"
          class="shield-indicator-btn"
          :class="{ 'has-quarantined': quarantinedCount > 0 }"
          title="Inspect quarantined AI content"
        >
          <div class="shield-pulse-dot"></div>
          <ShieldAlert v-if="quarantinedCount > 0" :size="18" class="text-amber-400" />
          <ShieldCheck v-else :size="18" class="text-emerald-400" />
          <span class="shield-label">
            <strong>{{ quarantinedCount }}</strong> Blocked
          </span>
          <span class="shield-ratio">({{ humanFeedRatio }}% Human)</span>
        </button>

        <!-- Filter Settings Button -->
        <button
          @click="openSettings"
          class="settings-btn"
          title="Configure AI Filter Sensitivity"
        >
          <SlidersHorizontal :size="19" />
          <span class="settings-label">Filter Settings</span>
        </button>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue';
import { useRouter } from 'vue-router';
import {
  ShieldCheck,
  ShieldAlert,
  Search,
  SlidersHorizontal,
  X,
} from 'lucide-vue-next';
import { useFilterStore } from '../stores/filterStore';

const router = useRouter();
const filterStore = useFilterStore();

const searchQuery = ref('');
const quarantinedCount = filterStore.quarantinedCount;
const humanFeedRatio = filterStore.humanFeedRatio;

function executeSearch() {
  const q = searchQuery.value.trim();
  if (q) {
    router.push({ name: 'search', query: { q } });
  }
}

function toggleQuarantine() {
  filterStore.toggleQuarantineDrawer(true);
}

function openSettings() {
  filterStore.toggleSettingsModal(true);
}
</script>

<style scoped>
.navbar-wrapper {
  position: sticky;
  top: 0;
  z-index: 40;
  background-color: rgba(10, 13, 20, 0.92);
  backdrop-filter: blur(12px);
  border-bottom: 1px solid var(--border-subtle);
  padding: 0.75rem 1.5rem;
}

.navbar-content {
  max-width: 1560px;
  margin: 0 auto;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
}

.brand-logo {
  display: flex;
  align-items: center;
  gap: 0.75rem;
  text-decoration: none;
}

.logo-shield {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 10px;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(6, 182, 212, 0.15));
  border: 1px solid rgba(16, 185, 129, 0.4);
  box-shadow: 0 0 15px rgba(16, 185, 129, 0.25);
}

.brand-text {
  display: flex;
  flex-direction: column;
}

.brand-title {
  font-size: 1.15rem;
  font-weight: 800;
  letter-spacing: -0.02em;
  color: #ffffff;
}

.brand-sub {
  color: #38bdf8;
}

.brand-tag {
  font-size: 0.68rem;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  color: var(--accent-shield);
}

.search-container {
  flex: 1;
  max-width: 640px;
}

.search-form {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.search-input-wrapper {
  position: relative;
  flex: 1;
  display: flex;
  align-items: center;
}

.search-icon {
  position: absolute;
  left: 0.9rem;
  color: var(--text-muted);
  pointer-events: none;
}

.search-input {
  width: 100%;
  padding: 0.6rem 2.2rem 0.6rem 2.6rem;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 9999px;
  color: var(--text-primary);
  font-size: 0.9rem;
  outline: none;
  transition: all 0.15s ease;
}

.search-input:focus {
  border-color: #3b82f6;
  background-color: #172233;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.2);
}

.clear-btn {
  position: absolute;
  right: 0.75rem;
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0.2rem;
  border-radius: 50%;
  display: flex;
}
.clear-btn:hover {
  color: var(--text-primary);
}

.search-submit-btn {
  padding: 0.6rem 1.25rem;
  background-color: var(--accent-primary);
  color: #ffffff;
  border: none;
  border-radius: 9999px;
  font-size: 0.875rem;
  font-weight: 600;
  cursor: pointer;
  transition: background-color 0.15s ease;
}
.search-submit-btn:hover {
  background-color: var(--accent-primary-hover);
}

.navbar-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.shield-indicator-btn {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  padding: 0.5rem 0.85rem;
  border-radius: 9999px;
  background-color: var(--accent-shield-bg);
  border: 1px solid var(--accent-shield-border);
  color: var(--text-primary);
  cursor: pointer;
  transition: all 0.15s ease;
  font-size: 0.8rem;
}
.shield-indicator-btn:hover {
  background-color: rgba(16, 185, 129, 0.2);
}

.shield-indicator-btn.has-quarantined {
  background-color: rgba(245, 158, 11, 0.14);
  border-color: rgba(245, 158, 11, 0.35);
}
.shield-indicator-btn.has-quarantined:hover {
  background-color: rgba(245, 158, 11, 0.22);
}

.shield-pulse-dot {
  width: 7px;
  height: 7px;
  border-radius: 50%;
  background-color: var(--accent-shield);
  box-shadow: 0 0 8px var(--accent-shield);
}

.shield-ratio {
  color: var(--text-muted);
  font-size: 0.75rem;
}

.settings-btn {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.5rem 0.85rem;
  border-radius: 0.5rem;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.15s ease;
}
.settings-btn:hover {
  background-color: var(--bg-card-hover);
  color: var(--text-primary);
  border-color: #3b4e6b;
}

@media (max-width: 860px) {
  .brand-tag, .settings-label, .shield-ratio {
    display: none;
  }
}
</style>
