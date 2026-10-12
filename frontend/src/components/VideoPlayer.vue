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

    <!-- Player Control & Mode Bar -->
    <div class="player-toolbar">
      <div class="player-mode-selector">
        <button
          v-if="filterStore.playerMode.value === 'premium'"
          @click="filterStore.setPlayerMode('privacy')"
          class="mode-badge premium-badge"
          title="Click to toggle to Anonymous Privacy Mode (youtube-nocookie.com)"
        >
          <Crown :size="14" class="text-amber-400" />
          <span class="badge-text font-medium">YouTube Premium Mode (Ad-Free Browser Session)</span>
          <span class="mode-switch-hint">Switch to Privacy</span>
        </button>

        <button
          v-else
          @click="filterStore.setPlayerMode('premium')"
          class="mode-badge privacy-badge"
          title="Click to toggle to YouTube Premium Mode (uses active browser Google session)"
        >
          <ShieldCheck :size="14" class="text-emerald-400" />
          <span class="badge-text font-medium">Privacy Mode (Cookie-Free)</span>
          <span class="mode-switch-hint">Switch to Premium</span>
        </button>
      </div>

      <div class="player-actions">
        <button
          @click="openInBrowser"
          class="external-yt-link"
          title="Watch directly on YouTube in your default browser"
        >
          <span>Watch on YouTube</span>
          <ExternalLink :size="13" />
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import { ShieldCheck, Crown, ExternalLink } from 'lucide-vue-next';
import { useFilterStore } from '../stores/filterStore';

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

const filterStore = useFilterStore();

const embedUrl = computed(() => {
  // When in YouTube Premium mode, we use youtube.com/embed instead of youtube-nocookie.
  // In a standard web browser (such as Vivaldi), this transmits your active Google login
  // session cookies, immediately unlocking 100% ad-free playback and Premium perks.
  if (filterStore.playerMode.value === 'premium') {
    return `https://www.youtube.com/embed/${props.videoId}?autoplay=1&rel=0&enablejsapi=1`;
  }
  // Privacy mode uses Google's cookie-stripped domain
  return `https://www.youtube-nocookie.com/embed/${props.videoId}?autoplay=1&rel=0&modestbranding=1&enablejsapi=1`;
});

function openInBrowser() {
  const url = `https://www.youtube.com/watch?v=${props.videoId}`;
  if (typeof window !== 'undefined' && window['go']?.['main']?.['App']?.['OpenInBrowser']) {
    window['go']['main']['App']['OpenInBrowser'](url);
  } else {
    window.open(url, '_blank');
  }
}
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
  gap: 1rem;
  flex-wrap: wrap;
}

.player-mode-selector {
  display: flex;
  align-items: center;
}

.mode-badge {
  display: inline-flex;
  align-items: center;
  gap: 0.45rem;
  padding: 0.3rem 0.65rem;
  border-radius: 9999px;
  border: 1px solid transparent;
  background: transparent;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.76rem;
}

.premium-badge {
  color: #fbbf24;
  background-color: rgba(245, 158, 11, 0.1);
  border-color: rgba(245, 158, 11, 0.3);
}

.premium-badge:hover {
  background-color: rgba(245, 158, 11, 0.2);
  border-color: rgba(245, 158, 11, 0.5);
}

.privacy-badge {
  color: #34d399;
  background-color: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.3);
}

.privacy-badge:hover {
  background-color: rgba(16, 185, 129, 0.2);
  border-color: rgba(16, 185, 129, 0.5);
}

.mode-switch-hint {
  font-size: 0.68rem;
  opacity: 0.6;
  margin-left: 0.35rem;
  padding-left: 0.35rem;
  border-left: 1px solid currentColor;
}

.player-actions {
  display: flex;
  align-items: center;
  gap: 0.75rem;
}

.external-yt-link {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  color: var(--text-muted);
  background: transparent;
  border: none;
  cursor: pointer;
  padding: 0;
  font-size: 0.76rem;
  transition: color 0.15s ease;
}

.external-yt-link:hover {
  color: #38bdf8;
}
</style>
