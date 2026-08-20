#!/usr/bin/env python3
"""Create low-text instructional PNGs for the Explore Research pages."""

from __future__ import annotations

import itertools
from pathlib import Path

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt
import numpy as np
from matplotlib.patches import (
    Arc,
    Circle,
    FancyArrowPatch,
    Polygon,
    Rectangle,
)
from PIL import Image, ImageDraw


ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "assets" / "img" / "projects"

BG = "#f7f5f0"
PAPER = "#fffdfa"
INK = "#17324d"
TEAL = "#168f82"
BLUE = "#3178bd"
CORAL = "#df684d"
GOLD = "#e3ad3d"
PALE_BLUE = "#dceaf2"
PALE_GOLD = "#f7e9bd"
GRID = "#c6d5de"


def new_figure():
    fig, ax = plt.subplots(figsize=(10, 5.625), facecolor=BG)
    ax.set_facecolor(BG)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 5.625)
    ax.set_aspect("equal")
    ax.set_axis_off()
    return fig, ax


def save(fig, name: str):
    fig.subplots_adjust(left=0, right=1, top=1, bottom=0)
    path = OUT / name
    fig.savefig(path, dpi=144, facecolor=BG)
    plt.close(fig)
    image = Image.open(path).convert("RGB")
    image.save(path, optimize=True)


def arrow(ax, start, end, color=INK, width=2.2, scale=17, curve=0):
    ax.add_patch(
        FancyArrowPatch(
            start,
            end,
            arrowstyle="-|>",
            mutation_scale=scale,
            color=color,
            linewidth=width,
            connectionstyle=f"arc3,rad={curve}",
            shrinkA=0,
            shrinkB=0,
        )
    )


def rounded_card(ax, x, y, width, height, edge=GRID, face=PAPER, linewidth=1.8):
    ax.add_patch(
        Rectangle(
            (x, y),
            width,
            height,
            facecolor=face,
            edgecolor=edge,
            linewidth=linewidth,
            joinstyle="round",
        )
    )


def draw_coloured_graph(ax, cx, cy, scale=1):
    positions = np.array(
        [
            [-0.7, 0.55],
            [0.55, 0.72],
            [0.82, -0.4],
            [-0.2, -0.75],
            [-0.85, -0.35],
        ]
    )
    edges = [(0, 1), (0, 2), (0, 4), (1, 2), (1, 3), (2, 3), (2, 4), (3, 4)]
    positions = positions * scale + np.array([cx, cy])
    for a, b in edges:
        ax.plot(
            positions[[a, b], 0],
            positions[[a, b], 1],
            color=GRID,
            linewidth=2.4,
            zorder=1,
        )
    for (x, y), color in zip(positions, [BLUE, CORAL, TEAL, BLUE, GOLD]):
        ax.add_patch(Circle((x, y), 0.18 * scale, facecolor=color, edgecolor=PAPER, linewidth=2, zorder=2))


def draw_pair_graph(ax, cx, cy, scale=1):
    positions = np.array(
        [
            [-0.85, 0.65],
            [0, 0.88],
            [0.85, 0.6],
            [-0.72, -0.2],
            [0.12, -0.1],
            [0.8, -0.45],
            [-0.25, -0.85],
        ]
    )
    edges = [(0, 1), (1, 2), (0, 3), (1, 4), (2, 5), (3, 4), (4, 5), (3, 6), (4, 6)]
    positions = positions * scale + np.array([cx, cy])
    for a, b in edges:
        ax.plot(
            positions[[a, b], 0],
            positions[[a, b], 1],
            color=TEAL,
            linewidth=2,
            zorder=1,
        )
    for index, (x, y) in enumerate(positions):
        color = GOLD if index >= 5 else BLUE
        ax.add_patch(Circle((x, y), 0.13 * scale, facecolor=color, edgecolor=PAPER, linewidth=1.5, zorder=2))


