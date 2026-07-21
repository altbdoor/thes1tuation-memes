---
layout: base
title: Home
description: |
  Small archive of the memes I made for the very cool dude,
  [TheS1tuation](https://www.twitch.tv/thes1tuation/).
---

<style>
  .alert:after {
    position: absolute;
    content: ' ';
    bottom: -1rem;
    left: calc(64px / 2 - 0.6rem);
    border-width: 1rem 0.6rem 0;
    border-style: solid;
    border-color: var(--bs-cyan) transparent transparent;
  }

  .partyhat.partyhat--active {
    display: block;
    position: absolute;
    top: -2rem;
    left: -0.6rem;
  }
</style>

<div class="alert alert-info position-relative" markdown="1">

Some of the memes are [NSFW](https://www.merriam-webster.com/dictionary/NSFW){: .alert-link} in nature. Please proceed at your own risk.

imgur album link is shadowbanned, [read more](./imgur-shadowban/){: .alert-link}.

</div>

<div class="position-relative d-inline-block" markdown="1">

![partyhat]({{ site.baseurl }}/assets/partyhat.png){: .partyhat.d-none.z-2 width="64" height="64"}
![ok]({{ site.baseurl }}/assets/ok.png){: .z-3.position-relative width="64" height="64"}

</div>

<script type="module">
  const today = new Date();
  const todayBirthday = today.getMonth() === 6 && today.getDate() === 21;
  const testBirthday = new URL(document.location.href).searchParams.has("bday");
  const partyhat = document.querySelector(".partyhat");
  const partyhatBox = partyhat.closest(".position-relative");

  async function init() {
    if (!todayBirthday && !testBirthday) {
      return;
    }

    partyhat.classList.add("partyhat--active");
    partyhat.classList.remove("d-none");
    partyhatBox.style.cursor = "pointer";

    const { default: confetti } =
      await import("https://cdn.jsdelivr.net/npm/canvas-confetti@1.9.4/+esm");
    const confettiInit = { startVelocity: 30, spread: 360, ticks: 60, zIndex: 0 };

    function randomInRange(min, max) {
      return Math.random() * (max - min) + min;
    }

    partyhatBox.onclick = () => {
      if (partyhat.classList.contains("partyhat--animating")) {
        return;
      }

      partyhat.classList.add("partyhat--animating");

      let intervalCounter = 0;
      const animatingInterval = setInterval(() => {
        confetti({
          ...confettiInit,
          origin: { x: Math.random(), y: Math.random() - 0.2 },
        });

        intervalCounter++;
        if (intervalCounter >= 60) {
          partyhat.classList.remove("partyhat--animating");
          clearInterval(animatingInterval);
        }
      }, 200);
    };
  }
  init();
</script>
