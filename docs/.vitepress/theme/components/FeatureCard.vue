<template>
  <div class="feature-card">
    <div class="feature-icon">
      <Icon :icon="icon" />
    </div>
    <h3 class="feature-title">{{ title }}</h3>
    <p class="feature-description">{{ description }}</p>
  </div>
</template>

<script setup>
import { Icon } from '@iconify/vue'

defineProps({
  icon: {
    type: String,
    required: true,
  },
  title: {
    type: String,
    required: true,
  },
  description: {
    type: String,
    required: true,
  },
})
</script>

<style scoped>
.feature-card {
  padding: 1.5rem;
  border-radius: 0.75rem;
  background: linear-gradient(135deg, rgba(var(--vp-c-bg-soft-rgb), 0.4), rgba(var(--vp-c-bg-soft-rgb), 0.1));
  border: 1px solid var(--vp-c-divider);
  transition: all 0.3s ease;
}

.feature-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  border-color: var(--vp-c-brand-1);
}

.feature-icon {
  width: 3rem;
  height: 3rem;
  margin-bottom: 1rem;
  color: var(--vp-c-brand-1);
  font-size: 2rem;
  display: flex;
  align-items: center;
  justify-content: center;
}

.feature-title {
  margin: 0 0 0.75rem 0;
  font-size: 1.125rem;
  font-weight: 600;
  color: var(--vp-c-text-1);
}

.feature-description {
  margin: 0;
  font-size: 0.9375rem;
  line-height: 1.6;
  color: var(--vp-c-text-2);
}

.dark .feature-card {
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.05), rgba(255, 255, 255, 0.02));
}

.dark .feature-card:hover {
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.3);
}
</style>
