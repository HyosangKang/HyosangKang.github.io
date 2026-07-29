#!/usr/bin/env python3
"""Create reproducible, topic-specific PNG illustrations for project pages."""

from __future__ import annotations

import hashlib
import math
import shutil
import sys
from pathlib import Path

import matplotlib

matplotlib.use("Agg")

import matplotlib.pyplot as plt
import networkx as nx
import numpy as np
from matplotlib.patches import Circle, Ellipse, FancyArrowPatch, Polygon, Rectangle
from mpl_toolkits.mplot3d import Axes3D  # noqa: F401
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


PROJECTS = {
    "calculus-calculator": ("Making a Calculator", "series"),
    "calculus-3d-graphing": ("3D Graphing", "surface"),
    "calculus-captcha": ("Making CAPTCHA", "captcha"),
    "calculus-polygon-area": ("Area of a Polygon", "polygon"),
    "calculus-curve-length": ("Length of a Curve", "curve"),
    "calculus-mnist": ("MNIST Data", "mnist"),
    "differential-equations-epidemic": ("Epidemic Simulation", "sir"),
    "differential-equations-stone-skipping": ("Stone-skipping Simulation", "skipping"),
    "differential-equations-projectile-motion": ("Projectile Motion", "projectile"),
    "differential-equations-orbital-motion": ("Orbital Motion", "orbit"),
    "differential-equations-ecosystem-populations": (
        "Ecosystem Population Dynamics",
        "predator-prey",
    ),
    "differential-equations-spring": ("Spring Oscillation", "spring"),
    "differential-equations-double-spring": ("Double-spring Motion", "double-spring"),
    "differential-equations-electric-circuit": ("Electric-circuit Simulation", "circuit"),
    "differential-equations-lanchester": ("Lanchester’s Law", "lanchester"),
    "linear-algebra-aes": ("AES Encryption", "aes"),
    "linear-algebra-3d-visualization": ("3D Visualization", "rotation"),
    "linear-algebra-pca": ("Image Classification with PCA", "pca"),
    "quantum-computing-simulator": ("Quantum Simulator", "bloch"),
    "quantum-computing-shor": ("Shor’s Algorithm", "period"),
    "quantum-computing-bb84": ("Quantum Key Distribution: BB84", "bb84"),
    "axiomatic-number-systems-mathematical-objects": (
        "Axiomatic Construction of Number Systems",
        "number-sets",
    ),
    "bernstein-vazirani-encoding": ("Bernstein–Vazirani Encoding", "bv"),
    "busy-beaver-quantum-algorithms": ("Busy Beaver Computation", "automaton"),
    "closed-geodesics-modular-surfaces": ("Closed Geodesics on Modular Surfaces", "hyperbolic"),
    "course-registration-school-timetabling": (
        "Course Registration and School Timetabling",
        "timetable",
    ),
    "covid-19-spread-modelling": ("COVID-19 Spread Modelling", "sir"),
    "escherian-images-hyperbolic-tilings": ("Escherian Hyperbolic Tilings", "hyperbolic-tiling"),
    "expander-graphs-pseudorandom-numbers": (
        "Expander Graphs for Pseudorandom Numbers",
        "expander",
    ),
    "fair-allocation-shared-resources": ("Fair Allocation of Shared Resources", "allocation"),
    "fair-division-rental-harmony-sperners-lemma": (
        "Rental Harmony and Sperner’s Lemma",
        "sperner",
    ),
    "gaussian-functions-protein-folding": ("Gaussian Functions in Protein Folding", "protein"),
    "graph-colouring-quantum-annealer": ("Graph Colouring with a Quantum Annealer", "annealing"),
    "graph-colouring-scheduling": ("Graph-colouring Algorithms for Scheduling", "colouring"),
    "graph-states-stabilizer-eigenstates": ("Graph States and Stabilizer Eigenstates", "graph-state"),
    "grovers-algorithm": ("Grover’s Algorithm", "grover"),
    "kanoodle-regular-polyhedra": ("Kanoodle Constructions of Regular Polyhedra", "polyhedron"),
    "lagrangian-geodesics-euler-lagrange-equations": (
        "Lagrangian Geodesics and Euler–Lagrange Equations",
        "geodesic-surface",
    ),
    "manifolds-natural-social-phenomena": ("Manifolds in Natural and Social Phenomena", "manifold"),
    "markov-numbers-uniqueness-conjecture": ("Markov Numbers and the Uniqueness Conjecture", "markov"),
    "n-body-dynamics-curved-spaces": ("N-body Dynamics on Curved Spaces", "nbody"),
    "neural-network-decoders-surface-codes": (
        "Neural-network Decoders for Surface Codes",
        "decoder",
    ),
    "programmable-mathematics": ("Programmable Mathematics", "programmable"),
    "quantum-annealing-combinatorial-optimization": (
        "Quantum Annealing for Combinatorial Optimization",
        "annealing",
    ),
    "quantum-coding-blocks": ("Quantum Coding Blocks", "code-blocks"),
    "quantum-error-correction-game": ("Quantum Error-correction Game", "qec-game"),
    "quantum-error-correction": ("Quantum Error Correction", "qec"),
    "quantum-simulators": ("Quantum Simulators", "bloch"),
    "rhythm-game-difficulty-classification": (
        "Rhythm-game Difficulty Classification",
        "rhythm",
    ),
    "rubiks-cube-turing-complete-computer": (
        "Rubik’s Cube as a Turing-complete Computer",
        "rubik",
    ),
    "school-building-energy-efficiency": ("School-building Energy Efficiency", "building"),
    "shors-algorithm": ("Shor’s Algorithm", "period"),
    "simons-algorithm": ("Simon’s Algorithm", "simon"),
    "simple-closed-geodesics-gamma-3-surfaces": (
        "Simple Closed Geodesics on Γ(3) Surfaces",
        "gamma",
    ),
    "smart-grid-simulation-optimization": ("Smart-grid Simulation and Optimization", "smart-grid"),
    "visual-tools-hyperbolic-geometry": ("Visual Tools for Hyperbolic Geometry", "hyperbolic-tiling"),
}


SOURCE_IMAGES = {
    "calculus-captcha-source.png": Path(
        "/Users/hk/github.com/book-engineering-math/sources/multivariable_calculus/original/img/ch-2/captcha.png"
    ),
    "calculus-3d-graphing-source.png": Path(
        "/Users/hk/github.com/book-advanced-calculus/archive/book-english/image/graph-sincos.png"
    ),
    "differential-equations-electric-circuit-source.png": Path(
        "/Users/hk/github.com/book-engineering-math/sources/engineering_math/img/125-6.png"
    ),
    "rauzy-fractal-1.png": Path("/Users/hk/github.com/rauzy/rauzy.png"),
    "rauzy-fractal-2.png": Path("/Users/hk/github.com/rauzy/rauzy4.png"),
    "qubo-1.png": Path("/Users/hk/github.com/qubo/img/graph_0.75.png"),
    "qubo-2.png": Path("/Users/hk/github.com/qubo/img/graph_0.90.png"),
    "qubo-3.png": Path("/Users/hk/github.com/qubo/img/graph_1.00.png"),
    "hyperbolic-surface-code-1.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/front-page-surface-code-hyperbolic-trace-and-vinberg-polytope.png"
    ),
    "hyperbolic-surface-code-2.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/2_3_7/hyperbolic_surface_code_layer_0.png"
    ),
    "hyperbolic-surface-code-3.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/2_3_7/threshold_02_03_07.png"
    ),
    "quantum-error-correction-1.png": Path(
        "/Users/hk/github.com/quantum-error-correction/ex-surf-code.png"
    ),
    "quantum-error-correction-2.png": Path(
        "/Users/hk/github.com/quantum-error-correction/ex-err-corr.png"
    ),
    "quantum-error-correction-3.png": Path(
        "/Users/hk/github.com/quantum-error-correction/ex-result.png"
    ),
    "neural-network-decoders-surface-codes-source.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/surface_code_threshold_30000_shots_per_point_with_zoom.png"
    ),
    "graph-states-stabilizer-eigenstates-source.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/xstab-diag.png"
    ),
    "visual-tools-hyperbolic-geometry-source.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/2_3_7/2_3_7_1_original_poinc.png"
    ),
    "escherian-images-hyperbolic-tilings-source.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/2_4_5/2_4_5_1_original_poinc.png"
    ),
    "closed-geodesics-modular-surfaces-source.png": Path(
        "/Users/hk/github.com/hyperbolic-surface-code/presentation/slides/img/2_3_7/2_3_7_triangle_poinc.png"
    ),
}


def rng_for(key: str) -> np.random.Generator:
    seed = int(hashlib.sha256(key.encode()).hexdigest()[:8], 16)
    return np.random.default_rng(seed)


