---
name: impl
description: Implementing any feature of this application
argument-hint: bounded context name
---

# Implementation

You are an expert go language programmer. Your task is to implement the features provided to you in the following todo list file. **MUST** follow the specification in the task list. **DO NOT** hallucinate, deviate from the specification.

## When to use this

Implementing any code in this repository. Make sure the bounded context is provided as an argument. This should map to the `.agents/progress/{bounded context}` folder. This folder should contain the todo instructions with `implementation-plan.md` for implementing the code for this agent or model. If this is not found, you might need to consider running `planimpl` skill.

## Workflow

bc should be derived by `$ARGUMENTS.`  
If stop instruction is there, you **MUST** stop. Then when user says proceed continue to next task.  
You can stop and prompt the user for any ambigous situation or any clarification needed.

1. Your file is at `.agents/progress/{bc}/implementation-plan.md`. List the `.agents/progress` directory only and find the file if not stop. You can't execute any task list in any other folder than `.agents/progress/{bc}`.
2. Then ask the user by

- Confirm: "Is this file {filepath} you asked me to implement?"
- Options: 
  - [ ] Yes  
  - [ ] No

1. If user says yes, Read the todo task specification at the path.
2. If user says no or you can't find the file tell the user and STOP.
3. Locate the first incompleted task in the list. Ignore all tasks which are marked as done. If all tasks completed STOP.
4. Perform the instruction in task then notifiying the user.
5. Mark the progress in the todo task file.
6. Move to the next task. (that is goto step 3). If all tasks completed STOP.
