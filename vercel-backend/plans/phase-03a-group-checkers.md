# Phase 03a: Base Group Checkers
Status: ✅ Completed
Dependencies: Phase 03

## Objective
Implement basic group checkers for vaccines with multiple alternative courses (e.g., Rota) or age-dependent alternative courses (e.g., Japanese Encephalitis, Hepatitis A).

## Requirements
### Functional
- [x] Implement `checkAlternativeCoursesGroup` logic for Rota.
    - Support finding the appropriate course based on age.
    - Support checking if the max completion age is exceeded.
- [x] Implement `checkAlternativeCoursesAgeRangeGroup` logic for JE_Group and HepA.
    - Warn when switching between courses (e.g., Jevax to Imojev).
    - Handle booster recommendations for Jevax.
- [x] Integrate these checkers into `engine.go`.

## Files to Create/Modify
- `internal/analyzer/group_alternative.go` - New file for alternative group logic.
- `internal/analyzer/engine.go` - Register the new checkers.
- `internal/analyzer/group_alternative_test.go` - Tests for the group checkers.

## Test Criteria
- [x] **Test Rota Multiple Courses**: A 2-month-old infant with no records should receive recommendations for *both* Rotarix and Rotateq as options.
- [x] **Test Rota Too Old**: An 8-month-old infant with no records should be marked as too old to start Rota.
- [x] **Test JE Switch Warning**: A patient with 1 Jevax and 1 Imojev should trigger a mixed course warning.

---
Next Phase: [Phase 03b: Specialized Checkers](./phase-03b-specialized-checkers.md)