def base_figure(title: str, projection: str | None = None):
    fig = plt.figure(figsize=(10, 5.625), facecolor=BG)
    ax = fig.add_subplot(111, projection=projection)
    ax.set_facecolor(BG)
    fig.text(0.045, 0.93, title, color=INK, fontsize=18, fontweight="semibold")
    return fig, ax


def finish(fig, ax, key: str, keep_axes: bool = False):
    if not keep_axes:
        ax.set_axis_off()
    fig.tight_layout(rect=(0.025, 0.025, 0.98, 0.89))
    path = OUT / f"{key}.png"
    fig.savefig(path, dpi=125, facecolor=BG)
    plt.close(fig)
    image = Image.open(path).convert("RGB")
    image.save(path, optimize=True)


def draw_series(key, title):
    fig, ax = base_figure(title)
    x = np.linspace(-2 * np.pi, 2 * np.pi, 500)
    exact = np.sin(x)
    approx = x - x**3 / math.factorial(3) + x**5 / math.factorial(5) - x**7 / math.factorial(7)
    ax.plot(x, exact, color=INK, lw=3, label="sin(x)")
    ax.plot(x, approx, color=CORAL, lw=2.5, ls="--", label="arithmetic approximation")
    ax.axhline(0, color=GRID, lw=1)
    ax.set_ylim(-2.2, 2.2)
    ax.set_xlabel("input x")
    ax.set_ylabel("value")
    ax.legend(frameon=False, loc="upper right")
    finish(fig, ax, key, True)


def draw_surface(key, title):
    fig, ax = base_figure(title, "3d")
    x = np.linspace(-5, 5, 90)
    y = np.linspace(-5, 5, 90)
    x, y = np.meshgrid(x, y)
    r = np.sqrt(x**2 + y**2) + 1e-9
    z = np.sin(r) / r
    ax.plot_surface(x, y, z, cmap="viridis", linewidth=0, antialiased=True, alpha=0.95)
    ax.set_xlabel("x")
    ax.set_ylabel("y")
    ax.set_zlabel("f(x, y)")
    ax.view_init(28, -55)
    finish(fig, ax, key, True)


def draw_captcha(key, title):
    fig, ax = base_figure(title)
    r = rng_for(key)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 5)
    for _ in range(16):
        ax.plot([r.uniform(0, 10), r.uniform(0, 10)], [r.uniform(0, 5), r.uniform(0, 5)],
                color=r.choice([TEAL, BLUE, GOLD]), alpha=0.35, lw=r.uniform(1, 3))
    for i, ch in enumerate("7K4M"):
        ax.text(1.6 + i * 1.9, 2.45 + r.uniform(-0.2, 0.2), ch, fontsize=62,
                rotation=r.uniform(-22, 22), color=INK, ha="center", va="center",
                fontweight="bold")
    ax.text(5, 0.35, "curves + transformations + controlled noise", ha="center", color=CORAL)
    finish(fig, ax, key)


def draw_polygon(key, title):
    fig, ax = base_figure(title)
    pts = np.array([[0.7, 0.8], [2.5, 0.25], [4.2, 1.0], [4.7, 2.8], [3.2, 4.0], [1.0, 3.4]])
    colors = [TEAL, GOLD, BLUE, CORAL]
    for i in range(1, len(pts) - 1):
        ax.add_patch(Polygon([pts[0], pts[i], pts[i + 1]], closed=True,
                             facecolor=colors[(i - 1) % len(colors)], alpha=0.38, edgecolor="white"))
    closed = np.vstack([pts, pts[0]])
    ax.plot(closed[:, 0], closed[:, 1], color=INK, lw=3)
    ax.scatter(pts[:, 0], pts[:, 1], s=80, color=INK, zorder=3)
    ax.text(6.0, 2.8, "signed edge sums", color=INK, fontsize=16)
    ax.text(6.0, 2.1, r"$A=\frac{1}{2}\left|\sum x_i y_{i+1}-y_i x_{i+1}\right|$",
            color=CORAL, fontsize=18)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 4.8)
    ax.set_aspect("equal")
    finish(fig, ax, key)


def draw_curve(key, title):
    fig, ax = base_figure(title)
    x = np.linspace(0, 8, 500)
    y = 0.45 * np.sin(1.3 * x) + 0.08 * x**1.5 + 1
    ax.plot(x, y, color=INK, lw=3, label="curve")
    for count, color, label in [(7, CORAL, "7 segments"), (18, TEAL, "18 segments")]:
        xp = np.linspace(0, 8, count)
        yp = 0.45 * np.sin(1.3 * xp) + 0.08 * xp**1.5 + 1
        ax.plot(xp, yp, marker="o", color=color, lw=1.8, alpha=0.85, label=label)
    ax.set_xlabel("x")
    ax.set_ylabel("y")
    ax.legend(frameon=False)
    finish(fig, ax, key, True)


def draw_mnist(key, title):
    fig, ax = base_figure(title)
    digit = np.zeros((14, 14))
    for t in np.linspace(0, 2 * np.pi, 150):
        x = int(np.clip(7 + 3.5 * np.cos(t), 0, 13))
        y = int(np.clip(6.5 + 5 * np.sin(t), 0, 13))
        digit[y, x] = 1
    for _ in range(2):
        digit = np.maximum.reduce([digit, np.roll(digit, 1, 0), np.roll(digit, -1, 0),
                                   np.roll(digit, 1, 1), np.roll(digit, -1, 1)])
    ax.imshow(digit, cmap="Blues", extent=(0, 4.5, 0, 4.5), origin="lower")
    layers = [(5.8, 7), (7.3, 5), (8.7, 3)]
    for (x0, n), (x1, m) in zip(layers[:-1], layers[1:]):
        ys0, ys1 = np.linspace(0.7, 3.8, n), np.linspace(1.0, 3.5, m)
        for y0 in ys0:
            for y1 in ys1:
                ax.plot([x0, x1], [y0, y1], color=GRID, lw=0.6)
    for x0, n in layers:
        ax.scatter(np.full(n, x0), np.linspace(0.7 if n == 7 else 1.0, 3.8 if n == 7 else 3.5, n),
                   s=110, color=TEAL if n > 3 else CORAL, edgecolor="white", zorder=3)
    ax.text(9.45, 2.3, "0–9", color=INK, fontsize=18, fontweight="bold")
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 4.5)
    finish(fig, ax, key)


def draw_sir(key, title):
    fig, ax = base_figure(title)
    t = np.linspace(0, 80, 800)
    s, i, r = np.empty_like(t), np.empty_like(t), np.empty_like(t)
    s[0], i[0], r[0] = 0.99, 0.01, 0
    beta, gamma = (0.28, 0.085) if "covid" in key else (0.34, 0.11)
    dt = t[1] - t[0]
    for k in range(len(t) - 1):
        inf = beta * s[k] * i[k]
        s[k + 1] = s[k] - dt * inf
        i[k + 1] = i[k] + dt * (inf - gamma * i[k])
        r[k + 1] = r[k] + dt * gamma * i[k]
    ax.plot(t, s, color=BLUE, lw=3, label="susceptible")
    ax.plot(t, i, color=CORAL, lw=3, label="infected")
    ax.plot(t, r, color=TEAL, lw=3, label="recovered")
    peak = np.argmax(i)
    ax.scatter([t[peak]], [i[peak]], color=CORAL, s=80, zorder=4)
    ax.annotate("peak", (t[peak], i[peak]), xytext=(t[peak] + 8, i[peak] + 0.12),
                arrowprops={"arrowstyle": "->", "color": INK}, color=INK)
    ax.set_xlabel("time")
    ax.set_ylabel("fraction of population")
    ax.set_ylim(0, 1.03)
    ax.legend(frameon=False, ncol=3, loc="upper center")
    finish(fig, ax, key, True)


