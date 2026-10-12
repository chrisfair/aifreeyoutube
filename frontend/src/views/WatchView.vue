<template>
  <div class="watch-view container">
    <div v-if="isLoading" class="watch-loading animate-pulse">
      <div class="skeleton-player"></div>
      <div class="skeleton-title"></div>
    </div>

    <div v-else-if="!video" class="watch-error">
      <h2>Unable to load video</h2>
      <router-link to="/" class="btn btn-primary mt-4">Return Home</router-link>
    </div>

    <div v-else class="watch-layout">
      <!-- Left / Main Column: Video Player & Metadata -->
      <div class="main-column">
        <!-- Privacy-Enhanced Video Player -->
        <VideoPlayer :video-id="video.id" :title="video.title" />

        <!-- Title & Primary Metrics -->
        <div class="video-info-box">
          <h1 class="video-title">{{ video.title }}</h1>

          <div class="creator-bar">
            <div class="creator-left">
              <div class="channel-avatar-fallback">
                {{ video.channelName ? video.channelName.charAt(0).toUpperCase() : 'Y' }}
              </div>
              <div class="creator-meta">
                <span class="creator-name">{{ video.channelName }}</span>
                <span class="creator-stats">{{ video.views }} • {{ video.published }}</span>
              </div>
            </div>

            <!-- Creator Action Controls -->
            <div class="creator-actions">
              <button @click="openExternal" class="btn btn-secondary btn-sm" title="Watch directly in system web browser">
                <ExternalLink :size="15" />
                <span>Open in Browser</span>
              </button>
              <button @click="whitelistCurrentChannel" class="btn btn-secondary btn-sm">
                <UserCheck :size="15" />
                <span>Whitelist Creator</span>
              </button>
              <button @click="blockCurrentChannel" class="btn btn-danger btn-sm">
                <Ban :size="15" />
                <span>Block Channel</span>
              </button>
            </div>
          </div>

          <!-- Video AI Shield Analysis Card -->
          <div class="ai-shield-card" :class="verdictCardClass">
            <div class="shield-header">
              <div class="shield-title-group">
                <ShieldCheck v-if="video.verdict === 'CLEAN'" :size="20" class="text-emerald-400" />
                <AlertTriangle v-else-if="video.verdict === 'SUSPICIOUS'" :size="20" class="text-amber-400" />
                <ShieldAlert v-else :size="20" class="text-rose-400" />
                <span class="shield-card-title">AI Content Shield Evaluation</span>
              </div>

              <div class="shield-score-box">
                <span class="score-val">{{ Math.round(video.humanScore * 100) }}%</span>
                <span class="score-lbl">Human Authenticity</span>
              </div>
            </div>

            <p class="shield-summary-text">
              <span v-if="video.verdict === 'CLEAN'">
                ✓ Verified authentic human production. No automated script templates, voice synth patterns, or content mill markers detected.
              </span>
              <span v-else-if="video.verdict === 'SUSPICIOUS'">
                ⚠ Borderline markers detected. This video contains some formulaic phrasing or clickbait patterns, but did not exceed the quarantine threshold.
              </span>
              <span v-else>
                ⛔ High risk of synthetic generation. Multiple heuristic markers triggered.
              </span>
            </p>

            <!-- Signals details list if present -->
            <div v-if="video.signals && video.signals.length > 0" class="signals-accordion">
              <div class="signals-accordion-title">Evaluation Details:</div>
              <div
                v-for="(sig, sIdx) in video.signals"
                :key="sIdx"
                class="watch-signal-item"
              >
                <span class="signal-tag">{{ sig.ruleName }}</span>
                <span class="signal-expl">{{ sig.description }}</span>
              </div>
            </div>
          </div>

          <!-- Collapsible Description Box -->
          <div v-if="video.description" class="description-box">
            <div class="desc-header">Description</div>
            <p class="desc-text" :class="{ 'desc-expanded': isDescExpanded }">
              {{ video.description }}
            </p>
            <button
              v-if="video.description.length > 200"
              @click="isDescExpanded = !isDescExpanded"
              class="desc-toggle-btn"
            >
              {{ isDescExpanded ? 'Show less' : 'Show more' }}
            </button>
          </div>
        </div>
      </div>

      <!-- Right Column: Clean Filtered Recommendations Sidebar -->
      <aside class="sidebar-column">
        <div class="sidebar-header">
          <div class="sidebar-title">
            <ShieldCheck :size="16" class="text-emerald-400" />
            <span>Human-Only Related Videos</span>
          </div>
          <span class="sidebar-count">{{ relatedVideos.length }} clean</span>
        </div>

        <!-- Quarantined recommendations notice if any were eliminated -->
        <div v-if="quarantinedRelated.length > 0" class="quarantine-sidebar-notice">
          <ShieldAlert :size="14" class="text-amber-400 flex-shrink-0" />
          <span>{{ quarantinedRelated.length }} AI recommendations blocked from this sidebar.</span>
        </div>

        <!-- Recommendations list -->
        <div class="related-list">
          <div
            v-for="rel in relatedVideos"
            :key="rel.id"
            class="related-card"
          >
            <router-link :to="{ name: 'watch', params: { id: rel.id } }" class="rel-thumb-link">
              <img :src="rel.thumbnail" :alt="rel.title" class="rel-thumb" />
              <span v-if="rel.duration" class="rel-duration">{{ rel.duration }}</span>
            </router-link>

            <div class="rel-meta">
              <router-link :to="{ name: 'watch', params: { id: rel.id } }" class="rel-title-link">
                <h4 class="rel-title" :title="rel.title">{{ rel.title }}</h4>
              </router-link>
              <span class="rel-channel">{{ rel.channelName }}</span>
              <div class="rel-score-row">
                <span class="badge badge-clean">
                  {{ Math.round(rel.humanScore * 100) }}% Human
                </span>
                <span class="rel-views">{{ rel.views }}</span>
              </div>
            </div>
          </div>
        </div>
      </aside>
    </div>
  </div>
