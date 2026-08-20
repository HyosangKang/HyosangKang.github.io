---
layout: page
title: AES Encryption
description: Classify the key strings and encrypted blocks from the BS203 AES challenge.
permalink: /learn/linear-algebra/aes-encryption/
subject: Linear Algebra
subject_slug: linear-algebra
project_order: 1
pathway: Learn Mathematics
back_url: /learn/linear-algebra/
back_label: Linear Algebra projects
question: Which of the given strings contains an AES key, and which one is ciphertext?
goals:
  - Distinguish a key-container string from a ciphertext block.
  - Connect a hidden 16-character key with one AES-128 decryption.
method_intro: Revisit Project I from the Fall 2023 BS203 Linear Algebra class as an eight-round simulator.
method:
  - Inspect one unused key string and one unused encrypted string in each round.
  - Choose the string that contains 16 consecutive characters used as the AES key.
  - Check the classification, then reveal the hidden key window and decoded message.
visual:
  - label: 16 strings
    text: 8 key containers and 8 ciphertexts
  - label: Classify
    text: Decide which string has which role
  - label: Reveal
    text: Recover the hidden fruit message
visual_caption: The original classroom challenge mixed key containers and ciphertext blocks that look alike.
---

{% include resource-page.liquid %}
