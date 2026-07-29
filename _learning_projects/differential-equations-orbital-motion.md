---
layout: page
title: Orbital Motion
description: Explore paths created by an inverse-square gravitational force.
permalink: /learn/differential-equations/orbital-motion/
subject: Differential Equations
subject_slug: differential-equations
project_order: 8
has_game: true
pathway: Learn Mathematics
back_url: /learn/differential-equations/
back_label: Differential Equations projects
question: What initial velocity makes a launched object collide with a planet, enter orbit, or escape?
goals:
  - Simulate motion under a central inverse-square force.
  - Compare collision, orbit, and escape paths produced by different launches.
method_intro: Recompute the gravitational acceleration from the object’s current position at every step.
method:
  - Point the acceleration toward the planet and scale its size by the inverse square of the distance.
  - Update velocity and position numerically while recording the trajectory.
  - Repeat with different initial directions and speeds and classify the resulting paths.
visual:
  - label: State
    text: Position and velocity
  - label: Force
    text: Central gravitational attraction
  - label: Path
    text: Collision, orbit, or escape
visual_caption: A small change in launch velocity can turn a collision path into an orbit or an escape path.
---

{% include resource-page.liquid %}
