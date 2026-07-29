#!/usr/bin/env python3
"""Create reproducible instructional PNGs for the Explore Research pages."""

from __future__ import annotations

import itertools
from pathlib import Path

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt
import networkx as nx
import numpy as np
from matplotlib.patches import Circle, FancyArrowPatch, Polygon, Rectangle
from PIL import Image


ROOT = Path(__file__).resolve().parents[1]
OUT = ROOT / "assets" / "img" / "projects"

BG = "#f7f5f0"
INK = "#16324f"
TEAL = "#1b998b"
BLUE = "#2d7dd2"
CORAL = "#e76f51"
GOLD = "#e9c46a"
PALE = "#dce8ef"
GRID = "#cbd8df"


def new_figure(title: str):
    fig, ax = plt.subplots(figsize=(10, 5.625), facecolor=BG)
    ax.set_facecolor(BG)
    fig.text(0.045, 0.93, title, color=INK, fontsize=18, fontweight="semibold")
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 5)
    ax.set_axis_off()
    return fig, ax


def save(fig, name: str):
    fig.tight_layout(rect=(0.025, 0.025, 0.98, 0.89))
    path = OUT / name
    fig.savefig(path, dpi=125, facecolor=BG)
    plt.close(fig)
    image = Image.open(path).convert("RGB")
    image.save(path, optimize=True)


def arrow(ax, start, end, label: str | None = None, color=INK):
    ax.add_patch(
        FancyArrowPatch(
            start,
            end,
            arrowstyle="-|>",
            mutation_scale=18,
            color=color,
            linewidth=2.2,
        )
    )
    if label:
        midpoint = (np.array(start) + np.array(end)) / 2
        ax.text(
            midpoint[0],
            midpoint[1] + 0.22,
            label,
            ha="center",
            color=color,
            fontsize=10,
        )


def card(ax, x, y, width, height, title, lines, color=BLUE):
    ax.add_patch(
        Rectangle(
            (x, y),
            width,
            height,
            facecolor="white",
            edgecolor=color,
            linewidth=2,
        )
    )
    ax.add_patch(
        Rectangle(
            (x, y + height - 0.48),
            width,
            0.48,
            facecolor=color,
            edgecolor=color,
        )
    )
    ax.text(
        x + width / 2,
        y + height - 0.24,
        title,
        ha="center",
        va="center",
        color="white",
        fontsize=11,
        fontweight="bold",
    )
    for index, line in enumerate(lines):
        ax.text(
            x + 0.16,
            y + height - 0.82 - index * 0.38,
            line,
            color=INK,
            fontsize=9.5,
            va="top",
        )


def draw_qubo_motivation():
    fig, ax = new_figure("Why graph colouring needs efficient degree reduction")

    graph = nx.Graph()
    graph.add_edges_from(
        [(0, 1), (0, 2), (0, 4), (1, 2), (1, 3), (2, 3), (2, 4), (3, 4)]
    )
    positions = {
        0: (1.3, 3.75),
        1: (0.75, 2.35),
        2: (2.15, 2.7),
        3: (2.75, 1.45),
        4: (1.25, 1.1),
    }
    node_colors = [BLUE, CORAL, TEAL, BLUE, GOLD]
    nx.draw_networkx_edges(graph, positions, ax=ax, edge_color=GRID, width=2)
    nx.draw_networkx_nodes(
        graph,
        positions,
        ax=ax,
        node_color=node_colors,
        node_size=420,
        edgecolors="white",
        linewidths=2,
    )
    ax.text(1.65, 4.55, "Graph colouring", ha="center", color=INK, fontweight="bold")
    ax.text(
        1.65,
        0.42,
        "penalize adjacent vertices\nthat receive the same colour",
        ha="center",
        color=INK,
        fontsize=10,
    )

    arrow(ax, (3.15, 2.65), (3.85, 2.65))

    card(
        ax,
        3.9,
        1.25,
        2.45,
        2.85,
        "Binary objective",
        [
            r"$x_ix_jx_k$",
            r"$x_ix_jx_kx_l$",
            "",
            "compact colour encoding",
            "creates terms above degree 2",
        ],
        CORAL,
    )

    arrow(ax, (6.45, 2.65), (7.1, 2.65), "reduce to pairs")

    card(
        ax,
        7.15,
        1.25,
        2.25,
        2.85,
        "Quantum annealer",
        [
            r"$Q(x)=\sum Q_{ij}x_ix_j$",
            "",
            "accepts a quadratic",
            "binary energy only",
        ],
        TEAL,
    )
    for row in range(3):
        for col in range(4):
            px = 7.52 + col * 0.43
            py = 1.62 + row * 0.3
            ax.add_patch(Circle((px, py), 0.055, color=BLUE))
            if col < 3:
                ax.plot([px + 0.055, px + 0.375], [py, py], color=GRID, lw=1)
            if row < 2:
                ax.plot([px, px], [py + 0.055, py + 0.245], color=GRID, lw=1)

    ax.text(
        5,
        0.12,
        "Every extra auxiliary bit consumes a qubit and enlarges the hardware-embedding problem.",
        ha="center",
        color=CORAL,
        fontsize=11,
        fontweight="semibold",
    )
    save(fig, "qubo-motivation.png")


