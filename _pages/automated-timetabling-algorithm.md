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
method_intro: DATA treats timetabling as a constraint-search problem, while leaving final review and adjustment in the hands of an administrator. It is a try-and-error search, but the retry strategy changes whenever the search becomes stuck.
method:
  - Import course, room, instructor, and student-plan data from a workbook.
  - Place highly constrained courses first and try an allowed time and suitable room for each one.
  - Reject clashes and other rule violations. When no valid choice remains, move the troublesome course earlier and retry from a previous point.
  - Review the result visually, move course blocks, check new conflicts, undo changes, and save the timetable.
core_idea: Ordinary try-and-error can loop if it repeatedly returns to the same recent decision. DATA keeps a backward-step value for the course that could not be placed. The first retry steps back a little; if the same position fails again, the step grows and the course is moved ahead of more previously placed courses. After reaching the start of the active search, the step cycles. This changing jump distance explores a different ordering instead of replaying one short loop.
core_steps:
  - Try a valid time-and-room candidate for the current course.
  - If every candidate conflicts, increase that course's backward step.
  - Jump back by that distance, move the difficult course earlier in the order, clear the affected later placements, and try again.
  - Continue until every included course has a valid assignment or the run is cancelled.
visual_caption: Course, room, instructor, and student constraints are combined in a backtracking search, producing a timetable that can still be checked and revised by a person.
demo_link_text: Open the interactive DATA demo ↗
demo_url: /demo/automated-timetabling/
demo_note: The sample workbook loads automatically in a full browser simulation of DATA, or you can choose your own. The timetable engine runs locally through WebAssembly; your workbook is not uploaded.
hide_source_images: true
---

{% include resource-page.liquid %}

<section class="resource-section">
  <h2>Source</h2>
  <p><a href="https://github.com/HyosangKang/dgist-automated-timetabling-algorithm" target="_blank" rel="noopener">DATA project repository ↗</a></p>
</section>
