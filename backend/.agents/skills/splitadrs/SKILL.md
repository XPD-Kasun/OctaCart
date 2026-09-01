---
name: splitadrs
description: Split a given monolithic ADR which is authored by the user to a bunch of ADRs with proper formatting.
argument-hint: bounded context name
disable-model-invocation: true
license: Apache 2.0 - NOAI
---

# Split the ADR which the architect authored.

bc=bounded context folder name from $ARGUMENTS

Your task is to read any adr file at `.agents/spec/{bc}/adr.md` and split to set of ADRs using following workflow.

## Workflow

1. Confirm whether the above file and print it to the user and ask the user to confirm.
2. If confirmed read the file.
3. Identify the separate points made by the architect and create a file for each point using `ADRXXX.md` format. XXX is a number in the order of points. Eg: ADR001.md,ADR002.md 
4. For each point write the ADR to the file using following template.
5. Template:

```markdown
Title
State the decision itself, not a description of the topic. For example, "Use Qwen2.5 1.5B Instruct for on-device translation".

The Problem
One sentence stating what we are trying to solve.

Options Considered
The alternatives that were on the table, with the selected option in bold.

Rationale
Only the decisive factors that led to the choice, not an exhaustive list of every pro and con.

Notes
Optional. Any additional context worth recording, such as constraints, assumptions, or links.
```