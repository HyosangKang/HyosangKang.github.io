---
layout: page
title: All topics
permalink: /browse/
description: Browse all Math Projects topics in one place.
nav: true
nav_order: 4
---

## Learn Mathematics

{% assign subject_pages = site.pages | where_exp: "item", "item.subject_landing == true" | sort: "subject_order" %}

{% for subject in subject_pages %}

### {{ subject.title }}

{% include project-list.liquid subject=subject.subject_slug %}
{% endfor %}

## Experience Research

### Geometry & Dynamics

{% include research-project-list.liquid area="geometry-dynamics" %}

### Algorithms & Discrete Mathematics

{% include research-project-list.liquid area="algorithms-discrete-mathematics" %}

### Quantum Computing

{% include research-project-list.liquid area="quantum-computing" %}

### Mathematical Modelling

{% include research-project-list.liquid area="mathematical-modelling" %}

### Mathematical Foundations

{% include research-project-list.liquid area="mathematical-foundations" %}

## Explore Research

<div class="subject-grid">
  <a class="subject-card" href="{{ '/explore/rauzy-fractal/' | relative_url }}">
    <span class="subject-count">Research Topic</span>
    <h3>Rauzy Fractal</h3>
    <span class="subject-link">Open topic →</span>
  </a>
  <a class="subject-card" href="{{ '/explore/qubo/' | relative_url }}">
    <span class="subject-count">Research Topic</span>
    <h3>QUBO</h3>
    <span class="subject-link">Open topic →</span>
  </a>
  <a class="subject-card" href="{{ '/explore/hyperbolic-surface-code/' | relative_url }}">
    <span class="subject-count">Research Topic</span>
    <h3>Hyperbolic Surface Code</h3>
    <span class="subject-link">Open topic →</span>
  </a>
  <a class="subject-card" href="{{ '/explore/automated-timetabling-algorithm/' | relative_url }}">
    <span class="subject-count">Research Topic</span>
    <h3>Automated Timetabling Algorithm</h3>
    <span class="subject-link">Open topic →</span>
  </a>
</div>
