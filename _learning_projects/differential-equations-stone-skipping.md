---
layout: page
title: Stone-skipping Simulation
description: Model and simulate the motion of a skipping stone.
permalink: /learn/differential-equations/stone-skipping/
subject: Differential Equations
subject_slug: differential-equations
project_order: 2
has_game: true
pathway: Learn Mathematics
back_url: /learn/differential-equations/
back_label: Differential Equations projects
question: What launch and impact conditions make a stone skip repeatedly across water?
goals:
  - Build a motion model with gravity in air and fluid forces during water contact.
  - Simulate the trajectory and measure how launch conditions change the number of skips.
method_intro: Integrate one motion model that changes its forces as the stone enters and leaves the water.
method:
  - Represent the stone as a shallow tilted segment with mass, length, position, and velocity.
  - Apply gravity above the surface; during partial immersion, compute lift and drag from speed, submerged length, and the two fluid coefficients.
  - Advance the state in short Euler steps, record each water contact, and compare trajectories for different launches.
visual:
  - label: Launch
    text: Speed, angle, and spin
  - label: Motion
    text: Flight and water impact
  - label: Result
    text: A sequence of skips
visual_caption: The MATLAB force model produces a shallow sequence of skips whose height and spacing decrease as fluid drag removes speed.
paper:
  title: The physics of stone skipping
  authors: Lyderic Bocquet
  meta: arXiv:physics/0210015
  url: https://arxiv.org/abs/physics/0210015
  description: A simplified physical model for water impact, repeated bounces, and the conditions behind a successful throw.
matlab_code:
  title: StoneSkipping.m
  url: /assets/code/StoneSkipping.m
  description: The original Euler-step simulation behind this project, with clarified comments and unchanged equations.
---

{% include resource-page.liquid %}
