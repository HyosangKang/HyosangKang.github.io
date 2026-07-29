---
layout: page
title: 3D Visualization
description: Rotate and zoom a three-dimensional object.
permalink: /learn/linear-algebra/3d-visualization/
subject: Linear Algebra
subject_slug: linear-algebra
project_order: 2
pathway: Learn Mathematics
back_url: /learn/linear-algebra/
back_label: Linear Algebra projects
question: How can matrices rotate, resize, move, and display a three-dimensional object?
goals:
  - Build matrix transformations for a 3D model.
  - Create interactive rotation and zoom with a screen projection.
method_intro: Store the object as coordinates and transform every point with matrices.
method:
  - Represent vertices and edges in homogeneous coordinates.
  - Multiply by rotation, scaling, translation, camera, and projection matrices.
  - Render the transformed points and connect controls for rotation and zoom.
visual:
  - label: Object
    text: Vertices and edges
  - label: Matrices
    text: Rotate, zoom, and project
  - label: Screen
    text: Interactive 3D view
visual_caption: Matrix multiplication carries each 3D point from the model to the screen.
---

{% include resource-page.liquid %}
