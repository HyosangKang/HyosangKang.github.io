---
layout: page
title: Ecosystem Population Dynamics
description: Model a food chain of coral, fish, and humans.
permalink: /learn/differential-equations/ecosystem-population-dynamics/
subject: Differential Equations
subject_slug: differential-equations
project_order: 9
has_game: true
pathway: Learn Mathematics
back_url: /learn/differential-equations/
back_label: Differential Equations projects
question: How can interactions among coral, fish, and humans produce changing population cycles?
goals:
  - Build a simple agent-based food-chain model.
  - Track how population counts respond when organisms are added or removed.
method_intro: Give each type of organism a small set of movement, feeding, and survival rules.
method:
  - Let fish depend on coral and let humans interact with the fish population.
  - Update every agent and population count in short time steps.
  - Change the starting populations and compare whether the ecosystem grows, cycles, or collapses.
visual:
  - label: Organisms
    text: Coral, fish, and humans
  - label: Interactions
    text: Feeding and survival rules
  - label: Result
    text: Changing population counts
visual_caption: A population that provides food rises first, followed later by the population that depends on it.
---

{% include resource-page.liquid %}
