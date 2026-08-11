# Project instructions

## Comments

Keep code comments lean. Default to no comment; add one only when the WHY is
truly non-obvious (a hidden constraint, a subtle invariant, a workaround, a
surprising precondition) — never to restate WHAT the code does. When a
comment would run more than 2-3 lines, that's a signal to cut it down to the
one sentence that matters, or move the fuller reasoning out of the code
entirely (PR description, review doc, tracker) rather than growing it in
place.

This project has many more commands left to implement using the same
schema-driven codegen pattern — comment bloat compounds across every one of
them, so err on the side of trimming even when a longer explanation feels
helpful in the moment.
