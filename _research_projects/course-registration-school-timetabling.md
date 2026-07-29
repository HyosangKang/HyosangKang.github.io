---
layout: page
title: Course Registration and School Timetabling
description: A guided research project in Algorithms & Discrete Mathematics.
area: Algorithms & Discrete Mathematics
area_slug: algorithms-discrete-mathematics
project_order: 4
pathway: Experience Research
back_url: "/research/#algorithms-discrete-mathematics"
back_label: Algorithms & Discrete Mathematics projects
question: How can courses be assigned to students, rooms, and times while satisfying many competing constraints?
goals:
  - Build a mathematical model of registration or timetabling constraints.
  - Create and evaluate schedules using real or realistic data.
method_intro: Express every possible assignment with decision variables, hard rules, and preferences.
method:
  - Record conflicts, capacities, availability, and student choices.
  - Use integer programming, graph methods, or a constructive heuristic.
  - Check feasibility and compare schedules by preference, balance, and running time.
visual:
  - label: Data
    text: Courses, students, rooms, times
  - label: Solver
    text: Constraints and priorities
  - label: Schedule
    text: Feasible assignments
visual_caption: The algorithm must satisfy fixed rules while balancing many requests that cannot all be met.
---

{% include resource-page.liquid %}
