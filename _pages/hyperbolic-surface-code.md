---
layout: page
title: Hyperbolic Surface Code
permalink: /explore/hyperbolic-surface-code/
description: Explore quantum codes built from tilings and topology on hyperbolic surfaces.
nav: false
pathway: Explore Research
back_url: /explore/
back_label: Explore Research
question: How can a hyperbolic tiling protect quantum information through topology?
goals:
  - Build a surface code from a finite hyperbolic tiling.
  - Study its stabilizers, logical operators, code rate, and distance.
method_intro: Place qubits and parity checks on a tiling whose non-trivial loops carry logical information.
method:
  - Construct a regular tiling and identify edges to form a closed hyperbolic surface.
  - Build the boundary maps or stabilizer-check matrices from vertices, edges, and faces.
  - Find logical cycles and test error patterns or decoders on the resulting code.
visual:
  - label: Tiling
    text: Hyperbolic faces and edges
  - label: Topology
    text: Checks and non-trivial cycles
  - label: Code
    text: Logical qubits and errors
visual_caption: Local stabilizer checks detect errors, while global cycles store the logical quantum information.
---

{% include resource-page.liquid %}
