---
layout: page
title: Making CAPTCHA
description: Create and transform a CAPTCHA image.
permalink: /learn/calculus/captcha/
subject: Calculus
subject_slug: calculus
project_order: 3
pathway: Learn Mathematics
back_url: /learn/calculus/
back_label: Calculus projects
question: How can mathematical transformations distort text while keeping it readable to a person?
goals:
  - Represent image positions with coordinates, vectors, and matrices.
  - Create a CAPTCHA image by transforming and recombining pixels.
method_intro: Begin with a clear text image and change the coordinates of its pixels.
method:
  - Store the image as a grid and identify the foreground pixels.
  - Apply rotations, shears, stretches, or other coordinate transformations.
  - Add controlled noise and check that the result remains readable.
visual:
  - label: Start
    text: Clear text image
  - label: Transform
    text: Move and distort pixels
  - label: Result
    text: CAPTCHA image
visual_caption: Coordinate transformations change the geometry of an image without changing the underlying text.
---

{% include resource-page.liquid %}
