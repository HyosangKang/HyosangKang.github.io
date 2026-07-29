---
layout: page
title: Busy Beaver Computation Using Quantum Algorithms
description: A guided research project in Algorithms & Discrete Mathematics.
area: Algorithms & Discrete Mathematics
area_slug: algorithms-discrete-mathematics
project_order: 8
pathway: Experience Research
back_url: "/research/#algorithms-discrete-mathematics"
back_label: Algorithms & Discrete Mathematics projects
question: Can quantum search ideas help explore small machines with extremely long halting times?
goals:
  - Encode and test small Turing machines in a bounded search.
  - Compare a classical search with a quantum-search formulation.
method_intro: Search a finite set of small machines while treating halting within a chosen limit as the test.
method:
  - Enumerate machine descriptions and simulate each one for a bounded number of steps.
  - Formulate the test as an oracle for amplitude amplification.
  - Compare query counts and verify all candidate record holders classically.
visual:
  - label: Machines
    text: Finite encoded candidates
  - label: Search
    text: Halting test and amplification
  - label: Record
    text: Longest verified run
visual_caption: The bounded experiment is computable, even though the general Busy Beaver problem is not.
---

{% include resource-page.liquid %}