def draw_qubo_motivation():
    fig, ax = new_figure()
    ax.text(
        0.38,
        5.23,
        "OBJECTIVE 01",
        color=TEAL,
        fontsize=10,
        fontweight="bold",
        va="center",
    )
    ax.text(
        0.38,
        4.91,
        "Convert graph colouring into a hardware-ready QUBO",
        color=INK,
        fontsize=21,
        fontweight="bold",
        va="center",
    )
    ax.text(
        0.38,
        4.57,
        "Preserve the colouring goal while replacing every 3- and 4-variable interaction with pairwise terms.",
        color="#536b7d",
        fontsize=10.5,
        va="center",
    )
    ax.plot([0.38, 9.62], [4.30, 4.30], color=GRID, linewidth=1.2)

    cards = [(0.38, BLUE), (3.63, CORAL), (6.88, TEAL)]
    for x, edge in cards:
        rounded_card(ax, x, 0.47, 2.74, 3.46, edge=edge, face=PAPER, linewidth=1.7)

    # Step 1: the original graph-colouring objective.
    ax.add_patch(Circle((0.72, 3.58), 0.18, facecolor=BLUE, edgecolor="none"))
    ax.text(0.72, 3.58, "1", ha="center", va="center", color=PAPER, fontsize=9, fontweight="bold")
    ax.text(1.00, 3.60, "INPUT", color=BLUE, fontsize=8.5, fontweight="bold", va="center")
    ax.text(0.68, 3.27, "Colour the graph", color=INK, fontsize=13, fontweight="bold", va="center")
    draw_coloured_graph(ax, 1.75, 2.28, 0.70)
    rounded_card(ax, 0.68, 0.77, 2.13, 0.48, edge=PALE_BLUE, face="#eef5f8", linewidth=1.0)
    ax.text(
        1.745,
        1.01,
        "Adjacent vertices ≠ same colour",
        ha="center",
        va="center",
        color=INK,
        fontsize=8.5,
        fontweight="semibold",
    )

    # Step 2: compact encoding creates unsupported higher-order interactions.
    ax.add_patch(Circle((3.97, 3.58), 0.18, facecolor=CORAL, edgecolor="none"))
    ax.text(3.97, 3.58, "2", ha="center", va="center", color=PAPER, fontsize=9, fontweight="bold")
    ax.text(4.25, 3.60, "COMPACT MODEL", color=CORAL, fontsize=8.5, fontweight="bold", va="center")
    ax.text(3.93, 3.27, "Encode with fewer bits", color=INK, fontsize=13, fontweight="bold", va="center")
    rounded_card(ax, 3.96, 2.27, 2.08, 0.62, edge=GRID, face=BG, linewidth=1.1)
    ax.text(5.00, 2.58, r"$x_i x_j x_k$", ha="center", va="center", color=INK, fontsize=15)
    rounded_card(ax, 3.96, 1.52, 2.08, 0.62, edge=GRID, face=BG, linewidth=1.1)
    ax.text(5.00, 1.83, r"$x_i x_j x_k x_\ell$", ha="center", va="center", color=INK, fontsize=15)
    rounded_card(ax, 4.12, 0.77, 1.76, 0.48, edge="#f3c6bc", face="#fff0ec", linewidth=1.0)
    ax.text(
        5.00,
        1.01,
        "Degree 3–4  •  unsupported",
        ha="center",
        va="center",
        color=CORAL,
        fontsize=8.3,
        fontweight="bold",
    )

    # Step 3: the destination is a pairwise QUBO.
    ax.add_patch(Circle((7.22, 3.58), 0.18, facecolor=TEAL, edgecolor="none"))
    ax.text(7.22, 3.58, "3", ha="center", va="center", color=PAPER, fontsize=9, fontweight="bold")
    ax.text(7.50, 3.60, "OUTPUT", color=TEAL, fontsize=8.5, fontweight="bold", va="center")
    ax.text(7.18, 3.27, "Reduce to pairwise terms", color=INK, fontsize=11.7, fontweight="bold", va="center")
    draw_pair_graph(ax, 8.25, 2.22, 0.67)
    rounded_card(ax, 7.22, 0.77, 2.06, 0.48, edge="#b6ddd7", face="#eaf6f3", linewidth=1.0)
    ax.text(
        8.25,
        1.01,
        r"Degree ≤ 2  •  annealer-ready",
        ha="center",
        va="center",
        color=TEAL,
        fontsize=8.3,
        fontweight="bold",
    )

    arrow(ax, (3.18, 2.17), (3.48, 2.17), color=BLUE, width=1.8, scale=12)
    arrow(ax, (6.43, 2.17), (6.73, 2.17), color=TEAL, width=1.8, scale=12)

    save(fig, "qubo-motivation.png")


