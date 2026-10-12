<template>
  <div class="home-view">
    <!-- Hero / Category Filter Bar -->
    <div class="hero-strip">
      <div class="hero-content">
        <div class="hero-header">
          <div class="hero-badge">
            <ShieldCheck :size="16" class="text-emerald-400" />
            <span>AI Shield: Active & Filtering</span>
          </div>
          <h1 class="hero-headline">Authentic Human Video Discovery</h1>
          <p class="hero-subtext">
            Browsing YouTube free from synthetic voice clones, formulaic ChatGPT scripts, and automated bot farms.
          </p>
        </div>

        <!-- Topic Chips -->
        <div class="topic-chips-row">
          <button
            v-for="topic in topics"
            :key="topic.id"
            type="button"
            class="topic-chip-btn"
            :class="{ active: activeTopic === topic.id }"
            @click="selectTopic(topic)"
          >
            <component :is="topic.icon" :size="14" />
            <span>{{ topic.label }}</span>
          </button>
        </div>
      </div>
    </div>

    <!-- Feed Content -->
    <main class="feed-main-container">
      <!-- Shield Filter Status Banner -->
      <div class="shield-status-strip">
        <div class="status-left">
          <span class="feed-label">Clean {{ activeTopicLabel }} Feed</span>
          <span class="feed-count">({{ videos.length }} Human-Verified Videos)</span>
        </div>

        <div class="status-right">
          <div class="metric-pill">
            <span class="metric-text">{{ humanRatio }}% Authentic Feed</span>
          </div>

          <button
            v-if="quarantinedVideos.length > 0"
            @click="openQuarantine"
            class="btn btn-danger btn-sm"
          >
            <ShieldAlert :size="14" />
            <span>{{ quarantinedVideos.length }} Quarantined</span>
          </button>
        </div>
      </div>

      <!-- Video Grid -->
      <VideoGrid
        :videos="videos"
        :loading="isLoading"
        @open-quarantine="openQuarantine"
        @open-settings="openSettings"
      />
    </main>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue';
import {
  ShieldCheck,
  ShieldAlert,
  Flame,
  Terminal,
  Cpu,
  Hammer,
  FlaskConical,
  Music,
  Compass,
} from 'lucide-vue-next';
import VideoGrid from '../components/VideoGrid.vue';
import { useFilterStore, setQuarantineResults } from '../stores/filterStore';

const filterStore = useFilterStore();

const topics = [
  { id: 'trending', label: 'Trending', icon: Flame, query: '' },
  { id: 'linux', label: 'Linux & Open Source', icon: Terminal, query: 'linux open source tutorial' },
  { id: 'tech', label: 'Tech & Engineering', icon: Cpu, query: 'computer science engineering' },
  { id: 'craft', label: 'Craftsmanship & DIY', icon: Hammer, query: 'woodworking handcrafted restored' },
  { id: 'science', label: 'Science & Math', icon: FlaskConical, query: 'physics mathematics documentary' },
  { id: 'music', label: 'Live Music & Sessions', icon: Music, query: 'live acoustic performance studio session' },
  { id: 'documentary', label: 'Longform Docs', icon: Compass, query: 'independent documentary deep dive' },
];

const activeTopic = ref('trending');
const videos = ref([]);
const quarantinedVideos = ref([]);
const humanRatio = ref(100);
const isLoading = ref(true);

const activeTopicLabel = computed(() => {
  const t = topics.find((item) => item.id === activeTopic.value);
  return t ? t.label : 'Video';
});

async function loadFeed(topic) {
  isLoading.value = true;
  videos.value = [];
  quarantinedVideos.value = [];

  try {
    let url = '/api/trending';
    if (topic.query) {
      url = `/api/search?q=${encodeURIComponent(topic.query)}`;
    }

    const res = await fetch(url);
    if (res.ok) {
      const data = await res.json();
      videos.value = data.videos || [];
      quarantinedVideos.value = data.quarantined || [];
      humanRatio.value = data.humanFeedRatio || 100;
      setQuarantineResults(quarantinedVideos.value, humanRatio.value);
    }
  } catch (err) {
    console.error('Failed to load feed:', err);
  } finally {
    isLoading.value = false;
  }
}

function selectTopic(topic) {
  activeTopic.value = topic.id;
  loadFeed(topic);
}

function openQuarantine() {
  filterStore.toggleQuarantineDrawer(true);
}

function openSettings() {
  filterStore.toggleSettingsModal(true);
}

onMounted(() => {
  loadFeed(topics[0]);
});
</script>

<style scoped>
.home-view {
  min-height: calc(100vh - 70px);
  padding-bottom: 4rem;
}

.hero-strip {
  background: radial-gradient(circle at 50% 0%, rgba(16, 185, 129, 0.08) 0%, transparent 70%);
  border-bottom: 1px solid var(--border-subtle);
  padding: 2.5rem 1.5rem 1.5rem;
}

.hero-content {
  max-width: 1560px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  text-align: center;
  gap: 1.25rem;
}

.hero-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  background-color: var(--accent-shield-bg);
  border: 1px solid var(--accent-shield-border);
  color: var(--accent-shield);
  padding: 0.25rem 0.75rem;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  margin-bottom: 0.5rem;
}

.hero-headline {
  font-size: 2.25rem;
  font-weight: 800;
  letter-spacing: -0.03em;
  color: #ffffff;
  max-width: 720px;
}

.hero-subtext {
  font-size: 1rem;
  color: var(--text-secondary);
  max-width: 580px;
  line-height: 1.5;
}

.topic-chips-row {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  flex-wrap: wrap;
  margin-top: 0.5rem;
}

.topic-chip-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.4rem;
  padding: 0.45rem 0.9rem;
  border-radius: 9999px;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  color: var(--text-secondary);
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.15s ease;
}
.topic-chip-btn:hover {
  background-color: var(--bg-card-hover);
  color: #ffffff;
  border-color: #3b506e;
}
.topic-chip-btn.active {
  background-color: rgba(16, 185, 129, 0.15);
  border-color: var(--accent-shield);
  color: var(--accent-shield);
}

.feed-main-container {
  max-width: 1560px;
  margin: 0 auto;
  padding: 1.5rem 1.5rem 0;
}

.shield-status-strip {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid rgba(255, 255, 255, 0.05);
}

.status-left {
  display: flex;
  align-items: baseline;
  gap: 0.5rem;
}

.feed-label {
  font-size: 1.15rem;
  font-weight: 700;
  color: #ffffff;
}

.feed-count {
  font-size: 0.85rem;
  color: var(--text-muted);
}

.status-right {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.metric-pill {
  background-color: var(--accent-shield-bg);
  border: 1px solid var(--accent-shield-border);
  color: var(--accent-shield);
  padding: 0.25rem 0.65rem;
  border-radius: 9999px;
  font-size: 0.78rem;
  font-weight: 700;
}

.btn-sm {
  padding: 0.3rem 0.7rem;
  font-size: 0.78rem;
}
</style>
