---
layout: page
title: Shor’s Algorithm
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 5
pathway: Experience Research
back_url: "/experience/quantum-computing/"
back_label: Quantum Computing projects
question: How can quantum period finding be used to factor a composite integer?
goals:
  - Derive the reduction from factoring to order finding.
  - Implement or simulate the quantum and classical stages on small cases.
method_intro: Find the period of modular powers and convert it into candidate factors.
method:
  - Choose a base coprime to \(N\) and construct modular exponentiation.
  - Apply phase estimation or the quantum Fourier transform to estimate its order.
  - Use continued fractions and greatest common divisors to recover and verify factors.
visual:
  - label: Number
    text: Composite integer \(N\)
  - label: Period
    text: Quantum order finding
  - label: Factors
    text: Classical recovery
visual_caption: The quantum circuit exposes a hidden period that ordinary arithmetic turns into factors.
---

{% include resource-page.liquid %}