def monomial_tile(ax, x, y, term, color):
    rounded_card(ax, x, y, 0.72, 0.45, edge=color, linewidth=1.25)
    ax.text(
        x + 0.36,
        y + 0.23,
        rf"$x_{{{term[0]}}}x_{{{term[1]}}}x_{{{term[2]}}}$",
        ha="center",
        va="center",
        color=INK,
        fontsize=7.7,
    )


def draw_qubo_symmetric_reduction():
    fig, ax = new_figure()
    terms = list(itertools.combinations(range(5), 3))
    ax.text(
        0.38,
        5.23,
        "OBJECTIVE 02",
        color=TEAL,
        fontsize=10,
        fontweight="bold",
        va="center",
    )
    ax.text(
        0.38,
        4.91,
        "Reduce a shared symmetric block — not every term",
        color=INK,
        fontsize=21,
        fontweight="bold",
        va="center",
    )
    ax.text(
        0.38,
        4.57,
        "Ten cubic terms form one homogeneous symmetric block — so reduce the block once.",
        color="#536b7d",
        fontsize=10.5,
        va="center",
    )
    ax.plot([0.38, 9.62], [4.30, 4.30], color=GRID, linewidth=1.2)

    # Conventional, term-by-term treatment.
    rounded_card(ax, 0.38, 0.47, 2.70, 3.46, edge=CORAL, face=PAPER, linewidth=1.7)
    ax.text(0.68, 3.58, "CONVENTIONAL", color=CORAL, fontsize=8.5, fontweight="bold")
    ax.text(0.68, 3.28, "Reduce terms one by one", color=INK, fontsize=12, fontweight="bold")
    shown_terms = terms[:4]
    for index, term in enumerate(shown_terms):
        y = 2.77 - index * 0.48
        rounded_card(ax, 0.68, y, 1.22, 0.34, edge="#efb3a6", face="#fff5f2", linewidth=0.9)
        ax.text(
            1.29,
            y + 0.17,
            rf"$x_{{{term[0]}}}x_{{{term[1]}}}x_{{{term[2]}}}$",
            ha="center",
            va="center",
            color=INK,
            fontsize=8.5,
        )
        arrow(ax, (1.99, y + 0.17), (2.22, y + 0.17), color=CORAL, width=1.1, scale=8)
        ax.add_patch(Circle((2.39, y + 0.17), 0.11, facecolor=CORAL, edgecolor=PAPER, linewidth=1.0))
    ax.text(1.51, 0.94, "+ 6 more independent reductions", ha="center", color="#536b7d", fontsize=8.2)
    rounded_card(ax, 0.68, 0.57, 2.10, 0.28, edge="#f3c6bc", face="#fff0ec", linewidth=1.0)
    ax.text(1.73, 0.71, "10 auxiliary bits", ha="center", va="center", color=CORAL, fontsize=9.5, fontweight="bold")

    # Recognize the common symmetric polynomial.
    rounded_card(ax, 3.58, 0.47, 2.84, 3.46, edge=BLUE, face=PAPER, linewidth=1.7)
    ax.text(3.88, 3.58, "RECOGNISE", color=BLUE, fontsize=8.5, fontweight="bold")
    ax.text(3.88, 3.28, "One repeated pattern", color=INK, fontsize=13, fontweight="bold")
    rounded_card(ax, 3.88, 2.42, 2.24, 0.58, edge=PALE_BLUE, face="#eef5f8", linewidth=1.0)
    ax.text(
        5.00,
        2.71,
        r"$e_3(x_0,\ldots,x_4)=\sum_{i<j<k}x_i x_j x_k$",
        ha="center",
        va="center",
        color=INK,
        fontsize=10.5,
    )
    for index, term in enumerate(terms):
        row, col = divmod(index, 5)
        x = 4.09 + col * 0.45
        y = 1.99 - row * 0.46
        ax.add_patch(Circle((x, y), 0.075, facecolor=BLUE if row == 0 else GOLD, edgecolor="none"))
    ax.add_patch(Arc((5.00, 1.74), 1.97, 0.82, theta1=190, theta2=350, color=TEAL, linewidth=1.8))
    rounded_card(ax, 3.88, 0.57, 2.24, 0.42, edge="#b6ddd7", face="#eaf6f3", linewidth=1.0)
    ax.text(5.00, 0.78, "1 shared symmetric block", ha="center", va="center", color=TEAL, fontsize=9.2, fontweight="bold")

    # Shared reduction output.
    rounded_card(ax, 6.92, 0.47, 2.70, 3.46, edge=TEAL, face=PAPER, linewidth=1.7)
    ax.text(7.22, 3.58, "SHARED REDUCTION", color=TEAL, fontsize=8.5, fontweight="bold")
    ax.text(7.22, 3.28, "Introduce two auxiliaries", color=INK, fontsize=11.7, fontweight="bold")
    ax.plot([7.46, 8.15], [2.34, 2.34], color=TEAL, linewidth=2.2)
    ax.plot([8.15, 8.78], [2.34, 2.72], color=TEAL, linewidth=2.2)
    ax.plot([8.15, 8.78], [2.34, 1.96], color=TEAL, linewidth=2.2)
    ax.add_patch(Circle((7.46, 2.34), 0.19, facecolor=BLUE, edgecolor=PAPER, linewidth=1.5))
    for index, y in enumerate((2.72, 1.96), start=1):
        ax.add_patch(Circle((8.85, y), 0.25, facecolor=TEAL, edgecolor=PAPER, linewidth=1.5))
        ax.text(8.85, y, rf"$y_{index}$", ha="center", va="center", color=PAPER, fontsize=9, fontweight="bold")
    rounded_card(ax, 7.22, 1.20, 2.10, 0.42, edge="#b6ddd7", face="#eaf6f3", linewidth=1.0)
    ax.text(8.27, 1.41, "2 auxiliary bits", ha="center", va="center", color=TEAL, fontsize=9.5, fontweight="bold")
    rounded_card(ax, 7.22, 0.57, 2.10, 0.42, edge=TEAL, face=TEAL, linewidth=1.0)
    ax.text(8.27, 0.78, "80% fewer", ha="center", va="center", color=PAPER, fontsize=10.2, fontweight="bold")

    arrow(ax, (3.16, 2.18), (3.45, 2.18), color=BLUE, width=1.8, scale=11)
    arrow(ax, (6.50, 2.18), (6.79, 2.18), color=TEAL, width=1.8, scale=11)

    save(fig, "qubo-symmetric-reduction.png")


