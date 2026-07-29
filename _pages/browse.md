---
layout: page
title: Find Resources
permalink: /browse/
description: Find all Math Projects resources in one place.
nav: true
nav_order: 4
---

{% assign subject_pages = site.pages | where_exp: "item", "item.subject_landing == true" | sort: "subject_order" %}

{% for subject in subject_pages %}

## {{ subject.title }}

{% include project-list.liquid subject=subject.subject_slug %}
{% endfor %}
