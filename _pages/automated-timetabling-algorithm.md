---
layout: page
title: Automated Timetabling Algorithm
permalink: /topics/automated-timetabling-algorithm/
description: Generate and revise a university timetable under course, room, instructor, and student constraints.
nav: false
pathway: Explore Research
back_url: /topics/
back_label: Explore Research
question: How can a university timetable be generated and revised while respecting courses, rooms, instructors, and students?
goals:
  - Produce a complete timetable without instructor, student, or classroom conflicts.
  - Respect course structure, room capacity, available times, preferences, and limits on teaching days.
  - Give administrators a practical app for review, revision, requests, saving, and distribution.
method_intro: DATA treats timetabling as a constraint-search problem, while leaving final review and adjustment in the hands of an administrator.
method:
  - Import course, room, instructor, and student-plan data from a workbook.
  - Place highly constrained courses first and try an allowed time and suitable room for each one.
  - Reject clashes and other rule violations, then backtrack or reorder courses when no valid choice remains.
  - Review the result visually, move course blocks, check new conflicts, undo changes, and save the timetable.
visual_caption: Course, room, instructor, and student constraints are combined in a backtracking search, producing a timetable that can still be checked and revised by a person.
demo_link_text: Open the interactive DATA demo ↗
demo_note: The browser-based WASM demo will be linked here when it is added.
---

{% include resource-page.liquid %}

<section class="resource-section">
  <h2>Source</h2>
  <p><a href="https://github.com/HyosangKang/dgist-automated-timetabling-algorithm" target="_blank" rel="noopener">DATA project repository ↗</a></p>
</section>
