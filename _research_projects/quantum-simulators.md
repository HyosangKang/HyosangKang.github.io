---
layout: page
title: Quantum Simulators
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 1
pathway: Experience Research
back_url: "/research/#quantum-computing"
back_label: Quantum Computing projects
question: How should a simulator represent quantum systems accurately and efficiently?
goals:
  - Implement states, gates, circuits, measurement, and tests.
  - Compare simulator designs as the number of qubits grows.
method_intro: Build the simulator from the mathematical rules for amplitudes and unitary operations.
method:
  - Represent multi-qubit states with complex vectors or another suitable data structure.
  - Apply gates, controlled operations, and probabilistic measurement.
  - Test known circuits and measure memory, speed, and numerical error.
visual:
  - label: State
    text: Complex amplitudes
  - label: Engine
    text: Gates and measurement
  - label: Study
    text: Correctness and scaling
visual_caption: Every added qubit doubles the state-vector size, making design choices part of the research.
---

{% include resource-page.liquid %}
