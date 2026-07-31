---
layout: page
title: Neural-network Decoders for Surface Codes
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 7
pathway: Experience Research
back_url: "/experience/quantum-computing/"
back_label: Quantum Computing projects
question: Can a neural network infer the right correction from a surface-code syndrome?
goals:
  - Generate error–syndrome data for a surface code.
  - Train a decoder and measure its logical error rate.
method_intro: Learn the pattern from stabilizer measurements to an effective correction class.
method:
  - Simulate physical errors and compute the resulting syndromes.
  - Train a suitable neural network on syndrome and correction pairs.
  - Compare its accuracy, speed, and logical failure rate with a standard decoder.
visual:
  - label: Errors
    text: Noise on physical qubits
  - label: Syndrome
    text: Stabilizer measurements
  - label: Decoder
    text: Predicted correction
visual_caption: The decoder sees only the syndrome and must choose a correction that preserves the logical state.
---

{% include resource-page.liquid %}