</template>

<script setup>
import { ref, watch, onMounted, computed } from 'vue';
import { useRoute } from 'vue-router';
import {
  ShieldCheck,
  ShieldAlert,
  AlertTriangle,
  UserCheck,
  Ban,
  ExternalLink,
} from 'lucide-vue-next';
import VideoPlayer from '../components/VideoPlayer.vue';
import { useFilterStore } from '../stores/filterStore';

const route = useRoute();
const filterStore = useFilterStore();

const video = ref(null);
const relatedVideos = ref([]);
const quarantinedRelated = ref([]);
const isLoading = ref(true);
const isDescExpanded = ref(false);

const verdictCardClass = computed(() => {
  if (video.value?.verdict === 'CLEAN') return 'shield-clean';
  if (video.value?.verdict === 'SUSPICIOUS') return 'shield-suspicious';
  return 'shield-slop';
});

async function loadVideoDetails(id) {
  if (!id) return;
  isLoading.value = true;
  video.value = null;
  relatedVideos.value = [];
  quarantinedRelated.value = [];

  try {
    const res = await fetch(`/api/video?id=${id}`);
    if (res.ok) {
      const data = await res.json();
      video.value = data.video;
      relatedVideos.value = data.relatedVideos || [];
      quarantinedRelated.value = data.quarantinedRelated || [];
    }
  } catch (err) {
    console.error('Failed to load video details:', err);
  } finally {
    isLoading.value = false;
  }
}

function openExternal() {
  if (!video.value?.id) return;
  const targetUrl = `https://www.youtube.com/watch?v=${video.value.id}`;
  if (typeof window !== 'undefined' && window['go']?.['main']?.['App']?.['OpenInBrowser']) {
    window['go']['main']['App']['OpenInBrowser'](targetUrl);
  } else {
    window.open(targetUrl, '_blank');
  }
}

async function whitelistCurrentChannel() {
  if (!video.value) return;
  const target = video.value.channelId || video.value.channelName;
  await filterStore.whitelistChannel(target);
  alert(`Channel "${video.value.channelName}" whitelisted.`);
}

async function blockCurrentChannel() {
  if (!video.value) return;
  const target = video.value.channelId || video.value.channelName;
  await filterStore.blockChannel(target);
  alert(`Channel "${video.value.channelName}" added to blocklist.`);
}

watch(
  () => route.params.id,
  (newId) => {
    if (newId) loadVideoDetails(newId);
  }
);

onMounted(() => {
  if (route.params.id) {
    loadVideoDetails(route.params.id);
  }
});
</script>

<style scoped>
.watch-view {
  padding-top: 1.5rem;
  padding-bottom: 4rem;
  max-width: 1560px;
}

.watch-layout {
  display: grid;
  grid-template-columns: 1fr;
  gap: 2rem;
}

@media (min-width: 1024px) {
  .watch-layout {
    grid-template-columns: 1fr 380px;
  }
}

@media (min-width: 1280px) {
  .watch-layout {
    grid-template-columns: 1fr 420px;
  }
}

.main-column {
  display: flex;
  flex-direction: column;
  gap: 1.25rem;
}

.video-info-box {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.video-title {
  font-size: 1.35rem;
  font-weight: 700;
  line-height: 1.35;
  color: #ffffff;
}

.creator-bar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  flex-wrap: wrap;
  gap: 1rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid var(--border-subtle);
}

