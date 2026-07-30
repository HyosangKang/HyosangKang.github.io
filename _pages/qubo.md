---
layout: page
title: QUBO
permalink: /topics/qubo/
description: Reduce graph-colouring models to efficient quadratic binary objectives for quantum annealing.
nav: false
pathway: Explore Research
back_url: /topics/
back_label: Explore Research
question: How can a graph-colouring model be reduced to QUBO form with far fewer auxiliary variables?
goals:
  - Explain why compact graph-colouring models create terms above degree two.
  - Reduce repeated symmetric terms together instead of handling every monomial separately.
  - Measure the saving in binary variables and quadratic interactions before quantum annealing.
method_intro: The key is to find homogeneous symmetric parts of the objective and reduce their shared structure as one block.
method:
  - Encode each vertex colour compactly with binary variables and penalize adjacent vertices that receive the same colour.
  - Identify large symmetric groups among the resulting higher-degree terms.
  - Reduce each symmetric group first, then apply ordinary monomial reduction only to the terms left over.
  - Compare auxiliary variables and quadratic terms on random graphs and complete graphs.
visual_caption: Efficient degree reduction makes a graph-colouring QUBO smaller before it is embedded on quantum hardware.
---

{% include resource-page.liquid %}

<section class="resource-section">
  <h2>Read more</h2>
  <ul>
    <li><a href="https://doi.org/10.7468/jksmeb.2024.31.1.57" target="_blank" rel="noopener">Research paper ↗</a></li>
    <li><a href="https://github.com/HyosangKang/symred" target="_blank" rel="noopener">Source code ↗</a></li>
  </ul>
</section>
