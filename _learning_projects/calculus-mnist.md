---
layout: page
title: MNIST Data
description: Build a neural-network model from scratch.
permalink: /learn/calculus/mnist/
subject: Calculus
subject_slug: calculus
project_order: 6
pathway: Learn Mathematics
back_url: /learn/calculus/
back_label: Calculus projects
question: How can gradients train a neural network to recognize handwritten digits?
goals:
  - Build the essential parts of a neural network without relying on a machine-learning framework.
  - Use a loss function and gradient-based updates to classify MNIST images.
method_intro: Treat each image as a vector of pixel values and learn a sequence of transformations.
method:
  - Normalize the MNIST data and implement a feed-forward network.
  - Calculate the loss and propagate its derivatives backward through the network.
  - Train with gradient descent, then examine accuracy and common classification errors.
visual:
  - label: Input
    text: 28 × 28 pixels
  - label: Learn
    text: Weights and gradients
  - label: Output
    text: Digit probabilities
visual_caption: Training adjusts the network weights so a pixel pattern becomes a probability for each digit.
---

{% include resource-page.liquid %}