def draw_skipping(key, title):
    fig, ax = base_figure(title)
    dt = 0.005
    theta = np.pi / 200
    x, y, vx, vy = 0.0, 0.5, 150.0, 0.5
    xs, ys, contacts = [x], [y], []
    in_water = False
    for _ in range(5000):
        ax_force, ay_force = 0.0, -10.0
        if y > 0:
            in_water = False
        else:
            if not in_water:
                contacts.append(x)
                in_water = True
            if y + 5 * np.sin(theta) > 0:
                submerged = min(abs(y), 5) / np.sin(theta)
                speed_squared = vx * vx + vy * vy
                ay_force = -10 + 0.5 * speed_squared * submerged * (
                    0.1 * np.cos(theta) - 0.3 * np.sin(theta)
                )
                ax_force = -0.5 * speed_squared * submerged * (
                    0.1 * np.sin(theta) + 0.3 * np.cos(theta)
                )
            else:
                ax_force = -0.3 * vx * vx
                ay_force = -10 + 0.3 * vy * vy
        ax_force = np.clip(ax_force, -7000, 7000)
        ay_force = np.clip(ay_force, -7000, 7000)
        vx += ax_force * dt
        vy += ay_force * dt
        x += vx * dt
        y += vy * dt
        if y > 0.04:
            in_water = False
        xs.append(x)
        ys.append(y)
        if y < -2 or vx < 3:
            break

    xs, ys = np.array(xs), np.array(ys)
    ax.fill_between([0, max(xs) * 1.02], -1.5, 0, color=BLUE, alpha=0.18)
    ax.axhline(0, color=BLUE, lw=2.6)
    ax.plot(xs, ys, color=CORAL, lw=2.6, label="computed stone trajectory")
    ax.scatter(contacts, np.zeros(len(contacts)), s=30, color=TEAL, edgecolor="white", zorder=4,
               label="water contacts")
    for i, contact_x in enumerate(contacts[:8]):
        radius = 2.2 + i * 0.8
        ax.add_patch(Ellipse((contact_x, 0), width=2 * radius, height=0.055 + i * 0.008,
                             fill=False, color=BLUE, alpha=0.38, lw=1.0))
    ax.add_patch(Polygon([[4, 0.48], [18, 0.7], [25, 0.56], [11, 0.42]],
                         color=INK, alpha=0.9))
    ax.annotate("gravity in air", xy=(60, 0.28), xytext=(72, 0.9),
                arrowprops={"arrowstyle": "->", "color": INK}, color=INK)
    ax.annotate("lift + drag while partly immersed", xy=(contacts[4], 0),
                xytext=(contacts[4] + 28, -0.85),
                arrowprops={"arrowstyle": "->", "color": TEAL}, color=TEAL)
    ax.set_xlim(0, max(xs) * 1.03)
    ax.set_ylim(-1.45, 1.15)
    ax.set_xlabel("distance")
    ax.set_ylabel("height above water")
    ax.legend(frameon=False, loc="upper right")
    finish(fig, ax, key, True)


def draw_projectile(key, title):
    fig, ax = base_figure(title)
    gravity = 9.8
    speed = 34
    angle = np.deg2rad(38)
    flight_time = 2 * speed * np.sin(angle) / gravity
    t = np.linspace(0, flight_time, 240)
    x = speed * np.cos(angle) * t
    y = speed * np.sin(angle) * t - 0.5 * gravity * t**2
    ax.fill_between(x, 0, y, color=BLUE, alpha=0.09)
    ax.plot(x, y, color=CORAL, lw=3, label="projectile path")
    ax.add_patch(FancyArrowPatch((0, 0), (11 * np.cos(angle), 11 * np.sin(angle)),
                                 arrowstyle="-|>", mutation_scale=18, color=TEAL, lw=3))
    ax.text(8, 6.7, "initial velocity", color=TEAL, fontsize=14)
    apex = np.argmax(y)
    ax.scatter([x[apex]], [y[apex]], color=GOLD, s=85, edgecolor="white", zorder=4)
    ax.annotate("vertical velocity = 0", (x[apex], y[apex]), xytext=(x[apex] + 8, y[apex] + 2),
                arrowprops={"arrowstyle": "->", "color": INK}, color=INK)
    ax.axhline(0, color=INK, lw=1.3)
    ax.set_xlabel("horizontal distance")
    ax.set_ylabel("height")
    ax.set_ylim(-1, max(y) * 1.23)
    finish(fig, ax, key, True)


def draw_orbit(key, title):
    fig, ax = base_figure(title)
    ax.add_patch(Circle((0, 0), 0.48, facecolor=BLUE, edgecolor="white", lw=2, zorder=5))
    ax.add_patch(Circle((-0.13, -0.08), 0.2, facecolor=TEAL, edgecolor="none", alpha=0.9, zorder=6))
    colors = [CORAL, GOLD, TEAL]
    labels = ["too slow", "near orbit", "escape path"]
    velocities = [0.72, 1.0, 1.42]
    for initial_velocity, c, label in zip(velocities, colors, labels):
        position = np.array([2.2, 0.0], dtype=float)
        velocity = np.array([0.0, initial_velocity], dtype=float)
        path = [position.copy()]
        for _ in range(1300):
            radius = np.linalg.norm(position)
            acceleration = -position / radius**3
            velocity += acceleration * 0.015
            position += velocity * 0.015
            path.append(position.copy())
            if radius < 0.5 or radius > 5.2:
                break
        path = np.array(path)
        ax.plot(path[:, 0], path[:, 1], color=c, lw=2.4, label=label)
    ax.add_patch(FancyArrowPatch((2.2, 0), (2.2, 1.1), arrowstyle="-|>",
                                 mutation_scale=16, color=INK, lw=2.5))
    ax.text(2.35, 0.55, "launch velocity", color=INK, fontsize=13)
    ax.set_aspect("equal")
    ax.set_xlim(-4.1, 4.8)
    ax.set_ylim(-3.4, 3.8)
    ax.legend(frameon=False, loc="upper left")
    finish(fig, ax, key)


def draw_predator_prey(key, title):
    fig, ax = base_figure(title)
    dt = 0.015
    prey, predator = 2.2, 0.7
    prey_values, predator_values = [prey], [predator]
    for _ in range(1500):
        dprey = 1.05 * prey - 0.58 * prey * predator
        dpredator = 0.34 * prey * predator - 0.82 * predator
        prey += dprey * dt
        predator += dpredator * dt
        prey_values.append(prey)
        predator_values.append(predator)
    time = np.arange(len(prey_values)) * dt
    ax.plot(time, prey_values, color=TEAL, lw=2.7, label="prey population")
    ax.plot(time, predator_values, color=CORAL, lw=2.7, label="predator population")
    ax.fill_between(time, prey_values, color=TEAL, alpha=0.08)
    ax.fill_between(time, predator_values, color=CORAL, alpha=0.08)
    ax.annotate("predators rise after prey", xy=(6.8, predator_values[int(6.8 / dt)]),
                xytext=(8.4, max(predator_values) * 0.88),
                arrowprops={"arrowstyle": "->", "color": INK}, color=INK)
    ax.set_xlabel("time")
    ax.set_ylabel("population")
    ax.legend(frameon=False, ncol=2, loc="upper right")
    finish(fig, ax, key, True)


def spring_curve(t, damping, frequency):
    return np.exp(-damping * t) * np.cos(frequency * t)


def draw_spring(key, title, double=False):
    fig, ax = base_figure(title)
    ax.set_xlim(0, 10)
    ax.set_ylim(-2.4, 2.4)
    t = np.linspace(0, 9.5, 500)
    if double:
        y1 = np.cos(2.1 * t) * np.cos(0.22 * t)
        y2 = np.sin(2.1 * t) * np.sin(0.22 * t)
        ax.plot(t, y1, color=BLUE, lw=2.6, label="mass 1")
        ax.plot(t, y2, color=CORAL, lw=2.6, label="mass 2")
        ax.legend(frameon=False, ncol=2)
    else:
        ax.plot(t, spring_curve(t, 0.13, 2.4), color=INK, lw=3, label="damped")
        ax.plot(t, 0.78 * np.sin(2.4 * t), color=CORAL, lw=2, alpha=0.8, label="undamped")
        ax.plot(t, 0.45 * np.sin(1.05 * t), color=TEAL, lw=2, ls="--", label="forcing")
        ax.legend(frameon=False, ncol=3)
    ax.axhline(0, color=GRID, lw=1)
    ax.set_xlabel("time")
    ax.set_ylabel("displacement")
    finish(fig, ax, key, True)


def draw_circuit(key, title):
    fig, ax = base_figure(title)
    ax.plot([1, 3], [3.7, 3.7], color=INK, lw=3)
    x = np.linspace(3, 5.4, 13)
    y = 3.7 + 0.32 * np.where(np.arange(13) % 2 == 0, -1, 1)
    ax.plot(x, y, color=INK, lw=3)
    ax.plot([5.4, 8, 8], [3.7, 3.7, 1], color=INK, lw=3)
    ax.plot([8, 1, 1], [1, 1, 3.7], color=INK, lw=3)
    for yy in [2.75, 2.35, 1.95]:
        ax.add_patch(Circle((1.1, yy), 0.45, fill=False, color=BLUE, lw=2))
    ax.plot([7.65, 8.35], [2.7, 2.7], color=TEAL, lw=4)
    ax.plot([7.65, 8.35], [2.35, 2.35], color=TEAL, lw=4)
    ax.text(4.0, 4.35, "R", color=CORAL, fontsize=20, ha="center")
    ax.text(0.25, 2.25, "L", color=BLUE, fontsize=20)
    ax.text(8.55, 2.4, "C", color=TEAL, fontsize=20)
    ax.text(3.0, 0.35, r"$Lq''+Rq'+\frac{1}{C}q=V(t)$", color=INK, fontsize=20)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 5)
    finish(fig, ax, key)


