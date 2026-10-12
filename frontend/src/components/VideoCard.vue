<template>
  <div class="video-card animate-fade-in" :class="{ 'card-suspicious': video.verdict === 'SUSPICIOUS' }">
    <!-- Thumbnail Container -->
    <router-link :to="{ name: 'watch', params: { id: video.id } }" class="thumbnail-wrapper">
      <img
        :src="video.thumbnail"
        :alt="video.title"
        loading="lazy"
        class="thumbnail-img"
        @error="handleImageError"
      />
      <!-- Duration Badge -->
      <span v-if="video.duration" class="duration-badge">
        {{ video.duration }}
      </span>

      <!-- Shorts Format Tag -->
      <span v-if="video.isShort" class="short-badge">
        Short
      </span>

      <!-- Play Overlay on Hover -->
      <div class="play-overlay">
        <Play :size="32" class="play-icon" />
      </div>
    </router-link>

    <!-- Content Metadata -->
    <div class="card-details">
      <!-- Human Score / AI Shield Badge -->
      <div class="card-shield-row">
        <button
          type="button"
          @click.stop="openDiagnostics"
          class="shield-badge-btn"
          :class="badgeClass"
          :title="'Click to inspect AI Shield diagnostics: ' + (video.signals?.length || 0) + ' signals evaluated'"
        >
          <ShieldCheck v-if="video.verdict === 'CLEAN'" :size="13" />
          <AlertTriangle v-else-if="video.verdict === 'SUSPICIOUS'" :size="13" />
          <ShieldAlert v-else :size="13" />
          <span>{{ Math.round(video.humanScore * 100) }}% Human Score</span>
          <Info :size="12" class="ml-1 opacity-70" />
        </button>

        <span v-if="video.published" class="published-text">
          {{ video.published }}
        </span>
      </div>

      <!-- Video Title -->
      <router-link :to="{ name: 'watch', params: { id: video.id } }" class="title-link">
        <h3 class="video-title" :title="video.title">
          {{ video.title }}
        </h3>
      </router-link>

      <!-- Channel and View Metrics -->
      <div class="channel-info">
        <span class="channel-name" :title="video.channelName">
          {{ video.channelName }}
        </span>
        <span v-if="video.views" class="view-count">
          • {{ video.views }}
        </span>
      </div>

      <!-- Quick Action Menu -->
      <div class="card-actions-row">
        <button
          type="button"
          @click.stop="quickWhitelist"
          class="card-action-btn"
          title="Whitelist this creator"
        >
          <UserCheck :size="13" />
          <span>Whitelist</span>
        </button>
        <button
          type="button"
          @click.stop="quickBlock"
          class="card-action-btn danger"
          title="Block this channel"
        >
          <Ban :size="13" />
          <span>Block</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import {
  Play,
  ShieldCheck,
  ShieldAlert,
  AlertTriangle,
  Info,
  UserCheck,
  Ban,
} from 'lucide-vue-next';
import { useFilterStore } from '../stores/filterStore';

const props = defineProps({
  video: {
    type: Object,
    required: true,
  },
});

const filterStore = useFilterStore();

const badgeClass = computed(() => {
  if (props.video.verdict === 'CLEAN') return 'badge-clean';
  if (props.video.verdict === 'SUSPICIOUS') return 'badge-suspicious';
  return 'badge-slop';
});

function handleImageError(e) {
  // Fallback to YouTube standard HQ thumbnail if highres fails
  e.target.src = `https://i.ytimg.com/vi/${props.video.id}/hqdefault.jpg`;
}

function openDiagnostics() {
  filterStore.inspectVideoSignals(props.video);
}

async function quickWhitelist() {
  const target = props.video.channelId || props.video.channelName;
  if (confirm(`Whitelist all content from "${props.video.channelName}"?`)) {
    await filterStore.whitelistChannel(target);
    alert(`Channel "${props.video.channelName}" has been whitelisted.`);
  }
}

async function quickBlock() {
  const target = props.video.channelId || props.video.channelName;
  if (confirm(`Block all content from "${props.video.channelName}"?`)) {
    await filterStore.blockChannel(target);
    alert(`Channel "${props.video.channelName}" has been added to your blocklist.`);
  }
}
</script>

<style scoped>
.video-card {
  display: flex;
  flex-direction: column;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.75rem;
  overflow: hidden;
  transition: transform 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

.video-card:hover {
  transform: translateY(-3px);
  border-color: #3b506e;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

.video-card.card-suspicious {
  border-color: rgba(245, 158, 11, 0.35);
}

.thumbnail-wrapper {
  position: relative;
  width: 100%;
  aspect-ratio: 16 / 9;
  background-color: #0d121c;
  overflow: hidden;
  display: block;
}

.thumbnail-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  transition: transform 0.3s ease;
}

.video-card:hover .thumbnail-img {
  transform: scale(1.03);
}

.duration-badge {
  position: absolute;
  bottom: 0.5rem;
  right: 0.5rem;
  background-color: rgba(0, 0, 0, 0.85);
  backdrop-filter: blur(4px);
  color: #ffffff;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 0.15rem 0.45rem;
  border-radius: 0.3rem;
}

.short-badge {
  position: absolute;
  top: 0.5rem;
  left: 0.5rem;
  background-color: rgba(244, 63, 94, 0.9);
  color: #ffffff;
  font-size: 0.7rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.05em;
  padding: 0.15rem 0.45rem;
  border-radius: 0.3rem;
}

.play-overlay {
  position: absolute;
  inset: 0;
  background-color: rgba(0, 0, 0, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.video-card:hover .play-overlay {
  opacity: 1;
}

.play-icon {
  color: #ffffff;
  filter: drop-shadow(0 2px 8px rgba(0, 0, 0, 0.5));
}

.card-details {
  padding: 0.85rem 1rem;
  display: flex;
  flex-direction: column;
  flex: 1;
  gap: 0.45rem;
}

.card-shield-row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.1rem;
}

.shield-badge-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  padding: 0.2rem 0.55rem;
  border-radius: 9999px;
  font-size: 0.72rem;
  font-weight: 700;
  cursor: pointer;
  transition: all 0.15s ease;
}
.shield-badge-btn:hover {
  filter: brightness(1.2);
  transform: scale(1.02);
}

.published-text {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.title-link {
  text-decoration: none;
}

.video-title {
  font-size: 0.95rem;
  font-weight: 600;
  line-height: 1.35;
  color: var(--text-primary);
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
  transition: color 0.15s ease;
}

.video-card:hover .video-title {
  color: #38bdf8;
}

.channel-info {
  display: flex;
  align-items: center;
  gap: 0.3rem;
  font-size: 0.8rem;
  color: var(--text-secondary);
}

.channel-name {
  font-weight: 500;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
  max-width: 180px;
}

.view-count {
  color: var(--text-muted);
  white-space: nowrap;
}

.card-actions-row {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-top: auto;
  padding-top: 0.4rem;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
}

.card-action-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  background: none;
  border: none;
  color: var(--text-muted);
  font-size: 0.72rem;
  cursor: pointer;
  padding: 0.2rem 0.4rem;
  border-radius: 0.25rem;
  transition: all 0.15s ease;
}
.card-action-btn:hover {
  color: var(--text-primary);
  background-color: rgba(255, 255, 255, 0.06);
}

.card-action-btn.danger:hover {
  color: var(--accent-danger);
  background-color: var(--accent-danger-bg);
}
</style>
