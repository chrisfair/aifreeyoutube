<template>
  <div class="search-view container">
    <!-- Search Query Banner -->
    <div class="search-header">
      <div class="header-left">
        <h2 class="search-title">
          Results for <span class="query-highlight">"{{ currentQuery }}"</span>
        </h2>
        <p class="search-subtext">
          Displaying human-verified content. Synthetic content mills and automated AI scripts have been quarantined.
        </p>
      </div>

      <div class="header-metrics">
        <div class="stat-bubble">
          <span class="bubble-val">{{ videos.length }}</span>
          <span class="bubble-lbl">Human Videos</span>
        </div>

        <button
          v-if="quarantinedVideos.length > 0"
          @click="openQuarantine"
          class="stat-bubble danger cursor-pointer"
          title="Inspect quarantined AI videos for this query"
        >
          <span class="bubble-val text-rose-400">{{ quarantinedVideos.length }}</span>
          <span class="bubble-lbl">AI Filtered</span>
        </button>
      </div>
    </div>

    <!-- Search Video Grid -->
    <VideoGrid
      :videos="videos"
      :loading="isLoading"
      @open-quarantine="openQuarantine"
      @open-settings="openSettings"
    />
  </div>
</template>

<script setup>
import { ref, watch, onMounted } from 'vue';
import { useRoute } from 'vue-router';
import VideoGrid from '../components/VideoGrid.vue';
import { useFilterStore, setQuarantineResults } from '../stores/filterStore';

const route = useRoute();
const filterStore = useFilterStore();

const currentQuery = ref('');
const videos = ref([]);
const quarantinedVideos = ref([]);
const isLoading = ref(true);

async function performSearch(q) {
  if (!q) return;
  currentQuery.value = q;
  isLoading.value = true;
  videos.value = [];
  quarantinedVideos.value = [];

  try {
    const res = await fetch(`/api/search?q=${encodeURIComponent(q)}`);
    if (res.ok) {
      const data = await res.json();
      videos.value = data.videos || [];
      quarantinedVideos.value = data.quarantined || [];
      setQuarantineResults(quarantinedVideos.value, data.humanFeedRatio || 100);
    }
  } catch (err) {
    console.error('Search request failed:', err);
  } finally {
    isLoading.value = false;
  }
}

function openQuarantine() {
  filterStore.toggleQuarantineDrawer(true);
}

function openSettings() {
  filterStore.toggleSettingsModal(true);
}

watch(
  () => route.query.q,
  (newQ) => {
    if (newQ) performSearch(newQ);
  }
);

onMounted(() => {
  if (route.query.q) {
    performSearch(route.query.q);
  }
});
</script>

<style scoped>
.search-view {
  padding-top: 2rem;
  padding-bottom: 4rem;
  max-width: 1560px;
}

.search-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 1.5rem;
  padding-bottom: 1.5rem;
  margin-bottom: 2rem;
  border-bottom: 1px solid var(--border-subtle);
}

.search-title {
  font-size: 1.6rem;
  font-weight: 800;
  color: #ffffff;
}

.query-highlight {
  color: #38bdf8;
}

.search-subtext {
  font-size: 0.85rem;
  color: var(--text-muted);
  margin-top: 0.25rem;
}

.header-metrics {
  display: flex;
  gap: 0.75rem;
}

.stat-bubble {
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.6rem;
  padding: 0.5rem 0.9rem;
  text-align: center;
  display: flex;
  flex-direction: column;
  min-width: 90px;
}

.stat-bubble.danger {
  border-color: rgba(244, 63, 94, 0.35);
  background-color: var(--accent-danger-bg);
}

.cursor-pointer {
  cursor: pointer;
  transition: transform 0.15s ease;
}
.cursor-pointer:hover {
  transform: translateY(-2px);
}

.bubble-val {
  font-size: 1.25rem;
  font-weight: 800;
  color: #ffffff;
  line-height: 1.1;
}

.bubble-lbl {
  font-size: 0.7rem;
  color: var(--text-muted);
  text-transform: uppercase;
  font-weight: 600;
  letter-spacing: 0.04em;
}
</style>
