<template>
  <div v-if="isOpen" class="modal-backdrop" @click.self="close">
    <div class="modal-dialog animate-fade-in">
      <!-- Modal Header -->
      <div class="modal-header">
        <div class="header-left">
          <div class="settings-avatar">
            <SlidersHorizontal :size="22" class="text-blue-400" />
          </div>
          <div>
            <h3 class="modal-title">AI Filter Engine Settings</h3>
            <p class="modal-subtitle">Customize AI Detection Sensitivity and Rules</p>
          </div>
        </div>
        <button @click="close" class="close-btn" title="Close Settings">
          <X :size="20" />
        </button>
      </div>

      <!-- Settings Body -->
      <div class="modal-body">
        <!-- Preset Sensitivity Tabs -->
        <div class="section-group">
          <label class="group-label">Detection Preset</label>
          <div class="presets-tabs">
            <button
              v-for="mode in presets"
              :key="mode.id"
              type="button"
              class="preset-tab-btn"
              :class="{ active: localSettings.sensitivityMode === mode.id }"
              @click="selectPreset(mode.id)"
            >
              <span class="preset-name">{{ mode.label }}</span>
              <span class="preset-desc">{{ mode.desc }}</span>
            </button>
          </div>
        </div>

        <!-- Custom Threshold Slider -->
        <div class="section-group">
          <div class="slider-header">
            <label class="group-label">Quarantine Threshold</label>
            <span class="threshold-display">{{ Math.round(localSettings.threshold * 100) }}% Cutoff</span>
          </div>
          <p class="field-hint">
            Videos with an AI score higher than this threshold will be quarantined from your search and trending feeds.
          </p>
          <input
            type="range"
            min="0.10"
            max="0.90"
            step="0.05"
            v-model.number="localSettings.threshold"
            @input="onSliderChange"
            class="range-slider"
          />
          <div class="slider-markers">
            <span>Aggressive (35%)</span>
            <span>Balanced (55%)</span>
            <span>Mild (75%)</span>
          </div>
        </div>

        <!-- Filter Rule Categories Toggles -->
        <div class="section-group">
          <label class="group-label">Rule Categories</label>
          <div class="toggle-list">
            <label class="toggle-item">
              <input type="checkbox" v-model="localSettings.filterShorts" class="toggle-checkbox" />
              <div class="toggle-text">
                <span class="toggle-title">Eliminate YouTube Shorts</span>
                <span class="toggle-sub">Blocks vertical Shorts format (overwhelmingly AI Reddit/TTS brainrot)</span>
              </div>
            </label>

            <label class="toggle-item">
              <input type="checkbox" v-model="localSettings.enabledCategories.lexical" class="toggle-checkbox" />
              <div class="toggle-text">
                <span class="toggle-title">ChatGPT & LLM Lexical Heuristics</span>
                <span class="toggle-sub">Filters distinctive LLM prose clichés ('delve into', 'tapestry of', 'game-changer')</span>
              </div>
            </label>

            <label class="toggle-item">
              <input type="checkbox" v-model="localSettings.enabledCategories.formula" class="toggle-checkbox" />
              <div class="toggle-text">
                <span class="toggle-title">Clickbait Alarmism Formulas</span>
                <span class="toggle-sub">Detects automated titles ('Scientists Just Discovered Terrifying', 'Why Nobody Talks About')</span>
              </div>
            </label>

            <label class="toggle-item">
              <input type="checkbox" v-model="localSettings.enabledCategories.channel" class="toggle-checkbox" />
              <div class="toggle-text">
                <span class="toggle-title">Synthetic Channel Handles</span>
                <span class="toggle-sub">Flags handles ending in '_ai', bot hashes, and automated faceless channel patterns</span>
              </div>
            </label>

            <label class="toggle-item">
              <input type="checkbox" v-model="localSettings.enabledCategories.blocklist" class="toggle-checkbox" />
              <div class="toggle-text">
                <span class="toggle-title">Community AI Farm Blocklist</span>
                <span class="toggle-sub">Enforces preloaded blocklist of confirmed mass-production AI mills</span>
              </div>
            </label>
          </div>
        </div>

        <!-- Custom Keyword Blacklist -->
        <div class="section-group">
          <label class="group-label">Custom Keyword Blocklist</label>
          <p class="field-hint">Any video containing these terms in its title or description will be immediately blocked.</p>
          <form @submit.prevent="addKeyword" class="keyword-add-row">
            <input
              type="text"
              v-model="newKeyword"
              placeholder="e.g. 'crypto bot', 'faceless affiliate'"
              class="keyword-input"
            />
            <button type="submit" class="btn btn-secondary">
              <Plus :size="16" />
              <span>Add Term</span>
            </button>
          </form>

          <div v-if="localSettings.customBlockedKeywords?.length" class="chips-cloud">
            <span
              v-for="(kw, idx) in localSettings.customBlockedKeywords"
              :key="idx"
              class="chip-item"
            >
              {{ kw }}
              <button type="button" @click="removeKeyword(idx)" class="chip-del-btn">×</button>
            </span>
          </div>
        </div>

        <!-- Channel Whitelist Management -->
        <div class="section-group">
          <label class="group-label">Channel Whitelist (Exempt From Filter)</label>
          <p class="field-hint">Channels that are permanently trusted and bypass all AI heuristics.</p>
          <form @submit.prevent="addWhitelistChannel" class="keyword-add-row">
            <input
              type="text"
              v-model="newWhitelistedChannel"
              placeholder="Channel Name or ID (e.g. UC...)"
              class="keyword-input"
            />
            <button type="submit" class="btn btn-shield">
              <Plus :size="16" />
              <span>Whitelist</span>
            </button>
          </form>

          <div v-if="localSettings.whitelistedChannels?.length" class="chips-cloud">
            <span
              v-for="(ch, idx) in localSettings.whitelistedChannels"
              :key="idx"
              class="chip-item chip-green"
            >
              {{ ch }}
              <button type="button" @click="removeWhitelistedChannel(idx)" class="chip-del-btn">×</button>
            </span>
          </div>
        </div>
      </div>

      <!-- Modal Footer -->
      <div class="modal-footer">
        <button type="button" @click="resetToDefaults" class="btn btn-secondary">
          <RotateCcw :size="16" />
          <span>Reset Defaults</span>
        </button>
        <button type="button" @click="save" class="btn btn-primary ml-auto" :disabled="isSaving">
          <Save :size="16" />
          <span>{{ isSaving ? 'Saving...' : 'Save & Apply Filter' }}</span>
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import {
  SlidersHorizontal,
  X,
  Plus,
  RotateCcw,
  Save,
} from 'lucide-vue-next';
import { useFilterStore } from '../stores/filterStore';