def draw_qubo_resources():
    fig, ax = new_figure()
    ax.text(
        0.38,
        5.23,
        "OBJECTIVE 03",
        color=TEAL,
        fontsize=10,
        fontweight="bold",
        va="center",
    )
    ax.text(
        0.38,
        4.91,
        "Cut the logical QUBO before it reaches quantum hardware",
        color=INK,
        fontsize=21,
        fontweight="bold",
        va="center",
    )
    ax.text(
        0.38,
        4.57,
        "Complete graph K₈ with eight colours  •  term-by-term reduction versus shared symmetric reduction",
        color="#536b7d",
        fontsize=10.5,
        va="center",
    )
    ax.plot([0.38, 9.62], [4.30, 4.30], color=GRID, linewidth=1.2)

    metrics = [
        {
            "x": 0.38,
            "title": "BINARY VARIABLES",
            "before": 1180,
            "after": 480,
            "saved": "700 fewer variables",
            "percent": "59% reduction",
        },
        {
            "x": 5.12,
            "title": "QUADRATIC INTERACTIONS",
            "before": 5849,
            "after": 2797,
            "saved": "3,052 fewer interactions",
            "percent": "52% reduction",
        },
    ]
    for metric in metrics:
        x = metric["x"]
        rounded_card(ax, x, 0.47, 4.50, 3.46, edge=GRID, face=PAPER, linewidth=1.5)
        ax.text(x + 0.30, 3.58, metric["title"], color=INK, fontsize=9.2, fontweight="bold")

        bar_x = x + 1.38
        bar_w = 2.72
        before_w = bar_w
        after_w = bar_w * metric["after"] / metric["before"]

        ax.text(x + 0.30, 2.89, "Term-by-term", color="#536b7d", fontsize=8.7, va="center")
        rounded_card(ax, bar_x, 2.65, before_w, 0.48, edge="#f3c6bc", face="#fbe3dd", linewidth=0.8)
        ax.add_patch(Rectangle((bar_x, 2.65), before_w, 0.48, facecolor=CORAL, edgecolor="none"))
        ax.text(
            bar_x + before_w - 0.12,
            2.89,
            f'{metric["before"]:,}',
            ha="right",
            va="center",
            color=PAPER,
            fontsize=10,
            fontweight="bold",
        )

        ax.text(x + 0.30, 2.14, "Shared", color="#536b7d", fontsize=8.7, va="center")
        rounded_card(ax, bar_x, 1.90, before_w, 0.48, edge="#d7ebe7", face="#eaf6f3", linewidth=0.8)
        ax.add_patch(Rectangle((bar_x, 1.90), after_w, 0.48, facecolor=TEAL, edgecolor="none"))
        ax.text(
            bar_x + after_w + 0.10,
            2.14,
            f'{metric["after"]:,}',
            ha="left",
            va="center",
            color=TEAL,
            fontsize=10,
            fontweight="bold",
        )

        ax.text(x + 0.30, 1.35, metric["saved"], color=INK, fontsize=10.4, fontweight="bold")
        rounded_card(ax, x + 0.30, 0.70, 1.74, 0.43, edge=TEAL, face=TEAL, linewidth=1.0)
        ax.text(
            x + 1.17,
            0.915,
            metric["percent"],
            ha="center",
            va="center",
            color=PAPER,
            fontsize=9.3,
            fontweight="bold",
        )

    rounded_card(ax, 2.29, 0.10, 5.42, 0.27, edge=PALE_BLUE, face="#eef5f8", linewidth=0.8)
    ax.text(
        5.00,
        0.235,
        "Smaller logical model  →  fewer physical qubits and couplers needed for embedding",
        ha="center",
        va="center",
        color=INK,
        fontsize=8.7,
        fontweight="semibold",
    )

    save(fig, "qubo-resource-comparison.png")


