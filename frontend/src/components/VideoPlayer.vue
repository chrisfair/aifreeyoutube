<template>
  <div class="player-container">
    <div class="iframe-wrapper">
      <iframe
        :src="embedUrl"
        :title="title"
        frameborder="0"
        allow="accelerometer; autoplay; clipboard-write; encrypted-media; gyroscope; picture-in-picture; web-share"
        allowfullscreen
        class="player-iframe"
      ></iframe>
    </div>

    <!-- Player Control Bar -->
    <div class="player-toolbar">
      <div class="privacy-badge">
        <ShieldCheck :size="14" class="text-emerald-400" />
        <span>Privacy-Enhanced Playback (Zero Google Ad Tracking & Algorithmic Clutter)</span>
      </div>

      <a
        :href="'https://youtu.be/' + videoId"
        target="_blank"
        rel="noopener noreferrer"
        class="external-yt-link"
        title="Open directly on YouTube"
      >
        <span>Open on YouTube</span>
        <ExternalLink :size="13" />
      </a>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { ShieldCheck, ExternalLink } from 'lucide-vue-next';

const props = defineProps({
  videoId: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    default: 'YouTube Video Player',
  },
});

const embedUrl = computed(() => {
  return `https://www.youtube-nocookie.com/embed/${props.videoId}?autoplay=1&rel=0&modestbranding=1&enablejsapi=1`;
});
</script>

<style scoped>
.player-container {
  width: 100%;
  border-radius: 0.75rem;
  overflow: hidden;
  background-color: #000000;
  border: 1px solid var(--border-subtle);
  box-shadow: 0 12px 32px rgba(0, 0, 0, 0.6);
}

.iframe-wrapper {
  position: relative;
  width: 100%;
  padding-bottom: 56.25%; /* 16:9 aspect ratio */
  height: 0;
}

.player-iframe {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  border: none;
}

.player-toolbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0.6rem 1rem;
  background-color: var(--bg-card);
  border-top: 1px solid var(--border-subtle);
  font-size: 0.78rem;
}

.privacy-badge {
  display: flex;
  align-items: center;
  gap: 0.4rem;
  color: var(--text-secondary);
}

.external-yt-link {
  display: flex;
  align-items: center;
  gap: 0.35rem;
  color: var(--text-muted);
  text-decoration: none;
  transition: color 0.15s ease;
}
.external-yt-link:hover {
  color: #38bdf8;
}
</style>
