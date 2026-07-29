---
layout: page
title: Graph Colouring with a Quantum Annealer
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 11
pathway: Experience Research
back_url: "/research/#quantum-computing"
back_label: Quantum Computing projects
question: Can a quantum annealer find valid colourings of a graph?
goals:
  - Encode graph colouring as a binary energy function.
  - Test how graph size and penalty weights affect valid solutions.
method_intro: Give every vertex–colour choice a binary variable and penalize broken colouring rules.
method:
  - Build QUBO terms for one colour per vertex and different colours across each edge.
  - Embed and sample the model on an annealer or simulator.
  - Decode the samples and compare validity and quality with a classical colouring method.
visual:
  - label: Graph
    text: Vertices and conflict edges
  - label: QUBO
    text: Colour variables and penalties
  - label: Sample
    text: Candidate colouring
visual_caption: A valid colouring appears as a low-energy state with no rule penalties.
---

{% include resource-page.liquid %}
