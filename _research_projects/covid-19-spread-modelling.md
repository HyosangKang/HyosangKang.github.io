---
layout: page
title: COVID-19 Spread Modelling with Differential Equations
description: A guided research project in Mathematical Modelling.
area: Mathematical Modelling
area_slug: mathematical-modelling
project_order: 1
pathway: Experience Research
back_url: "/research/#mathematical-modelling"
back_label: Mathematical Modelling projects
question: How well can a modified compartment model explain real COVID-19 case data?
goals:
  - Build an SIR-based model suited to the chosen data.
  - Estimate parameters and compare policy or transmission scenarios.
method_intro: Represent disease spread as flows between population groups and fit those flows to observations.
method:
  - Select compartments, rates, initial values, and a reliable time series.
  - Solve the differential equations and estimate parameters from the data.
  - Compare simulated curves with observations and test alternative assumptions.
visual:
  - label: Data
    text: Cases over time
  - label: Model
    text: Modified SIR equations
  - label: Study
    text: Fit and scenario curves
visual_caption: The gap between simulated and observed curves shows both useful structure and model limits.
---

{% include resource-page.liquid %}