def draw_qubo_symmetric_reduction():
    fig, ax = new_figure("Reduce the shared structure, not each monomial separately")
    terms = list(itertools.combinations(range(5), 3))

    ax.text(2.1, 4.55, "Conventional monomial reduction", ha="center", color=INK, fontweight="bold")
    for index, term in enumerate(terms):
        row, col = divmod(index, 2)
        x, y = 0.45 + col * 1.45, 3.95 - row * 0.58
        ax.add_patch(Rectangle((x, y), 1.18, 0.4, facecolor="white", edgecolor=GRID))
        ax.text(
            x + 0.59,
            y + 0.2,
            "$x_{}x_{}x_{}$".format(*term),
            ha="center",
            va="center",
            color=INK,
            fontsize=10,
        )
        ax.add_patch(Circle((3.22, y + 0.2), 0.13, facecolor=CORAL, edgecolor="white"))
        ax.text(3.22, y + 0.2, f"$w_{{{index}}}$", ha="center", va="center", color="white", fontsize=7)
    ax.text(
        2.05,
        0.62,
        "10 cubic terms → 10 auxiliary bits",
        ha="center",
        color=CORAL,
        fontsize=12,
        fontweight="semibold",
    )

    arrow(ax, (4.05, 2.55), (5.05, 2.55), "recognize symmetry", TEAL)

    ax.text(7.55, 4.55, "Symmetric reduction", ha="center", color=INK, fontweight="bold")
    ax.add_patch(
        Rectangle((5.25, 1.05), 4.35, 3.05, facecolor="white", edgecolor=TEAL, linewidth=2.5)
    )
    ax.text(
        7.43,
        3.7,
        r"$P_5^{(3)}=\sum_{0\leq i<j<k\leq4}x_ix_jx_k$",
        ha="center",
        color=INK,
        fontsize=14,
    )
    for index, term in enumerate(terms):
        row, col = divmod(index, 5)
        x, y = 5.58 + col * 0.72, 3.05 - row * 0.56
        ax.text(
            x,
            y,
            "$x_{}x_{}x_{}$".format(*term),
            color=BLUE if row == 0 else GOLD,
            fontsize=9,
            ha="center",
        )
    ax.text(
        7.43,
        1.82,
        "one homogeneous symmetric block",
        ha="center",
        color=INK,
        fontsize=11,
    )
    for index, x in enumerate([6.9, 7.95]):
        ax.add_patch(Circle((x, 1.35), 0.22, facecolor=TEAL, edgecolor="white"))
        ax.text(x, 1.35, f"$w_{index}$", ha="center", va="center", color="white", fontsize=10)
    ax.text(
        7.43,
        0.62,
        "same values → only 2 auxiliary bits",
        ha="center",
        color=TEAL,
        fontsize=12,
        fontweight="semibold",
    )
    save(fig, "qubo-symmetric-reduction.png")


