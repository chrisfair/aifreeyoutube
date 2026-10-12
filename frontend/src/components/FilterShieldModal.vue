<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="close">
    <div class="modal-dialog animate-fade-in">
      <!-- Modal Header -->
      <div class="modal-header">
        <div class="header-left">
          <div class="shield-avatar" :class="verdictColorClass">
            <ShieldCheck v-if="video.verdict === 'CLEAN'" :size="24" />
            <AlertTriangle v-else-if="video.verdict === 'SUSPICIOUS'" :size="24" />
            <ShieldAlert v-else :size="24" />
          </div>
          <div>
            <h3 class="modal-title">AI Shield Diagnostics</h3>
            <p class="modal-subtitle">Algorithmic Authenticity Breakdown</p>
          </div>
        </div>
        <button @click="close" class="close-btn" title="Close Diagnostics">
          <X :size="20" />
        </button>
      </div>

      <!-- Video Reference Strip -->
      <div class="video-ref-strip">
        <img :src="video.thumbnail" :alt="video.title" class="ref-thumb" />
        <div class="ref-meta">
          <span class="ref-title">{{ video.title }}</span>
          <span class="ref-channel">By {{ video.channelName }}</span>
        </div>
      </div>

      <!-- Probability Gauges -->
      <div class="metrics-grid">
        <div class="metric-card human">
          <span class="metric-label">Human Authenticity</span>
          <span class="metric-val">{{ Math.round(video.humanScore * 100) }}%</span>
          <div class="progress-track">
            <div class="progress-fill human-fill" :style="{ width: (video.humanScore * 100) + '%' }"></div>
          </div>
        </div>

        <div class="metric-card ai">
          <span class="metric-label">Synthetic / AI Risk</span>
          <span class="metric-val">{{ Math.round(video.aiScore * 100) }}%</span>
          <div class="progress-track">
            <div class="progress-fill ai-fill" :style="{ width: (video.aiScore * 100) + '%' }"></div>
          </div>
        </div>
      </div>

      <!-- Signals Breakdown List -->
      <div class="signals-section">
        <div class="signals-header">
          <h4 class="section-title">Triggered Heuristic Signals</h4>
          <span class="signal-count-tag">{{ video.signals?.length || 0 }} Signals Found</span>
        </div>

        <!-- If no signals triggered -->
        <div v-if="!video.signals || video.signals.length === 0" class="clean-pass-box">
          <CheckCircle2 :size="20" class="text-emerald-400" />
          <div>
            <strong>100% Clean Pass</strong>
            <p>No synthetic prose, clickbait alarmism templates, or automated bot markers were triggered.</p>
          </div>
        </div>

        <!-- List of signals -->
        <div v-else class="signals-list">
          <div
            v-for="(sig, idx) in video.signals"
            :key="idx"
            class="signal-item"
          >
            <div class="signal-top">
              <span class="signal-name">{{ sig.ruleName }}</span>
              <span class="signal-cat-badge">{{ sig.category }}</span>
              <span class="signal-weight-badge">+{{ Math.round(sig.weight * 100) }}% Impact</span>
            </div>
            <p class="signal-desc">{{ sig.description }}</p>
            <div v-if="sig.matchedText" class="matched-snippet">
              <code>Matched: "{{ sig.matchedText }}"</code>
            </div>
          </div>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="modal-footer">
        <button @click="whitelist" class="btn btn-secondary">
          <UserCheck :size="16" />
          <span>Whitelist Channel</span>
        </button>
        <button @click="block" class="btn btn-danger">
          <Ban :size="16" />
          <span>Block Channel</span>
        </button>
        <button @click="close" class="btn btn-primary ml-auto">
          Done
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue';
import {
  ShieldCheck,
  ShieldAlert,
  AlertTriangle,
  CheckCircle2,
  X,
  UserCheck,
  Ban,
} from 'lucide-vue-next';
import { useFilterStore } from '../stores/filterStore';

const filterStore = useFilterStore();
const isOpen = filterStore.isShieldModalOpen;
const video = filterStore.inspectedVideo;

const verdictColorClass = computed(() => {
  if (video.value?.verdict === 'CLEAN') return 'shield-clean';
  if (video.value?.verdict === 'SUSPICIOUS') return 'shield-suspicious';
  return 'shield-slop';
});

function close() {
  filterStore.closeShieldModal();
}

async function whitelist() {
  if (!video.value) return;
  const target = video.value.channelId || video.value.channelName;
  await filterStore.whitelistChannel(target);
  alert(`Channel "${video.value.channelName}" whitelisted.`);
  close();
}

async function block() {
  if (!video.value) return;
  const target = video.value.channelId || video.value.channelName;
  await filterStore.blockChannel(target);
  alert(`Channel "${video.value.channelName}" added to blocklist.`);
  close();
}
</script>

