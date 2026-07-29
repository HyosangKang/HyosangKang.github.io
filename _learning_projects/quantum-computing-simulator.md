---
layout: page
title: Quantum Simulator
description: Build a simulator for quantum states, gates, and circuits.
permalink: /learn/quantum-computing/quantum-simulator/
subject: Quantum Computing
subject_slug: quantum-computing
project_order: 1
pathway: Learn Mathematics
back_url: /learn/quantum-computing/
back_label: Quantum Computing projects
question: How can a classical program reproduce quantum states, gates, and measurements?
goals:
  - Represent single and multiple qubits with complex amplitudes.
  - Build gates, circuits, and probabilistic measurement in a simulator.
method_intro: Use vectors for states and matrices for the operations applied to them.
method:
  - Store a normalized state vector and combine qubits with tensor products.
  - Apply unitary gate matrices and sample measurement outcomes from their probabilities.
  - Test small circuits such as Deutsch's or Simon's algorithm against expected results.
visual:
  - label: State
    text: Qubit amplitudes
  - label: Circuit
    text: Quantum gate matrices
  - label: Measure
    text: Classical outcomes
visual_caption: The simulator updates amplitudes through each gate and turns them into measurement probabilities.
---

{% include resource-page.liquid %}
