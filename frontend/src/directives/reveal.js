const reducedMotion = () => window.matchMedia('(prefers-reduced-motion: reduce)').matches;

const revealObserver = typeof window !== 'undefined'
  ? new IntersectionObserver((entries, observer) => {
      entries.forEach((entry) => {
        if (!entry.isIntersecting) return;
        entry.target.classList.add('motion-visible');
        observer.unobserve(entry.target);
      });
    }, { threshold: 0.12, rootMargin: '0px 0px -48px' })
  : null;

export const reveal = {
  mounted(el, binding) {
    const options = binding.value || {};
    el.classList.add('motion-reveal');
    el.dataset.motion = options.from || 'up';
    el.style.setProperty('--motion-delay', `${options.delay || 0}ms`);
    el.style.setProperty('--motion-distance', `${options.distance || 28}px`);

    if (reducedMotion() || !revealObserver) {
      el.classList.add('motion-visible');
      return;
    }

    revealObserver.observe(el);
  },
  unmounted(el) {
    revealObserver?.unobserve(el);
  }
};
