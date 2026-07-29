---
layout: page
title: Epidemic Simulation
description: Model and simulate the spread of an epidemic.
permalink: /learn/differential-equations/epidemic-simulation/
subject: Differential Equations
subject_slug: differential-equations
project_order: 1
has_game: true
pathway: Learn Mathematics
back_url: /learn/differential-equations/
back_label: Differential Equations projects
question: How can the flow between population groups help us predict an epidemic?
goals:
  - Build and modify an SIR model for the spread of disease.
  - Compare simulated curves with real data and test changes in transmission.
method_intro: Turn the movement of people between health groups into a system of equations.
method:
  - Define susceptible, infected, and recovered groups with transmission and recovery rates.
  - Solve the system numerically with Euler's method or a Runge–Kutta method.
  - Change the rates, compare the results with data, and discuss what the model leaves out.
visual:
  - label: Population
    text: Susceptible people
  - label: Model
    text: Infection and recovery flows
  - label: Result
    text: S, I, and R curves
visual_caption: The model follows people as they move from susceptible to infected to recovered.
---

{% include resource-page.liquid %}
