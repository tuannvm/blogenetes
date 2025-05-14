---
title: "More predictable benchmarking with testing.B.Loop"
date: "2025-05-14T16:19:27-07:00"
---

The article discusses the introduction of `testing.B.Loop` in Go 1.24 as a new way to write benchmarks. It highlights the benefits of using `testing.B.Loop` over the traditional `b.N` benchmark loop, such as preventing unwanted compiler optimizations, excluding setup/cleanup code from timing, and avoiding accidental dependencies on iteration counts.

`testing.B.Loop` simplifies benchmarking by integrating timer management for setup/cleanup code and preventing dead code elimination within the loop. It also offers a one-shot ramp-up approach for more efficient benchmarking.

The article explains common pitfalls of the old benchmark loop structure and demonstrates how `testing.B.Loop` helps address these issues. It also provides guidelines on when to use `testing.B.Loop` and acknowledges contributors who provided feedback on the feature.

Overall, `testing.B.Loop` provides a faster, more accurate, and intuitive way to write benchmarks in Go.
