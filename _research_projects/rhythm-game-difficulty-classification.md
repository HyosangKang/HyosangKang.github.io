---
layout: page
title: Rhythm-game Difficulty Classification with Machine Learning
description: A guided research project in Algorithms & Discrete Mathematics.
area: Algorithms & Discrete Mathematics
area_slug: algorithms-discrete-mathematics
project_order: 6
pathway: Experience Research
back_url: "/research/#algorithms-discrete-mathematics"
back_label: Algorithms & Discrete Mathematics projects
question: Which measurable patterns in a rhythm-game chart determine its difficulty?
goals:
  - Extract mathematical features from rhythm-game charts.
  - Train and evaluate a model that predicts difficulty.
method_intro: Turn note timing and movement patterns into data that a classifier can compare.
method:
  - Measure features such as note density, interval variation, jumps, and repeated patterns.
  - Pair charts with difficulty labels and split them into training and test sets.
  - Train a classifier and inspect both accuracy and the charts it misclassifies.
visual:
  - label: Chart
    text: Notes, timing, and movement
  - label: Features
    text: Density and pattern measures
  - label: Model
    text: Predicted difficulty
visual_caption: The model learns from numerical features rather than seeing the rhythm chart as a player does.
---

{% include resource-page.liquid %}
