/**
 * @file filterStore.js
 * @description Centralized reactive state store managing AI Shield settings,
 * lifetime quarantine statistics, user preferences, and transparency drawers.
 */

import { reactive, ref, computed } from 'vue';

// Reactive state container
const state = reactive({
  /** @type {import('../models').FilterSettings} */
  settings: {
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
  },

  /** @type {Array<any>} Quarantined videos from the latest query or trending load */
  quarantinedVideos: [],

  /** @type {number} Total count of videos currently filtered out */
  quarantinedCount: 0,

  /** @type {number} Percentage of human authentic content on current view */
  humanFeedRatio: 100.0,

  /** @type {object|null} Global lifetime statistics */
  globalStats: null,

  /** UI Dialog/Drawer controls */
  isSettingsOpen: false,
  isQuarantineOpen: false,
  isShieldModalOpen: false,

  /** @type {object|null} Video currently inspected in the AI Shield diagnostic modal */
  inspectedVideo: null,

  /** Loading states */
  isSavingSettings: false,
  isLoadingStats: false,
});

/**
 * Initializes and retrieves active filter preferences from the Go backend.
 */
export async function fetchFilterSettings() {
  try {
    const res = await fetch('/api/filter/settings');
    if (res.ok) {
      const data = await res.json();
      state.settings = data;
    }
  } catch (err) {
    console.error('Failed to fetch filter settings from server:', err);
  }
}

/**
 * Persists updated filter preferences to the Go backend.
 * @param {object} newSettings
 */
export async function saveFilterSettings(newSettings) {
  state.isSavingSettings = true;
  try {
    const res = await fetch('/api/filter/settings', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(newSettings),
    });
    if (res.ok) {
      const data = await res.json();
      state.settings = data;
    }
  } catch (err) {
    console.error('Failed to save filter settings:', err);
  } finally {
    state.isSavingSettings = false;
  }
}

/**
 * Fetches lifetime scan and quarantine metrics from the backend.
 */
export async function fetchGlobalStats() {
  state.isLoadingStats = true;
  try {
    const res = await fetch('/api/filter/stats');
    if (res.ok) {
      state.globalStats = await res.json();
    }
  } catch (err) {
    console.error('Failed to fetch global stats:', err);
  } finally {
    state.isLoadingStats = false;
  }
}

/**
 * Updates quarantined video list and human feed ratio from an API search/trending response.
 * @param {Array<any>} quarantined
 * @param {number} ratio
 */
export function setQuarantineResults(quarantined = [], ratio = 100.0) {
  state.quarantinedVideos = quarantined;
  state.quarantinedCount = quarantined.length;
  state.humanFeedRatio = ratio;
}

/**
 * Adds a channel to the custom blocklist on the backend.
 * @param {string} channelIdentifier - Channel Name, ID, or @handle
 */
export async function blockChannel(channelIdentifier) {
  try {
    const res = await fetch('/api/filter/block-channel', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ channel: channelIdentifier }),
    });
    if (res.ok) {
      await fetchFilterSettings();
    }
  } catch (err) {
    console.error('Failed to block channel:', err);
  }
}

/**
 * Whitelists a creator channel, exempting them from AI filtering.
 * @param {string} channelIdentifier - Channel Name or ID
 */
export async function whitelistChannel(channelIdentifier) {
  try {
    const res = await fetch('/api/filter/whitelist-channel', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ channel: channelIdentifier }),
    });
    if (res.ok) {
      await fetchFilterSettings();
    }
  } catch (err) {
    console.error('Failed to whitelist channel:', err);
  }
}

/**
 * Opens the AI Shield breakdown modal for a specific video.
 * @param {object} video
 */
export function inspectVideoSignals(video) {
  state.inspectedVideo = video;
  state.isShieldModalOpen = true;
}

export function closeShieldModal() {
  state.isShieldModalOpen = false;
  state.inspectedVideo = null;
}

export function toggleSettingsModal(open) {
  state.isSettingsOpen = open !== undefined ? open : !state.isSettingsOpen;
}

export function toggleQuarantineDrawer(open) {
  state.isQuarantineOpen = open !== undefined ? open : !state.isQuarantineOpen;
}

/**
 * Export state and computed properties
 */
export function useFilterStore() {
  return {
    state,
    sensitivityMode: computed(() => state.settings.sensitivityMode),
    thresholdPercent: computed(() => Math.round(state.settings.threshold * 100)),
    quarantinedCount: computed(() => state.quarantinedCount),
    humanFeedRatio: computed(() => state.humanFeedRatio),
    quarantinedVideos: computed(() => state.quarantinedVideos),
    globalStats: computed(() => state.globalStats),
    inspectedVideo: computed(() => state.inspectedVideo),
    isSettingsOpen: computed(() => state.isSettingsOpen),
    isQuarantineOpen: computed(() => state.isQuarantineOpen),
    isShieldModalOpen: computed(() => state.isShieldModalOpen),

    fetchFilterSettings,
    saveFilterSettings,
    fetchGlobalStats,
    setQuarantineResults,
    blockChannel,
    whitelistChannel,
    inspectVideoSignals,
    closeShieldModal,
    toggleSettingsModal,
    toggleQuarantineDrawer,
  };
}
