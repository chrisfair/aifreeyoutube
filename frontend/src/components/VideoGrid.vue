<template>
  <div class="video-grid-wrapper">
    <!-- Loading Skeletons -->
    <div v-if="loading" class="video-grid">
      <div v-for="i in 8" :key="i" class="skeleton-card animate-pulse">
        <div class="skeleton-thumb"></div>
        <div class="skeleton-body">
          <div class="skeleton-line line-badge"></div>
          <div class="skeleton-line line-title"></div>
          <div class="skeleton-line line-sub"></div>
        </div>
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="!videos || videos.length === 0" class="empty-state">
      <div class="empty-icon-shield">
        <ShieldCheck :size="48" class="text-emerald-400" />
      </div>
      <h3 class="empty-title">All Content Filtered or No Results</h3>
      <p class="empty-desc">
        No videos matching human authenticity criteria were found. If search results were overly aggressive,
        try adjusting your filter threshold or inspect quarantined content.
      </p>
      <div class="empty-actions">
        <button @click="$emit('open-quarantine')" class="btn btn-secondary">
          <ShieldAlert :size="16" />
          <span>View Quarantined Videos</span>
        </button>
        <button @click="$emit('open-settings')" class="btn btn-primary">
          <SlidersHorizontal :size="16" />
          <span>Adjust Sensitivity</span>
        </button>
      </div>
    </div>

    <!-- Video Grid -->
    <div v-else class="video-grid">
      <VideoCard
        v-for="video in videos"
        :key="video.id"
        :video="video"
      />
    </div>
  </div>
</template>

<script setup>
import { ShieldCheck, ShieldAlert, SlidersHorizontal } from 'lucide-vue-next';
import VideoCard from './VideoCard.vue';

defineProps({
  videos: {
    type: Array,
    default: () => [],
  },
  loading: {
    type: Boolean,
    default: false,
  },
});

defineEmits(['open-quarantine', 'open-settings']);
</script>

<style scoped>
.video-grid-wrapper {
  width: 100%;
}

.video-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 1.5rem;
}

@media (min-width: 1280px) {
  .video-grid {
    grid-template-columns: repeat(4, 1fr);
  }
}

/* Skeletons */
.skeleton-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.75rem;
  overflow: hidden;
  height: 310px;
}

.skeleton-thumb {
  width: 100%;
  aspect-ratio: 16 / 9;
  background-color: #1a2436;
}

.skeleton-body {
  padding: 1rem;
  display: flex;
  flex-direction: column;
  gap: 0.6rem;
}

.skeleton-line {
  background-color: #1f2c42;
  border-radius: 4px;
}

.line-badge {
  width: 35%;
  height: 14px;
}

.line-title {
  width: 90%;
  height: 18px;
}

.line-sub {
  width: 60%;
  height: 12px;
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 5rem 1.5rem;
  max-width: 520px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.empty-icon-shield {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 80px;
  height: 80px;
  border-radius: 20px;
  background: rgba(16, 185, 129, 0.1);
  border: 1px solid rgba(16, 185, 129, 0.25);
  margin-bottom: 1.5rem;
}

.empty-title {
  font-size: 1.35rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
}

.empty-desc {
  font-size: 0.9rem;
  color: var(--text-secondary);
  line-height: 1.6;
  margin-bottom: 1.75rem;
}

.empty-actions {
  display: flex;
  gap: 0.75rem;
  flex-wrap: wrap;
  justify-content: center;
}
</style>
