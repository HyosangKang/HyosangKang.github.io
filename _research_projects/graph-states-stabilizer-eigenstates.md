---
layout: page
title: Graph States and Stabilizer Eigenstates
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 8
pathway: Experience Research
back_url: "/research/#quantum-computing"
back_label: Quantum Computing projects
question: How does a graph describe an entangled state and its stabilizer equations?
goals:
  - Construct graph states from vertices and edges.
  - Derive and verify their stabilizer eigenvalue relations.
method_intro: Associate one qubit with each vertex and an entangling gate with each edge.
method:
  - Prepare all qubits in \(|+\rangle\) and apply controlled-Z gates along the graph edges.
  - Build the stabilizer generator for every vertex.
  - Verify the state algebraically or in a simulator and explore graph transformations.
visual:
  - label: Graph
    text: Vertices and edges
  - label: Circuit
    text: Qubits and CZ gates
  - label: State
    text: Stabilizer eigenstate
visual_caption: The neighbourhood of each vertex determines one stabilizer that fixes the whole graph state.
---

{% include resource-page.liquid %}
