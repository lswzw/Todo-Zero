import type { Directive } from 'vue'

const ripple: Directive = {
  mounted(el: HTMLElement) {
    el.style.position = 'relative'
    el.style.overflow = 'hidden'

    el.addEventListener('click', (e: MouseEvent) => {
      const rect = el.getBoundingClientRect()
      const size = Math.max(rect.width, rect.height) * 2
      const x = e.clientX - rect.left - size / 2
      const y = e.clientY - rect.top - size / 2

      const ripple = document.createElement('span')
      ripple.style.cssText = `
        position: absolute;
        width: ${size}px;
        height: ${size}px;
        left: ${x}px;
        top: ${y}px;
        border-radius: 50%;
        background: rgba(255, 255, 255, 0.3);
        transform: scale(0);
        pointer-events: none;
        animation: ripple-effect 0.6s ease-out forwards;
      `

      el.appendChild(ripple)
      ripple.addEventListener('animationend', () => ripple.remove())
    })
  },
}

export default ripple
