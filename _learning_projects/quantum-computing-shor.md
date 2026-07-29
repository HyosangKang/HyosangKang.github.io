---
layout: page
title: "Quantum Cryptography: Shor’s Algorithm"
description: Explore Shor’s algorithm and its implications for RSA.
permalink: /learn/quantum-computing/shors-algorithm/
subject: Quantum Computing
subject_slug: quantum-computing
project_order: 2
pathway: Learn Mathematics
back_url: /learn/quantum-computing/
back_label: Quantum Computing projects
question: How does quantum period finding turn integer factorization into a tractable task?
goals:
  - Explain the link between period finding, factorization, and RSA security.
  - Run the main stages of Shor's algorithm on small integers.
method_intro: Use a quantum routine to find a period, then recover factors with classical arithmetic.
method:
  - Choose a number \(N\) and a base \(a\), then study the periodic powers of \(a\) modulo \(N\).
  - Model modular exponentiation and the quantum Fourier transform in a small circuit or simulator.
  - Estimate the period and use greatest common divisors to obtain candidate factors.
visual:
  - label: Integer
    text: Number to factor
  - label: Quantum step
    text: Find a modular period
  - label: Classical step
    text: Recover the factors
visual_caption: Shor's algorithm uses a quantum period to expose the arithmetic structure of a composite number.
---

{% include resource-page.liquid %}
