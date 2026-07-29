---
layout: page
title: Area of a Polygon
description: Compute the area of a polygon.
permalink: /learn/calculus/polygon-area/
subject: Calculus
subject_slug: calculus
project_order: 4
pathway: Learn Mathematics
back_url: /learn/calculus/
back_label: Calculus projects
question: How can the area of an irregular polygon be computed from the coordinates of its vertices?
goals:
  - Implement a reliable area calculation for simple polygons.
  - Connect coordinate formulas with triangulation, determinants, and integration.
method_intro: Describe the boundary by listing its vertices in order.
method:
  - Draw the polygon and check that the vertices follow one orientation around its boundary.
  - Compute signed areas using coordinate products or divide the polygon into triangles.
  - Test the method on convex and non-convex examples and compare the results.
visual:
  - label: Input
    text: Ordered vertices
  - label: Combine
    text: Signed triangle areas
  - label: Output
    text: Polygon area
visual_caption: An irregular polygon can be measured by combining the oriented areas determined by its edges.
---

{% include resource-page.liquid %}
