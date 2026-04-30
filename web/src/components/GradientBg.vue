<template>
  <div :class="['gradient-bg', `gradient-bg--${variant}`]">
    <div class="gradient-bg__layer" />
    <div class="gradient-bg__layer gradient-bg__layer--secondary" />
    <div v-if="showDecor" class="gradient-bg__decor">
      <div class="gradient-bg__bubble gradient-bg__bubble--1" />
      <div class="gradient-bg__bubble gradient-bg__bubble--2" />
      <div class="gradient-bg__bubble gradient-bg__bubble--3" />
    </div>
    <div class="gradient-bg__content">
      <slot />
    </div>
  </div>
</template>

<script setup lang="ts">
type GradientVariant = 'primary' | 'secondary' | 'dark'

withDefaults(
  defineProps<{
    /** 渐变变体 */
    variant?: GradientVariant
    /** 是否显示装饰气泡 */
    showDecor?: boolean
  }>(),
  {
    variant: 'primary',
    showDecor: false,
  }
)
</script>

<style scoped>
.gradient-bg {
  position: relative;
  min-height: 100vh;
  overflow: hidden;
}

.gradient-bg__layer {
  position: absolute;
  inset: 0;
  z-index: 0;
}

/* === primary 变体 === */
.gradient-bg--primary .gradient-bg__layer {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  opacity: 0.85;
}

.gradient-bg--primary .gradient-bg__layer--secondary {
  background: radial-gradient(
    circle at 20% 80%,
    rgba(120, 80, 200, 0.3) 0%,
    transparent 50%
  );
  animation: float 15s ease-in-out infinite;
}

/* === secondary 变体（浅色） === */
.gradient-bg--secondary .gradient-bg__layer {
  background: linear-gradient(135deg, #f5f7fa 0%, #c3cfe2 100%);
  opacity: 1;
}

.gradient-bg--secondary .gradient-bg__layer--secondary {
  background: radial-gradient(
    circle at 80% 20%,
    rgba(102, 126, 234, 0.08) 0%,
    transparent 50%
  );
  animation: float 20s ease-in-out infinite;
}

/* === dark 变体 === */
.gradient-bg--dark .gradient-bg__layer {
  background: linear-gradient(135deg, #1a1a2e 0%, #16213e 100%);
  opacity: 1;
}

.gradient-bg--dark .gradient-bg__layer--secondary {
  background: radial-gradient(
    circle at 30% 70%,
    rgba(102, 126, 234, 0.15) 0%,
    transparent 50%
  );
  animation: float 18s ease-in-out infinite;
}

/* === 装饰气泡 === */
.gradient-bg__decor {
  position: absolute;
  inset: 0;
  z-index: 1;
  pointer-events: none;
}

.gradient-bg__bubble {
  position: absolute;
  border-radius: 50%;
  opacity: 0.12;
}

.gradient-bg__bubble--1 {
  width: 300px;
  height: 300px;
  background: #667eea;
  top: 10%;
  left: -5%;
  animation: bubble-float 12s ease-in-out infinite;
}

.gradient-bg__bubble--2 {
  width: 200px;
  height: 200px;
  background: #764ba2;
  top: 60%;
  right: -3%;
  animation: bubble-float 10s ease-in-out infinite 3s;
}

.gradient-bg__bubble--3 {
  width: 150px;
  height: 150px;
  background: #f093fb;
  bottom: 10%;
  left: 30%;
  animation: bubble-float 14s ease-in-out infinite 6s;
}

.gradient-bg--secondary .gradient-bg__bubble {
  opacity: 0.06;
}

.gradient-bg--dark .gradient-bg__bubble {
  opacity: 0.08;
}

/* === 内容层 === */
.gradient-bg__content {
  position: relative;
  z-index: 2;
  min-height: 100vh;
}

/* === 动画 === */
@keyframes float {
  0%,
  100% {
    transform: translateX(0) translateY(0);
  }
  50% {
    transform: translateX(5%) translateY(-5%);
  }
}

@keyframes bubble-float {
  0%,
  100% {
    transform: translateY(0) scale(1);
  }
  33% {
    transform: translateY(-20px) scale(1.05);
  }
  66% {
    transform: translateY(10px) scale(0.95);
  }
}
</style>
