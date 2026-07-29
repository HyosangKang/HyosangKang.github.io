---
layout: page
title: Grover’s Algorithm
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 4
pathway: Experience Research
back_url: "/research/#quantum-computing"
back_label: Quantum Computing projects
question: How does amplitude amplification find a marked item with fewer queries?
goals:
  - Build Grover's oracle and diffusion steps.
  - Measure how iteration count and noise affect success.
method_intro: Repeatedly flip the marked phase and reflect amplitudes toward the target.
method:
  - Prepare an equal superposition and define an oracle for the marked states.
  - Alternate the oracle with the diffusion operator.
  - Track amplitudes and compare observed success with the ideal iteration count.
visual:
  - label: Search space
    text: Equal amplitudes
  - label: Amplify
    text: Oracle and diffusion
  - label: Measure
    text: Marked result
visual_caption: Each Grover iteration rotates probability toward the marked part of the search space.
---

{% include resource-page.liquid %}
