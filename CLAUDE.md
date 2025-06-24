# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is **pug** (µcompiler) - a learning-oriented compiler implementation project written in Japanese. The project aims to teach compiler construction by building from a simple lexer to a full LLVM-connected optimizing compiler through 4 distinct phases:

- **Phase 1**: Basic language processing (lexer, parser, interpreter)
- **Phase 2**: Compiler foundation (code generation, type system, control structures)  
- **Phase 3**: Optimization engine (IR, SSA, optimization passes)
- **Phase 4**: Industrial-grade compiler (LLVM integration, multi-target support)

## Project Language and Documentation

**IMPORTANT**: All documentation, commit messages, and GitHub issues must be written in Japanese (日本語) as specified in the project README.

## Project Structure (Planned)

The project follows a phase-based structure under `µcompiler/`:
- `phase1/` - Basic language processing (lexer.go, parser.go, ast.go, interpreter.go)
- `phase2/` - Compiler foundation (codegen.go, types.go, symbols.go, control.go)
- `phase3/` - IR and optimization (ir/, optimizer/, backend/)
- `phase4/` - LLVM integration (llvm/, runtime/, tools/)
- `cmd/` - CLI tools (µc, µinterp, µtools)
- `examples/` - Sample programs (.µ files)
- `benchmark/` - Performance measurement tools

## Build System (Planned)

The project plans to use a comprehensive Makefile-based build system with these key targets:

**Development Commands:**
- `make help` - Show available commands
- `make dev` - Development environment setup
- `make phase1-build` - Build Phase 1 interpreter
- `make bench-compile` - Compile time benchmarks
- `make bench-runtime` - Runtime performance benchmarks
- `make bench-vs-gcc` - Compare with GCC
- `make bench-vs-rust` - Compare with Rust

**Phase-specific targets:**
- Phase 1: `./bin/interp` - Interpreter with REPL support
- Phase 2: `./bin/pug` - Basic compiler with assembly output
- Phase 3: Optimization flags (-O0, -O1, -O2) and IR emission
- Phase 4: `--backend=llvm` with multi-target support

## Language Implementation (.µ files)

The µ language syntax is designed to be Rust-like with progressive feature addition:
- Basic types: int, float, string, bool
- Functions with `fn name(params) -> type` syntax
- Control flow: if/while/for with Rust-like syntax
- Advanced features planned: structs, generics, traits, async/await

## Performance Goals

The project has specific performance targets comparing each phase:
- Phase 1 (Interpreter): Baseline performance
- Phase 2 (Basic Compiler): 10x faster than interpreter
- Phase 3 (Optimizing): 50x faster than interpreter  
- Phase 4 (LLVM): 100x faster than interpreter, 50% smaller code size

## Development Status

**Current Status**: Very early stage - only README exists, no implementation yet.

The project is designed for AI-assisted development with:
- Comprehensive documentation in Japanese
- Phase-based progressive complexity
- Extensive testing and benchmarking infrastructure
- Clear separation of concerns across phases

## File Extensions

- `.µ` or `.dog` - Source files for the µ language
- `.go` - Implementation files (Go language)
- `*_test.go` - Test files following Go conventions

## Performance Analysis

The project includes comprehensive benchmarking comparing against industry standards (GCC, Rust, Clang) with automated CI/CD performance regression testing.