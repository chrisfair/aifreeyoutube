<template>
  <div>
    <!-- Backdrop overlay -->
    <div
      v-if="isOpen"
      class="drawer-backdrop"
      @click="close"
    ></div>

    <!-- Slide-over drawer container -->
    <aside class="drawer-container" :class="{ 'drawer-open': isOpen }">
      <!-- Drawer Header -->
      <div class="drawer-header">
        <div class="header-info">
          <div class="quarantine-icon-box">
            <ShieldAlert :size="22" class="text-rose-400" />
          </div>
          <div>
            <h3 class="drawer-title">Quarantined Content</h3>
            <p class="drawer-subtitle">
              {{ videos.length }} Synthetic / AI Videos Blocked
            </p>
          </div>
        </div>
        <button @click="close" class="close-btn" title="Close Drawer">
          <X :size="20" />
        </button>
      </div>

      <!-- Informational Callout -->
      <div class="transparency-note">
        <Info :size="16" class="text-sky-400 flex-shrink-0" />
        <span>
          These videos exceeded your active AI threshold. We preserve transparency by allowing you to review the exact signals and restore any false positives.
        </span>
      </div>

      <!-- Quarantined Video Items List -->
      <div class="drawer-body">
        <div v-if="videos.length === 0" class="drawer-empty">
          <CheckCircle2 :size="36" class="text-emerald-400" />
          <span class="empty-title">Quarantine Clean</span>
          <span class="empty-sub">No videos have been intercepted on this view.</span>
        </div>

        <div v-else class="quarantined-list">
          <div
            v-for="(video, idx) in videos"
            :key="video.id || idx"
            class="quarantined-card animate-fade-in"
          >
            <!-- Thumbnail and Risk Score Tag -->
            <div class="thumb-box">
              <img :src="video.thumbnail" :alt="video.title" class="q-thumb" />
              <span class="ai-risk-tag">
                {{ Math.round(video.aiScore * 100) }}% AI Risk
              </span>
            </div>

            <!-- Title and Creator -->
            <div class="q-details">
              <h4 class="q-title" :title="video.title">
                {{ video.title }}
              </h4>
              <span class="q-channel">Channel: {{ video.channelName }}</span>

              <!-- Primary Triggered Rules Pill List -->
              <div v-if="video.signals?.length" class="q-signals-preview">
                <span
                  v-for="(sig, sIdx) in video.signals.slice(0, 2)"
                  :key="sIdx"
                  class="q-signal-pill"
                >
                  {{ sig.ruleName }}
                </span>
                <span v-if="video.signals.length > 2" class="q-more-signals">
                  +{{ video.signals.length - 2 }} more
                </span>
              </div>

              <!-- Action Controls -->
              <div class="q-actions">
                <button
                  type="button"
                  @click="inspect(video)"
                  class="q-btn"
                  title="View full AI diagnostics breakdown"
                >
                  <Eye :size="13" />
                  <span>Inspect</span>
                </button>

                <button
                  type="button"
                  @click="whitelist(video)"
                  class="q-btn text-emerald-400"
                  title="Whitelist this channel"
                >
                  <UserCheck :size="13" />
                  <span>Whitelist</span>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </aside>
  </div>
</template>

<script setup>
import {
  ShieldAlert,
  CheckCircle2,
  Info,
  X,
  Eye,
  UserCheck,
} from 'lucide-vue-next';
import { useFilterStore } from '../stores/filterStore';

const filterStore = useFilterStore();
const isOpen = filterStore.isQuarantineOpen;
const videos = filterStore.quarantinedVideos;

function close() {
  filterStore.toggleQuarantineDrawer(false);
}

function inspect(video) {
  filterStore.inspectVideoSignals(video);
}

async function whitelist(video) {
  const target = video.channelId || video.channelName;
  await filterStore.whitelistChannel(target);
  alert(`Channel "${video.channelName}" has been whitelisted.`);
}
</script>

<style scoped>
.drawer-backdrop {
  position: fixed;
  inset: 0;
  z-index: 50;
  background-color: rgba(0, 0, 0, 0.6);
  backdrop-filter: blur(4px);
}