const filterStore = useFilterStore();
const isOpen = filterStore.isSettingsOpen;
const isSaving = ref(false);

const presets = [
  { id: 'aggressive', label: 'Aggressive', desc: '35% Cutoff (Strict clean)' },
  { id: 'balanced', label: 'Balanced', desc: '55% Cutoff (Default)' },
  { id: 'mild', label: 'Mild', desc: '75% Cutoff (Obvious farms)' },
  { id: 'custom', label: 'Custom', desc: 'Manual Slider' },
];

const localSettings = ref(JSON.parse(JSON.stringify(filterStore.state.settings)));
const newKeyword = ref('');
const newWhitelistedChannel = ref('');

watch(isOpen, (newVal) => {
  if (newVal) {
    localSettings.value = JSON.parse(JSON.stringify(filterStore.state.settings));
  }
});

function selectPreset(presetId) {
  localSettings.value.sensitivityMode = presetId;
  if (presetId === 'aggressive') localSettings.value.threshold = 0.35;
  if (presetId === 'balanced') localSettings.value.threshold = 0.55;
  if (presetId === 'mild') localSettings.value.threshold = 0.75;
}

function onSliderChange() {
  localSettings.value.sensitivityMode = 'custom';
}

function addKeyword() {
  const kw = newKeyword.value.trim();
  if (kw && !localSettings.value.customBlockedKeywords.includes(kw)) {
    localSettings.value.customBlockedKeywords.push(kw);
    newKeyword.value = '';
  }
}

function removeKeyword(idx) {
  localSettings.value.customBlockedKeywords.splice(idx, 1);
}

function addWhitelistChannel() {
  const ch = newWhitelistedChannel.value.trim();
  if (ch && !localSettings.value.whitelistedChannels.includes(ch)) {
    localSettings.value.whitelistedChannels.push(ch);
    newWhitelistedChannel.value = '';
  }
}

function removeWhitelistedChannel(idx) {
  localSettings.value.whitelistedChannels.splice(idx, 1);
}

