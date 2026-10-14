# ODD Roadmap

Coordinate product direction with an ODD roadmap: one living index issue, result
milestones, ODD-carrying work units, and readiness gates. This practice keeps planning
honest, reviewable, and free of invented dates.

## What it is

A single coordination index issue (a metadocument, never an implementation issue)
holds the approved order of work: a table of milestones, the ODD work units inside
each, and each unit's status and readiness conditions. Milestones group approved
outcomes, not dates. Each work unit carries its own plan in its issue: scope,
acceptance criteria, checks, and what still blocks readiness.

## Hard rules

- Approved revisions are appended to the index (v1, v2, ...), each with its own
  proposal marker; never rewrite or delete an approved revision.
- No dates and no promised versions anywhere in the index or milestone descriptions.
- Ideas without an approved milestone stay in a "Not prioritized" section and never
  create issues automatically.
- A unit is contribution-ready only when its plan is complete and its readiness
  conditions are resolved; never announce readiness while a decision is pending.
- Bugs and small tasks keep the normal issue-first flow; do not force ODD treatment
  by size.

## Workflow

1. Read the current index before planning; never contradict an approved revision.
2. Get maintainer approval for outcomes (milestones) and unit scope before creating
   anything.
3. Create outcome-named milestones and assign existing issues.
4. Create or append an index revision: objective, approved-order table
   (milestone → work units → status and readiness), how-to-contribute, out-of-flow
   section, not-prioritized section, and a proposal marker line.
5. Point any static roadmap document at the index as the living source for sequencing
   and readiness.
6. Write each work unit's full plan in its own issue before linking it from the index.

## Live reference

The practice in production: the Cortex repository's roadmap index and its operating
guide's ODD roadmap workflow section.
