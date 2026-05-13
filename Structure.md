# Project Structure - TraCuuTiemChung

This document outlines the directory structure and the purpose of each module in the **TraCuuTiemChung** Android application. The project follows **Clean Architecture** principles to ensure separation of concerns and maintainability.

## Directory Tree

```text
TraCuuTiemChung/
├── app/                        # Main Android Application module
│   ├── build.gradle.kts        # Module-level build configuration
│   └── src/
│       ├── main/
│       │   ├── java/com/tracuutiemchung/app/
│       │   │   ├── core/           # Common utilities and base classes
│       │   │   ├── data/           # Data Layer: Repositories, API, Storage
│       │   │   │   ├── credentials/ # Secure credential management (DataStore + Crypto)
│       │   │   │   ├── model/       # Data transfer objects (DTOs)
│       │   │   │   ├── portal/      # VNCDC Portal client and HTML parsers (Jsoup)
│       │   │   │   └── rules/       # Business logic rules for data mapping
│       │   │   ├── domain/         # Domain Layer: Use Cases and Business Entities
│       │   │   │   ├── analyzer/    # Logic for analyzing vaccination records
│       │   │   │   └── usecase/     # Application-specific business logic
│       │   │   ├── ui/             # UI Layer: Screens, ViewModels, Theme
│       │   │   │   ├── login/       # Login screen and authentication UI
│       │   │   │   ├── lookup/      # Phone-based lookup screen
│       │   │   │   ├── result/      # Vaccination history display
│       │   │   │   └── theme/       # Design system (Compose Material 3)
│       │   │   └── MainActivity.kt # Entry point of the app
│       │   └── res/                # Android resources (drawables, layouts, strings)
│       └── test/               # Unit tests for business logic
│           └── java/com/tracuutiemchung/app/
│               ├── data/           # Tests for data parsing and storage
│               └── domain/         # Tests for use cases
├── gradle/                     # Gradle wrapper and build scripts
├── build.gradle.kts            # Project-level build configuration
├── settings.gradle.kts         # Project settings and module inclusion
└── README.md                   # Project overview and setup guide
```

## Module Responsibilities

### 1. Data Layer (`data/`)
- **Portal Client**: Handles HTTP requests to the VNCDC portal using OkHttp.
- **Parser**: Uses Jsoup to extract vaccination data from raw HTML responses.
- **Credentials**: Implements secure storage for user credentials using Android DataStore and Keystore-backed encryption (`CredentialCrypto`).
- **Models**: Defines JSON-serializable data classes for network and local storage.

### 2. Domain Layer (`domain/`)
- **Use Cases**: Encapsulates single pieces of business logic (e.g., `LookupVaccinationByPhoneUseCase`).
- **Entities**: Pure Kotlin classes representing the core data models used by the UI.
- **Analyzer**: Contains the core logic for analyzing vaccination records against medical rules. It is split into several specialized checkers:
    - `VaccineAnalysisEngine`: The main dispatcher that coordinates the analysis process.
    - `SingleSeriesChecker`: Handles standard vaccination series and boosters.
    - `AgeDependentSeriesChecker`: Manages vaccines with regimens that change based on start age (e.g., Prevenar 13).
    - `MmrEquivalentGroupChecker`: Handles measles-containing vaccines and their complex interactions.
    - `AlternativeCoursesChecker`: Manages vaccines with multiple valid courses (e.g., Rota, JE).
    - `Specialized Checkers`: `FluGroupChecker`, `CumulativeGroupChecker`, `MeningococcalAcywGroupChecker`, `PneumococcalSpecialChecker`.
    - `VaccineDateUtils` & `AnalysisRuleUtils`: Common utilities for date arithmetic and rule validation.

### 3. UI Layer (`ui/`)
- **Compose UI**: Built entirely with Jetpack Compose for a modern, reactive interface.
- **ViewModels**: Manages UI state and communicates with the Domain layer using Coroutines.
- **Navigation**: Uses Navigation Compose to handle transitions between screens.

## Design Patterns
- **Clean Architecture**: Decouples business logic from implementation details (UI/DB).
- **MVVM (Model-View-ViewModel)**: Decouples UI from data handling.
- **Repository Pattern**: Abstracts data sources (Network/Local) from the rest of the app.
- **Dependency Injection**: Manual or Hilt-based (pending project scale).
