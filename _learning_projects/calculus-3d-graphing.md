---
layout: page
title: 3D Graphing
description: Draw the surface defined by a two-variable function.
permalink: /learn/calculus/3d-graphing/
subject: Calculus
subject_slug: calculus
project_order: 2
pathway: Learn Mathematics
back_url: /learn/calculus/
back_label: Calculus projects
question: How can a two-variable formula be turned into a readable picture of a surface?
goals:
  - Sample a function of two variables and represent its values as points in three-dimensional space.
  - Render a surface that can be rotated, scaled, and viewed from different directions.
method_intro: Treat the formula as a height above each point of a grid.
method:
  - Evaluate the function on a rectangular mesh of input points.
  - Join neighbouring points to form a wireframe or shaded surface.
  - Apply projection and coordinate transformations to produce the final image.
visual:
  - label: Formula
    text: z = f(x, y)
  - label: Sample
    text: Grid of 3D points
  - label: Draw
    text: Surface image
visual_caption: A surface appears when sampled function values are joined across a two-dimensional grid.
---

{% include resource-page.liquid %}
