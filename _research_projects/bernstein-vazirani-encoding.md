---
layout: page
title: Bernstein–Vazirani Encoding
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 2
pathway: Experience Research
back_url: "/experience/quantum-computing/"
back_label: Quantum Computing projects
question: How can one quantum query recover a hidden binary string?
goals:
  - Derive and implement the Bernstein–Vazirani circuit.
  - Explain how phase encoding and interference reveal the string.
method_intro: Encode the hidden string in an oracle and let interference collect all of its bits.
method:
  - Define the oracle \(f(x)=s\cdot x \pmod 2\).
  - Build the Hadamard–oracle–Hadamard circuit and simulate or run it.
  - Measure the output and test how noise changes the recovered string.
visual:
  - label: Secret
    text: Hidden bit string
  - label: Oracle
    text: Phase-encoded inner product
  - label: Measure
    text: Recover all bits
visual_caption: Quantum interference turns the oracle's phase pattern directly into the hidden string.
---

{% include resource-page.liquid %}