def draw_qubo_resources():
    fig, axes = plt.subplots(1, 2, figsize=(10, 5.625), facecolor=BG)
    fig.suptitle(
        "A smaller QUBO for graph colouring on the complete graph K₈",
        x=0.045,
        y=0.965,
        ha="left",
        color=INK,
        fontsize=18,
        fontweight="semibold",
    )
    labels = ["Monomial\nreduction", "Symmetric\nreduction"]
    colors = [CORAL, TEAL]
    values = [(1180, 480), (5849, 2797)]
    titles = ["Total binary variables", "Quadratic terms"]
    for ax, pair, title in zip(axes, values, titles):
        ax.set_facecolor(BG)
        bars = ax.bar(labels, pair, color=colors, width=0.58)
        ax.set_title(title, color=INK, fontsize=14, fontweight="semibold", pad=14)
        ax.spines[["top", "right", "left"]].set_visible(False)
        ax.spines["bottom"].set_color(GRID)
        ax.tick_params(axis="y", left=False, labelleft=False)
        ax.tick_params(axis="x", colors=INK)
        ax.set_ylim(0, max(pair) * 1.2)
        for bar, value in zip(bars, pair):
            ax.text(
                bar.get_x() + bar.get_width() / 2,
                value + max(pair) * 0.035,
                f"{value:,}",
                ha="center",
                color=INK,
                fontsize=13,
                fontweight="bold",
            )
    axes[0].text(
        0.5,
        0.84,
        "456 rather than 1,156\nauxiliary variables",
        transform=axes[0].transAxes,
        ha="center",
        color=TEAL,
        fontsize=11,
    )
    axes[1].text(
        0.5,
        0.84,
        "less than half as many\nquadratic terms",
        transform=axes[1].transAxes,
        ha="center",
        color=TEAL,
        fontsize=11,
    )
    fig.text(
        0.5,
        0.025,
        "Fewer variables and interactions make the graph-colouring model easier to embed on limited quantum hardware.",
        ha="center",
        color=INK,
        fontsize=10.5,
    )
    fig.tight_layout(rect=(0.04, 0.14, 0.98, 0.84), w_pad=3)
    path = OUT / "qubo-resource-comparison.png"
    fig.savefig(path, dpi=125, facecolor=BG)
    plt.close(fig)
    image = Image.open(path).convert("RGB")
    image.save(path, optimize=True)


def rauzy_points(iterations=16):
    rules = {0: [0, 1], 1: [0, 2], 2: [0]}
    word = [0]
    for _ in range(iterations):
        word = [symbol for item in word for symbol in rules[item]]
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
    return np.array(word), projected


def draw_rauzy_method():
    fig, ax = new_figure("From symbolic substitution to a Rauzy fractal")
    card(ax, 0.35, 1.0, 2.35, 3.15, "1. Grow a word", ["0 → 01", "1 → 02", "2 → 0"], BLUE)
    words = ["0", "01", "0102", "0102010", "0102010010201"]
    for index, word in enumerate(words):
        ax.text(0.62, 2.55 - index * 0.36, word, color=INK, fontsize=9.5)

    arrow(ax, (2.78, 2.58), (3.38, 2.58))

    card(
        ax,
        3.45,
        1.0,
        2.45,
        3.15,
        "2. Count prefixes",
        ["walk in a lattice:", "", "0: one step in x", "1: one step in y", "2: one step in z"],
        GOLD,
    )
    path = np.array([[0, 0], [0.4, 0], [0.4, 0.45], [0.8, 0.45], [0.8, 0.9], [1.2, 0.9], [1.5, 0.68]])
    ax.plot(3.88 + path[:, 0], 1.35 + path[:, 1], color=INK, lw=2.3)
    ax.scatter(3.88 + path[:, 0], 1.35 + path[:, 1], color=[BLUE, TEAL, BLUE, GOLD, BLUE, TEAL, GOLD], s=38)

    arrow(ax, (5.98, 2.58), (6.55, 2.58), "remove growth")

    word, points = rauzy_points()
    colors = np.array([BLUE, TEAL, CORAL])
    px = 8.05 + points[:, 0] * 1.25
    py = 2.58 + points[:, 1] * 1.55
    ax.scatter(px, py, c=colors[word], s=2.4, alpha=0.72, linewidths=0)
    ax.text(8.05, 4.55, "3. Project to the contracting plane", ha="center", color=INK, fontweight="bold")
    ax.text(
        8.05,
        0.63,
        "the expanding direction disappears;\nthe remaining fluctuations stay bounded",
        ha="center",
        color=INK,
        fontsize=10,
    )
    save(fig, "rauzy-method.png")


