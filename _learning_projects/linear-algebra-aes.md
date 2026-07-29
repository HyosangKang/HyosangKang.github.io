---
layout: page
title: AES Encryption
description: Investigate the mathematics used by AES encryption.
permalink: /learn/linear-algebra/aes-encryption/
subject: Linear Algebra
subject_slug: linear-algebra
project_order: 1
pathway: Learn Mathematics
back_url: /learn/linear-algebra/
back_label: Linear Algebra projects
question: How do matrix and finite-field operations turn a readable message into encrypted data?
goals:
  - Represent a data block as a matrix of bytes.
  - Implement and trace the main transformations in AES encryption.
method_intro: Follow one block of data through a sequence of reversible rounds.
method:
  - Arrange the input bytes in the AES state matrix and create the round keys.
  - Apply byte substitution, row shifting, column mixing, and key addition.
  - Test encryption with known examples and inspect how a small input change spreads.
visual:
  - label: Input
    text: Plaintext byte matrix
  - label: Rounds
    text: Substitute, shift, mix, add key
  - label: Output
    text: Ciphertext block
visual_caption: Each AES round rearranges and combines the bytes so the original pattern disappears.
---

{% include resource-page.liquid %}
