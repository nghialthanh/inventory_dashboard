# Issue tracker: Local Markdown

Issues and specs for this repo live as markdown files in `.scratch/`.

## Conventions

- One feature per directory: `.scratch/<feature-slug>/`
- The spec is `.scratch/<feature-slug>/spec.md`
- Implementation issues are one file per ticket at `.scratch/<feature-slug>/issues/<NN>-<slug>.md`, numbered from `01`, never a single combined tickets file
- Triage state is recorded as a `Status:` line near the top of each issue file (see `triage-labels.md` for the role strings)
- Comments and conversation history append to the bottom of the file under a `## Comments` heading

## When a skill says "publish to the issue tracker"

Create a new file under `.scratch/<feature-slug>/` (creating the directory if needed).

## When a skill says "fetch the relevant ticket"

Read the file at the referenced path. The user will normally pass the path or the issue number directly.

## Wayfinding

`/wayfinder` stores a map at `.scratch/<effort>/map.md` and one child per question at `.scratch/<effort>/issues/NN-<slug>.md`.

- `Type:` is `research`, `prototype`, `grilling`, or `task`. `Status:` is `claimed` or `resolved`.
- `Blocked by: NN, NN` — unblocked when every listed file is `resolved`.
- Frontier: open, unblocked, unclaimed files; lowest number wins.
- Claim by setting `Status: claimed` before work. Resolve by appending `## Answer`, setting `Status: resolved`, and adding a gist plus link to the map's Decisions-so-far.
