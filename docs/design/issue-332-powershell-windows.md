# Scoped plan: PowerShell runtime target and Windows platform (#332)

Status: proposal for maintainer approval. This document is the scoped plan required
by roadmap policy before either phase can become a milestone in the ODD roadmap index
(#288). It creates no issues and promises no dates or versions.

Refs #332 #303 #78 #288

## Evaluation

The issue identifies two independent gaps. Both are new roadmap directions: per the
ODD roadmap hard rules they become milestones only through a maintainer-approved
revision appended to the index, and stay in "Not prioritized" until then.

Evidence from the current tree:

- Runtime targets are defined in `internal/runtimematrix/matrix.go` (currently `pi`,
  `opencode`, `claude-code`); the install path validates pairs in
  `internal/installplan` (`validPair`, `supportedLocalTarget`), with wiring across
  `internal/artifact/bundle.go` (expected runtimes), `internal/ownership`,
  `internal/skillprojection`/`internal/projection`, and `internal/runtimecompat`.
- QA routes and backends live in `internal/qapi` (`backendselect.go`, `runtime.go`,
  `probes.go`) with the binding contract used for the pi/claude/opencode backends.
- Windows-adjacent semantics are not validated today: file modes in
  `internal/atomicfile` and `internal/safepath` assume POSIX modes, and
  `internal/runtimeprobe` probes are designed around POSIX shells.
- Releases are evaluation prereleases for darwin_arm64 and linux_amd64
  (`internal/releasecatalog`); there is no windows_amd64 target.

## Phase A — Directory projection target (scoped plan, revised v2)

Outcome: PowerShell-based agent surfaces receive curated skills and QA actors by
projecting them into a user-selected directory, with the same transactional
guarantees as the existing runtimes. Revision v2 (2026-10-19, maintainer decision):
the target is a directory projection mode, not a new runtime family.

### Evaluation behind the revision

- Surface research (2026-10-19) found no dominant PowerShell-native agent runtime
  with its own skills directory: AIShell is archived since January 2026; existing
  CLIs on Windows reuse the same destinations as their POSIX RuntimeIDs; niche
  PowerShell harnesses have no skills contract.
- The open consumption contract is the Agent Skills specification
  (agentskills.io), implemented by Microsoft Agent Framework: a directory of
  `SKILL.md` files with `name`, `description`, `license`, and arbitrary `metadata`
  frontmatter — compatible with Cortex's projection, and ignoring unknown fields
  (`metadata.cortex-*` passes through).
- Copying files by hand works for landing content (the TARS precedent) but loses
  catalog provenance, manifest hashes, ownership, and transactionality — the gap
  #303 and #294 document.

### Scope

1. **Directory destination mode.** `cortex install --target-dir <path>` projects
   skills and QA actors as plain `SKILL.md` files into the chosen directory, inside
   the existing transaction (manifest, ownership, backup, rollback, readback). No
   new `RuntimeID` family and no `runtimematrix` entry: this is a destination mode,
   OS-agnostic, usable today on POSIX.
2. **Update and uninstall.** The directory target participates in `update` and
   conservative `uninstall` through the recorded ownership, without touching files
   Cortex does not own.
3. **Diagnostics.** `cortex doctor` reports the configured directory target
   (presence, writability) without version certification — there is no host binary
   to probe.
4. **QA actors.** QA actors render as skills (the four lenses are self-contained;
   already proven by the TARS vendoring). No QA backend binding is added: this
   target does not run `qa run` backends.
5. **Consumption contract.** The Agent Skills specification is documented as the
   consumption contract for the output (any host reading `SKILL.md` from
   directories can consume it), without coupling Cortex to a specific SDK.

### Resolved decisions

- Target: directory projection mode (this revision), replacing the full runtime
  target. Milestone "PowerShell runtime target" retains its outcome name; the scope
  is this mode.
- QA actor handling: render as skills; no QA backend binding.
- Emit-to-directory: it is the whole of Phase A, not an optional extra; it is also
  the natural prerequisite for Phase B on Windows.

### Readiness gate

Phase A is ready for implementation when the maintainer approves this revised
scope. Destination evidence and the former open decisions are resolved above.

## Phase B — Windows platform (scoped plan)

Outcome: the CLI installs, updates, and uninstalls with verified semantics on
windows_amd64, with release binaries available. The Phase A directory target is a
prerequisite: a Windows user installs skills into a directory, so that mode exists
and is proven on POSIX first.

### Scope

1. **Release binaries.** windows_amd64 targets in the release pipeline
   (`internal/releasecatalog` and release docs).
2. **Path and mode semantics.** Canonicalize and verify path handling, file modes,
   shadow-root validation, backup/state layout, and install state on Windows:
   `internal/atomicfile`, `internal/safepath`, `internal/filetxn`, backup journals,
   `internal/installstate`.
3. **Runtime probes.** Validate `internal/runtimeprobe` behavior on PowerShell for
   every supported runtime family.
4. **Lifecycle verification.** Isolated install/update/uninstall and readback on a
   real Windows environment, covering absent runtimes, unidentified versions, and
   known-incompatible adapters per the verification rules.

### Readiness gate

Phase B is ready for implementation when: a real Windows verification environment is
available (this is the dominant condition), the mode/backup semantics design for
Windows is written and approved, and the maintainer approves this scope.

## Sequencing

Phase A does not require Phase B (S3): the directory target lands while the CLI
remains POSIX-only, and Phase B consumes it as its install-mode prerequisite. Each
phase carries its own milestone and work units.

## Proposed roadmap revision (for issue #288, maintainer approval required)

The following revision would be appended to the ODD roadmap index once the maintainer
approves the outcomes. Milestones would be created outcome-named, issues created only
after approval, and each work unit's plan written in its own issue.

> **Approved order addition (revision v3):**
>
> | Milestone | Work / active ODD | Status and readiness |
> | --- | --- | --- |
> | [PowerShell runtime target](https://github.com/refactor-ia/cortex/milestone/4) | Scoped plan: docs/design/issue-332-powershell-windows.md Phase A, revised v2 — directory projection mode (#332) | Readiness gate resolved (revision v2, 2026-10-19): directory-mode scope approved by the maintainer; work-unit issues follow. |
> | [Windows platform] | Scoped plan: docs/design/issue-332-powershell-windows.md Phase B (#332) | Blocked on readiness gate: Windows verification environment, Windows mode/backup semantics design, maintainer scope approval. No issues created yet. |
>
> Proposal marker: `<!-- odd-roadmap: proposal=cortex-powershell-windows-v1 operation=roadmap-revision -->`

## Non-goals

- No dates and no promised versions anywhere (roadmap rule).
- No QA lifecycle changes; the directory target adds no QA backend.
- No new capability content in either phase.
- No Windows-specific behavior in Phase A: the directory target is OS-agnostic.
