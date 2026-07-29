---
layout: page
title: Rubik’s Cube as a Turing-complete Computer
description: A guided research project in Algorithms & Discrete Mathematics.
area: Algorithms & Discrete Mathematics
area_slug: algorithms-discrete-mathematics
project_order: 7
pathway: Experience Research
back_url: "/research/#algorithms-discrete-mathematics"
back_label: Algorithms & Discrete Mathematics projects
question: Can sequences of Rubik's Cube moves represent and carry out a general computation?
goals:
  - Encode computational states in cube configurations.
  - Construct move sequences that act as logical or state-transition rules.
method_intro: Treat the cube as a finite state system whose legal moves perform controlled changes.
method:
  - Describe configurations and moves using permutations or group operations.
  - Design small gadgets that store symbols and update them through move sequences.
  - Trace a computation and verify that unintended cube pieces are restored.
visual:
  - label: Encode
    text: Data in cube positions
  - label: Moves
    text: Controlled state changes
  - label: Compute
    text: Simulated machine steps
visual_caption: Carefully designed move sequences can change one encoded state while preserving the rest.
---

{% include resource-page.liquid %}