def draw_lanchester(key, title):
    fig, ax = base_figure(title)
    t = np.linspace(0, 10, 250)
    blue = np.maximum(0, 100 * np.cos(0.105 * t) - 8 * t)
    red = np.maximum(0, 85 * np.cos(0.12 * t) - 6.5 * t)
    ax.plot(t, blue, color=BLUE, lw=3, label="force A")
    ax.plot(t, red, color=CORAL, lw=3, label="force B")
    ax.fill_between(t, blue, alpha=0.12, color=BLUE)
    ax.fill_between(t, red, alpha=0.12, color=CORAL)
    ax.set_xlabel("time")
    ax.set_ylabel("remaining strength")
    ax.legend(frameon=False)
    finish(fig, ax, key, True)


def draw_aes(key, title):
    fig, ax = base_figure(title)
    r = rng_for(key)
    state = r.integers(0, 256, (4, 4))
    for offset, label in [(0.4, "state"), (5.6, "mixed state")]:
        for i in range(4):
            for j in range(4):
                value = state[i, j] if offset < 1 else state[(i + j) % 4, j] ^ (17 * (i + 1))
                ax.add_patch(Rectangle((offset + j * 0.8, 0.7 + (3 - i) * 0.8), 0.72, 0.72,
                                       facecolor=PALE if (i + j) % 2 else TEAL, alpha=0.8,
                                       edgecolor="white"))
                ax.text(offset + j * 0.8 + 0.36, 1.06 + (3 - i) * 0.8, f"{value:02X}",
                        ha="center", va="center", fontsize=10, color=INK)
        ax.text(offset + 1.55, 4.3, label, ha="center", color=INK, fontsize=15)
    ax.add_patch(FancyArrowPatch((4.0, 2.3), (5.25, 2.3), arrowstyle="-|>",
                                 mutation_scale=18, color=CORAL, lw=2.5))
    ax.text(4.62, 2.65, "round", color=CORAL, ha="center")
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 5)
    finish(fig, ax, key)


def draw_rotation(key, title):
    fig, ax = base_figure(title, "3d")
    vertices = np.array([[x, y, z] for x in [-1, 1] for y in [-1, 1] for z in [-1, 1]])
    theta = np.deg2rad(32)
    phi = np.deg2rad(20)
    rz = np.array([[np.cos(theta), -np.sin(theta), 0], [np.sin(theta), np.cos(theta), 0], [0, 0, 1]])
    ry = np.array([[np.cos(phi), 0, np.sin(phi)], [0, 1, 0], [-np.sin(phi), 0, np.cos(phi)]])
    rotated = vertices @ (rz @ ry).T
    for pts, color, alpha in [(vertices, GRID, 0.6), (rotated, BLUE, 1.0)]:
        for i, a in enumerate(pts):
            for j, b in enumerate(pts):
                if j > i and np.sum(np.abs(vertices[i] - vertices[j])) == 2:
                    ax.plot(*zip(a, b), color=color, lw=2.3, alpha=alpha)
    ax.scatter(rotated[:, 0], rotated[:, 1], rotated[:, 2], s=42, color=CORAL)
    ax.set_xlabel("x")
    ax.set_ylabel("y")
    ax.set_zlabel("z")
    ax.view_init(22, -58)
    finish(fig, ax, key, True)


def draw_pca(key, title):
    fig, ax = base_figure(title)
    r = rng_for(key)
    angle = np.deg2rad(32)
    rot = np.array([[np.cos(angle), -np.sin(angle)], [np.sin(angle), np.cos(angle)]])
    pts = r.normal(size=(240, 2)) @ np.diag([2.4, 0.65]) @ rot.T
    ax.scatter(pts[:, 0], pts[:, 1], s=18, color=BLUE, alpha=0.38)
    mean = pts.mean(0)
    for scale, color, label in [(3.2, CORAL, "PC 1"), (1.5, TEAL, "PC 2")]:
        vec = (rot[:, 0] if label == "PC 1" else rot[:, 1]) * scale
        ax.add_patch(FancyArrowPatch(mean - vec, mean + vec, arrowstyle="-|>",
                                     mutation_scale=18, color=color, lw=3))
        ax.text(*(mean + vec * 1.05), label, color=color, fontsize=14)
    ax.set_aspect("equal")
    ax.set_xlabel("pixel feature 1")
    ax.set_ylabel("pixel feature 2")
    finish(fig, ax, key, True)


def draw_bloch(key, title):
    fig, ax = base_figure(title, "3d")
    u = np.linspace(0, 2 * np.pi, 70)
    v = np.linspace(0, np.pi, 35)
    x = np.outer(np.cos(u), np.sin(v))
    y = np.outer(np.sin(u), np.sin(v))
    z = np.outer(np.ones_like(u), np.cos(v))
    ax.plot_wireframe(x, y, z, rstride=7, cstride=5, color=GRID, lw=0.6)
    vec = np.array([0.6, 0.48, 0.64])
    vec /= np.linalg.norm(vec)
    ax.quiver(0, 0, 0, *vec, color=CORAL, lw=3, arrow_length_ratio=0.15)
    ax.text(0, 0, 1.14, "|0⟩", ha="center", color=INK)
    ax.text(0, 0, -1.2, "|1⟩", ha="center", color=INK)
    ax.text(*(vec * 1.15), "|ψ⟩", color=CORAL, fontsize=15)
    ax.set_box_aspect((1, 1, 1))
    ax.view_init(20, -55)
    finish(fig, ax, key)


def draw_period(key, title):
    fig, ax = base_figure(title)
    x = np.arange(0, 31)
    period = 6 if "shor" in key else 5
    y = np.mod(2**x, 21 if period == 6 else 15)
    ax.vlines(x, 0, y, color=BLUE, lw=1.5, alpha=0.7)
    ax.scatter(x, y, c=np.mod(x, period), cmap="viridis", s=70, edgecolor="white", zorder=3)
    for start in range(0, 31, period):
        ax.axvspan(start - 0.35, start + period - 0.65, color=GOLD, alpha=0.08)
    ax.annotate(f"period r = {period}", xy=(period, y[period]), xytext=(period + 3, max(y) * 0.83),
                arrowprops={"arrowstyle": "->", "color": CORAL}, color=CORAL, fontsize=15)
    ax.set_xlabel("exponent x")
    ax.set_ylabel("modular value")
    finish(fig, ax, key, True)


def draw_bb84(key, title):
    fig, ax = base_figure(title)
    bits = [1, 0, 1, 1, 0, 0, 1, 0]
    bases = ["×", "+", "+", "×", "+", "×", "×", "+"]
    measured = [1, 0, 0, 1, 0, 1, 1, 0]
    for i, (bit, basis, result) in enumerate(zip(bits, bases, measured)):
        x = i + 1
        ax.add_patch(Circle((x, 3.3), 0.28, color=TEAL if bit else BLUE))
        ax.text(x, 3.3, str(bit), ha="center", va="center", color="white", fontweight="bold")
        ax.text(x, 2.4, basis, ha="center", va="center", color=INK, fontsize=20)
        ax.add_patch(Circle((x, 1.35), 0.28, fill=False, lw=2.5,
                            edgecolor=CORAL if result != bit else TEAL))
        ax.text(x, 1.35, str(result), ha="center", va="center", color=INK)
        if result == bit:
            ax.plot([x - 0.22, x + 0.22], [0.65, 0.65], color=TEAL, lw=4)
    ax.text(0.1, 3.3, "sent", color=INK, va="center")
    ax.text(0.1, 2.4, "basis", color=INK, va="center")
    ax.text(0.1, 1.35, "read", color=INK, va="center")
    ax.text(4.5, 0.28, "matching bases keep the shared key bits", ha="center", color=TEAL)
    ax.set_xlim(-0.5, 9)
    ax.set_ylim(0, 4.2)
    finish(fig, ax, key)


def draw_number_sets(key, title):
    fig, ax = base_figure(title)
    sets = [(4.7, 2.35, 4.2, "ℂ", PALE), (4.7, 2.35, 3.35, "ℝ", "#cce7e3"),
            (4.7, 2.35, 2.5, "ℚ", "#f4dfb4"), (4.7, 2.35, 1.65, "ℤ", "#f4c9bd"),
            (4.7, 2.35, 0.8, "ℕ", "#b9d5f1")]
    for x, y, radius, label, color in sets:
        ax.add_patch(Circle((x, y), radius, facecolor=color, edgecolor="white", lw=2))
        ax.text(x, y + radius - 0.35, label, ha="center", va="top", color=INK,
                fontsize=18, fontweight="bold")
    ax.text(8.4, 3.6, "objects built from", color=INK, fontsize=14)
    ax.text(8.4, 3.15, "definitions + axioms", color=CORAL, fontsize=16, fontweight="bold")
    ax.set_xlim(0, 10)
    ax.set_ylim(-2.1, 6.8)
    ax.set_aspect("equal")
    finish(fig, ax, key)


