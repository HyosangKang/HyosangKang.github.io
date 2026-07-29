---
layout: page
title: QUBO
permalink: /topics/qubo/
description: Explore how discrete optimization problems become quadratic functions of binary variables.
nav: false
pathway: Explore Research
back_url: /topics/
back_label: Explore Research
question: How can a constrained discrete problem be written as a quadratic energy over binary variables?
goals:
  - Formulate objectives and constraints in QUBO form.
  - Connect the model with scheduling, graph colouring, and quantum annealing.
method_intro: Encode every decision as 0 or 1 and make invalid choices increase the objective value.
method:
  - Define the binary variables and write the original objective.
  - Add quadratic penalty terms for each constraint.
  - Inspect the energy landscape and solve small cases with classical or annealing methods.
visual:
  - label: Problem
    text: Choices, costs, and rules
  - label: QUBO
    text: Quadratic binary energy
  - label: Solution
    text: Lowest-energy bit string
visual_caption: A successful QUBO makes feasible, high-quality solutions appear at the bottom of its energy landscape.
---

{% include resource-page.liquid %}