.drawer-container {
  position: fixed;
  top: 0;
  right: 0;
  bottom: 0;
  width: 100%;
  max-width: 440px;
  background-color: #101622;
  border-left: 1px solid var(--border-subtle);
  z-index: 55;
  display: flex;
  flex-direction: column;
  box-shadow: -10px 0 35px rgba(0, 0, 0, 0.7);
  transform: translateX(100%);
  transition: transform 0.3s cubic-bezier(0.16, 1, 0.3, 1);
}

.drawer-container.drawer-open {
  transform: translateX(0);
}

.drawer-header {
  padding: 1.25rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-subtle);
}

.header-info {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.quarantine-icon-box {
  width: 42px;
  height: 42px;
  border-radius: 10px;
  background-color: var(--accent-danger-bg);
  border: 1px solid var(--accent-danger-border);
  display: flex;
  align-items: center;
  justify-content: center;
}

.drawer-title {
  font-size: 1.05rem;
  font-weight: 700;
  color: #ffffff;
}

.drawer-subtitle {
  font-size: 0.78rem;
  color: var(--accent-danger);
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: var(--text-muted);
  cursor: pointer;
  padding: 0.3rem;
  border-radius: 0.35rem;
}
.close-btn:hover {
  color: #ffffff;
  background-color: var(--bg-card);
}

.transparency-note {
  display: flex;
  align-items: flex-start;
  gap: 0.6rem;
  padding: 0.85rem 1.25rem;
  background-color: rgba(56, 189, 248, 0.08);
  border-bottom: 1px solid rgba(56, 189, 248, 0.2);
  font-size: 0.76rem;
  color: #bae6fd;
  line-height: 1.45;
}

.drawer-body {
  padding: 1rem 1.25rem;
  overflow-y: auto;
  flex: 1;
}

.drawer-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 5rem 1rem;
  text-align: center;
  gap: 0.5rem;
}

.empty-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #ffffff;
}

.empty-sub {
  font-size: 0.82rem;
  color: var(--text-muted);
}

.quarantined-list {
  display: flex;
  flex-direction: column;
  gap: 0.85rem;
}

.quarantined-card {
  display: flex;
  gap: 0.85rem;
  background-color: var(--bg-card);
  border: 1px solid rgba(244, 63, 94, 0.25);
  border-radius: 0.65rem;
  padding: 0.75rem;
  transition: border-color 0.15s ease;
}
.quarantined-card:hover {
  border-color: rgba(244, 63, 94, 0.5);
}

.thumb-box {
  position: relative;
  width: 120px;
  height: 68px;
  flex-shrink: 0;
  border-radius: 0.4rem;
  overflow: hidden;
  background-color: #0b0f17;
}

.q-thumb {
  width: 100%;
  height: 100%;
  object-fit: cover;
  filter: grayscale(40%);
}

.ai-risk-tag {
  position: absolute;
  bottom: 0.25rem;
  right: 0.25rem;
  background-color: rgba(244, 63, 94, 0.95);
  color: #ffffff;
  font-size: 0.65rem;
  font-weight: 700;
  padding: 0.1rem 0.35rem;
  border-radius: 3px;
}

.q-details {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 0.25rem;
  overflow: hidden;
}

.q-title {
  font-size: 0.82rem;
  font-weight: 600;
  line-height: 1.3;
  color: #f1f5f9;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.q-channel {
  font-size: 0.72rem;
  color: var(--text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.q-signals-preview {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
  margin-top: 0.15rem;
}

.q-signal-pill {
  font-size: 0.65rem;
  background-color: rgba(244, 63, 94, 0.15);
  color: #fda4af;
  padding: 0.1rem 0.4rem;
  border-radius: 3px;
  white-space: nowrap;
}

.q-more-signals {
  font-size: 0.65rem;
  color: var(--text-muted);
}

.q-actions {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.35rem;
}

.q-btn {
  display: inline-flex;
  align-items: center;
  gap: 0.25rem;
  background: none;
  border: 1px solid var(--border-subtle);
  border-radius: 4px;
  color: var(--text-secondary);
  font-size: 0.7rem;
  padding: 0.2rem 0.45rem;
  cursor: pointer;
  transition: all 0.15s ease;
}
.q-btn:hover {
  background-color: rgba(255, 255, 255, 0.05);
  color: #ffffff;
}
</style>
