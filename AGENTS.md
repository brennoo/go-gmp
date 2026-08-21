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

## API stability and code generation

The library is alpha and its public API is not stable. Make necessary breaking
changes now rather than preserving handwritten APIs that conflict with the
schema-driven design. Do not add compatibility aliases, deprecated duplicate
fields, or legacy wire models solely to avoid an API break.

For generated command families, prioritize:

- accurate GMP XML wire behavior;
- deterministic schema-derived naming and family-local generation;
- consistent generated APIs across command families;
- correct presence semantics, including omitted versus explicit-empty values;
- fail-closed handling of ambiguous or unsupported schema shapes;
- a clean, maintainable generated surface suitable for future stabilization.

Update all repository call sites, tests, fixtures, and documentation when an
adoption changes the API. Record deliberate breaking changes where relevant,
but do not compromise generator correctness or future maintainability to keep
an alpha API compatible.