def draw_quantum_circuit(key, title, algorithm):
    fig, ax = base_figure(title)
    n = 4
    for q in range(n):
        y = 3.7 - q * 0.8
        ax.plot([0.8, 9.1], [y, y], color=INK, lw=1.5)
        ax.text(0.35, y, f"q{q}", va="center", color=INK)
    gate_color = {"bv": TEAL, "simon": BLUE, "code-blocks": CORAL}.get(algorithm, TEAL)
    for q in range(n):
        ax.add_patch(Rectangle((1.3, 3.42 - q * 0.8), 0.55, 0.55,
                               facecolor=PALE, edgecolor=INK))
        ax.text(1.575, 3.7 - q * 0.8, "H", ha="center", va="center", color=INK)
    ax.add_patch(Rectangle((3.1, 1.05), 2.2, 3.0, facecolor=gate_color, alpha=0.18,
                           edgecolor=gate_color, lw=2))
    oracle = {"bv": "Uₛ  hidden string", "simon": "Uƒ  paired inputs", "code-blocks": "ENCODE"}.get(algorithm, "U")
    ax.text(4.2, 2.55, oracle, ha="center", va="center", color=INK, fontsize=15)
    for q in range(n - 1):
        ax.add_patch(Rectangle((6.1, 3.42 - q * 0.8), 0.55, 0.55,
                               facecolor=PALE, edgecolor=INK))
        ax.text(6.375, 3.7 - q * 0.8, "H", ha="center", va="center", color=INK)
    for q in range(n - 1):
        y = 3.7 - q * 0.8
        ax.add_patch(Rectangle((7.7, y - 0.25), 0.5, 0.5, fill=False, edgecolor=CORAL, lw=2))
        ax.plot([7.82, 8.05], [y - 0.02, y + 0.16], color=CORAL, lw=1.5)
    ax.text(8.55, 2.9, "measure", color=CORAL, rotation=90, va="center")
    ax.set_xlim(0, 10)
    ax.set_ylim(0.6, 4.4)
    finish(fig, ax, key)


def draw_automaton(key, title):
    fig, ax = base_figure(title)
    pos = {"A": (1.5, 2.5), "B": (4, 3.7), "C": (4, 1.3), "D": (6.8, 2.5), "HALT": (8.7, 2.5)}
    edges = [("A", "B", "1/R"), ("B", "C", "0/L"), ("C", "A", "1/R"),
             ("B", "D", "1/R"), ("C", "D", "0/R"), ("D", "HALT", "1/–")]
    for a, b, label in edges:
        ax.add_patch(FancyArrowPatch(pos[a], pos[b], arrowstyle="-|>", mutation_scale=14,
                                     color=GRID, lw=2, connectionstyle="arc3,rad=0.08"))
        mx, my = (np.array(pos[a]) + np.array(pos[b])) / 2
        ax.text(mx, my + 0.18, label, color=INK, fontsize=10)
    for label, (x, y) in pos.items():
        color = CORAL if label == "HALT" else TEAL
        ax.add_patch(Circle((x, y), 0.48, facecolor=color, edgecolor="white", lw=2))
        ax.text(x, y, label, ha="center", va="center", color="white", fontweight="bold")
    ax.text(5, 0.35, "How many marks can a small machine write before it halts?",
            ha="center", color=INK, fontsize=14)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 4.6)
    finish(fig, ax, key)


def draw_hyperbolic(key, title, tiled=False):
    fig, ax = base_figure(title)
    ax.add_patch(Circle((0, 0), 1, fill=False, edgecolor=INK, lw=3))
    r = rng_for(key)
    if tiled:
        for ring in [0.25, 0.48, 0.69, 0.84]:
            count = int(8 + ring * 20)
            for j in range(count):
                angle = 2 * np.pi * j / count + ring
                size = 150 * (1 - ring) + 18
                ax.scatter([ring * np.cos(angle)], [ring * np.sin(angle)], s=size,
                           color=[TEAL, BLUE, CORAL, GOLD][j % 4], alpha=0.82,
                           edgecolor="white", lw=0.5)
        for theta in np.linspace(0, 2 * np.pi, 11, endpoint=False):
            ax.plot([0, 0.96 * np.cos(theta)], [0, 0.96 * np.sin(theta)],
                    color=GRID, lw=0.7, zorder=0)
    else:
        for _ in range(11):
            center = r.uniform(-0.65, 0.65, 2)
            rad = r.uniform(0.25, 0.85)
            theta = np.linspace(0, 2 * np.pi, 300)
            x = center[0] + rad * np.cos(theta)
            y = center[1] + rad * np.sin(theta)
            mask = x**2 + y**2 < 0.99
            ax.plot(np.where(mask, x, np.nan), np.where(mask, y, np.nan),
                    color=r.choice([TEAL, BLUE, CORAL]), lw=1.6, alpha=0.85)
    ax.set_xlim(-1.08, 1.08)
    ax.set_ylim(-1.08, 1.08)
    ax.set_aspect("equal")
    finish(fig, ax, key)


def draw_timetable(key, title):
    fig, ax = base_figure(title)
    days = ["Mon", "Tue", "Wed", "Thu", "Fri"]
    times = ["09", "10", "11", "13", "14", "15"]
    r = rng_for(key)
    colors = [BLUE, TEAL, CORAL, GOLD, "#8e7dbe"]
    for i, day in enumerate(days):
        ax.text(1.9 + i * 1.45, 4.55, day, ha="center", color=INK, fontweight="bold")
    for j, time in enumerate(times):
        ax.text(0.85, 3.95 - j * 0.62, time, ha="center", va="center", color=INK)
    for i in range(5):
        for j in range(6):
            ax.add_patch(Rectangle((1.2 + i * 1.45, 3.65 - j * 0.62), 1.35, 0.54,
                                   facecolor="white", edgecolor=GRID))
    for label, i, j, span, c in [("MATH", 0, 0, 2, 0), ("PHYS", 1, 2, 1, 1),
                                  ("COMP", 2, 1, 2, 2), ("LAB", 3, 3, 2, 3),
                                  ("SEMINAR", 4, 0, 1, 4), ("ALG", 1, 5, 1, 0)]:
        ax.add_patch(Rectangle((1.24 + i * 1.45, 3.69 - j * 0.62 - (span - 1) * 0.62),
                               1.27, 0.46 + (span - 1) * 0.62,
                               facecolor=colors[c], edgecolor="white", alpha=0.85))
        ax.text(1.875 + i * 1.45, 3.92 - j * 0.62 - (span - 1) * 0.31, label,
                ha="center", va="center", color="white", fontsize=9, fontweight="bold")
    ax.text(5.1, 0.2, "no clashes • balanced loads • required rooms", ha="center", color=INK)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 5)
    finish(fig, ax, key)


def draw_expander(key, title):
    fig, ax = base_figure(title)
    graph = nx.random_regular_graph(3, 22, seed=int(rng_for(key).integers(10000)))
    pos = nx.spring_layout(graph, seed=7, k=0.75)
    nx.draw_networkx_edges(graph, pos, ax=ax, edge_color=GRID, width=1.5)
    side = [n for n in graph if pos[n][0] < 0]
    nx.draw_networkx_nodes(graph, pos, nodelist=side, ax=ax, node_color=BLUE, node_size=130)
    nx.draw_networkx_nodes(graph, pos, nodelist=[n for n in graph if n not in side],
                           ax=ax, node_color=TEAL, node_size=130)
    boundary = [(a, b) for a, b in graph.edges if (a in side) != (b in side)]
    nx.draw_networkx_edges(graph, pos, edgelist=boundary, ax=ax, edge_color=CORAL, width=3)
    ax.text(0, -1.16, f"{len(boundary)} edges cross this cut", ha="center", color=CORAL, fontsize=14)
    finish(fig, ax, key)


def draw_allocation(key, title):
    fig, ax = base_figure(title)
    resources = ["room", "equipment", "time", "budget"]
    shares = np.array([[0.38, 0.27, 0.22, 0.13], [0.24, 0.29, 0.31, 0.16],
                       [0.21, 0.26, 0.25, 0.28]])
    colors = [BLUE, TEAL, GOLD, CORAL]
    y = np.arange(3)
    left = np.zeros(3)
    for j, resource in enumerate(resources):
        ax.barh(y, shares[:, j], left=left, color=colors[j], label=resource, height=0.55)
        left += shares[:, j]
    ax.set_yticks(y, ["group A", "group B", "group C"])
    ax.set_xlim(0, 1)
    ax.set_xlabel("share of available resources")
    ax.legend(frameon=False, ncol=4, loc="lower center", bbox_to_anchor=(0.5, -0.28))
    finish(fig, ax, key, True)


