---
layout: page
title: Quantum Error-correction Game
description: A guided research project in Quantum Computing.
area: Quantum Computing
area_slug: quantum-computing
project_order: 10
pathway: Experience Research
back_url: "/research/#quantum-computing"
back_label: Quantum Computing projects
question: Can game play make quantum error detection and correction easier to understand?
goals:
  - Turn syndromes, errors, and corrections into clear game actions.
  - Build a prototype that remains faithful to the underlying code.
method_intro: Let players diagnose hidden errors from the same limited clues available to a decoder.
method:
  - Choose a small code and map qubits, stabilizers, and errors to visual game elements.
  - Design levels that introduce one idea at a time and give useful feedback.
  - Test whether player strategies match valid decoding and revise confusing rules.
visual:
  - label: Hidden event
    text: Error on encoded qubits
  - label: Clues
    text: Visible syndrome pattern
  - label: Player
    text: Choose a correction
visual_caption: The player never sees the error directly and must reason from the syndrome.
---

{% include resource-page.liquid %}