def draw_rauzy_self_similarity():
    fig, ax = new_figure("The fractal is assembled from smaller copies of itself")
    word, points = rauzy_points()
    palette = np.array([BLUE, TEAL, CORAL])

    ax.scatter(
        2.35 + points[:, 0] * 1.8,
        2.55 + points[:, 1] * 1.95,
        c=palette[word],
        s=3,
        alpha=0.78,
        linewidths=0,
    )
    ax.text(2.35, 4.58, "Complete domain", ha="center", color=INK, fontweight="bold")

    arrow(ax, (4.35, 2.55), (5.2, 2.55), "separate by symbol", TEAL)

    centres = [(6.35, 3.75), (8.25, 3.75), (7.3, 1.55)]
    for symbol, (cx, cy) in enumerate(centres):
        subset = points[word == symbol]
        subset = subset - subset.mean(axis=0)
        scale = max(np.ptp(subset[:, 0]), np.ptp(subset[:, 1]))
        ax.scatter(
            cx + subset[:, 0] / scale * 1.25,
            cy + subset[:, 1] / scale * 1.05,
            color=palette[symbol],
            s=3,
            alpha=0.8,
            linewidths=0,
        )
        ax.text(cx, cy - 0.86, f"piece {symbol}", ha="center", color=INK, fontsize=10)
    ax.text(
        7.3,
        0.38,
        "Each coloured piece is transformed and fitted back into the whole domain.",
        ha="center",
        color=INK,
        fontsize=10.5,
    )
    save(fig, "rauzy-self-similarity.png")


def hyperbolic_network(ax, centre, radius):
    cx, cy = centre
    rings = [0.0, 0.28, 0.53, 0.75, 0.9]
    counts = [1, 7, 14, 28, 42]
    ring_nodes = []
    for ring, count in zip(rings, counts):
        nodes = []
        for index in range(count):
            angle = 2 * np.pi * index / count + ring
            node = (cx + radius * ring * np.cos(angle), cy + radius * ring * np.sin(angle))
            nodes.append(node)
        ring_nodes.append(nodes)
    for level in range(1, len(ring_nodes)):
        previous = ring_nodes[level - 1]
        current = ring_nodes[level]
        for index, node in enumerate(current):
            parent = previous[int(index * len(previous) / len(current))]
            ax.plot([parent[0], node[0]], [parent[1], node[1]], color=GRID, lw=0.7)
            neighbour = current[(index + 1) % len(current)]
            ax.plot([node[0], neighbour[0]], [node[1], neighbour[1]], color=GRID, lw=0.6)
    for nodes in ring_nodes:
        if nodes:
            ax.scatter(*np.array(nodes).T, s=10, color=TEAL, zorder=2)
    ax.add_patch(Circle((cx, cy), radius, fill=False, edgecolor=INK, linewidth=2.2))


def draw_hyperbolic_local_global():
    fig, ax = new_figure("Local checks detect errors; global cycles store information")
    hyperbolic_network(ax, (3.0, 2.55), 2.0)
    local = np.array([[2.42, 2.25], [2.78, 2.65], [3.18, 2.42], [2.95, 2.02]])
    ax.add_patch(Polygon(local, closed=True, fill=False, edgecolor=CORAL, linewidth=3))
    ax.annotate(
        "local stabilizer check",
        xy=(2.82, 2.35),
        xytext=(0.3, 4.4),
        arrowprops={"arrowstyle": "->", "color": CORAL},
        color=CORAL,
        fontsize=10,
    )

    theta = np.linspace(-1.25, 1.25, 180)
    ax.plot(3.0 + 1.88 * np.sin(theta), 2.55 + 0.72 * np.cos(theta), color=BLUE, lw=3.3)
    ax.annotate(
        "non-trivial logical cycle",
        xy=(3.0, 3.27),
        xytext=(0.35, 0.45),
        arrowprops={"arrowstyle": "->", "color": BLUE},
        color=BLUE,
        fontsize=10,
    )

    card(
        ax,
        5.65,
        2.85,
        3.8,
        1.45,
        "Error correction",
        ["short error chain", "→ syndrome at its boundary", "→ decoder proposes a correction"],
        CORAL,
    )
    card(
        ax,
        5.65,
        0.75,
        3.8,
        1.45,
        "Logical information",
        ["a loop that cannot shrink away", "crosses the surface globally", "and acts on an encoded qubit"],
        BLUE,
    )
    arrow(ax, (7.55, 2.72), (7.55, 2.3), "same tiling", TEAL)
    save(fig, "hyperbolic-surface-code-local-global.png")