def draw_sperner(key, title):
    fig, ax = base_figure(title)
    n = 8
    tri = np.array([[0, 0], [1, 0], [0.5, np.sqrt(3) / 2]])
    colors = [BLUE, CORAL, TEAL]
    for i in range(n + 1):
        for j in range(n + 1 - i):
            p = (i * tri[1] + j * tri[2] + (n - i - j) * tri[0]) / n
            if i == 0:
                c = 2
            elif j == 0:
                c = 1
            elif i + j == n:
                c = 0
            else:
                c = (2 * i + 3 * j) % 3
            ax.scatter(*p, s=72, color=colors[c], edgecolor="white", zorder=3)
    ax.plot(*np.vstack([tri, tri[0]]).T, color=INK, lw=2)
    ax.text(0.5, -0.11, "a fully labelled small triangle marks a balanced choice",
            ha="center", color=INK, fontsize=13)
    ax.set_xlim(-0.08, 1.08)
    ax.set_ylim(-0.15, 0.98)
    ax.set_aspect("equal")
    finish(fig, ax, key)


def draw_protein(key, title):
    fig, ax = base_figure(title)
    x = np.linspace(-3, 3, 300)
    y = np.linspace(-2.2, 2.2, 220)
    xx, yy = np.meshgrid(x, y)
    z = (-1.1 * np.exp(-((xx + 1.2) ** 2 + (yy - 0.5) ** 2) / 0.7)
         - 0.8 * np.exp(-((xx - 1.1) ** 2 + (yy + 0.7) ** 2) / 0.4)
         + 0.16 * (xx**2 + yy**2))
    contour = ax.contourf(xx, yy, z, levels=18, cmap="viridis")
    path_x = np.linspace(-2.6, 1.1, 90)
    path_y = 0.5 * np.sin(2.2 * path_x) * np.exp(-0.15 * (path_x + 1.2) ** 2)
    ax.plot(path_x, path_y, color="white", lw=3)
    ax.scatter([path_x[-1]], [path_y[-1]], s=100, color=CORAL, edgecolor="white")
    ax.text(1.25, -0.55, "low-energy fold", color="white", fontsize=13)
    fig.colorbar(contour, ax=ax, shrink=0.72, label="model energy")
    ax.set_xlabel("shape coordinate 1")
    ax.set_ylabel("shape coordinate 2")
    finish(fig, ax, key, True)


def draw_colouring(key, title, annealing=False):
    fig, ax = base_figure(title)
    graph = nx.watts_strogatz_graph(16, 4, 0.25, seed=4)
    pos = nx.spring_layout(graph, seed=5)
    greedy = nx.coloring.greedy_color(graph, strategy="largest_first")
    palette = [BLUE, CORAL, TEAL, GOLD, "#8e7dbe"]
    nx.draw_networkx_edges(graph, pos, ax=ax, edge_color=GRID, width=1.4)
    nx.draw_networkx_nodes(graph, pos, ax=ax,
                           node_color=[palette[greedy[n] % len(palette)] for n in graph],
                           node_size=270, edgecolors="white", linewidths=1.2)
    if annealing:
        x = np.linspace(-1.2, 1.2, 160)
        curve = 0.45 * (x**2 - 0.7) ** 2 - 0.15 * x
        ax.plot(1.15 + x * 0.42, -1.05 + curve * 0.8, color=INK, lw=2.5)
        ax.scatter([1.15 + x[np.argmin(curve)] * 0.42],
                   [-1.05 + curve.min() * 0.8], color=CORAL, s=65)
    finish(fig, ax, key)


def draw_graph_state(key, title):
    fig, ax = base_figure(title)
    graph = nx.grid_2d_graph(4, 4)
    pos = {node: (node[0], node[1]) for node in graph}
    nx.draw_networkx_edges(graph, pos, ax=ax, edge_color=GRID, width=2)
    nx.draw_networkx_nodes(graph, pos, ax=ax, node_color=TEAL, node_size=310,
                           edgecolors="white", linewidths=2)
    for node in [(1, 1), (2, 2)]:
        ax.add_patch(Circle(pos[node], 0.26, fill=False, edgecolor=CORAL, lw=3))
        ax.text(*pos[node], "K", ha="center", va="center", color="white", fontweight="bold")
    ax.text(4.1, 2.25, r"$K_v=X_v\prod_{u\in N(v)} Z_u$", color=INK, fontsize=18)
    ax.text(4.1, 1.55, "one stabilizer couples", color=INK, fontsize=14)
    ax.text(4.1, 1.15, "a vertex to its neighbours", color=CORAL, fontsize=14)
    ax.set_xlim(-0.5, 7.3)
    ax.set_ylim(-0.5, 3.6)
    finish(fig, ax, key)


def draw_grover(key, title):
    fig, ax = base_figure(title)
    x = np.arange(16)
    initial = np.full(16, 0.25)
    final = np.full(16, 0.08)
    final[11] = 0.91
    ax.bar(x - 0.18, initial, width=0.34, color=GRID, label="before")
    ax.bar(x + 0.18, final, width=0.34, color=[CORAL if i == 11 else TEAL for i in x], label="after")
    ax.annotate("marked state", (11.18, 0.91), xytext=(8.2, 1.02),
                arrowprops={"arrowstyle": "->", "color": CORAL}, color=CORAL)
    ax.set_xlabel("basis state")
    ax.set_ylabel("amplitude")
    ax.set_xticks(x[::2])
    ax.legend(frameon=False)
    finish(fig, ax, key, True)


def draw_polyhedron(key, title):
    fig, ax = base_figure(title, "3d")
    phi = (1 + np.sqrt(5)) / 2
    verts = []
    for a in [-1, 1]:
        for b in [-phi, phi]:
            verts.extend([(0, a, b), (a, b, 0), (b, 0, a)])
    verts = np.array(verts)
    distances = np.linalg.norm(verts[:, None, :] - verts[None, :, :], axis=2)
    edge = np.min(distances[distances > 1e-6])
    for i in range(len(verts)):
        for j in range(i + 1, len(verts)):
            if abs(distances[i, j] - edge) < 1e-5:
                ax.plot(*zip(verts[i], verts[j]), color=BLUE, lw=2)
    ax.scatter(verts[:, 0], verts[:, 1], verts[:, 2], s=95, color=CORAL, edgecolor="white")
    ax.set_box_aspect((1, 1, 1))
    ax.view_init(18, 35)
    finish(fig, ax, key)


def draw_geodesic_surface(key, title):
    fig, ax = base_figure(title, "3d")
    x = np.linspace(-2.5, 2.5, 80)
    y = np.linspace(-2.5, 2.5, 80)
    xx, yy = np.meshgrid(x, y)
    zz = 0.28 * (xx**2 - yy**2)
    ax.plot_surface(xx, yy, zz, cmap="Blues", alpha=0.62, linewidth=0)
    t = np.linspace(-2.2, 2.2, 180)
    gy = 0.65 * np.sin(0.9 * t)
    gz = 0.28 * (t**2 - gy**2) + 0.05
    ax.plot(t, gy, gz, color=CORAL, lw=4)
    ax.scatter([t[0], t[-1]], [gy[0], gy[-1]], [gz[0], gz[-1]], color=INK, s=55)
    ax.view_init(28, -60)
    finish(fig, ax, key)


def draw_manifold(key, title):
    fig, ax = base_figure(title, "3d")
    t = 1.5 * np.pi * (1 + 2 * np.linspace(0, 1, 500))
    h = 18 * np.linspace(0, 1, 500)
    x = t * np.cos(t)
    z = t * np.sin(t)
    points = np.column_stack([x, h, z])
    ax.scatter(points[:, 0], points[:, 1], points[:, 2], c=t, cmap="viridis", s=8)
    ax.plot(points[70:440, 0], points[70:440, 1], points[70:440, 2], color=CORAL, lw=2.5)
    ax.set_xlabel("observed feature")
    ax.set_ylabel("observed feature")
    ax.set_zlabel("hidden geometry")
    ax.view_init(18, -64)
    finish(fig, ax, key, True)


