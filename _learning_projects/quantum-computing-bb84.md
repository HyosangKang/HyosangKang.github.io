---
layout: page
title: "Quantum Key Distribution: BB84"
description: Investigate quantum key distribution through the BB84 protocol.
permalink: /learn/quantum-computing/bb84/
subject: Quantum Computing
subject_slug: quantum-computing
project_order: 3
pathway: Learn Mathematics
back_url: /learn/quantum-computing/
back_label: Quantum Computing projects
question: How can two people create a secret key and detect someone listening to the quantum channel?
goals:
  - Simulate state preparation, measurement, and key sifting in BB84.
  - Measure the errors caused by an intercept-and-resend attack.
method_intro: Send random quantum states in two incompatible measurement bases.
method:
  - Let Alice choose random bits and bases, and let Bob measure in random bases.
  - Keep only positions where their bases agree and compare a sample for errors.
  - Insert an eavesdropper and observe how the quantum bit error rate changes.
visual:
  - label: Alice
    text: Random bits and bases
  - label: Channel
    text: Quantum states and possible Eve
  - label: Bob
    text: Sifted secret key
visual_caption: An eavesdropper must measure the states, and the disturbance appears as errors in the shared key.
---

{% include resource-page.liquid %}
