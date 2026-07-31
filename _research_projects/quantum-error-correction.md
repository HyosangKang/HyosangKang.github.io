---
layout: page
title: Quantum Error Correction
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 9
pathway: Experience Research
back_url: "/experience/quantum-computing/"
back_label: Quantum Computing projects
question: How can a logical qubit be protected without directly measuring its unknown state?
goals:
  - Encode a logical qubit and detect selected physical errors.
  - Apply corrections and measure the code's failure rate.
method_intro: Store logical information across several qubits and measure only error syndromes.
method:
  - Choose a repetition, stabilizer, or small surface code.
  - Build the encoding, error, syndrome, and recovery circuits.
  - Simulate different noise levels and compare physical and logical error rates.
visual:
  - label: Encode
    text: One logical qubit across many
  - label: Detect
    text: Syndrome without reading data
  - label: Recover
    text: Correct the logical state
visual_caption: Redundancy reveals the error pattern while leaving the protected quantum information hidden.
---

{% include resource-page.liquid %}
