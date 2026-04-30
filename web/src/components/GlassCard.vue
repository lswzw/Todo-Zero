<template>
  <div :class="glassClasses">
    <slot />
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

type GlassVariant = 'default' | 'dark' | 'colored'
type GlassPadding = 'none' | 'sm' | 'md' | 'lg'

const props = withDefaults(
  defineProps<{
    /** 变体样式 */
    variant?: GlassVariant
    /** 是否可悬停 */
    hoverable?: boolean
    /** 内边距 */
    padding?: GlassPadding
    /** 自定义背景渐变 */
    gradient?: string
  }>(),
  {
    variant: 'default',
    hoverable: false,
    padding: 'md',
  }
)

const glassClasses = computed(() => [
  'glass-card',
  `glass-card--${props.variant}`,
  `glass-card--p-${props.padding}`,
  {
    'glass-card--hoverable': props.hoverable,
  },
])

const gradientStyle = computed(() =>
  props.gradient ? { background: props.gradient } : {}
)
</script>

<style scoped>
.glass-card {
  background: rgba(255, 255, 255, 0.7);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.3);
  border-radius: 16px;
  box-shadow: 0 4px 24px rgba(0, 0, 0, 0.08);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  transform: translateZ(0);
  will-change: transform;
}

/* 变体 */
.glass-card--default {
  background: rgba(255, 255, 255, 0.7);
}

.glass-card--dark {
  background: rgba(255, 255, 255, 0.5);
}

.glass-card--colored {
  background: linear-gradient(135deg, rgba(102, 126, 234, 0.1) 0%, rgba(118, 75, 162, 0.1) 100%);
}

/* 内边距 */
.glass-card--p-none {
  padding: 0;
}

.glass-card--p-sm {
  padding: 12px;
}

.glass-card--p-md {
  padding: 20px;
}

.glass-card--p-lg {
  padding: 32px;
}

/* 悬停效果 */
.glass-card--hoverable {
  cursor: pointer;
  will-change: transform;
}

.glass-card--hoverable:hover {
  transform: translateY(-4px) scale(1.01);
  box-shadow: 0 12px 40px rgba(0, 0, 0, 0.12);
}

.glass-card--hoverable:active {
  transform: translateY(-2px) scale(1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.1);
}
</style>
