---
layout: page
title: Graph-colouring Algorithms for Scheduling
description: A guided research project in Algorithms & Discrete Mathematics.
area: Algorithms & Discrete Mathematics
area_slug: algorithms-discrete-mathematics
project_order: 5
pathway: Experience Research
back_url: "/research/#algorithms-discrete-mathematics"
back_label: Algorithms & Discrete Mathematics projects
question: How many time slots are needed when conflicting activities cannot occur together?
goals:
  - Convert a scheduling problem into graph colouring.
  - Compare exact and approximate colouring algorithms.
method_intro: Make each activity a vertex and join two vertices whenever their times must differ.
method:
  - Build the conflict graph from the scheduling data.
  - Apply greedy, backtracking, or other colouring methods.
  - Measure the number of colours, running time, and effect of the vertex order.
visual:
  - label: Conflicts
    text: Activities joined by edges
  - label: Colouring
    text: Assign non-conflicting colours
  - label: Schedule
    text: Colours become time slots
visual_caption: A proper graph colouring is a schedule in which every conflict is separated.
---

{% include resource-page.liquid %}