def draw_markov(key, title):
    fig, ax = base_figure(title)
    nodes = {
        "(1,1,1)": (5, 4.1), "(1,1,2)": (3.0, 3.1), "(1,2,5)": (1.6, 1.9),
        "(1,5,13)": (0.7, 0.7), "(2,5,29)": (2.5, 0.7),
        "(1,2,5)′": (4.3, 1.9), "(1,5,13)′": (3.7, 0.7), "(2,5,29)′": (5.0, 0.7),
        "(1,1,2)′": (7.0, 3.1), "(1,2,5)″": (6.0, 1.9), "(2,5,29)″": (5.9, 0.7),
        "(1,5,13)″": (7.2, 0.7), "(1,2,5)‴": (8.4, 1.9), "(2,5,29)‴": (8.3, 0.7),
        "(1,5,13)‴": (9.3, 0.7),
    }
    for label, (x, y) in nodes.items():
        if y > 0.7:
            children = sorted([(l, p) for l, p in nodes.items() if abs(p[1] - (y - 1.2)) < 0.1],
                              key=lambda item: abs(item[1][0] - x))[:2]
            for _, p in children:
                if abs(p[0] - x) < 2.2:
                    ax.plot([x, p[0]], [y, p[1]], color=GRID, lw=1.6)
    for label, (x, y) in nodes.items():
        ax.add_patch(Circle((x, y), 0.36, color=TEAL if y > 1 else BLUE, alpha=0.9))
        ax.text(x, y, label.rstrip("′″‴"), ha="center", va="center", color="white", fontsize=8)
    ax.text(5, 0.12, r"$x^2+y^2+z^2=3xyz$", ha="center", color=CORAL, fontsize=18)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 4.7)
    finish(fig, ax, key)


def draw_nbody(key, title):
    fig, ax = base_figure(title)
    ax.add_patch(Circle((5, 2.35), 1.85, facecolor=PALE, edgecolor=INK, lw=2))
    colors = [CORAL, TEAL, BLUE]
    phases = [0, 2.05, 4.2]
    for color, phase in zip(colors, phases):
        t = np.linspace(0, 2.2 * np.pi, 260)
        x = 5 + 1.62 * np.cos(t + phase)
        y = 2.35 + 0.72 * np.sin(2 * t + phase) * np.sqrt(np.maximum(0.08, 1 - 0.18 * np.cos(t)))
        ax.plot(x, y, color=color, lw=2.2, alpha=0.8)
        ax.scatter([x[-1]], [y[-1]], color=color, s=90, edgecolor="white", zorder=3)
    ax.text(7.4, 3.8, "curvature changes", color=INK, fontsize=14)
    ax.text(7.4, 3.4, "every trajectory", color=CORAL, fontsize=14)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 4.7)
    ax.set_aspect("equal")
    finish(fig, ax, key)


def draw_decoder(key, title):
    fig, ax = base_figure(title)
    graph = nx.grid_2d_graph(5, 5)
    pos = {n: (n[0] * 0.65, n[1] * 0.65) for n in graph}
    nx.draw_networkx_edges(graph, pos, ax=ax, edge_color=GRID, width=1.2)
    nx.draw_networkx_nodes(graph, pos, ax=ax, node_color=BLUE, node_size=95)
    errors = [(1, 2), (3, 1), (4, 3)]
    nx.draw_networkx_nodes(graph, pos, nodelist=errors, ax=ax, node_color=CORAL, node_size=150)
    ax.add_patch(FancyArrowPatch((3.5, 1.3), (4.5, 1.3), arrowstyle="-|>",
                                 mutation_scale=18, color=INK, lw=2.5))
    layers = [(5.2, 6), (6.4, 4), (7.6, 2)]
    for (x0, n), (x1, m) in zip(layers[:-1], layers[1:]):
        for y0 in np.linspace(0.45, 3.7, n):
            for y1 in np.linspace(1.0, 3.2, m):
                ax.plot([x0, x1], [y0, y1], color=GRID, lw=0.6)
    for x0, n in layers:
        ax.scatter(np.full(n, x0), np.linspace(0.45 if n == 6 else 1.0, 3.7 if n == 6 else 3.2, n),
                   s=105, color=TEAL if n > 2 else CORAL, edgecolor="white", zorder=3)
    ax.text(8.35, 2.1, "correction", color=INK, fontsize=14)
    ax.set_xlim(-0.4, 9.5)
    ax.set_ylim(-0.2, 4.5)
    finish(fig, ax, key)


def draw_programmable(key, title):
    fig, ax = base_figure(title)
    expressions = [("x", 1.0, 2.5, BLUE), ("sin(x)", 3.0, 3.4, TEAL),
                   ("x²", 3.0, 1.6, GOLD), ("sin(x)+x²", 5.5, 2.5, CORAL),
                   ("plot", 8.1, 2.5, INK)]
    for label, x, y, color in expressions:
        ax.add_patch(Rectangle((x - 0.65, y - 0.38), 1.3, 0.76,
                               facecolor=color, edgecolor="white", lw=2, alpha=0.88))
        ax.text(x, y, label, ha="center", va="center", color="white", fontweight="bold")
    for a, b in [((1.65, 2.5), (2.35, 3.25)), ((1.65, 2.5), (2.35, 1.75)),
                 ((3.65, 3.4), (4.85, 2.65)), ((3.65, 1.6), (4.85, 2.35)),
                 ((6.15, 2.5), (7.45, 2.5))]:
        ax.add_patch(FancyArrowPatch(a, b, arrowstyle="-|>", mutation_scale=14, color=GRID, lw=2))
    ax.text(5, 0.45, "mathematical objects become inspectable and reusable", ha="center", color=INK)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 4.6)
    finish(fig, ax, key)


def draw_annealing(key, title):
    fig, ax = base_figure(title)
    x = np.linspace(-3, 3, 600)
    energy = 0.09 * x**6 - 0.75 * x**4 + 1.35 * x**2 + 0.25 * np.sin(5 * x)
    ax.plot(x, energy, color=INK, lw=3)
    steps_x = np.array([-2.7, -2.25, -1.78, -1.42, -1.18, -0.98, -0.83])
    steps_y = np.interp(steps_x, x, energy) + np.array([1.3, 1.0, 0.7, 0.45, 0.28, 0.13, 0])
    ax.plot(steps_x, steps_y, "o-", color=CORAL, lw=2.4, ms=7)
    ax.annotate("low-energy solution", (steps_x[-1], steps_y[-1]),
                xytext=(0.4, energy.min() + 0.8),
                arrowprops={"arrowstyle": "->", "color": TEAL}, color=TEAL, fontsize=14)
    ax.set_xlabel("candidate configuration")
    ax.set_ylabel("QUBO energy")
    finish(fig, ax, key, True)


def draw_qec(key, title, game=False):
    fig, ax = base_figure(title)
    for i in range(7):
        for j in range(5):
            color = PALE if (i + j) % 2 else "#cce7e3"
            ax.add_patch(Rectangle((0.8 + i * 0.65, 0.75 + j * 0.65), 0.61, 0.61,
                                   facecolor=color, edgecolor="white"))
            ax.scatter([1.105 + i * 0.65], [1.055 + j * 0.65], s=34, color=BLUE)
    error_nodes = [(2, 1), (3, 1), (4, 2), (4, 3)]
    for i, j in error_nodes:
        ax.scatter([1.105 + i * 0.65], [1.055 + j * 0.65], s=125, color=CORAL,
                   edgecolor="white", zorder=4)
    if game:
        ax.text(6.1, 3.25, "syndrome", color=INK, fontsize=14)
        ax.text(6.1, 2.75, "find a correction", color=CORAL, fontsize=18, fontweight="bold")
        ax.text(6.1, 2.15, "without creating", color=INK, fontsize=14)
        ax.text(6.1, 1.7, "a logical error", color=TEAL, fontsize=18, fontweight="bold")
    else:
        ax.add_patch(FancyArrowPatch((5.7, 2.2), (7.2, 2.2), arrowstyle="-|>",
                                     mutation_scale=18, color=INK, lw=2.5))
        ax.text(8.05, 2.65, "decode", color=CORAL, ha="center", fontsize=15)
        ax.text(8.05, 2.05, "correct", color=TEAL, ha="center", fontsize=15)
        ax.text(8.05, 1.45, "recover |ψ⟩", color=INK, ha="center", fontsize=15)
    ax.set_xlim(0, 10)
    ax.set_ylim(0, 4.7)
    finish(fig, ax, key)


def draw_rhythm(key, title):
    fig, ax = base_figure(title)
    r = rng_for(key)
    times = np.sort(r.uniform(0.4, 9.3, 42))
    lanes = r.integers(0, 4, len(times))
    colors = [BLUE, TEAL, CORAL, GOLD]
    for lane in range(4):
        ax.axhline(lane, color=GRID, lw=1)
    for t, lane in zip(times, lanes):
        ax.add_patch(Rectangle((t - 0.06, lane - 0.27), 0.12, 0.54,
                               facecolor=colors[lane], edgecolor="none"))
    density = np.convolve(np.histogram(times, bins=60, range=(0, 10))[0], np.ones(7), mode="same")
    ax.plot(np.linspace(0, 10, 60), 3.55 + 0.09 * density, color=INK, lw=2.5)
    ax.text(8.1, 4.05, "difficulty signal", color=INK)
    ax.set_xlim(0, 10)
    ax.set_ylim(-0.55, 4.65)
    ax.set_xlabel("song time")
    ax.set_yticks(range(4), ["lane 1", "lane 2", "lane 3", "lane 4"])
    finish(fig, ax, key, True)


