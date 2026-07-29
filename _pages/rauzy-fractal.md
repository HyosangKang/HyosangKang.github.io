---
layout: page
title: Rauzy Fractal
permalink: /topics/rauzy-fractal/
description: Explore how substitutions and eigenvectors produce self-similar geometric fractals.
nav: false
pathway: Explore Research
back_url: /topics/
back_label: Explore Research
question: How does a symbolic substitution create a self-similar geometric fractal?
goals:
  - Construct a Rauzy fractal from a Pisot substitution.
  - Study its projection, self-similarity, and two- or three-dimensional form.
method_intro: Turn an expanding symbolic sequence into a lattice path, then project away its main growth direction.
method:
  - Choose a substitution and compute its substitution matrix and eigenvectors.
  - Iterate the substitution and record the prefix-count path.
  - Project the path into the contracting space and compare the resulting pieces under substitution.
visual:
  - label: Symbols
    text: Iterated substitution word
  - label: Algebra
    text: Matrix and eigenprojection
  - label: Geometry
    text: Self-similar Rauzy fractal
visual_caption: The long symbolic word grows in one direction, while its projected fluctuations fill a bounded fractal.
---

{% include resource-page.liquid %}
