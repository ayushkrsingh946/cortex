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

## Phase A — PowerShell runtime target (scoped plan)

Outcome: a PowerShell-based agent runtime can receive curated skills and QA actors
through the same transactional guarantees as the existing three runtimes.

### Scope

1. **Runtime definition and admission.** New `RuntimeID` for the PowerShell runtime
   in `internal/runtimematrix`, wired through `installplan.validPair`,
   `internal/artifact/bundle.go` expected runtimes, `internal/ownership`,
   projection, and `internal/runtimecompat` (compatible / uncertified /
   known_incompatible admission, following the #303 pattern).
2. **Projection and destinations.** Define where skills land for a PowerShell runtime
   (root kind in `internal/skilldest`), including multi-profile handling analogous to
   the #303 proposal, and what the install plan emits per target.
3. **QA route/backend binding.** Add the QA route/backend binding for the runtime in
   `internal/qapi` following the existing contract. QA lifecycle rules do not change;
   the binding decides how (or whether) QA actors project for this runtime.
4. **Diagnostics.** `cortex doctor` reports presence, compatibility, and install
   action for the runtime like the others.
5. **Verification evidence.** Structural conformance, isolated install and readback,
   and at least one real smoke invocation for the new runtime family, recorded per
   the architecture's verification rules.

### Open decisions for the scoped plan (resolved before readiness)

- Exact destination layout and profile selection mechanism for the target runtime
  (needs the same first-hand evidence #303 collected for Hermes).
- QA actor handling: render as skill, native equivalent, or skip — decided with the
  destination evidence, and documented.
- Whether an emit-to-directory mode for vendoring is in scope or deferred.

### Readiness gate

Phase A is ready for implementation when: destination evidence exists (layout,
profiles, probe command), the four open decisions above are resolved and recorded,
and the maintainer approves this scope.

## Phase B — Windows platform (scoped plan)

Outcome: the CLI installs, updates, and uninstalls with verified semantics on
windows_amd64, with release binaries available.

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

Phase A does not require Phase B (S3): a PowerShell target can land while the CLI
remains POSIX-only. Phase B benefits from Phase A's probe validation but has no hard
dependency on it. Each phase carries its own milestone and work units.

## Proposed roadmap revision (for issue #288, maintainer approval required)

The following revision would be appended to the ODD roadmap index once the maintainer
approves the outcomes. Milestones would be created outcome-named, issues created only
after approval, and each work unit's plan written in its own issue.

> **Approved order addition (revision v3):**
>
> | Milestone | Work / active ODD | Status and readiness |
> | --- | --- | --- |
> | [PowerShell runtime target] | Scoped plan: docs/design/issue-332-powershell-windows.md Phase A (#332) | Blocked on readiness gate: destination evidence, open decisions resolved, maintainer scope approval. No issues created yet. |
> | [Windows platform] | Scoped plan: docs/design/issue-332-powershell-windows.md Phase B (#332) | Blocked on readiness gate: Windows verification environment, Windows mode/backup semantics design, maintainer scope approval. No issues created yet. |
>
> Proposal marker: `<!-- odd-roadmap: proposal=cortex-powershell-windows-v1 operation=roadmap-revision -->`

## Non-goals

- No dates and no promised versions anywhere (roadmap rule).
- No QA lifecycle changes beyond the Phase A backend binding.
- No new capability content in either phase.
