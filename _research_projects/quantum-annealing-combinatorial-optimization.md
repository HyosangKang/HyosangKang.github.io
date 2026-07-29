---
layout: page
title: Quantum Annealing for Combinatorial Optimization
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 6
pathway: Experience Research
back_url: "/research/#quantum-computing"
back_label: Quantum Computing projects
question: How can a discrete optimization problem be encoded as the lowest energy of a quantum annealer?
goals:
  - Convert a chosen combinatorial problem into an Ising or QUBO model.
  - Compare annealing results with a classical solution method.
method_intro: Make low-energy binary states represent good solutions and high-energy states represent violations.
method:
  - Define binary variables, the objective, and penalty terms.
  - Map the model to available qubits and choose penalty or chain strengths.
  - Sample solutions and compare validity, quality, and running time.
visual:
  - label: Problem
    text: Discrete choices and rules
  - label: Energy
    text: Ising or QUBO model
  - label: Samples
    text: Candidate solutions
visual_caption: Optimization becomes a search for the binary state with the smallest designed energy.
---

{% include resource-page.liquid %}
