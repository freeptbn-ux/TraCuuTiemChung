# Supported Rule Types
This document lists the vaccination rule types supported by the Kotlin Analysis Engine. These types map directly to the `type` field in `vaccine_rules.json`.

## Core Series
- `single_series`: Standard multi-dose series (e.g., 5-in-1, Hepatitis B). Supports basic boosters.
- `series`: Alias for `single_series`.
- `single_dose_min_age`: Vaccines requiring only one dose after a minimum age (e.g., BCG, if tiêm muộn).
- `single_series_min_age`: Multi-dose series with a minimum age for the first dose.

## Advanced Regimens
- `age_dependent_series`: Regimen (number of doses/intervals) changes based on the age at the first dose (e.g., Prevenar 13, Synflorix).
- `group_alternative_courses`: Supports switching between different valid courses (e.g., JE_Group mixing, Rota 2-dose vs 3-dose).
- `group_alternative_courses_min_age`: Alternative courses with minimum age constraints.
- `group_alternative_courses_age_range`: Alternative courses where the valid course depends on current age.

## Group Interactions
- `mmr_equivalent_group`: Handles complex interactions between Measles, MR, and MMR vaccines.
- `flu_group`: Seasonal flu logic, including different regimens for children vs. adults.
- `meningococcal_acyw_group`: Logic for MenACWY vaccines (e.g., Menactra, Menveo).
- `group_cumulative_unique_doses`: Combines multiple vaccine keys to count total unique doses (e.g., 6-in-1 + 5-in-1 + HepB).
- `group_cumulative_unique_doses_min_age`: Cumulative doses with minimum age constraints.

## Specialized Logic
- `pneumococcal_special`: Specialized logic for Pneumococcal mixing and Pneumovax 23 recommendations.

## Known Limitations & Intentional Differences
- **VA-MENGOC-BC**: Includes a reverse interaction check where it's not recommended if the patient is over 24 months and has received ACYW.
- **Grace Period**: The Kotlin engine uses a configurable `GRACE_PERIOD_DAYS` (currently set to 0) for strict adherence to intervals, unlike some systems that allow 4-day grace periods.
- **Rules Not Ported**: All major rule types from the legacy Python system have been ported as of Phase 11.