<style scoped>
.modal-backdrop {
  position: fixed;
  inset: 0;
  z-index: 60;
  background-color: rgba(5, 8, 14, 0.75);
  backdrop-filter: blur(8px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 1.5rem;
}

.modal-dialog {
  background-color: #121926;
  border: 1px solid var(--border-subtle);
  border-radius: 1rem;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  display: flex;
  flex-direction: column;
  box-shadow: 0 20px 48px rgba(0, 0, 0, 0.7);
  overflow: hidden;
}

.modal-header {
  padding: 1.25rem 1.5rem;
  display: flex;
  align-items: center;
  justify-content: space-between;
  border-bottom: 1px solid var(--border-subtle);
}

.header-left {
  display: flex;
  align-items: center;
  gap: 0.85rem;
}

.shield-avatar {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}
.shield-clean {
  background: var(--accent-shield-bg);
  border: 1px solid var(--accent-shield-border);
  color: var(--accent-shield);
}
.shield-suspicious {
  background: var(--accent-warning-bg);
  border: 1px solid rgba(245, 158, 11, 0.3);
  color: var(--accent-warning);
}
.shield-slop {
  background: var(--accent-danger-bg);
  border: 1px solid var(--accent-danger-border);
  color: var(--accent-danger);
}

.modal-title {
  font-size: 1.1rem;
  font-weight: 700;
  color: #ffffff;
}

.modal-subtitle {
  font-size: 0.8rem;
  color: var(--text-muted);
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

.video-ref-strip {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 0.85rem 1.5rem;
  background-color: #0d121c;
  border-bottom: 1px solid var(--border-subtle);
}

.ref-thumb {
  width: 80px;
  height: 45px;
  object-fit: cover;
  border-radius: 0.4rem;
}

.ref-meta {
  flex: 1;
  overflow: hidden;
}

.ref-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-primary);
  display: block;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.ref-channel {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.metrics-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1rem;
  padding: 1.25rem 1.5rem 0.5rem;
}

.metric-card {
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.75rem;
  padding: 0.85rem 1rem;
}

.metric-label {
  font-size: 0.75rem;
  font-weight: 500;
  color: var(--text-muted);
  display: block;
}

.metric-val {
  font-size: 1.6rem;
  font-weight: 800;
  display: block;
  margin: 0.2rem 0 0.5rem;
}
.metric-card.human .metric-val {
  color: var(--accent-shield);
}
.metric-card.ai .metric-val {
  color: var(--accent-danger);
}

.progress-track {
  width: 100%;
  height: 6px;
  background-color: #1f2c42;
  border-radius: 3px;
  overflow: hidden;
}

.progress-fill {
  height: 100%;
  border-radius: 3px;
  transition: width 0.4s ease;
}
.human-fill {
  background-color: var(--accent-shield);
}
.ai-fill {
  background-color: var(--accent-danger);
}

.signals-section {
  padding: 1rem 1.5rem 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.signals-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  margin-bottom: 0.85rem;
}

.section-title {
  font-size: 0.9rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--text-secondary);
}

.signal-count-tag {
  font-size: 0.72rem;
  background-color: var(--bg-card);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  color: var(--text-muted);
}

.clean-pass-box {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  padding: 1rem;
  background-color: var(--accent-shield-bg);
  border: 1px solid var(--accent-shield-border);
  border-radius: 0.6rem;
  font-size: 0.85rem;
  color: var(--text-primary);
}

.signals-list {
  display: flex;
  flex-direction: column;
  gap: 0.75rem;
}

.signal-item {
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.6rem;
  padding: 0.75rem 0.9rem;
}

.signal-top {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  margin-bottom: 0.35rem;
}

.signal-name {
  font-size: 0.82rem;
  font-weight: 700;
  color: #ffffff;
}

.signal-cat-badge {
  font-size: 0.68rem;
  background-color: #243247;
  color: #94a3b8;
  padding: 0.1rem 0.4rem;
  border-radius: 3px;
  text-transform: uppercase;
}

.signal-weight-badge {
  margin-left: auto;
  font-size: 0.72rem;
  font-weight: 700;
  color: var(--accent-warning);
}

.signal-desc {
  font-size: 0.8rem;
  color: var(--text-secondary);
  line-height: 1.4;
}

.matched-snippet {
  margin-top: 0.4rem;
  background-color: #0b0f17;
  padding: 0.3rem 0.5rem;
  border-radius: 4px;
  font-size: 0.75rem;
  color: #38bdf8;
}

.modal-footer {
  padding: 1rem 1.5rem;
  background-color: #0d121c;
  border-top: 1px solid var(--border-subtle);
  display: flex;
  align-items: center;
  gap: 0.75rem;
}
.ml-auto {
  margin-left: auto;
}
</style>