.creator-left {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.channel-avatar-fallback {
  width: 44px;
  height: 44px;
  border-radius: 50%;
  background: linear-gradient(135deg, #2563eb, #38bdf8);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1.2rem;
  font-weight: 700;
  color: #ffffff;
}

.creator-meta {
  display: flex;
  flex-direction: column;
}

.creator-name {
  font-size: 1rem;
  font-weight: 600;
  color: #ffffff;
}

.creator-stats {
  font-size: 0.8rem;
  color: var(--text-muted);
}

.creator-actions {
  display: flex;
  gap: 0.5rem;
}

.ai-shield-card {
  border-radius: 0.75rem;
  padding: 1.15rem 1.25rem;
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.ai-shield-card.shield-clean {
  background-color: var(--accent-shield-bg);
  border: 1px solid var(--accent-shield-border);
}
.ai-shield-card.shield-suspicious {
  background-color: var(--accent-warning-bg);
  border: 1px solid rgba(245, 158, 11, 0.35);
}
.ai-shield-card.shield-slop {
  background-color: var(--accent-danger-bg);
  border: 1px solid var(--accent-danger-border);
}

.shield-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.shield-title-group {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.shield-card-title {
  font-size: 0.95rem;
  font-weight: 700;
  color: #ffffff;
}

.shield-score-box {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
}

.score-val {
  font-size: 1.25rem;
  font-weight: 800;
  color: var(--accent-shield);
  line-height: 1;
}

.score-lbl {
  font-size: 0.68rem;
  color: var(--text-muted);
  text-transform: uppercase;
}

.shield-summary-text {
  font-size: 0.85rem;
  color: var(--text-primary);
  line-height: 1.45;
}

.signals-accordion {
  background-color: rgba(0, 0, 0, 0.25);
  border-radius: 0.5rem;
  padding: 0.75rem 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.4rem;
}

.signals-accordion-title {
  font-size: 0.75rem;
  font-weight: 700;
  text-transform: uppercase;
  color: var(--text-muted);
  margin-bottom: 0.2rem;
}

.watch-signal-item {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  font-size: 0.78rem;
}

.signal-tag {
  background-color: #1e293b;
  color: #94a3b8;
  padding: 0.15rem 0.45rem;
  border-radius: 4px;
  font-weight: 600;
}

.signal-expl {
  color: var(--text-secondary);
}

.description-box {
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.75rem;
  padding: 1rem 1.25rem;
}

.desc-header {
  font-size: 0.85rem;
  font-weight: 700;
  color: var(--text-muted);
  text-transform: uppercase;
  margin-bottom: 0.5rem;
}

.desc-text {
  font-size: 0.88rem;
  line-height: 1.55;
  color: var(--text-secondary);
  white-space: pre-line;
  max-height: 100px;
  overflow: hidden;
}

.desc-text.desc-expanded {
  max-height: none;
}

.desc-toggle-btn {
  background: none;
  border: none;
  color: #38bdf8;
  font-size: 0.82rem;
  font-weight: 600;
  cursor: pointer;
  margin-top: 0.5rem;
}

/* Sidebar */
.sidebar-column {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.sidebar-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding-bottom: 0.5rem;
  border-bottom: 1px solid var(--border-subtle);
}

.sidebar-title {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  font-size: 0.95rem;
  font-weight: 700;
  color: #ffffff;
}

.sidebar-count {
  font-size: 0.75rem;
  color: var(--accent-shield);
  font-weight: 600;
}

.quarantine-sidebar-notice {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  background-color: rgba(245, 158, 11, 0.1);
  border: 1px solid rgba(245, 158, 11, 0.3);
  padding: 0.6rem 0.75rem;
  border-radius: 0.5rem;
  font-size: 0.75rem;
  color: #fde68a;
}

.related-list {
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.related-card {
  display: flex;
  gap: 0.75rem;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.6rem;
  padding: 0.6rem;
  transition: all 0.15s ease;
}
.related-card:hover {
  border-color: #3b506e;
  transform: translateY(-2px);
}

.rel-thumb-link {
  position: relative;
  width: 120px;
  height: 68px;
  flex-shrink: 0;
  border-radius: 0.35rem;
  overflow: hidden;
  background-color: #0b0f17;
}

.rel-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.rel-duration {
  position: absolute;
  bottom: 0.2rem;
  right: 0.2rem;
  background-color: rgba(0, 0, 0, 0.85);
  color: #ffffff;
  font-size: 0.68rem;
  font-weight: 600;
  padding: 0.05rem 0.35rem;
  border-radius: 2px;
}

.rel-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.2rem;
  overflow: hidden;
}

.rel-title-link {
  text-decoration: none;
}

.rel-title {
  font-size: 0.82rem;
  font-weight: 600;
  line-height: 1.3;
  color: #ffffff;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.rel-channel {
  font-size: 0.72rem;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.rel-score-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: auto;
}

.rel-views {
  font-size: 0.68rem;
  color: var(--text-muted);
}
</style>
