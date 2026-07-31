---
layout: page
title: Gaussian Functions in Protein Folding
description: A guided research project in Mathematical Modelling.
area: Mathematical Modelling
area_slug: mathematical-modelling
project_order: 2
pathway: Experience Research
back_url: "/experience/mathematical-modelling/"
back_label: Mathematical Modelling projects
question: Can Gaussian functions describe useful interactions or shapes in a protein-folding model?
goals:
  - Build a Gaussian-based energy or similarity model.
  - Test whether it distinguishes or guides candidate protein shapes.
method_intro: Use smooth Gaussian terms to assign stronger influence to nearby structural features.
method:
  - Choose a simplified representation of amino-acid positions or contacts.
  - Define Gaussian features or an energy function from distances.
  - Fit or simulate the model and compare scores with known or generated structures.
visual:
  - label: Structure
    text: Amino-acid positions
  - label: Function
    text: Gaussian distance terms
  - label: Result
    text: Energy or fold score
visual_caption: Each Gaussian term creates a smooth region of influence around a structural feature.
---

{% include resource-page.liquid %}
