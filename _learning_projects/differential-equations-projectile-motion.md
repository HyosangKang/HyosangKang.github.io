---
layout: page
title: Projectile Motion
description: Launch an object and predict its path under gravity.
permalink: /learn/differential-equations/projectile-motion/
subject: Differential Equations
subject_slug: differential-equations
project_order: 7
has_game: true
pathway: Learn Mathematics
back_url: /learn/differential-equations/
back_label: Differential Equations projects
question: How do launch speed and direction determine the path and landing point of a thrown object?
goals:
  - Simulate two-dimensional motion from an initial position and velocity.
  - Predict the maximum height, flight time, and landing distance.
method_intro: Separate the motion into horizontal and vertical components.
method:
  - Keep horizontal velocity constant while gravity changes vertical velocity.
  - Update position and velocity in short time steps until the object reaches the ground.
  - Compare the numerical path with the familiar parabolic formula.
visual:
  - label: Launch
    text: Initial position and velocity
  - label: Motion
    text: Horizontal motion and gravity
  - label: Result
    text: Trajectory and landing point
visual_caption: Gravity steadily changes the vertical velocity while the horizontal component remains constant.
---

{% include resource-page.liquid %}
