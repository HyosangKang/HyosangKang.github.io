---
layout: page
title: Image Classification with PCA
description: Use principal component analysis for image classification.
permalink: /learn/linear-algebra/image-classification-pca/
subject: Linear Algebra
subject_slug: linear-algebra
project_order: 3
pathway: Learn Mathematics
back_url: /learn/linear-algebra/
back_label: Linear Algebra projects
question: How can an image be compressed while keeping the features needed for classification?
goals:
  - Find principal components from a collection of images.
  - Project, reconstruct, and classify images in a smaller space.
method_intro: Replace thousands of pixel values with a few directions that explain the most variation.
method:
  - Turn each image into a vector, center the data, and compute a covariance matrix or SVD.
  - Keep the leading eigenvectors and project every image onto them.
  - Train a simple classifier and measure how the number of components affects accuracy.
visual:
  - label: Images
    text: High-dimensional pixels
  - label: PCA
    text: Leading variation directions
  - label: Result
    text: Compact features and classes
visual_caption: PCA keeps the strongest patterns in the data and removes many less useful dimensions.
---

{% include resource-page.liquid %}