def substitution_word(iterations=16):
    rules = {0: [0, 1], 1: [0, 2], 2: [0]}
    word = [0]
    for _ in range(iterations):
        word = [symbol for item in word for symbol in rules[item]]
    return np.array(word)


def rauzy_points(iterations=16):
    word = substitution_word(iterations)
    counts = np.zeros((len(word), 3))
    running = np.zeros(3)
    for index, symbol in enumerate(word):
        running[symbol] += 1
        counts[index] = running
    substitution = np.array([[1, 1, 1], [1, 0, 0], [0, 1, 0]], dtype=float)
    eigenvalues, eigenvectors = np.linalg.eig(substitution)
    growth = np.real(eigenvectors[:, np.argmax(np.real(eigenvalues))])
    growth /= np.linalg.norm(growth)
    basis = np.linalg.svd(growth.reshape(1, 3))[2][1:].T
    projected = counts @ basis
    projected -= projected.mean(axis=0)
    projected /= np.max(np.linalg.norm(projected, axis=1))
    return word, projected


def draw_rauzy_method():
    fig, ax = new_figure()
    palette = np.array([BLUE, TEAL, CORAL])

    # Symbolic substitution.
    rules = [(0, [0, 1]), (1, [0, 2]), (2, [0])]
    for row, (source, result) in enumerate(rules):
        y = 3.9 - row * 1.0
        ax.add_patch(Circle((0.65, y), 0.21, facecolor=palette[source], edgecolor=PAPER, linewidth=1.5))
        arrow(ax, (0.95, y), (1.4, y), color=GRID, width=1.6, scale=11)
        for index, symbol in enumerate(result):
            ax.add_patch(
                Circle((1.72 + index * 0.5, y), 0.21, facecolor=palette[symbol], edgecolor=PAPER, linewidth=1.5)
            )

    arrow(ax, (2.6, 2.8), (3.35, 2.8), color=BLUE)

    # Prefix counts as a three-direction lattice walk.
    short_word = substitution_word(5)[:18]
    directions = np.array([[0.38, 0], [0.12, 0.34], [-0.26, 0.22]])
    walk = np.vstack([np.zeros(2), np.cumsum(directions[short_word], axis=0)])
    walk -= walk.mean(axis=0)
    walk *= 0.78
    walk += np.array([4.65, 2.8])
    for index, symbol in enumerate(short_word):
        ax.plot(walk[index : index + 2, 0], walk[index : index + 2, 1], color=palette[symbol], linewidth=3)
    ax.scatter(walk[:, 0], walk[:, 1], s=14, color=INK, zorder=3)

    arrow(ax, (5.9, 2.8), (6.65, 2.8), color=TEAL)

    # Projection to the contracting plane.
    word, points = rauzy_points()
    ax.scatter(
        8.15 + points[:, 0] * 1.5,
        2.8 + points[:, 1] * 1.75,
        c=palette[word],
        s=3,
        alpha=0.78,
        linewidths=0,
    )

    save(fig, "rauzy-method.png")


