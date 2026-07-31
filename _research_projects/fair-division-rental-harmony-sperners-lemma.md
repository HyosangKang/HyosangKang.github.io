---
layout: page
title: Fair Division, Rental Harmony, and Sperner’s Lemma
description: A guided research project in Algorithms & Discrete Mathematics.
area: Algorithms & Discrete Mathematics
area_slug: algorithms-discrete-mathematics
project_order: 2
pathway: Experience Research
back_url: "/experience/algorithms-discrete-mathematics/"
back_label: Algorithms & Discrete Mathematics projects
question: Can rooms and rent be divided so that every person prefers their own choice?
goals:
  - Model an envy-free rental division.
  - Use Sperner's lemma to find or approximate a fair solution.
method_intro: Represent possible rent divisions as points in a simplex and label them by preferences.
method:
  - Collect or define each person's room preferences at different prices.
  - Triangulate the price simplex and assign preference labels to its vertices.
  - Search for a fully labelled small simplex and interpret it as a rental division.
visual:
  - label: Preferences
    text: Rooms at different rents
  - label: Simplex
    text: Sperner labels
  - label: Division
    text: Envy-free assignment
visual_caption: A fully labelled region points to prices where all rooms can be assigned without envy.
---

{% include resource-page.liquid %}
