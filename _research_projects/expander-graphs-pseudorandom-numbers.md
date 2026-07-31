---
layout: page
title: Expander Graphs for Pseudorandom-number Generation
description: A guided research project in Algorithms & Discrete Mathematics.
area: Algorithms & Discrete Mathematics
area_slug: algorithms-discrete-mathematics
project_order: 1
pathway: Experience Research
back_url: "/experience/algorithms-discrete-mathematics/"
back_label: Algorithms & Discrete Mathematics projects
question: How can a sparse graph produce a sequence that behaves almost randomly?
goals:
  - Build or study an expander graph and its random walks.
  - Measure mixing and use the walk as a pseudorandom generator.
method_intro: Move through a highly connected sparse graph using a small random choice at each step.
method:
  - Construct a regular graph and examine its expansion or spectral gap.
  - Generate sequences from random walks on the graph.
  - Compare their distribution and correlations with truly random and ordinary graph walks.
visual:
  - label: Seed
    text: A few random bits
  - label: Expander
    text: Rapidly mixing graph walk
  - label: Output
    text: Pseudorandom sequence
visual_caption: Strong expansion spreads a local random choice quickly across the whole graph.
---

{% include resource-page.liquid %}
