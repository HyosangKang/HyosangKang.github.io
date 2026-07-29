---
layout: page
title: Length of a Curve
description: Compute the length of curves.
permalink: /learn/calculus/curve-length/
subject: Calculus
subject_slug: calculus
project_order: 5
pathway: Learn Mathematics
back_url: /learn/calculus/
back_label: Calculus projects
question: How can the length of a curved path be estimated and then expressed as an integral?
goals:
  - Approximate a curve by short line segments and observe convergence as the sampling becomes finer.
  - Connect the numerical approximation with the arc-length formula.
method_intro: Replace the curve by a polygonal path and progressively shorten its segments.
method:
  - Sample points from a graph or parametric curve and add the distances between consecutive points.
  - Repeat with finer samples and compare the estimates.
  - Evaluate the corresponding arc-length integral numerically and discuss the error.
visual:
  - label: Curve
    text: Smooth path
  - label: Approximate
    text: Short line segments
  - label: Measure
    text: Limiting length
visual_caption: The length estimate improves as a polygonal path follows the curve more closely.
---

{% include resource-page.liquid %}