def draw_rubik(key, title):
    fig, ax = base_figure(title)
    palette = [BLUE, CORAL, GOLD, TEAL, "#f4f1de", "#8e7dbe"]
    for face, ox, oy in [(0, 1.2, 1.4), (1, 4.05, 1.4), (2, 6.9, 1.4)]:
        for i in range(3):
            for j in range(3):
                color = palette[(face * 2 + i + 2 * j) % len(palette)]
                ax.add_patch(Rectangle((ox + i * 0.72, oy + j * 0.72), 0.66, 0.66,
                                       facecolor=color, edgecolor=INK, lw=1))
        if face < 2:
            ax.add_patch(FancyArrowPatch((ox + 2.25, 2.45), (ox + 2.75, 2.45),
                                         arrowstyle="-|>", mutation_scale=14, color=CORAL, lw=2))
    ax.text(2.25, 0.75, "state", ha="center", color=INK)
    ax.text(5.1, 0.75, "transition", ha="center", color=INK)
    ax.text(8.0, 0.75, "encoded tape", ha="center", color=INK)
    ax.set_xlim(0, 10)
    ax.set_ylim(0.3, 4.3)
    finish(fig, ax, key)


def draw_building(key, title):
    fig, ax = base_figure(title)
    hours = np.linspace(0, 24, 241)
    outside = 16 + 7 * np.sin((hours - 7) * np.pi / 12)
    baseline = 3.2 + 2.7 * np.exp(-((hours - 14) / 3.7) ** 2)
    improved = baseline * 0.66
    ax2 = ax.twinx()
    ax.plot(hours, baseline, color=CORAL, lw=3, label="baseline energy")
    ax.plot(hours, improved, color=TEAL, lw=3, label="improved envelope")
    ax.fill_between(hours, improved, baseline, color=GOLD, alpha=0.25, label="saving")
    ax2.plot(hours, outside, color=BLUE, lw=1.8, ls="--", alpha=0.7, label="outside temperature")
    ax.set_xlabel("hour")
    ax.set_ylabel("energy demand")
    ax2.set_ylabel("outside temperature")
    ax.legend(frameon=False, loc="upper left")
    finish(fig, ax, key, True)


def draw_smart_grid(key, title):
    fig, ax = base_figure(title)
    graph = nx.random_geometric_graph(18, radius=0.42, seed=7)
    pos = nx.get_node_attributes(graph, "pos")
    source = min(graph, key=lambda n: pos[n][0])
    nx.draw_networkx_edges(graph, pos, ax=ax, edge_color=GRID, width=1.5)
    nx.draw_networkx_nodes(graph, pos, ax=ax, node_color=TEAL, node_size=150,
                           edgecolors="white")
    nx.draw_networkx_nodes(graph, pos, nodelist=[source], ax=ax, node_color=CORAL,
                           node_size=260, edgecolors="white")
    path_lengths = nx.single_source_shortest_path_length(graph, source)
    far = max(path_lengths, key=path_lengths.get)
    route = nx.shortest_path(graph, source, far)
    nx.draw_networkx_edges(graph, pos, edgelist=list(zip(route, route[1:])), ax=ax,
                           edge_color=CORAL, width=3)
    ax.text(0.5, -0.09, "generation • storage • demand • optimized flow",
            ha="center", color=INK, transform=ax.transAxes, fontsize=14)
    finish(fig, ax, key)


def draw_gamma(key, title):
    fig, ax = base_figure(title)
    ax.axhline(0, color=INK, lw=2)
    for center, radius, color in [(-2, 1.1, BLUE), (0, 1.0, TEAL), (2, 1.1, CORAL)]:
        theta = np.linspace(0, np.pi, 220)
        ax.plot(center + radius * np.cos(theta), radius * np.sin(theta), color=color, lw=3)
    for x in [-3.1, -1, 1, 3.1]:
        ax.plot([x, x], [0, 3.7], color=GRID, lw=1.4)
    t = np.linspace(-2.9, 2.9, 300)
    geodesic = 0.55 + 1.9 / (1 + np.exp(-2.2 * t)) + 0.18 * np.sin(3 * t)
    ax.plot(t, geodesic, color=INK, lw=3)
    ax.scatter([t[0], t[-1]], [geodesic[0], geodesic[-1]], color=GOLD, s=75)
    ax.text(0, 3.25, "quotient identifications close the path", ha="center", color=CORAL)
    ax.set_xlim(-3.5, 3.5)
    ax.set_ylim(-0.1, 3.8)
    finish(fig, ax, key)


def draw_project(key: str, title: str, kind: str):
    if kind == "series":
        draw_series(key, title)
    elif kind == "surface":
        draw_surface(key, title)
    elif kind == "captcha":
        draw_captcha(key, title)
    elif kind == "polygon":
        draw_polygon(key, title)
    elif kind == "curve":
        draw_curve(key, title)
    elif kind == "mnist":
        draw_mnist(key, title)
    elif kind == "sir":
        draw_sir(key, title)
    elif kind == "skipping":
        draw_skipping(key, title)
    elif kind == "projectile":
        draw_projectile(key, title)
    elif kind == "orbit":
        draw_orbit(key, title)
    elif kind == "predator-prey":
        draw_predator_prey(key, title)
    elif kind == "spring":
        draw_spring(key, title)
    elif kind == "double-spring":
        draw_spring(key, title, True)
    elif kind == "circuit":
        draw_circuit(key, title)
    elif kind == "lanchester":
        draw_lanchester(key, title)
    elif kind == "aes":
        draw_aes(key, title)
    elif kind == "rotation":
        draw_rotation(key, title)
    elif kind == "pca":
        draw_pca(key, title)
    elif kind == "bloch":
        draw_bloch(key, title)
    elif kind == "period":
        draw_period(key, title)
    elif kind == "bb84":
        draw_bb84(key, title)
    elif kind == "number-sets":
        draw_number_sets(key, title)
    elif kind in {"bv", "simon", "code-blocks"}:
        draw_quantum_circuit(key, title, kind)
    elif kind == "automaton":
        draw_automaton(key, title)
    elif kind == "hyperbolic":
        draw_hyperbolic(key, title)
    elif kind == "hyperbolic-tiling":
        draw_hyperbolic(key, title, True)
    elif kind == "timetable":
        draw_timetable(key, title)
    elif kind == "expander":
        draw_expander(key, title)
    elif kind == "allocation":
        draw_allocation(key, title)
    elif kind == "sperner":
        draw_sperner(key, title)
    elif kind == "protein":
        draw_protein(key, title)
    elif kind == "colouring":
        draw_colouring(key, title)
    elif kind == "annealing":
        if "graph-colouring" in key:
            draw_colouring(key, title, True)
        else:
            draw_annealing(key, title)
    elif kind == "graph-state":
        draw_graph_state(key, title)
    elif kind == "grover":
        draw_grover(key, title)
    elif kind == "polyhedron":
        draw_polyhedron(key, title)
    elif kind == "geodesic-surface":
        draw_geodesic_surface(key, title)
    elif kind == "manifold":
        draw_manifold(key, title)
    elif kind == "markov":
        draw_markov(key, title)
    elif kind == "nbody":
        draw_nbody(key, title)
    elif kind == "decoder":
        draw_decoder(key, title)
    elif kind == "programmable":
        draw_programmable(key, title)
    elif kind == "qec":
        draw_qec(key, title)
    elif kind == "qec-game":
        draw_qec(key, title, True)
    elif kind == "rhythm":
        draw_rhythm(key, title)
    elif kind == "rubik":
        draw_rubik(key, title)
    elif kind == "building":
        draw_building(key, title)
    elif kind == "smart-grid":
        draw_smart_grid(key, title)
    elif kind == "gamma":
        draw_gamma(key, title)
    else:
        raise ValueError(f"Unknown visual kind: {kind}")


def copy_source_images():
    for output_name, source in SOURCE_IMAGES.items():
        if not source.exists():
            print(f"Skipping missing source image: {source}")
            continue
        target = OUT / output_name
        try:
            image = Image.open(source)
            image.load()
            if image.mode not in {"RGB", "RGBA"}:
                image = image.convert("RGB")
            if max(image.size) > 1800:
                image.thumbnail((1800, 1800), Image.Resampling.LANCZOS)
            image.save(target, optimize=True)
        except Exception:
            shutil.copy2(source, target)


def main():
    OUT.mkdir(parents=True, exist_ok=True)
    requested = set(sys.argv[1:])
    for key, (title, kind) in PROJECTS.items():
        if requested and key not in requested:
            continue
        draw_project(key, title, kind)
    if not requested:
        copy_source_images()
    count = len(requested) if requested else len(PROJECTS)
    print(f"Created {count} project illustrations in {OUT}")


if __name__ == "__main__":
    main()