def draw_rauzy_self_similarity():
    fig, ax = new_figure()
    word, points = rauzy_points()
    palette = np.array([BLUE, TEAL, CORAL])

    ax.scatter(
        2.4 + points[:, 0] * 1.75,
        2.82 + points[:, 1] * 1.95,
        c=palette[word],
        s=3.2,
        alpha=0.82,
        linewidths=0,
    )
    arrow(ax, (4.45, 2.82), (5.2, 2.82), color=TEAL)

    centres = [(6.25, 3.75), (8.35, 3.75), (7.3, 1.65)]
    for symbol, (cx, cy) in enumerate(centres):
        subset = points[word == symbol]
        subset -= subset.mean(axis=0)
        scale = max(np.ptp(subset[:, 0]), np.ptp(subset[:, 1]))
        subset /= scale
        ax.scatter(
            cx + subset[:, 0] * 1.35,
            cy + subset[:, 1] * 1.1,
            color=palette[symbol],
            s=3.2,
            alpha=0.82,
            linewidths=0,
        )

    save(fig, "rauzy-self-similarity.png")


def draw_hyperbolic_motivation():
    """Compare two verified source figures without redrawing either lattice."""
    source_dir = OUT / "source"
    euclidean = Image.open(source_dir / "hyperbolic-surface-code-euclidean-source.png").convert("RGBA")
    hyperbolic = Image.open(source_dir / "hyperbolic-surface-code-2-3-7-source.png").convert("RGBA")

    canvas = Image.new("RGB", (1440, 810), BG)
    euclidean.thumbnail((510, 510), Image.Resampling.LANCZOS)
    hyperbolic.thumbnail((650, 650), Image.Resampling.LANCZOS)
    canvas.paste(euclidean, (70, 150), euclidean)
    canvas.paste(hyperbolic, (735, 80), hyperbolic)

    draw = ImageDraw.Draw(canvas)
    draw.line((615, 405, 700, 405), fill=TEAL, width=8)
    draw.polygon([(700, 405), (672, 386), (672, 424)], fill=TEAL)

    path = OUT / "hyperbolic-surface-code-motivation.png"
    canvas.save(path, optimize=True)


def draw_book_icon(ax, cx, cy, color):
    ax.add_patch(Polygon([(cx - 0.42, cy + 0.28), (cx - 0.03, cy + 0.18), (cx - 0.03, cy - 0.38), (cx - 0.42, cy - 0.28)], closed=True, facecolor=PAPER, edgecolor=color, linewidth=1.8))
    ax.add_patch(Polygon([(cx + 0.42, cy + 0.28), (cx + 0.03, cy + 0.18), (cx + 0.03, cy - 0.38), (cx + 0.42, cy - 0.28)], closed=True, facecolor=PAPER, edgecolor=color, linewidth=1.8))


def draw_building_icon(ax, cx, cy, color):
    ax.add_patch(Rectangle((cx - 0.42, cy - 0.36), 0.84, 0.68, facecolor=PAPER, edgecolor=color, linewidth=1.8))
    ax.add_patch(Polygon([(cx - 0.5, cy + 0.32), (cx, cy + 0.58), (cx + 0.5, cy + 0.32)], closed=True, facecolor=color, edgecolor=color))
    for row in range(2):
        for col in range(3):
            ax.add_patch(Rectangle((cx - 0.29 + col * 0.25, cy - 0.18 + row * 0.27), 0.11, 0.13, facecolor=PALE_BLUE, edgecolor="none"))


def draw_person_icon(ax, cx, cy, color, group=False):
    offsets = [-0.28, 0.28] if group else [0]
    for offset in offsets:
        ax.add_patch(Circle((cx + offset, cy + 0.26), 0.16, facecolor=color, edgecolor=PAPER, linewidth=1.2))
        ax.add_patch(Arc((cx + offset, cy - 0.25), 0.58, 0.75, theta1=20, theta2=160, color=color, linewidth=3))