function resetToDefaults() {
  localSettings.value = {
    sensitivityMode: 'balanced',
    threshold: 0.55,
    filterShorts: true,
    enabledCategories: {
      lexical: true,
      formula: true,
      metadata: true,
      channel: true,
      blocklist: true,
    },
    customBlockedKeywords: [],
    customBlockedChannels: [],
    whitelistedChannels: [],
  };
}

async function save() {
  isSaving.value = true;
  await filterStore.saveFilterSettings(localSettings.value);
  isSaving.value = false;
  close();
  // Reload current window to apply new filter parameters to view
  window.location.reload();
}

function close() {
  filterStore.toggleSettingsModal(false);
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
  max-width: 680px;
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

.settings-avatar {
  width: 44px;
  height: 44px;
  border-radius: 12px;
  background-color: rgba(59, 130, 246, 0.15);
  border: 1px solid rgba(59, 130, 246, 0.35);
  display: flex;
  align-items: center;
  justify-content: center;
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

.modal-body {
  padding: 1.5rem;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
  gap: 1.75rem;
  flex: 1;
}

.section-group {
  display: flex;
  flex-direction: column;
  gap: 0.5rem;
}

.group-label {
  font-size: 0.85rem;
  font-weight: 700;
  text-transform: uppercase;
  letter-spacing: 0.04em;
  color: var(--text-secondary);
}

.field-hint {
  font-size: 0.8rem;
  color: var(--text-muted);
  margin-top: -0.2rem;
}

.presets-tabs {
  display: grid;
  grid-template-columns: repeat(4, 1fr);
  gap: 0.5rem;
}

.preset-tab-btn {
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.5rem;
  padding: 0.6rem 0.5rem;
  text-align: center;
  cursor: pointer;
  transition: all 0.15s ease;
  display: flex;
  flex-direction: column;
  gap: 0.15rem;
}
.preset-tab-btn:hover {
  border-color: #3b506e;
}
.preset-tab-btn.active {
  background-color: rgba(59, 130, 246, 0.15);
  border-color: #3b82f6;
}

.preset-name {
  font-size: 0.85rem;
  font-weight: 700;
  color: #ffffff;
}

.preset-desc {
  font-size: 0.68rem;
  color: var(--text-muted);
}

.slider-header {
  display: flex;
  align-items: center;
  justify-content: space-between;
}

.threshold-display {
  font-size: 0.9rem;
  font-weight: 700;
  color: var(--accent-shield);
}

.range-slider {
  width: 100%;
  accent-color: #3b82f6;
  margin-top: 0.5rem;
}

.slider-markers {
  display: flex;
  justify-content: space-between;
  font-size: 0.72rem;
  color: var(--text-muted);
}

.toggle-list {
  display: flex;
  flex-direction: column;
  gap: 0.65rem;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.6rem;
  padding: 0.75rem 1rem;
}

.toggle-item {
  display: flex;
  align-items: flex-start;
  gap: 0.75rem;
  cursor: pointer;
  user-select: none;
}

.toggle-checkbox {
  margin-top: 0.25rem;
  accent-color: #10b981;
  width: 16px;
  height: 16px;
}

.toggle-text {
  display: flex;
  flex-direction: column;
}

.toggle-title {
  font-size: 0.85rem;
  font-weight: 600;
  color: var(--text-primary);
}

.toggle-sub {
  font-size: 0.75rem;
  color: var(--text-muted);
}

.keyword-add-row {
  display: flex;
  gap: 0.5rem;
  margin-top: 0.25rem;
}

.keyword-input {
  flex: 1;
  background-color: var(--bg-card);
  border: 1px solid var(--border-subtle);
  border-radius: 0.5rem;
  padding: 0.5rem 0.8rem;
  color: #ffffff;
  font-size: 0.85rem;
  outline: none;
}
.keyword-input:focus {
  border-color: #3b82f6;
}

.chips-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 0.4rem;
  margin-top: 0.5rem;
}

.chip-item {
  display: inline-flex;
  align-items: center;
  gap: 0.35rem;
  background-color: rgba(244, 63, 94, 0.12);
  border: 1px solid rgba(244, 63, 94, 0.3);
  color: #f43f5e;
  padding: 0.2rem 0.55rem;
  border-radius: 9999px;
  font-size: 0.75rem;
  font-weight: 600;
}

.chip-item.chip-green {
  background-color: var(--accent-shield-bg);
  border-color: var(--accent-shield-border);
  color: var(--accent-shield);
}

.chip-del-btn {
  background: none;
  border: none;
  color: inherit;
  font-size: 1rem;
  line-height: 1;
  cursor: pointer;
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
