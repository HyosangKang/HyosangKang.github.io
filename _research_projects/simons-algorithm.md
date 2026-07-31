---
layout: page
title: Simon’s Algorithm
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 3
pathway: Experience Research
back_url: "/experience/quantum-computing/"
back_label: Quantum Computing projects
question: How can a quantum algorithm reveal a hidden XOR period?
goals:
  - Implement Simon's circuit and collect its measurement equations.
  - Recover the hidden string with linear algebra over \(\mathbb F_2\).
method_intro: Use interference to sample vectors perpendicular to the unknown period.
method:
  - Construct an oracle satisfying \(f(x)=f(x\oplus s)\).
  - Run the circuit repeatedly and record the measured bit strings.
  - Solve the resulting binary linear system and verify the recovered period.
visual:
  - label: Oracle
    text: Two-to-one hidden period
  - label: Samples
    text: Orthogonal bit equations
  - label: Solve
    text: Hidden XOR string
visual_caption: Each measurement gives one linear restriction on the secret, and enough restrictions determine it.
---

{% include resource-page.liquid %}