def input_card(ax, x, y, title, color, icon):
    rounded_card(ax, x, y, 1.75, 1.03, edge=color, linewidth=1.8)
    icon(ax, x + 0.42, y + 0.51, color)
    ax.text(
        x + 1.28,
        y + 0.51,
        title,
        ha="center",
        va="center",
        color=INK,
        fontsize=9.2,
        fontweight="semibold",
    )


def draw_clock(ax, cx, cy):
    ax.add_patch(Circle((cx, cy), 0.18, fill=False, edgecolor=INK, linewidth=1.5))
    ax.plot([cx, cx], [cy, cy + 0.1], color=INK, linewidth=1.4)
    ax.plot([cx, cx + 0.08], [cy, cy - 0.06], color=INK, linewidth=1.4)


def draw_timetable(ax, x, y, width, height):
    cols, rows = 5, 6
    cell_w, cell_h = width / cols, height / rows
    for col in range(cols):
        for row in range(rows):
            ax.add_patch(
                Rectangle(
                    (x + col * cell_w, y + row * cell_h),
                    cell_w,
                    cell_h,
                    facecolor=PAPER,
                    edgecolor=GRID,
                    linewidth=0.7,
                )
            )
    blocks = [
        (0, 4, 1, 1, BLUE),
        (1, 2, 2, 1, GOLD),
        (3, 4, 1, 2, TEAL),
        (0, 0, 2, 1, CORAL),
        (2, 0, 1, 2, BLUE),
        (4, 1, 1, 2, GOLD),
    ]
    for col, row, span_x, span_y, color in blocks:
        ax.add_patch(
            Rectangle(
                (x + col * cell_w + 0.04, y + row * cell_h + 0.04),
                span_x * cell_w - 0.08,
                span_y * cell_h - 0.08,
                facecolor=color,
                edgecolor=PAPER,
                linewidth=1,
                alpha=0.9,
            )
        )


def draw_timetabling_method():
    fig, ax = new_figure()

    input_card(ax, 0.3, 4.05, "Courses", BLUE, draw_book_icon)
    input_card(ax, 0.3, 2.72, "Rooms", GOLD, draw_building_icon)
    input_card(ax, 0.3, 1.39, "Instructors", CORAL, draw_person_icon)
    input_card(
        ax,
        0.3,
        0.06,
        "Students",
        TEAL,
        lambda plot, cx, cy, color: draw_person_icon(plot, cx, cy, color, group=True),
    )

    for y, color in zip([4.56, 3.23, 1.90, 0.57], [BLUE, GOLD, CORAL, TEAL]):
        arrow(ax, (2.15, y), (3.1, 2.81), color=color, width=1.65, scale=13)

    rounded_card(ax, 3.18, 1.35, 3.05, 2.92, edge=INK, linewidth=2.2)
    ax.text(4.7, 3.94, "Constraints", ha="center", va="center", color=INK, fontsize=13, fontweight="semibold")

    factors = [
        (3.48, 2.82, "time", BLUE),
        (4.83, 2.82, "capacity", GOLD),
        (3.48, 1.78, "conflict", CORAL),
        (4.83, 1.78, "preference", TEAL),
    ]
    for x, y, label, color in factors:
        rounded_card(ax, x, y, 1.13, 0.72, edge=color, face=BG, linewidth=1.3)
        if label == "time":
            draw_clock(ax, x + 0.22, y + 0.36)
            text_x = x + 0.72
        else:
            text_x = x + 0.56
        ax.text(text_x, y + 0.36, label, ha="center", va="center", color=INK, fontsize=8.5)

    arrow(ax, (6.32, 2.81), (7.08, 2.81), color=TEAL, width=2.5)
    draw_timetable(ax, 7.28, 1.2, 2.35, 3.22)

    save(fig, "automated-timetabling-method.png")


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    draw_qubo_motivation()
    draw_qubo_symmetric_reduction()
    draw_qubo_resources()
    draw_rauzy_method()
    draw_rauzy_self_similarity()
    draw_hyperbolic_motivation()
    draw_timetabling_method()
    print("Created 7 low-text instructional topic images")


if __name__ == "__main__":
    main()