def draw_hyperbolic_motivation():
    fig, ax = new_figure("Why place a surface code on hyperbolic geometry?")

    ax.text(2.3, 4.55, "Euclidean surface code", ha="center", color=INK, fontweight="bold")
    for row in range(6):
        for col in range(6):
            x, y = 0.85 + col * 0.55, 1.2 + row * 0.55
            ax.scatter(x, y, s=30, color=BLUE)
            if col < 5:
                ax.plot([x, x + 0.55], [y, y], color=GRID, lw=1.2)
            if row < 5:
                ax.plot([x, x], [y, y + 0.55], color=GRID, lw=1.2)
    ax.text(
        2.3,
        0.48,
        "excellent distance,\nbut only a small number of logical qubits",
        ha="center",
        color=INK,
        fontsize=10.5,
    )

    arrow(ax, (4.55, 2.65), (5.35, 2.65), "change the geometry", TEAL)

    hyperbolic_network(ax, (7.45, 2.65), 1.85)
    ax.text(7.45, 4.78, "Hyperbolic surface code", ha="center", color=INK, fontweight="bold")
    for angle in np.linspace(0, 2 * np.pi, 4, endpoint=False):
        ax.add_patch(
            Circle(
                (7.45 + 0.65 * np.cos(angle), 2.65 + 0.65 * np.sin(angle)),
                0.13,
                fill=False,
                edgecolor=CORAL,
                linewidth=2,
            )
        )
    ax.text(
        7.45,
        0.35,
        "many logical qubits grow with code size,\nwhile distance grows more slowly",
        ha="center",
        color=INK,
        fontsize=10.5,
    )
    ax.text(
        9.42,
        3.9,
        "constant\nencoding rate",
        ha="center",
        color=CORAL,
        fontsize=10,
        fontweight="semibold",
    )
    save(fig, "hyperbolic-surface-code-motivation.png")


def draw_timetabling_method():
    fig, ax = new_figure("How DATA turns university constraints into a working timetable")

    inputs = [
        ("Courses", "hours, sections,\npossible times"),
        ("Rooms", "capacity and\nspecial facilities"),
        ("Instructors", "availability,\npreferred days"),
        ("Students", "planned courses\nand overlaps"),
    ]
    colors = [BLUE, GOLD, CORAL, TEAL]
    for index, ((title, detail), color) in enumerate(zip(inputs, colors)):
        y = 3.85 - index * 0.95
        card(ax, 0.25, y, 2.15, 0.78, title, [detail], color)
        arrow(ax, (2.48, y + 0.38), (3.0, 2.6), color=color)

    card(
        ax,
        3.05,
        1.0,
        3.55,
        3.3,
        "Constraint search",
        [
            "1. place the most constrained course",
            "2. try an allowed time and room",
            "3. reject instructor, student,",
            "   and classroom conflicts",
            "4. backtrack when no option remains",
            "5. continue until every course fits",
        ],
        INK,
    )

    arrow(ax, (6.72, 2.6), (7.2, 2.6), "complete assignment", TEAL)

    ax.add_patch(Rectangle((7.25, 1.0), 2.4, 3.3, facecolor="white", edgecolor=TEAL, linewidth=2))
    ax.text(8.45, 4.02, "Reviewable timetable", ha="center", color=INK, fontweight="bold")
    schedule = [
        (0, 0, 1, BLUE),
        (1, 0, 2, GOLD),
        (2, 1, 1, TEAL),
        (0, 2, 2, CORAL),
        (2, 3, 1, BLUE),
    ]
    for col in range(3):
        for row in range(5):
            ax.add_patch(
                Rectangle(
                    (7.52 + col * 0.61, 1.48 + row * 0.42),
                    0.54,
                    0.35,
                    facecolor="white",
                    edgecolor=GRID,
                    linewidth=0.8,
                )
            )
    for col, row, span, color in schedule:
        ax.add_patch(
            Rectangle(
                (7.54 + col * 0.61, 1.5 + row * 0.42),
                0.5,
                0.31 + (span - 1) * 0.42,
                facecolor=color,
                edgecolor="white",
                alpha=0.88,
            )
        )
    ax.text(
        8.45,
        1.2,
        "inspect • move • undo • save",
        ha="center",
        color=INK,
        fontsize=9.5,
    )
    ax.text(
        5,
        0.25,
        "The algorithm produces the first complete assignment; the app supports review, requests, revision, and distribution.",
        ha="center",
        color=INK,
        fontsize=10.5,
    )
    save(fig, "automated-timetabling-method.png")


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    draw_qubo_motivation()
    draw_qubo_symmetric_reduction()
    draw_qubo_resources()
    draw_rauzy_method()
    draw_rauzy_self_similarity()
    draw_hyperbolic_local_global()
    draw_hyperbolic_motivation()
    draw_timetabling_method()
    print("Created 8 instructional topic images")


if __name__ == "__main__":
    main()
